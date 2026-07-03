from typing import Optional, Any
from google.adk.models.lite_llm import LiteLlm

from adk_router.config import DEFAULT_PROVIDERS, get_env_var, detect_default_endpoint
from adk_router.exceptions import ConfigurationError, ProviderError


class LiteLLMProvider:
    _instance = None
    
    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            cls._instance._providers = {}
            cls._instance._register_default_providers()
        return cls._instance
    
    def _register_default_providers(self):
        for provider_name, config in DEFAULT_PROVIDERS.items():
            self.register_provider(provider_name, config)
    
    def register_provider(self, provider_name: str, config: dict):
        self._providers[provider_name.lower()] = config
    
    def create_model(
        self,
        model_name: str,
        provider: Optional[str] = None,
        api_key: Optional[str] = None,
        api_base: Optional[str] = None,
        **kwargs,
    ) -> LiteLlm:
        provider = provider or self._infer_provider(model_name)
        config = self._providers.get(provider)
        
        if not config:
            raise ProviderError(f"Provider '{provider}' is not registered. Available: {list(self._providers.keys())}")
        
        model_key = model_name
        if provider != "openai" and provider != "anthropic":
            model_key = f"{provider}/{model_name}"
        
        final_api_key = api_key or get_env_var(config["env_key"])
        if not final_api_key and config["env_key"] not in kwargs:
            raise ConfigurationError(
                f"API key for provider '{provider}' not provided. "
                f"Set environment variable '{config['env_key']}' or pass api_key parameter."
            )
        
        final_api_base = api_base or detect_default_endpoint(provider)
        
        lite_llm_kwargs = {}
        if final_api_key:
            lite_llm_kwargs["api_key"] = final_api_key
        if final_api_base:
            lite_llm_kwargs["api_base"] = final_api_base
        lite_llm_kwargs.update(kwargs)
        
        return LiteLlm(model=model_key, **lite_llm_kwargs)
    
    def _infer_provider(self, model_name: str) -> str:
        if "/" in model_name:
            return model_name.split("/")[0].lower()
        return "openai"
    
    def get_model_string(self, model_name: str, provider: str) -> str:
        if "/" in model_name:
            return model_name
        prefix = self._providers.get(provider, {}).get("name_prefix", provider)
        return f"{prefix}/{model_name}"