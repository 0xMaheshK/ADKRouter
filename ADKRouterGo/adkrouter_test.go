package adkrouter

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
)

// ---------------------------------------------------------------------------
// Config tests
// ---------------------------------------------------------------------------

func TestDefaultProvidersExist(t *testing.T) {
	for _, key := range []string{"openrouter", "openai", "anthropic"} {
		if _, ok := DEFAULT_PROVIDERS[key]; !ok {
			t.Errorf("expected DEFAULT_PROVIDERS to contain %q", key)
		}
	}
}

func TestOpenRouterDefaults(t *testing.T) {
	cfg, ok := DEFAULT_PROVIDERS["openrouter"]
	if !ok {
		t.Fatal("openrouter provider missing")
	}
	if cfg.BaseURL != "https://openrouter.ai/api/v1" {
		t.Errorf("openrouter base_url = %q, want %q", cfg.BaseURL, "https://openrouter.ai/api/v1")
	}
	if cfg.EnvKey != "OPENROUTER_API_KEY" {
		t.Errorf("openrouter env_key = %q, want %q", cfg.EnvKey, "OPENROUTER_API_KEY")
	}
	if cfg.NamePrefix != "openrouter" {
		t.Errorf("openrouter name_prefix = %q, want %q", cfg.NamePrefix, "openrouter")
	}
}

func TestProvidersByCategory(t *testing.T) {
	gateways := GetProvidersByCategory(CategoryGateway)
	if _, ok := gateways["openrouter"]; !ok {
		t.Error("expected openrouter in gateway category")
	}
	custom := GetProvidersByCategory(CategoryCustom)
	if _, ok := custom["custom_openai"]; !ok {
		t.Error("expected custom_openai in custom category")
	}
}

func TestGetAllProvidersCopy(t *testing.T) {
	all := GetAllProviders()
	all["openai"] = ProviderConfig{}
	if _, ok := DEFAULT_PROVIDERS["openai"]; !ok {
		t.Fatal("openai should still exist in DEFAULT_PROVIDERS")
	}
	if DEFAULT_PROVIDERS["openai"].BaseURL == "" {
		t.Error("mutating returned copy should not affect built-ins")
	}
}

func TestDetectDefaultEndpoint(t *testing.T) {
	if got := DetectDefaultEndpoint("openrouter"); got != "https://openrouter.ai/api/v1" {
		t.Errorf("DetectDefaultEndpoint(openrouter) = %q", got)
	}
	if got := DetectDefaultEndpoint("unknown"); got != "" {
		t.Errorf("DetectDefaultEndpoint(unknown) = %q, want empty", got)
	}
}

func TestGetEnvVar(t *testing.T) {
	t.Setenv("ADKROUTER_TEST_KEY", "test-value")
	if got := GetEnvVar("ADKROUTER_TEST_KEY", "default"); got != "test-value" {
		t.Errorf("GetEnvVar = %q, want test-value", got)
	}
	t.Setenv("ADKROUTER_TEST_KEY", "")
	if got := GetEnvVar("ADKROUTER_TEST_KEY", "default"); got != "default" {
		t.Errorf("GetEnvVar with empty = %q, want default", got)
	}
}

// ---------------------------------------------------------------------------
// LiteLLMProvider tests
// ---------------------------------------------------------------------------

func TestLiteLLMProviderSingleton(t *testing.T) {
	a := LiteLLMProviderInstance()
	b := LiteLLMProviderInstance()
	if a != b {
		t.Error("LiteLLMProviderInstance should return the same instance")
	}
	p := a.GetAvailableProviders()
	if !contains(p, "openrouter") || !contains(p, "openai") || !contains(p, "anthropic") {
		t.Errorf("expected core providers, got %v", p)
	}
}

