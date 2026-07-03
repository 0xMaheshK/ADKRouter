# ADK Router

A provider-agnostic module for extending Google ADK (Agent Development Kit) to support any OpenAI-compatible LLM API.

## Overview

ADK Router allows ADK projects to connect to LLM providers that are not natively supported by ADK, including:
- OpenRouter
- OpenAI compatible APIs
- Azure OpenAI
- Custom model endpoints
- And any provider with OpenAI-compatible API format

## Installation

```bash
pip install adk-router
```

Or manually add to your project:

```bash
pip install google-adk litellm python-dotenv
```

## Quick Start

### Basic Usage with OpenRouter

```python
from adk_router import get_model, OpenRouterProvider

# Create a model instance for OpenRouter
model = get_model(
    model_name="xiaomi/mimo-v2-flash:free",
    provider="openrouter",
    api_key="YOUR_OPENROUTER_API_KEY",
)

# Use with ADK Agent
from google.adk.agents.lm_agent import LlmAgent

agent = LlmAgent(
    model=model,
    name="rover_agent",
    instruction="You are a helpful assistant.",
)
```

### Using Environment Variables

Set environment variables:

```bash
# .env file
OPENROUTER_API_KEY=your_api_key_here
OPENROUTER_BASE_URL=https://openrouter.ai/api/v1
```

```python
from adk_router import OpenRouterProvider

# Auto-load from environment
provider = OpenRouterProvider.from_env()
model = provider.create_model("xiaomi/mimo-v2-flash:free")
```

### Using with Any OpenAI-Compatible Provider

```python
from adk_router import get_model

# Custom provider with base URL
model = get_model(
    model_name="your-model-name",
    provider="openai",
    api_key="YOUR_API_KEY",
    api_base="https://your-custom-endpoint.com/v1",
)

agent = LlmAgent(
    model=model,
    name="custom_agent",
    instruction="You are a helpful assistant.",
)
```

## API Reference

### `get_model(model_name: str, provider: str, api_key: str, api_base: str) -> LiteLlm`

Creates a LiteLlm model instance for ADK agents.

**Parameters:**
- `model_name` (str): The model identifier (e.g., `"xiaomi/mimo-v2-flash:free"` or just `"mimo-v2-flash:free"`)
- `provider` (str, optional): Provider name (`"openrouter"`, `"openai"`, `"anthropic"`, `"azure"`, `"gemini"`, or a custom provider)
- `api_key` (str, optional): API key for authentication
- `api_base` (str, optional): Base URL for the API endpoint

### `OpenRouterProvider(api_key: str, api_base: str)`

A dedicated provider class for OpenRouter integration.

**Methods:**
- `create_model(model_name: str) -> LiteLlm`: Creates a model instance
- `from_env()`: Class method to create from environment variables

### `register_provider(name: str, config: dict)`

Register a custom provider.

**Config structure:**
```python
{
    "base_url": "https://api.example.com/v1",
    "env_key": "EXAMPLE_API_KEY",
    "name_prefix": "example",
}
```

## Supported Providers

| Provider | Environment Variable | Default Base URL |
|----------|---------------------|------------------|
| OpenRouter | `OPENROUTER_API_KEY` | `https://openrouter.ai/api/v1` |
| OpenAI | `OPENAI_API_KEY` | `https://api.openai.com/v1` |
| Anthropic | `ANTHROPIC_API_KEY` | `https://api.anthropic.com` |
| Azure | `AZURE_OPENAI_API_KEY` | None (required) |
| Gemini | `GOOGLE_API_KEY` | None (required) |

## Model Naming Convention

Models are specified with provider prefixes:

```python
# OpenRouter
model="openrouter/xiaomi/mimo-v2-flash:free"

# OpenAI
model="openai/gpt-4o"

# Anthropic
model="anthropic/claude-3-haiku-20240307"
```

Or let the module infer the prefix:

```python
model="xiaomi/mimo-v2-flash:free"  # Will use "openrouter" provider
```

## Development

### Running Tests

```bash
pytest tests/
```

### Project Structure

```
adk_router/
├── __init__.py           # Public API exports
├── model_registry.py     # Model registry and factory
├── lite_llm_provider.py  # Core LiteLLM integration
├── openrouter.py         # OpenRouter-specific helper
├── config.py            # Configuration management
└── exceptions.py        # Custom exceptions
```

## Security Advisory

When using LiteLLM, be aware of potential supply chain risks. Always verify:
1. Use official ADK releases
2. Rotate API keys periodically
3. Review the [LiteLLM security advisory](https://github.com/BerriAI/litellm/security)
4. Pin specific versions in requirements.txt

## License

Apache License 2.0