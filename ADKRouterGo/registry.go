package adkrouter

import (
	"sort"
	"strings"
	"sync"
)

// customRegistry is the module-level registry used by the package-level
// RegisterProvider / GetAvailableProviders functions, mirroring the
// _registry dict in the Python and JavaScript implementations.
var customRegistry = struct {
	sync.RWMutex
	providers map[string]ProviderConfig
}{
	providers: map[string]ProviderConfig{},
}

// RegisterProvider registers a custom provider config under the given name.
// Custom providers take precedence over built-in defaults for the same name.
// Names are normalized to lowercase.
//
// The provider is made available to both [GetAvailableProviders] and [GetModel].
func RegisterProvider(name string, config ProviderConfig) {
	customRegistry.Lock()
	customRegistry.providers[strings.ToLower(name)] = config
	customRegistry.Unlock()
	// Also register with the shared singleton so GetModel can resolve it.
	LiteLLMProviderInstance().RegisterProvider(name, config)
}

// GetAvailableProviders returns the names of all registered providers. When
// no custom providers are registered, the built-in provider names are
// returned; otherwise the union of custom and built-in names is returned.
func GetAvailableProviders() []string {
	customRegistry.RLock()
	customNames := make([]string, 0, len(customRegistry.providers))
	for name := range customRegistry.providers {
		customNames = append(customNames, name)
	}
	customRegistry.RUnlock()

	builtin := getAllDefaultProviderNames()

	seen := make(map[string]struct{}, len(customNames)+len(builtin))
	names := make([]string, 0, len(customNames)+len(builtin))
	for _, n := range customNames {
		if _, ok := seen[n]; !ok {
			seen[n] = struct{}{}
			names = append(names, n)
		}
	}
	for _, n := range builtin {
		if _, ok := seen[n]; !ok {
			seen[n] = struct{}{}
			names = append(names, n)
		}
	}
	sort.Strings(names)
	return names
}

// getAllDefaultProviderNames returns the name of every built-in provider.
func getAllDefaultProviderNames() []string {
	names := make([]string, 0, len(DEFAULT_PROVIDERS))
	for name := range DEFAULT_PROVIDERS {
		names = append(names, name)
	}
	return names
}

// GetModel creates a configured [Model] for the given model name. Options
// control the provider, API key, endpoint, and sampling parameters.
//
// This is the Go equivalent of the Python get_model helper.
//
// Example:
//
//	model, err := adkrouter.GetModel("gpt-4o",
//	    adkrouter.WithProvider("openai"),
//	    adkrouter.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
//	    adkrouter.WithTemperature(0.7),
//	)
func GetModel(modelName string, opts ...Option) (*Model, error) {
	return LiteLLMProviderInstance().CreateModel(modelName, opts...)
}
