# ADK Router Go

A provider-agnostic Go module for extending Google ADK (Agent Development Kit) integrations to any OpenAI-compatible LLM API.

## Overview

ADK Router Go bridges Go applications with 100+ LLM providers through a unified, OpenAI-compatible interface, mirroring the [Python implementation](../adk_router/README.md) and [TypeScript implementation](../ADKrouterjs/README.md) of ADK Router. It enables support for:

- OpenRouter (unified gateway with free models)
- OpenAI, Azure OpenAI, Anthropic
- Groq, Mistral, xAI, DeepSeek, Cohere
- AWS Bedrock, Google Vertex AI, HuggingFace
- Local models (Ollama, vLLM)
- Custom OpenAI-compatible endpoints

Unlike Google ADK (which is Python-only), the Go implementation acts as a **standalone client**: model configurations created through this module can also make real chat completion and streaming requests directly.

## Installation

```bash
go get github.com/adkrouter/adkroutergo
```

No external dependencies are required — the module is self-contained (stdlib only).

## Quick Start

### Method 1: Top-level `GetModel` helper

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/adkrouter/adkroutergo"
)

func main() {
    model, err := adkrouter.GetModel(
        "xiaomi/mimo-v2-flash:free",
        adkrouter.WithProvider("openrouter"),
        adkrouter.WithAPIKey("YOUR_OPENROUTER_API_KEY"),
        adkrouter.WithTemperature(0.7),
        adkrouter.WithMaxTokens(4096),
    )
    if err != nil {
        log.Fatal(err)
    }

    resp, err := model.CreateChatCompletion(context.Background(), []adkrouter.ChatMessage{
        {Role: "user", Content: "What is ADK Router?"},
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(resp.Choices[0].Message.Content)
}
```

### Method 2: `OpenRouterProvider` from environment variables

```go
provider, err := adkrouter.NewOpenRouterProviderFromEnv() // loads OPENROUTER_API_KEY
if err != nil {
    log.Fatal(err)
}

model, err := provider.CreateModel("xiaomi/mimo-v2-flash:free")
if err != nil {
    log.Fatal(err)
}
```

### Method 3: Custom OpenAI-compatible endpoint

```go
model, err := adkrouter.GetModel(
    "my-model",
    adkrouter.WithProvider("openai"),
    adkrouter.WithAPIKey("YOUR_KEY"),
    adkrouter.WithAPIBase("https://custom-endpoint.com/v1"),
)
```

### Streaming

```go
stream, err := model.CreateChatCompletionStream(ctx, messages)
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

scanner := bufio.NewScanner(stream)
for scanner.Scan() {
    line := strings.TrimPrefix(strings.TrimSpace(scanner.Text()), "data:")
    if line == "[DONE]" || line == "" {
        continue
    }
    var chunk adkrouter.StreamChunk
    if err := json.Unmarshal([]byte(line), &chunk); err != nil {
        log.Fatal(err)
    }
    if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
        fmt.Print(chunk.Choices[0].Delta.Content)
    }
}
```

## API Reference

### `GetModel(modelName string, opts ...Option) (*Model, error)`

Creates a configured model for an ADK agent / client. The provider is resolved in this order:

1. Explicit provider via `WithProvider(...)`
2. Inferred from the model name prefix (e.g. `qwen/...` → `openrouter` if registered)
3. Defaults to `openai`

Options:

| Option | Description |
|--------|-------------|
| `WithProvider(name)` | Provider name (e.g. `"openrouter"`, `"openai"`, custom) |
| `WithAPIKey(key)` | API key (falls back to the provider's env var) |
| `WithAPIBase(base)` | Base URL override |
| `WithTemperature(v)` | Sampling temperature |
| `WithMaxTokens(v)` | Max tokens to generate |
| `WithTopP(v)` | Nucleus sampling probability |
| `WithFrequencyPenalty(v)` | Frequency penalty |
| `WithPresencePenalty(v)` | Presence penalty |
| `WithStop(seqs)` | Stop sequences |
| `WithStream(bool)` | Enable streaming mode |
| `WithTimeout(d)` | HTTP client timeout |
| `WithHTTPClient(c)` | Custom `http.Client` |
| `WithExtra(key, value)` | Arbitrary provider-specific parameter |

### `(*Model).CreateChatCompletion(ctx, messages, ...opts) (*ChatCompletionResponse, error)`

Sends a non-streaming chat completion request.

### `(*Model).CreateChatCompletionStream(ctx, messages, ...opts) (io.ReadCloser, error)`

Sends a streaming request; returns a reader of Server-Sent Events (one `StreamChunk` JSON object per `data:` line).

### `OpenRouterProvider`

| Function/Method | Description |
|-----------------|-------------|
| `NewOpenRouterProvider(apiKey)` | Provider with explicit key |
| `NewOpenRouterProviderFromEnv()` | From `OPENROUTER_API_KEY` / `OPENROUTER_BASE_URL` |
| `(*OpenRouterProvider).CreateModel(modelName)` | Create an OpenRouter model |
| `(*OpenRouterProvider).APIKey()` / `APIBase()` | Configured values |

### `LiteLLMProvider`

The process-wide singleton (`LiteLLMProviderInstance()`) that backs `GetModel`. Methods:

- `RegisterProvider(name, config)` — register a custom provider
- `GetAvailableProviders() []string`
- `GetProvider(name) (ProviderConfig, bool)`
- `InferProvider(modelName) string`
- `GetModelString(modelName, provider) string`

### `RegisterProvider` / `GetAvailableProviders`

Module-level helpers that register a custom provider and list providers.

### Provider catalog helpers

- `GetAllProviders() map[string]ProviderConfig`
- `GetProvidersByCategory(category ProviderCategory) map[string]ProviderConfig`
- `DetectDefaultEndpoint(providerName) string`
- `GetEnvVar(key, default) string`
- `LoadEnvConfig(envPath) error`

## Environment Variables

| Variable | Purpose |
|----------|---------|
| `OPENROUTER_API_KEY` | OpenRouter key |
| `OPENROUTER_BASE_URL` | Optional OpenRouter override |
| `OPENAI_API_KEY` | OpenAI key |
| `OPENAI_BASE_URL` | Custom OpenAI-compatible endpoint |
| `ANTHROPIC_API_KEY` | Anthropic key |
| `AZURE_OPENAI_API_KEY` | Azure key |
| `GOOGLE_API_KEY` | Gemini key |
| `GROQ_API_KEY` | Groq key |
| `MISTRAL_API_KEY` | Mistral key |
| `DEEPSEEK_API_KEY` | DeepSeek key |
| ... | See `config.go` for the full list |

When an API key is not passed explicitly, the module reads the provider's environment variable. A `ConfigurationError` is returned if neither is available.

## Custom Provider Registration

```go
adkrouter.RegisterProvider("my-endpoint", adkrouter.ProviderConfig{
    BaseURL:     "https://my-api.com/v1",
    EnvKey:      "MY_API_KEY",
    NamePrefix:  "custom",
    Description: "My custom OpenAI-compatible API",
    Category:    adkrouter.CategoryCustom,
})

model, err := adkrouter.GetModel("model-name", adkrouter.WithProvider("my-endpoint"))
```

## Model Naming Convention

Like the Python and JS implementations, the model string is provider-prefixed for routing:

- OpenRouter: `openrouter/xiaomi/mimo-v2-flash:free`
- OpenAI: `openai/gpt-4o` (no prefix — `gpt-4o` is used verbatim)
- Groq: `groq/llama-3.3-70b-versatile`

The provider prefix is used internally for routing; the actual API request uses the model name without the forcing prefix.

## Errors

The module defines four error types, all of which satisfy the `ADKRouterError` interface:

| Type | Raised when |
|------|-------------|
| `ADKRouterError` | Base interface for every adkrouter error |
| `ConfigurationError` | Missing/invalid configuration (e.g. no API key) |
| `ProviderError` | Unknown/unregistered provider |
| `ModelError` | Invalid model name |

Use `errors.As` for type-specific handling:

```go
var cfgErr adkrouter.ConfigurationError
if errors.As(err, &cfgErr) {
    // handle configuration error
}
```

## Development

```bash
go build ./...      # compile
go vet ./...        # static analysis
go test ./...       # run tests
go test ./... -race # with race detector
```

### Project Structure

```
ADKRouterGo/
├── doc.go            # Package documentation
├── config.go         # Provider catalog + config helpers
├── dotenv.go         # Minimal .env loader
├── errors.go         # Error types
├── types.go          # Chat completion request/response types
├── provider.go       # Model, LiteLLMProvider, options
├── registry.go       # RegisterProvider, GetModel, GetAvailableProviders
├── openrouter.go     # OpenRouterProvider helper
├── adkrouter.go      # Version constant
├── examples/         # Usage example
└── adkrouter_test.go # Tests
```

## Documentation

- Python (canonical): [../adk_router/README.md](../adk_router/README.md)
- JavaScript/TypeScript: [../ADKrouterjs/README.md](../ADKrouterjs/README.md)
- Usage guide (Python): [../usage.html](../usage.html)
- Provider catalog (Python): [../providers.html](../providers.html)

## License

Apache License 2.0