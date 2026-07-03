class ADKRouterError(Exception):
    """Base exception for adk_router module errors."""
    pass


class ConfigurationError(ADKRouterError):
    """Raised when configuration is missing or invalid."""
    pass


class ProviderError(ADKRouterError):
    """Raised when provider configuration or model name is invalid."""
    pass


class ModelError(ADKRouterError):
    """Raised when model specification is invalid."""
    pass