package adkrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// defaultTimeout is the default HTTP timeout for chat completion requests.
const defaultTimeout = 2 * time.Minute

// Option configures a [Model] while it is being created or before a request.
type Option interface {
	apply(*Model)
}

type optionFunc func(*Model)

func (f optionFunc) apply(m *Model) { f(m) }

// WithAPIKey sets the API key used to authenticate requests.
func WithAPIKey(key string) Option {
	return optionFunc(func(m *Model) { m.apiKey = key })
}

// WithProvider sets the provider name explicitly (e.g. "openrouter").
// When omitted, the provider is inferred from the model name or defaults
// to "openai".
func WithProvider(provider string) Option {
	return optionFunc(func(m *Model) { m.provider = provider })
}

// WithAPIBase overrides the base URL of the API endpoint. When omitted, the
// provider's default endpoint is used.
func WithAPIBase(base string) Option {
	return optionFunc(func(m *Model) { m.apiBase = base })
}

// WithTemperature sets the sampling temperature (0..2).
func WithTemperature(temperature float64) Option {
	return optionFunc(func(m *Model) { m.temperature = &temperature })
}

// WithMaxTokens sets the maximum number of tokens to generate.
func WithMaxTokens(maxTokens int) Option {
	return optionFunc(func(m *Model) { m.maxTokens = &maxTokens })
}

// WithTopP sets the nucleus sampling probability (0..1).
func WithTopP(topP float64) Option {
	return optionFunc(func(m *Model) { m.topP = &topP })
}

// WithFrequencyPenalty sets the frequency penalty (-2..2).
func WithFrequencyPenalty(penalty float64) Option {
	return optionFunc(func(m *Model) { m.frequencyPenalty = &penalty })
}

// WithPresencePenalty sets the presence penalty (-2..2).
func WithPresencePenalty(penalty float64) Option {
	return optionFunc(func(m *Model) { m.presencePenalty = &penalty })
}

// WithStop sets the stop sequences. When multiple stop sequences are
// supported, the generation stops at the first one encountered.
func WithStop(stop []string) Option {
	return optionFunc(func(m *Model) {
		if stop != nil {
			m.stop = append([]string(nil), stop...)
		} else {
			m.stop = nil
		}
	})
}

// WithStream enables streaming mode so that CreateChatCompletion yields
// chunks as they are generated.
func WithStream(stream bool) Option {
	return optionFunc(func(m *Model) { m.stream = stream })
}

// WithTimeout sets the HTTP client timeout for outgoing requests.
func WithTimeout(timeout time.Duration) Option {
	return optionFunc(func(m *Model) {
		m.timeout = timeout
		m.client = &http.Client{Timeout: timeout}
	})
}

// WithHTTPClient sets a custom HTTP client. This is useful for injecting
// transports, proxies, or TLS configuration.
func WithHTTPClient(client *http.Client) Option {
	return optionFunc(func(m *Model) { m.client = client })
}

// WithExtra stores an arbitrary key/value pair that is forwarded to the API
// request body. Unknown keys are forwarded verbatim, useful for provider
// specific parameters.
func WithExtra(key string, value any) Option {
	return optionFunc(func(m *Model) {
		if m.extra == nil {
			m.extra = map[string]any{}
		}
		m.extra[key] = value
	})
}

// LiteLLMProvider holds the provider registry and creates [Model] values.
// It is a process-wide singleton; call [LiteLLMProviderInstance] to access it.
type LiteLLMProvider struct {
	mu        sync.RWMutex
	providers map[string]ProviderConfig
	client    *http.Client
}

var (
	liteLLMOnce sync.Once
	liteLLMInst *LiteLLMProvider
)

// LiteLLMProviderInstance returns the package-level singleton provider.
func LiteLLMProviderInstance() *LiteLLMProvider {
	liteLLMOnce.Do(func() {
		liteLLMInst = newLiteLLMProvider()
	})
	return liteLLMInst
}

func newLiteLLMProvider() *LiteLLMProvider {
	p := &LiteLLMProvider{
		providers: map[string]ProviderConfig{},
		client:    &http.Client{Timeout: defaultTimeout},
	}
	p.RegisterDefaultProviders()
	return p
}

// RegisterDefaultProviders registers all built-in providers.
func (p *LiteLLMProvider) RegisterDefaultProviders() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for name, config := range DEFAULT_PROVIDERS {
		p.providers[strings.ToLower(name)] = config
	}
}