func TestCreateModelWithProvider(t *testing.T) {
	p := LiteLLMProviderInstance()
	m, err := p.CreateModel("test-model",
		WithProvider("openrouter"),
		WithAPIKey("test-key"),
	)
	if err != nil {
		t.Fatalf("CreateModel: %v", err)
	}
	if m.String() != "openrouter/test-model" {
		t.Errorf("model key = %q, want %q", m.String(), "openrouter/test-model")
	}
	if m.apiKey != "test-key" {
		t.Errorf("apiKey = %q", m.apiKey)
	}
	if m.apiBase != "https://openrouter.ai/api/v1" {
		t.Errorf("apiBase = %q", m.apiBase)
	}
}

func TestCreateModelInferredProvider(t *testing.T) {
	p := LiteLLMProviderInstance()
	m, err := p.CreateModel("qwen/qwen3-preview:free", WithProvider("openrouter"), WithAPIKey("key"))
	if err != nil {
		t.Fatalf("CreateModel: %v", err)
	}
	if m.String() != "openrouter/qwen/qwen3-preview:free" {
		t.Errorf("model key = %q", m.String())
	}
}

func TestCreateModelOpenAIIsUnprefixed(t *testing.T) {
	p := LiteLLMProviderInstance()
	m, err := p.CreateModel("gpt-4o", WithProvider("openai"), WithAPIKey("key"))
	if err != nil {
		t.Fatalf("CreateModel: %v", err)
	}
	if m.String() != "gpt-4o" {
		t.Errorf("openai model key = %q, want gpt-4o", m.String())
	}
}

func TestCreateModelFromEnvVar(t *testing.T) {
	t.Setenv("GROQ_API_KEY", "env-key")
	p := LiteLLMProviderInstance()
	m, err := p.CreateModel("llama-3.3-70b-versatile", WithProvider("groq"))
	if err != nil {
		t.Fatalf("CreateModel: %v", err)
	}
	if m.apiKey != "env-key" {
		t.Errorf("apiKey = %q, want env-key", m.apiKey)
	}
}

func TestCreateModelMissingAPIKey(t *testing.T) {
	// Ensure the env var is not set (this test doesn't modify user environment).
	if os.Getenv("OPENAI_API_KEY") != "" {
		t.Skip("OPENAI_API_KEY is set in this environment")
	}
	p := LiteLLMProviderInstance()
	_, err := p.CreateModel("gpt-4o", WithProvider("openai"))
	if err == nil {
		t.Fatal("expected error for missing API key")
	}
	var cfgErr ConfigurationError
	if !errors.As(err, &cfgErr) {
		t.Errorf("expected ConfigurationError, got %T", err)
	}
}

func TestCreateModelUnregisteredProvider(t *testing.T) {
	p := LiteLLMProviderInstance()
	_, err := p.CreateModel("anything", WithProvider("does-not-exist"), WithAPIKey("k"))
	if err == nil {
		t.Fatal("expected error for unregistered provider")
	}
	var provErr ProviderError
	if !errors.As(err, &provErr) {
		t.Errorf("expected ProviderError, got %T", err)
	}
	var baseErr ADKRouterError
	if !errors.As(err, &baseErr) {
		t.Error("ProviderError should satisfy the ADKRouterError interface")
	}
}

func TestCreateModelEmptyName(t *testing.T) {
	p := LiteLLMProviderInstance()
	_, err := p.CreateModel("  ", WithProvider("openai"), WithAPIKey("k"))
	if err == nil {
		t.Fatal("expected error for empty model name")
	}
	var modelErr ModelError
	if !errors.As(err, &modelErr) {
		t.Errorf("expected ModelError, got %T", err)
	}
}

func TestInferProvider(t *testing.T) {
	p := LiteLLMProviderInstance()
	if got := p.InferProvider("openrouter/xiaomi/mimo-v2-flash:free"); got != "openrouter" {
		t.Errorf("InferProvider = %q", got)
	}
	if got := p.InferProvider("gpt-4o"); got != "openai" {
		t.Errorf("InferProvider = %q, want openai", got)
	}
}

