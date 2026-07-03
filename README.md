# ADK Router

A provider-agnostic Python module for extending Google ADK (Agent Development Kit) to support any OpenAI-compatible LLM API.

## Overview

ADK Router bridges ADK with 100+ LLM providers through LiteLLM, enabling support for:
- OpenRouter (unified gateway with free models)
- OpenAI, Azure OpenAI, Anthropic
- Groq, Mistral, xAI, DeepSeek, Cohere
- AWS Bedrock, Google Vertex AI, HuggingFace
- Local models (Ollama, vLLM)
- Custom OpenAI-compatible endpoints

## Quick Start

```python
from adk_router import get_model, OpenRouterProvider
from google.adk.agents.lm_agent import LlmAgent

# Method 1: Direct with get_model
model = get_model(
    model_name="xiaomi/mimo-v2-flash:free",
    provider="openrouter",
    api_key="YOUR_OPENROUTER_API_KEY",
)

agent = LlmAgent(
    model=model,
    name="agent",
    instruction="You are a helpful assistant.",
)

# Method 2: From environment variables
provider = OpenRouterProvider.from_env()  # Loads OPENROUTER_API_KEY
model = provider.create_model("xiaomi/mimo-v2-flash:free")

# Method 3: Custom endpoint
model = get_model(
    model_name="my-model",
    provider="openai",
    api_key="YOUR_KEY",
    api_base="https://custom-endpoint.com/v1",
)
```

## Installation

```bash
pip install -e .
# or
pip install adk-router
```

## Requirements

- Python >= 3.9
- google-adk >= 0.1.0
- litellm >= 1.0.0
- python-dotenv >= 1.0.0

## Environment Variables

### OpenRouter
```bash
OPENROUTER_API_KEY=your_openrouter_api_key
OPENROUTER_BASE_URL=https://openrouter.ai/api/v1  # optional
```

### OpenAI
```bash
OPENAI_API_KEY=sk-...
```

### Azure OpenAI
```bash
AZURE_OPENAI_API_KEY=your_key
AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com/
```

### Anthropic
```bash
ANTHROPIC_API_KEY=sk-ant-...
```

## Usage Examples

### Using Free Models (OpenRouter)
```python
from adk_router import get_model

# Free models via OpenRouter
model = get_model(
    model_name="qwen/qwen3-preview:free",
    provider="openrouter",
    api_key=os.getenv("OPENROUTER_API_KEY"),
)

# Xiaomi MiMo free model
model = get_model("xiaomi/mimo-v2-flash:free", provider="openrouter")
```

### With LiteLLM Parameters
```python
model = get_model(
    model_name="openai/gpt-4o",
    api_key=os.getenv("OPENAI_API_KEY"),
    temperature=0.7,
    max_tokens=4096,
    top_p=0.9,
    stop=["\n\n"],
)
```

### Custom Provider Registration
```python
from adk_router import register_provider

register_provider("my-endpoint", {
    "base_url": "https://my-api.com/v1",
    "env_key": "MY_API_KEY",
    "name_prefix": "custom",
})

model = get_model("model-name", provider="my-endpoint")
```

## API Reference

### `get_model(model_name, provider, api_key, api_base, **kwargs)`
Creates a LiteLlm model for ADK agents.

### `OpenRouterProvider(api_key, api_base)`
Dedicated OpenRouter helper class.

### `register_provider(name, config)`
Register custom OpenAI-compatible providers.

### `get_all_providers()` / `get_providers_by_category(category)`
Query available providers.

## Documentation

- [usage.html](usage.html) - Detailed usage guide
- [providers.html](providers.html) - Complete provider catalog with examples

## License

Apache License 2.0