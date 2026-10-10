package adkrouter

// ADKRouterError is the base error type for all errors in the adkrouter package.
// Specific error types ([ConfigurationError], [ProviderError], [ModelError])
// implement this interface, allowing callers to use [errors.As] to check for
// the base type — mirroring the Python class hierarchy where all exceptions
// subclass ADKRouterError.
//
// Example:
//
//	var baseErr adkrouter.ADKRouterError
//	if errors.As(err, &baseErr) {
//	    // err is any adkrouter error
//	}
type ADKRouterError interface {
	error
	isADKRouterError()
}

// ConfigurationError is raised when configuration is missing or invalid,
// such as when a required API key is not provided or a provider is
// not registered.
type ConfigurationError struct {
	message string
}

func (e ConfigurationError) Error() string     { return e.message }
func (e ConfigurationError) isADKRouterError() {}

// NewConfigurationError creates a new ConfigurationError with the given message.
func NewConfigurationError(message string) error {
	return ConfigurationError{message: message}
}

// ProviderError is raised when a provider configuration or model name is invalid,
// such as when a provider is not registered or the model name is empty.
type ProviderError struct {
	message string
}

func (e ProviderError) Error() string     { return e.message }
func (e ProviderError) isADKRouterError() {}

// NewProviderError creates a new ProviderError with the given message.
func NewProviderError(message string) error {
	return ProviderError{message: message}
}

// ModelError is raised when a model specification is invalid.
type ModelError struct {
	message string
}

func (e ModelError) Error() string     { return e.message }
func (e ModelError) isADKRouterError() {}

// NewModelError creates a new ModelError with the given message.
func NewModelError(message string) error {
	return ModelError{message: message}
}