// RegisterProvider registers (or replaces) a provider config under the given
// name. Names are normalized to lowercase.
func (p *LiteLLMProvider) RegisterProvider(name string, config ProviderConfig) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.providers[strings.ToLower(name)] = config
}

// GetProvider returns a copy of the provider config for the given name.
func (p *LiteLLMProvider) GetProvider(name string) (ProviderConfig, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	cfg, ok := p.providers[strings.ToLower(name)]
	return cfg, ok
}

// GetAvailableProviders returns the names of all registered providers,
// sorted alphabetically.
func (p *LiteLLMProvider) GetAvailableProviders() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	names := make([]string, 0, len(p.providers))
	for name := range p.providers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// InferProvider determines the provider from a model name. When the model
// name contains "/", the leading segment is treated as the provider name;
// otherwise "openai" is assumed, matching the original implementations.
func (p *LiteLLMProvider) InferProvider(modelName string) string {
	if strings.Contains(modelName, "/") {
		return strings.ToLower(strings.SplitN(modelName, "/", 2)[0])
	}
	return "openai"
}

// GetModelString builds the model key for a provider/model pair. When the
// model name already contains "/", it is returned unchanged; otherwise the
// provider's name prefix is prepended.
func (p *LiteLLMProvider) GetModelString(modelName, provider string) string {
	if strings.Contains(modelName, "/") {
		return modelName
	}
	prefix := provider
	if cfg, ok := p.GetProvider(provider); ok && cfg.NamePrefix != "" {
		prefix = cfg.NamePrefix
	}
	return fmt.Sprintf("%s/%s", prefix, modelName)
}

// CreateModel creates a configured [Model] for the given model name and
// options. The provider is resolved in this order: explicit provider from
// [WithProvider], provider inferred from the model name prefix, otherwise
// "openai".
//
// The returned model can be used to make chat completion requests directly.
func (p *LiteLLMProvider) CreateModel(modelName string, opts ...Option) (*Model, error) {
	if strings.TrimSpace(modelName) == "" {
		return nil, ModelError{message: "model name cannot be empty"}
	}

	m := &Model{extra: map[string]any{}}
	for _, opt := range opts {
		if opt != nil {
			opt.apply(m)
		}
	}

	provider := strings.ToLower(m.provider)
	if provider == "" {
		provider = p.InferProvider(modelName)
	}
	m.provider = provider

	config, ok := p.GetProvider(provider)
	if !ok {
		return nil, ProviderError{
			message: fmt.Sprintf(
				"Provider '%s' is not registered. Available: %v",
				provider, p.GetAvailableProviders(),
			),
		}
	}

	key := modelName
	if provider != "openai" && provider != "anthropic" {
		key = fmt.Sprintf("%s/%s", provider, modelName)
	}
	m.key = key
	m.name = modelName

	if m.apiKey == "" && config.EnvKey != "" {
		m.apiKey = GetEnvVar(config.EnvKey, "")
	}
	if m.apiKey == "" && config.EnvKey != "" {
		return nil, ConfigurationError{
			message: fmt.Sprintf(
				"API key for provider '%s' not provided. Set environment variable '%s' or pass an explicit API key.",
				provider, config.EnvKey,
			),
		}
	}

	if m.apiBase == "" {
		m.apiBase = config.BaseURL
	}

	if m.client == nil {
		m.client = p.client
	}
	return m, nil
}

// buildRequest constructs the chat completion request from model settings.
func (m *Model) buildRequest(messages []ChatMessage, stream bool) *ChatCompletionRequest {
	req := &ChatCompletionRequest{
		Model:    m.name,
		Messages: messages,
		Stream:   stream,
		Extra:    m.extra,
	}
	if m.temperature != nil {
		req.Temperature = m.temperature
	}
	if m.maxTokens != nil {
		req.MaxTokens = m.maxTokens
	}
	if m.topP != nil {
		req.TopP = m.topP
	}
	if m.frequencyPenalty != nil {
		req.FrequencyPenalty = m.frequencyPenalty
	}
	if m.presencePenalty != nil {
		req.PresencePenalty = m.presencePenalty
	}
	if len(m.stop) > 0 {
		req.Stop = m.stop
	}
	return req
}

// endpoint returns the full chat completions URL for the model.
func (m *Model) endpoint() string {
	base := strings.TrimRight(m.apiBase, "/")
	if strings.HasSuffix(base, "/chat/completions") {
		return base
	}
	return fmt.Sprintf("%s/chat/completions", base)
}

