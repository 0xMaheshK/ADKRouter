from typing import Optional
from google.adk.models.lite_llm import LiteLlm

from adk_router.lite_llm_provider import LiteLLMProvider
from adk_router.config import OPENROUTER_DEFAULT_BASE_URL


class OpenRouterProvider:
    DEFAULT_BASE_URL = OPENROUTER_DEFAULT_BASE_URL
    
    def __init__(
        self,
        api_key: Optional[str] = None,
        api_base: Optional[str] = None,
    ):
        self.api_key = api_key
        self.api_base = api_base or self.DEFAULT_BASE_URL
        self._provider = LiteLLMProvider()
    
    def create_model(
        self,
        model_name: str,
        **kwargs,
    ) -> LiteLlm:
        provider_key = self.api_key or kwargs.pop("api_key", None)
        provider_base = self.api_base or kwargs.pop("api_base", None)
        
        return self._provider.create_model(
            model_name=model_name,
            provider="openrouter",
            api_key=provider_key,
            api_base=provider_base,
            **kwargs,
        )
    
    @classmethod
    def from_env(cls) -> "OpenRouterProvider":
        import os
        from dotenv import load_dotenv
        load_dotenv()
        return cls(
            api_key=os.getenv("OPENROUTER_API_KEY"),
            api_base=os.getenv("OPENROUTER_BASE_URL", cls.DEFAULT_BASE_URL),
        )