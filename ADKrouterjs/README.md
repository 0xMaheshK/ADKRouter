# ADK Router JS

A provider-agnostic module for extending Google ADK (Agent Development Kit) to support any OpenAI-compatible LLM API in TypeScript/JavaScript.

## Overview

ADK Router JS bridges ADK with 100+ LLM providers through LiteLLM, enabling support for:
- OpenRouter (unified gateway with free models)
- OpenAI, Azure OpenAI, Anthropic
- Groq, Mistral, xAI, DeepSeek, Cohere
- AWS Bedrock, Google Vertex AI, HuggingFace
- Local models (Ollama, vLLM)
- Custom OpenAI-compatible endpoints

## Installation

```bash
npm install adk-router-js
# or
yarn add adk-router-js
```

## Quick Start

```typescript
import { getModel, OpenRouterProvider, LiteLLMProvider } from 'adk-router-js';

// Method 1: Using getModel helper
const model = getModel({
    modelName: "xiaomi/mimo-v2-flash:free",
    provider: "openrouter",
    apiKey: process.env.OPENROUTER_API_KEY
});

// Method 2: Using OpenRouterProvider
const provider = OpenRouterProvider.fromEnv();
const model = await provider.createModel("xiaomi/mimo-v2-flash:free");

// Method 3: Using LiteLLMProvider directly
const liteProvider = LiteLLMProvider.getInstance();
const model = liteProvider.createModel(
    "qwen/qwen3-preview:free",
    "openrouter",
    process.env.OPENROUTER_API_KEY
);
```

## Supported Providers

| Category | Providers |
|----------|-----------|
| **Gateways** | OpenRouter |
| **Major** | OpenAI, Anthropic, Azure, Vertex AI, Gemini |
| **Fast** | Groq, Together AI, Fireworks |
| **Open** | Mistral, xAI, DeepSeek, Qwen, Cohere, Meta Llama |
| **Cloud** | Bedrock, HuggingFace, Replicate |
| **Local** | Ollama |

## API

### `getModel(options: ModelOptions): LiteLlmOptions`
Creates a LiteLLm model configuration for ADK agents.

### `OpenRouterProvider.fromEnv(): OpenRouterProvider`
Creates provider from environment variables.

### `registerProvider(name: string, config: ProviderConfig): void`
Register custom OpenAI-compatible providers.

### `getAllProviders(): ProviderRegistry`
Get all available providers.

### `getProvidersByCategory(category: string): ProviderRegistry`
Get providers filtered by category.

## Environment Variables

```bash
OPENROUTER_API_KEY=your_openrouter_api_key
OPENROUTER_BASE_URL=https://openrouter.ai/api/v1  # optional

OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=sk-ant-...
AZURE_OPENAI_API_KEY=...
```

## Documentation

- [usage.html](usage.html) - Detailed usage guide
- [providers.html](providers.html) - Complete provider catalog

## License

Apache License 2.0