func TestGetModelString(t *testing.T) {
	p := LiteLLMProviderInstance()
	if got := p.GetModelString("mimo-v2-flash:free", "openrouter"); got != "openrouter/mimo-v2-flash:free" {
		t.Errorf("GetModelString = %q", got)
	}
	if got := p.GetModelString("openrouter/mimo-v2-flash:free", "openrouter"); got != "openrouter/mimo-v2-flash:free" {
		t.Errorf("GetModelString with / = %q", got)
	}
	if got := p.GetModelString("llama", "groq"); got != "groq/llama" {
		t.Errorf("GetModelString groq = %q", got)
	}
}

// ---------------------------------------------------------------------------
// OpenRouterProvider tests
// ---------------------------------------------------------------------------

func TestOpenRouterDefaultBaseURL(t *testing.T) {
	if OPENROUTER_DEFAULT_BASE_URL != "https://openrouter.ai/api/v1" {
		t.Errorf("OPENROUTER_DEFAULT_BASE_URL = %q", OPENROUTER_DEFAULT_BASE_URL)
	}
	p := NewOpenRouterProvider("test-key")
	if p.APIBase() != "https://openrouter.ai/api/v1" {
		t.Errorf("apiBase = %q", p.APIBase())
	}
	if p.APIKey() != "test-key" {
		t.Errorf("apiKey = %q", p.APIKey())
	}
}

func TestOpenRouterCreateModel(t *testing.T) {
	p := NewOpenRouterProvider("test-key")
	m, err := p.CreateModel("xiaomi/mimo-v2-flash:free")
	if err != nil {
		t.Fatalf("CreateModel: %v", err)
	}
	if m.String() != "openrouter/xiaomi/mimo-v2-flash:free" {
		t.Errorf("model key = %q", m.String())
	}
	if m.apiKey != "test-key" {
		t.Errorf("apiKey = %q", m.apiKey)
	}
}

func TestOpenRouterFromEnv(t *testing.T) {
	t.Setenv("OPENROUTER_API_KEY", "env-key")
	p, err := NewOpenRouterProviderFromEnv()
	if err != nil {
		t.Fatalf("fromEnv: %v", err)
	}
	if p.APIKey() != "env-key" {
		t.Errorf("apiKey = %q", p.APIKey())
	}
}

func TestOpenRouterFromEnvMissingKey(t *testing.T) {
	t.Setenv("OPENROUTER_API_KEY", "")
	_, err := NewOpenRouterProviderFromEnv()
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	var cfgErr ConfigurationError
	if !errors.As(err, &cfgErr) {
		t.Errorf("expected ConfigurationError, got %T", err)
	}
}

// ---------------------------------------------------------------------------
// Registry tests
// ---------------------------------------------------------------------------

func TestRegisterProvider(t *testing.T) {
	RegisterProvider("my-custom", ProviderConfig{
		BaseURL:     "https://custom.example.com/v1",
		EnvKey:      "MY_CUSTOM_API_KEY",
		NamePrefix:  "custom",
		Description: "test",
		Category:    CategoryCustom,
	})
	providers := GetAvailableProviders()
	if !contains(providers, "my-custom") {
		t.Errorf("expected my-custom in providers, got %v", providers)
	}

	m, err := GetModel("my-model", WithProvider("my-custom"), WithAPIKey("k"))
	if err != nil {
		t.Fatalf("GetModel custom: %v", err)
	}
	if m.String() != "my-custom/my-model" {
		t.Errorf("custom model key = %q", m.String())
	}
	if m.apiBase != "https://custom.example.com/v1" {
		t.Errorf("custom apiBase = %q", m.apiBase)
	}
}

func TestGetModelTopLevel(t *testing.T) {
	m, err := GetModel("gpt-4o",
		WithProvider("openai"),
		WithAPIKey("test-key"),
		WithTemperature(0.7),
		WithMaxTokens(4096),
	)
	if err != nil {
		t.Fatalf("GetModel: %v", err)
	}
	if m.String() != "gpt-4o" {
		t.Errorf("key = %q", m.String())
	}
	if m.temperature == nil || *m.temperature != 0.7 {
		t.Errorf("temperature = %v", m.temperature)
	}
	if m.maxTokens == nil || *m.maxTokens != 4096 {
		t.Errorf("maxTokens = %v", m.maxTokens)
	}
}

