import pytest
from unittest.mock import patch, MagicMock


class TestOpenRouterProvider:
    def test_openrouter_provider_init(self):
        from adk_router import OpenRouterProvider
        provider = OpenRouterProvider(api_key="test_key", api_base="https://test.com")
        assert provider.api_key == "test_key"
        assert provider.api_base == "https://test.com"
    
    def test_openrouter_default_base_url(self):
        from adk_router import OpenRouterProvider
        provider = OpenRouterProvider(api_key="test_key")
        assert provider.api_base == "https://openrouter.ai/api/v1"


class TestGetModel:
    def test_get_model_with_provider(self):
        from adk_router import get_model
        with patch("adk_router.lite_llm_provider.LiteLLMProvider.create_model") as mock_create:
            mock_create.return_value = MagicMock()
            result = get_model("test-model", provider="openrouter", api_key="key")
            mock_create.assert_called_once_with(
                model_name="test-model",
                provider="openrouter",
                api_key="key",
                api_base=None,
            )


class TestConfig:
    def test_default_providers_exist(self):
        from adk_router.config import DEFAULT_PROVIDERS
        assert "openrouter" in DEFAULT_PROVIDERS
        assert "openai" in DEFAULT_PROVIDERS
        assert "anthropic" in DEFAULT_PROVIDERS
    
    def test_openrouter_defaults(self):
        from adk_router.config import DEFAULT_PROVIDERS
        config = DEFAULT_PROVIDERS["openrouter"]
        assert config["base_url"] == "https://openrouter.ai/api/v1"
        assert config["env_key"] == "OPENROUTER_API_KEY"
        assert config["name_prefix"] == "openrouter"


class TestExceptions:
    def test_configuration_error(self):
        from adk_router.exceptions import ConfigurationError, ADKRouterError
        assert issubclass(ConfigurationError, ADKRouterError)
    
    def test_provider_error(self):
        from adk_router.exceptions import ProviderError, ADKRouterError
        assert issubclass(ProviderError, ADKRouterError)