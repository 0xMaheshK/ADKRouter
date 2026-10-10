package adkrouter

import (
	"encoding/json"
	"net/http"
	"time"
)

// ChatMessage represents a single message in a chat conversation.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ToolCallFunction holds the arguments of a tool call request by the model.
type ToolCallFunction struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

// ToolCall is a request by the model to invoke a tool.
type ToolCall struct {
	ID       string           `json:"id,omitempty"`
	Type     string           `json:"type,omitempty"`
	Function ToolCallFunction `json:"function,omitempty"`
}

// ChatMessageToolCall is a message variant that carries tool calls.
type ChatMessageToolCall struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	Name       string     `json:"name,omitempty"`
}

// Usage reports token consumption for a request.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionRequest is the OpenAI-compatible chat completion request body.
type ChatCompletionRequest struct {
	Model            string         `json:"model"`
	Messages         []ChatMessage  `json:"messages"`
	Temperature      *float64       `json:"temperature,omitempty"`
	MaxTokens        *int           `json:"max_tokens,omitempty"`
	TopP             *float64       `json:"top_p,omitempty"`
	FrequencyPenalty *float64       `json:"frequency_penalty,omitempty"`
	PresencePenalty  *float64       `json:"presence_penalty,omitempty"`
	Stop             []string       `json:"stop,omitempty"`
	N                *int           `json:"n,omitempty"`
	User             string         `json:"user,omitempty"`
	Stream           bool           `json:"stream,omitempty"`
	Extra            map[string]any `json:"-"`
	Tools            []ToolCall     `json:"tools,omitempty"`
}

// MarshalJSON satisfies custom serialization used by the client.
func (r *ChatCompletionRequest) MarshalJSON() ([]byte, error) {
	type alias ChatCompletionRequest
	body, err := json.Marshal((*alias)(r))
	if err != nil {
		return nil, err
	}
	if len(r.Extra) == 0 {
		return body, nil
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	for k, v := range r.Extra {
		payload[k] = v
	}
	return json.Marshal(payload)
}

// Choice is a single completion choice in a non-streaming response.
type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// ChatCompletionResponse is the OpenAI-compatible chat completion response.
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   *Usage   `json:"usage,omitempty"`
}

// StreamChoice is a single chunk of a streaming response.
type StreamChoice struct {
	Index        int          `json:"index"`
	Delta        *ChatMessage `json:"delta,omitempty"`
	FinishReason string       `json:"finish_reason"`
}

// StreamChunk is one Server-Sent Event payload of a streaming response.
type StreamChunk struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
}

// APIErrorBody carries the error payload returned on non-2xx HTTP responses.
type APIErrorBody struct {
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// ModelOption configures a [Model] when it is created.
type ModelOption func(*Model)

// Model is a fully configured model that can make chat completion requests
// to an OpenAI-compatible endpoint.
type Model struct {
	provider string // provider name (e.g. "openrouter")
	key      string // model key including provider prefix (e.g. "openrouter/xiaomi/mimo-v2-flash:free")
	name     string // model name sent to the API (e.g. "xiaomi/mimo-v2-flash:free")

	apiKey  string
	apiBase string

	temperature      *float64
	maxTokens        *int
	topP             *float64
	frequencyPenalty *float64
	presencePenalty  *float64
	stop             []string
	timeout          time.Duration
	stream           bool
	extra            map[string]any
	client           *http.Client
}

// String returns the model key (including the provider prefix when applicable).
func (m *Model) String() string { return m.key }

// ToMap returns a copy of the model configuration as a map. Sensitive values
// (API keys) are excluded unless includeSecrets is true.
func (m *Model) ToMap(includeSecrets bool) map[string]any {
	out := map[string]any{
		"model": m.key,
	}
	if includeSecrets {
		out["api_key"] = m.apiKey
	}
	if m.apiBase != "" {
		out["api_base"] = m.apiBase
	}
	if m.temperature != nil {
		out["temperature"] = *m.temperature
	}
	if m.maxTokens != nil {
		out["max_tokens"] = *m.maxTokens
	}
	if m.topP != nil {
		out["top_p"] = *m.topP
	}
	if m.frequencyPenalty != nil {
		out["frequency_penalty"] = *m.frequencyPenalty
	}
	if m.presencePenalty != nil {
		out["presence_penalty"] = *m.presencePenalty
	}
	if len(m.stop) > 0 {
		out["stop"] = m.stop
	}
	if m.stream {
		out["stream"] = true
	}
	for k, v := range m.extra {
		out[k] = v
	}
	return out
}
