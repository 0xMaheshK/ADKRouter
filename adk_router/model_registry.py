from typing import Optional, Dict, Any, List
from google.adk.models.lite_llm import LiteLlm

from adk_router.lite_llm_provider import LiteLLMProvider
from adk_router.exceptions import ConfigurationError, ProviderError


_registry: Dict[str, Dict[str, Any]] = {}


def register_provider(provider_name: str, config: dict) -> None:
    global _registry
    _registry[provider_name.lower()] = config


def get_available_providers() -> List[str]:
    return list(_registry.keys())


def get_model(
    model_name: str,
    provider: Optional[str] = None,
    api_key: Optional[str] = None,
    api_base: Optional[str] = None,
    **kwargs,
) -> LiteLlm:
    lite_provider = LiteLLMProvider()
    return lite_provider.create_model(
        model_name=model_name,
        provider=provider,
        api_key=api_key,
        api_base=api_base,
        **kwargs,
    )