// ---------------------------------------------------------------------------
// Chat completion HTTP tests
// ---------------------------------------------------------------------------

func TestCreateChatCompletion(t *testing.T) {
	var (
		gotBody string
		gotAuth string
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		gotAuth = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "chatcmpl-123",
			"object": "chat.completion",
			"created": 1720000000,
			"model": "gpt-4o",
			"choices": [{"index": 0, "message": {"role": "assistant", "content": "Hello!"}, "finish_reason": "stop"}],
			"usage": {"prompt_tokens": 5, "completion_tokens": 3, "total_tokens": 8}
		}`))
	}))
	defer srv.Close()

	p := LiteLLMProviderInstance()
	m, err := p.CreateModel("gpt-4o",
		WithProvider("openai"),
		WithAPIKey("secret-key"),
		WithAPIBase(srv.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("CreateModel: %v", err)
	}

	resp, err := m.CreateChatCompletion(context.Background(), []ChatMessage{
		{Role: "user", Content: "Hi"},
	})
	if err != nil {
		t.Fatalf("CreateChatCompletion: %v", err)
	}
	if len(resp.Choices) != 1 || resp.Choices[0].Message.Content != "Hello!" {
		t.Errorf("unexpected response: %+v", resp.Choices)
	}
	if resp.Usage == nil || resp.Usage.TotalTokens != 8 {
		t.Errorf("unexpected usage: %+v", resp.Usage)
	}
	if gotAuth != "Bearer secret-key" {
		t.Errorf("Authorization = %q", gotAuth)
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(gotBody), &parsed); err != nil {
		t.Fatalf("request body not JSON: %v", err)
	}
	if parsed["model"] != "gpt-4o" {
		t.Errorf("request model = %v", parsed["model"])
	}
}

func TestCreateChatCompletionStream(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "data: {\"id\":\"1\",\"object\":\"chat.completion.chunk\",\"created\":1,\"model\":\"gpt-4o\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\"},\"finish_reason\":null}]}\n\n")
		fmt.Fprintf(w, "data: {\"id\":\"1\",\"object\":\"chat.completion.chunk\",\"created\":1,\"model\":\"gpt-4o\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\"Hel\"},\"finish_reason\":null}]}\n\n")
		fmt.Fprintf(w, "data: {\"id\":\"1\",\"object\":\"chat.completion.chunk\",\"created\":1,\"model\":\"gpt-4o\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\"lo\"},\"finish_reason\":null}]}\n\n")
		fmt.Fprintf(w, "data: {\"id\":\"1\",\"object\":\"chat.completion.chunk\",\"created\":1,\"model\":\"gpt-4o\",\"choices\":[{\"index\":0,\"delta\":{},\"finish_reason\":\"stop\"}]}\n\n")
		fmt.Fprintf(w, "data: [DONE]\n\n")
	}))
	defer srv.Close()

	p := LiteLLMProviderInstance()
	m, err := p.CreateModel("gpt-4o",
		WithProvider("openai"),
		WithAPIKey("key"),
		WithAPIBase(srv.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("CreateModel: %v", err)
	}

	stream, err := m.CreateChatCompletionStream(context.Background(), []ChatMessage{{Role: "user", Content: "Hi"}})
	if err != nil {
		t.Fatalf("CreateChatCompletionStream: %v", err)
	}
	defer stream.Close()

	var chunks []StreamChunk
	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			break
		}
		var chunk StreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			t.Fatalf("unmarshal chunk: %v", err)
		}
		chunks = append(chunks, chunk)
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scanner: %v", err)
	}
	if len(chunks) != 4 {
		t.Fatalf("got %d chunks, want 4", len(chunks))
	}
	var content strings.Builder
	for _, c := range chunks {
		if len(c.Choices) > 0 && c.Choices[0].Delta != nil {
			content.WriteString(c.Choices[0].Delta.Content)
		}
	}
	if content.String() != "Hello" {
		t.Errorf("streamed content = %q, want Hello", content.String())
	}
}

func TestCreateChatCompletionHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": {"message": "Invalid API key", "type": "invalid_request_error"}}`))
	}))
	defer srv.Close()

	p := LiteLLMProviderInstance()
	m, err := p.CreateModel("gpt-4o",
		WithProvider("openai"),
		WithAPIKey("bad-key"),
		WithAPIBase(srv.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("CreateModel: %v", err)
	}
	_, err = m.CreateChatCompletion(context.Background(), []ChatMessage{{Role: "user", Content: "Hi"}})
	if err == nil {
		t.Fatal("expected error for 401")
	}
	if !strings.Contains(err.Error(), "Invalid API key") {
		t.Errorf("error = %v, want message about invalid API key", err)
	}
}