// applyHeaders sets the standard headers for an OpenAI-compatible request.
func (m *Model) applyHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	if m.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+m.apiKey)
	}
	if m.provider == "openrouter" {
		req.Header.Set("HTTP-Referer", "https://github.com/adkrouter/adkroutergo")
		req.Header.Set("X-Title", "ADK Router Go")
	}
}

// do sends the request using the model's HTTP client.
func (m *Model) do(req *http.Request) (*http.Response, error) {
	client := m.client
	if client == nil {
		client = &http.Client{Timeout: defaultTimeout}
	}
	return client.Do(req)
}

// clone returns a copy of the model with its mutable fields duplicated,
// so per-request options do not mutate the shared model instance.
func (m *Model) clone() *Model {
	c := *m
	if m.temperature != nil {
		v := *m.temperature
		c.temperature = &v
	}
	if m.maxTokens != nil {
		v := *m.maxTokens
		c.maxTokens = &v
	}
	if m.topP != nil {
		v := *m.topP
		c.topP = &v
	}
	if m.frequencyPenalty != nil {
		v := *m.frequencyPenalty
		c.frequencyPenalty = &v
	}
	if m.presencePenalty != nil {
		v := *m.presencePenalty
		c.presencePenalty = &v
	}
	if m.stop != nil {
		c.stop = append([]string(nil), m.stop...)
	}
	if m.extra != nil {
		c.extra = make(map[string]any, len(m.extra))
		for k, v := range m.extra {
			c.extra[k] = v
		}
	}
	return &c
}

// CreateChatCompletion sends a chat completion request for the given messages
// and returns the parsed response.
//
// Example:
//
//	resp, err := model.CreateChatCompletion(ctx, []adkrouter.ChatMessage{
//	    {Role: "user", Content: "What is ADKRouter?"},
//	})
//	text := resp.Choices[0].Message.Content
func (m *Model) CreateChatCompletion(ctx context.Context, messages []ChatMessage, opts ...Option) (*ChatCompletionResponse, error) {
	if len(opts) == 0 {
		return m.createChatCompletion(ctx, messages)
	}
	c := m.clone()
	for _, opt := range opts {
		if opt != nil {
			opt.apply(c)
		}
	}
	return c.createChatCompletion(ctx, messages)
}

func (m *Model) createChatCompletion(ctx context.Context, messages []ChatMessage) (*ChatCompletionResponse, error) {
	req := m.buildRequest(messages, m.stream)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("adkrouter: marshal chat completion request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, m.endpoint(), strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("adkrouter: build request: %w", err)
	}
	m.applyHeaders(httpReq)

	resp, err := m.do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("adkrouter: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, m.decodeError(resp, "chat completion")
	}

	var parsed ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("adkrouter: decode response: %w", err)
	}
	return &parsed, nil
}

// CreateChatCompletionStream sends a streaming chat completion request and
// returns an [io.ReadCloser] that yields Server-Sent Events. Each event chunk
// is a JSON object matching [StreamChunk]. The caller must close the reader.
//
// Parsing example:
//
//	stream, err := model.CreateChatCompletionStream(ctx, messages)
//	defer stream.Close()
//	reader := bufio.NewScanner(stream)
//	for reader.Scan() {
//	    // skip "data:" prefix, parse StreamChunk, accumulate deltas
//	}
func (m *Model) CreateChatCompletionStream(ctx context.Context, messages []ChatMessage, opts ...Option) (io.ReadCloser, error) {
	if len(opts) == 0 {
		return m.createChatCompletionStream(ctx, messages)
	}
	c := m.clone()
	for _, opt := range opts {
		if opt != nil {
			opt.apply(c)
		}
	}
	return c.createChatCompletionStream(ctx, messages)
}

func (m *Model) createChatCompletionStream(ctx context.Context, messages []ChatMessage) (io.ReadCloser, error) {
	req := m.buildRequest(messages, true)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("adkrouter: marshal streaming request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, m.endpoint(), strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("adkrouter: build streaming request: %w", err)
	}
	m.applyHeaders(httpReq)

	resp, err := m.do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("adkrouter: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, m.decodeError(resp, "streaming chat completion")
	}
	return resp.Body, nil
}

// decodeError parses an OpenAI-style error body into a Go error.
func (m *Model) decodeError(resp *http.Response, action string) error {
	var apiErr APIErrorBody
	raw, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(raw, &apiErr)
	msg := fmt.Sprintf("adkrouter: %s request failed with status %d", action, resp.StatusCode)
	if apiErr.Error != nil && apiErr.Error.Message != "" {
		msg = fmt.Sprintf("adkrouter: %s", apiErr.Error.Message)
	}
	return fmt.Errorf("%s", msg)
}
