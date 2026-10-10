// Package adkrouter provides a provider-agnostic way to create model configurations
// and make API calls to 100+ LLM providers through OpenAI-compatible APIs.
//
// # Overview
//
// ADK Router bridges Go applications with OpenAI-compatible LLM APIs, supporting
// providers like OpenRouter, OpenAI, Anthropic, Groq, Mistral, DeepSeek, Cohere,
// xAI, HuggingFace, AWS Bedrock, Google Vertex AI, Ollama, and many more.
//
// # Quick Start
//
// Using OpenRouterProvider:
//
//	provider, err := adkrouter.NewOpenRouterProviderWithOptions(
//	    adkrouter.WithAPIKey("your-openrouter-key"),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	model, err := provider.CreateModel("xiaomi/mimo-v2-flash:free")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	resp, err := model.CreateChatCompletion(ctx, []adkrouter.ChatMessage{
//	    {Role: "user", Content: "Hello!"},
//	})
//
// Or using the top-level GetModel helper:
//
//	model, err := adkrouter.GetModel("gpt-4o",
//	    adkrouter.WithProvider("openai"),
//	    adkrouter.WithAPIKey("your-openai-key"),
//	    adkrouter.WithTemperature(0.7),
//	    adkrouter.WithMaxTokens(4096),
//	)
//
// # Streaming
//
//	model.CreateChatCompletionStream(ctx, messages) returns an
//	[io.ReadCloser] that yields Server-Sent Events (SSE) chunks.
//
// # Environment Variables
//
// Each provider has an associated environment variable for the API key.
// Set the environment variable or pass the key explicitly to avoid
// [ConfigurationError].
//
// # Supported Providers
//
// The package includes built-in configuration for 25+ LLM providers across
// categories: gateway, provider, local, custom, and alias.
package adkrouter
