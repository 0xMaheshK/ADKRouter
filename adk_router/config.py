import os
from typing import Optional, Dict, Any

DEFAULT_PROVIDERS = {
    "openrouter": {
        "base_url": "https://openrouter.ai/api/v1",
        "env_key": "OPENROUTER_API_KEY",
        "name_prefix": "openrouter",
        "description": "OpenRouter - Unified API gateway for models",
        "category": "gateway",
    },
    "openai": {
        "base_url": "https://api.openai.com/v1",
        "env_key": "OPENAI_API_KEY",
        "name_prefix": "openai",
        "description": "OpenAI GPT models",
        "category": "provider",
    },
    "anthropic": {
        "base_url": "https://api.anthropic.com",
        "env_key": "ANTHROPIC_API_KEY",
        "name_prefix": "anthropic",
        "description": "Anthropic Claude models",
        "category": "provider",
    },
    "azure": {
        "base_url": None,
        "env_key": "AZURE_OPENAI_API_KEY",
        "name_prefix": "azure",
        "description": "Azure OpenAI",
        "category": "provider",
    },
    "azure_ai": {
        "base_url": None,
        "env_key": "AZURE_OPENAI_API_KEY",
        "name_prefix": "azure_ai",
        "description": "Azure AI Studio",
        "category": "provider",
    },
    "gemini": {
        "base_url": None,
        "env_key": "GOOGLE_API_KEY",
        "name_prefix": "gemini",
        "description": "Google Gemini",
        "category": "provider",
    },
    "vertex": {
        "base_url": None,
        "env_key": "GOOGLE_APPLICATION_CREDENTIALS",
        "name_prefix": "vertex",
        "description": "Google Vertex AI",
        "category": "provider",
    },
    "claude": {
        "base_url": None,
        "env_key": "ANTHROPIC_API_KEY",
        "name_prefix": "anthropic",
        "description": "Claude via Bedrock (alias for anthropic)",
        "category": "alias",
    },
    "bedrock": {
        "base_url": None,
        "env_key": "AWS_ACCESS_KEY_ID",
        "name_prefix": "bedrock",
        "description": "AWS Bedrock",
        "category": "provider",
    },
    "groq": {
        "base_url": "https://api.groq.com/v1",
        "env_key": "GROQ_API_KEY",
        "name_prefix": "groq",
        "description": "Groq Llama models",
        "category": "provider",
    },
    "mistral": {
        "base_url": "https://api.mistral.ai/v1",
        "env_key": "MISTRAL_API_KEY",
        "name_prefix": "mistral",
        "description": "Mistral AI models",
        "category": "provider",
    },
    "ollama": {
        "base_url": "http://localhost:11434",
        "env_key": None,
        "name_prefix": "ollama",
        "description": "Ollama local models",
        "category": "local",
    },
    "together": {
        "base_url": "https://api.together.ai/v1",
        "env_key": "TOGETHER_API_KEY",
        "name_prefix": "together",
        "description": "Together AI models",
        "category": "provider",
    },
    "xai": {
        "base_url": "https://api.x.ai/v1",
        "env_key": "XAI_API_KEY",
        "name_prefix": "xai",
        "description": "xAI Grok models",
        "category": "provider",
    },
    "deepseek": {
        "base_url": "https://api.deepseek.com/v1",
        "env_key": "DEEPSEEK_API_KEY",
        "name_prefix": "deepseek",
        "description": "DeepSeek models",
        "category": "provider",
    },
    "cohere": {
        "base_url": "https://api.cohere.com/v1",
        "env_key": "COHERE_API_KEY",
        "name_prefix": "cohere",
        "description": "Cohere models",
        "category": "provider",
    },
    "nvidia_nim": {
        "base_url": None,
        "env_key": "NVIDIA_NIM_API_KEY",
        "name_prefix": "nvidia_nim",
        "description": "NVIDIA NIM models",
        "category": "provider",
    },
    "nebius": {
        "base_url": "https://api.nebius.ai/v1",
        "env_key": "NEBIUS_API_KEY",
        "name_prefix": "nebius",
        "description": "Nebius AI Studio",
        "category": "provider",
    },
    "meta_llama": {
        "base_url": None,
        "env_key": "METALAMA_API_KEY",
        "name_prefix": "meta_llama",
        "description": "Meta Llama models",
        "category": "provider",
    },
    "anyscale": {
        "base_url": "https://api.endpoints.anyscale.com/v1",
        "env_key": "ANYSCALE_API_KEY",
        "name_prefix": "anyscale",
        "description": "Anyscale endpoints",
        "category": "provider",
    },
    "fireworks": {
        "base_url": "https://api.fireworks.ai/inference/v1",
        "env_key": "FIREWORKS_API_KEY",
        "name_prefix": "fireworks_ai",
        "description": "Fireworks AI",
        "category": "provider",
    },
    "moonshot": {
        "base_url": "https://api.moonshot.ai/v1",
        "env_key": "MOONSHOT_API_KEY",
        "name_prefix": "moonshot",
        "description": "Moonshot AI Kimi models",
        "category": "provider",
    },
    "minimax": {
        "base_url": "https://api.minimax.chat/v1",
        "env_key": "MINIMAX_API_KEY",
        "name_prefix": "minimax",
        "description": "Minimax models",
        "category": "provider",
    },
    "replicate": {
        "base_url": "https://api.replicate.com/v1",
        "env_key": "REPLICATE_API_TOKEN",
        "name_prefix": "replicate",
        "description": "Replicate models",
        "category": "provider",
    },
    "huggingface": {
        "base_url": "https://api-inference.huggingface.co",
        "env_key": "HUGGING_FACE_HUB_TOKEN",
        "name_prefix": "huggingface",
        "description": "HuggingFace models",
        "category": "provider",
    },
    "qwen": {
        "base_url": "https://dashscope.ali-bab.com/api/v1",
        "env_key": "DASHSCOPE_API_KEY",
        "name_prefix": "dashscope",
        "description": "Qwen via DashScope",
        "category": "provider",
    },
    "custom_openai": {
        "base_url": None,
        "env_key": "OPENAI_API_KEY",
        "name_prefix": "openai",
        "description": "Custom OpenAI-compatible endpoint (set OPENAI_BASE_URL)",
        "category": "custom",
    },
}


def load_env_config(env_path: Optional[str] = None) -> Dict[str, Any]:
    try:
        from dotenv import load_dotenv
        load_dotenv(env_path)
    except ImportError:
        pass
    return {}


def get_env_var(key: str, default: Optional[str] = None) -> Optional[str]:
    return os.getenv(key, default)


def detect_default_endpoint(provider_name: str) -> Optional[str]:
    provider_config = DEFAULT_PROVIDERS.get(provider_name.lower())
    if provider_config:
        return provider_config.get("base_url")
    return None


def get_providers_by_category(category: str) -> Dict[str, Any]:
    return {
        name: config for name, config in DEFAULT_PROVIDERS.items()
        if config.get("category") == category
    }


def get_all_providers() -> Dict[str, Any]:
    return DEFAULT_PROVIDERS.copy()