// ---------------------------------------------------------------------------
// Misc
// ---------------------------------------------------------------------------

func TestDotEnvLoading(t *testing.T) {
	dir := t.TempDir()
	path := dir + string(os.PathSeparator) + "test.env"
	content := "# comment\nEXPORT_FLAG=unused\nexport TEST_DOTENV_KEY=value-from-file\nQUOTED=\"hello world\"\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TEST_DOTENV_KEY", "existing")

	if err := LoadEnvConfig(path); err != nil {
		t.Fatalf("LoadEnvConfig: %v", err)
	}
	if got := os.Getenv("QUOTED"); got != "hello world" {
		t.Errorf("QUOTED = %q, want hello world", got)
	}
	// Existing env must not be overridden.
	if got := os.Getenv("TEST_DOTENV_KEY"); got != "existing" {
		t.Errorf("existing env overridden: %q", got)
	}
}

func TestModelToMap(t *testing.T) {
	p := LiteLLMProviderInstance()
	m, err := p.CreateModel("gpt-4o",
		WithProvider("openai"),
		WithAPIKey("key"),
		WithTemperature(0.5),
		WithMaxTokens(100),
	)
	if err != nil {
		t.Fatalf("CreateModel: %v", err)
	}
	anonymous, err := json.Marshal(m.ToMap(false))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(anonymous), "key") {
		t.Error("ToMap(false) should not include the API key")
	}
	secretMap := m.ToMap(true)
	if secretMap["api_key"] != "key" {
		t.Errorf("ToMap(true) api_key = %v", secretMap["api_key"])
	}
}

func TestModelCloneDoesNotShareState(t *testing.T) {
	p := LiteLLMProviderInstance()
	m, err := p.CreateModel("gpt-4o", WithProvider("openai"), WithAPIKey("k"))
	if err != nil {
		t.Fatalf("CreateModel: %v", err)
	}
	c := m.clone()
	WithTemperature(0.9).apply(c)
	if m.temperature != nil {
		t.Error("mutating clone should not affect original")
	}
}

func TestConcurrentCreateModel(t *testing.T) {
	p := LiteLLMProviderInstance()
	var wg sync.WaitGroup
	errCh := make(chan error, 20)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := p.CreateModel("gpt-4o", WithProvider("openai"), WithAPIKey("k"))
			if err != nil {
				errCh <- err
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Errorf("concurrent CreateModel: %v", err)
	}
}

func TestEndpointConstruction(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()

	cases := []struct {
		apiBase string
		want    string
	}{
		{srv.URL, srv.URL + "/chat/completions"},
		{srv.URL + "/v1", srv.URL + "/v1/chat/completions"},
		{srv.URL + "/v1/chat/completions", srv.URL + "/v1/chat/completions"},
	}
	for _, tc := range cases {
		m := &Model{apiBase: tc.apiBase, client: &http.Client{}}
		got := m.endpoint()
		if got != tc.want {
			t.Errorf("endpoint(%q) = %q, want %q", tc.apiBase, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func contains(list []string, target string) bool {
	for _, x := range list {
		if x == target {
			return true
		}
	}
	return false
}
