package adkrouter

// OpenRouterProvider is a dedicated helper for the OpenRouter gateway.
// It hard-codes the OpenRouter provider name and default endpoint.
type OpenRouterProvider struct {
	apiKey   string
	apiBase  string
	provider *LiteLLMProvider
}

// NewOpenRouterProvider creates a new OpenRouterProvider with the given
// API key. It uses the built-in default base URL unless overridden.
func NewOpenRouterProvider(apiKey string) *OpenRouterProvider {
	return &OpenRouterProvider{
		apiKey:   apiKey,
		apiBase:  OPENROUTER_DEFAULT_BASE_URL,
		provider: LiteLLMProviderInstance(),
	}
}

// NewOpenRouterProviderWithOptions creates a new OpenRouterProvider using
// the provided options. Only [WithAPIKey] and [WithAPIBase] are honored
// for the provider itself; other options are ignored.
func NewOpenRouterProviderWithOptions(opts ...Option) *OpenRouterProvider {
	p := &OpenRouterProvider{
		apiBase:  OPENROUTER_DEFAULT_BASE_URL,
		provider: LiteLLMProviderInstance(),
	}
	m := &Model{}
	for _, opt := range opts {
		if opt != nil {
			opt.apply(m)
		}
	}
	if m.apiKey != "" {
		p.apiKey = m.apiKey
	}
	if m.apiBase != "" {
		p.apiBase = m.apiBase
	}
	return p
}

// NewOpenRouterProviderFromEnv creates a provider from the environment
// variables OPENROUTER_API_KEY and OPENROUTER_BASE_URL. It loads a .env
// file from the current directory if present. Returns a
// [ConfigurationError] if the API key is not set.
func NewOpenRouterProviderFromEnv(envPath ...string) (*OpenRouterProvider, error) {
	if len(envPath) > 0 && envPath[0] != "" {
		_ = LoadEnvConfig(envPath[0])
	} else {
		_ = LoadEnvConfig("")
	}
	apiKey := GetEnvVar("OPENROUTER_API_KEY", "")
	if apiKey == "" {
		return nil, NewConfigurationError(
			"OpenRouter API key not found. Set OPENROUTER_API_KEY environment variable or pass an explicit key.",
		)
	}
	apiBase := GetEnvVar("OPENROUTER_BASE_URL", OPENROUTER_DEFAULT_BASE_URL)
	return &OpenRouterProvider{
		apiKey:   apiKey,
		apiBase:  apiBase,
		provider: LiteLLMProviderInstance(),
	}, nil
}

// APIKey returns the configured API key.
func (p *OpenRouterProvider) APIKey() string { return p.apiKey }

// APIBase returns the configured base URL.
func (p *OpenRouterProvider) APIBase() string { return p.apiBase }

// CreateModel creates a [Model] for the given model name through OpenRouter.
// The model name is prefixed with "openrouter/" for routing compatibility.
func (p *OpenRouterProvider) CreateModel(modelName string, opts ...Option) (*Model, error) {
	all := make([]Option, 0, len(opts)+2)
	if p.apiKey != "" {
		all = append(all, WithAPIKey(p.apiKey))
	}
	if p.apiBase != "" {
		all = append(all, WithAPIBase(p.apiBase))
	}
	all = append(all, opts...)
	all = append(all, WithProvider("openrouter"))
	return p.provider.CreateModel(modelName, all...)
}

// GetModelString returns the model key for the given model name.
func (p *OpenRouterProvider) GetModelString(modelName string) string {
	return p.provider.GetModelString(modelName, "openrouter")
}
