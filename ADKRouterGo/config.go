package adkrouter

import (
	"os"
	"strings"
)

// OPENROUTER_DEFAULT_BASE_URL is the default base URL for OpenRouter.
const OPENROUTER_DEFAULT_BASE_URL = "https://openrouter.ai/api/v1"

// ProviderCategory groups providers by their usage characteristic.
type ProviderCategory string

// Provider categories.
const (
	CategoryProvider ProviderCategory = "provider"
	CategoryGateway  ProviderCategory = "gateway"
	CategoryLocal    ProviderCategory = "local"
	CategoryCustom   ProviderCategory = "custom"
	CategoryAlias    ProviderCategory = "alias"
)

// ProviderConfig describes how to connect to a model provider.
type ProviderConfig struct {
	// BaseURL is the default API endpoint. Empty when the endpoint must be
	// supplied by the caller (e.g. Azure, Gemini).
	BaseURL string `json:"base_url"`

	// EnvKey is the environment variable holding the API key.
	// Empty for providers that do not require a key (e.g. Ollama).
	EnvKey string `json:"env_key"`

	// NamePrefix is the prefix used when building the model string.
	NamePrefix string `json:"name_prefix"`

	// Description is a human-readable description of the provider.
	Description string `json:"description"`

	// Category groups the provider (provider, gateway, local, custom, alias).
	Category ProviderCategory `json:"category"`
}

// DEFAULT_PROVIDERS is the built-in registry of OpenAI-compatible providers.
var DEFAULT_PROVIDERS = map[string]ProviderConfig{
	"openrouter": {
		BaseURL:     "https://openrouter.ai/api/v1",
		EnvKey:      "OPENROUTER_API_KEY",
		NamePrefix:  "openrouter",
		Description: "OpenRouter - Unified API gateway for models",
		Category:    CategoryGateway,
	},
	"openai": {
		BaseURL:     "https://api.openai.com/v1",
		EnvKey:      "OPENAI_API_KEY",
		NamePrefix:  "openai",
		Description: "OpenAI GPT models",
		Category:    CategoryProvider,
	},
	"anthropic": {
		BaseURL:     "https://api.anthropic.com",
		EnvKey:      "ANTHROPIC_API_KEY",
		NamePrefix:  "anthropic",
		Description: "Anthropic Claude models",
		Category:    CategoryProvider,
	},
	"azure": {
		EnvKey:      "AZURE_OPENAI_API_KEY",
		NamePrefix:  "azure",
		Description: "Azure OpenAI",
		Category:    CategoryProvider,
	},
	"azure_ai": {
		EnvKey:      "AZURE_OPENAI_API_KEY",
		NamePrefix:  "azure_ai",
		Description: "Azure AI Studio",
		Category:    CategoryProvider,
	},
	"gemini": {
		EnvKey:      "GOOGLE_API_KEY",
		NamePrefix:  "gemini",
		Description: "Google Gemini",
		Category:    CategoryProvider,
	},
	"vertex": {
		EnvKey:      "GOOGLE_APPLICATION_CREDENTIALS",
		NamePrefix:  "vertex",
		Description: "Google Vertex AI",
		Category:    CategoryProvider,
	},
	"claude": {
		EnvKey:      "ANTHROPIC_API_KEY",
		NamePrefix:  "anthropic",
		Description: "Claude via Bedrock (alias for anthropic)",
		Category:    CategoryAlias,
	},
	"bedrock": {
		EnvKey:      "AWS_ACCESS_KEY_ID",
		NamePrefix:  "bedrock",
		Description: "AWS Bedrock",
		Category:    CategoryProvider,
	},
	"groq": {
		BaseURL:     "https://api.groq.com/v1",
		EnvKey:      "GROQ_API_KEY",
		NamePrefix:  "groq",
		Description: "Groq Llama models",
		Category:    CategoryProvider,
	},
	"mistral": {
		BaseURL:     "https://api.mistral.ai/v1",
		EnvKey:      "MISTRAL_API_KEY",
		NamePrefix:  "mistral",
		Description: "Mistral AI models",
		Category:    CategoryProvider,
	},
	"ollama": {
		BaseURL:     "http://localhost:11434",
		NamePrefix:  "ollama",
		Description: "Ollama local models",
		Category:    CategoryLocal,
	},
	"together": {
		BaseURL:     "https://api.together.ai/v1",
		EnvKey:      "TOGETHER_API_KEY",
		NamePrefix:  "together",
		Description: "Together AI models",
		Category:    CategoryProvider,
	},
	"xai": {
		BaseURL:     "https://api.x.ai/v1",
		EnvKey:      "XAI_API_KEY",
		NamePrefix:  "xai",
		Description: "xAI Grok models",
		Category:    CategoryProvider,
	},
	"deepseek": {
		BaseURL:     "https://api.deepseek.com/v1",
		EnvKey:      "DEEPSEEK_API_KEY",
		NamePrefix:  "deepseek",
		Description: "DeepSeek models",
		Category:    CategoryProvider,
	},
	"cohere": {
		BaseURL:     "https://api.cohere.com/v1",
		EnvKey:      "COHERE_API_KEY",
		NamePrefix:  "cohere",
		Description: "Cohere models",
		Category:    CategoryProvider,
	},
	"nvidia_nim": {
		EnvKey:      "NVIDIA_NIM_API_KEY",
		NamePrefix:  "nvidia_nim",
		Description: "NVIDIA NIM models",
		Category:    CategoryProvider,
	},
	"nebius": {
		BaseURL:     "https://api.nebius.ai/v1",
		EnvKey:      "NEBIUS_API_KEY",
		NamePrefix:  "nebius",
		Description: "Nebius AI Studio",
		Category:    CategoryProvider,
	},
	"meta_llama": {
		EnvKey:      "METALAMA_API_KEY",
		NamePrefix:  "meta_llama",
		Description: "Meta Llama models",
		Category:    CategoryProvider,
	},
	"anyscale": {
		BaseURL:     "https://api.endpoints.anyscale.com/v1",
		EnvKey:      "ANYSCALE_API_KEY",
		NamePrefix:  "anyscale",
		Description: "Anyscale endpoints",
		Category:    CategoryProvider,
	},
	"fireworks": {
		BaseURL:     "https://api.fireworks.ai/inference/v1",
		EnvKey:      "FIREWORKS_API_KEY",
		NamePrefix:  "fireworks_ai",
		Description: "Fireworks AI",
		Category:    CategoryProvider,
	},
	"moonshot": {
		BaseURL:     "https://api.moonshot.ai/v1",
		EnvKey:      "MOONSHOT_API_KEY",
		NamePrefix:  "moonshot",
		Description: "Moonshot AI Kimi models",
		Category:    CategoryProvider,
	},
	"minimax": {
		BaseURL:     "https://api.minimax.chat/v1",
		EnvKey:      "MINIMAX_API_KEY",
		NamePrefix:  "minimax",
		Description: "Minimax models",
		Category:    CategoryProvider,
	},
	"replicate": {
		BaseURL:     "https://api.replicate.com/v1",
		EnvKey:      "REPLICATE_API_TOKEN",
		NamePrefix:  "replicate",
		Description: "Replicate models",
		Category:    CategoryProvider,
	},
	"huggingface": {
		BaseURL:     "https://api-inference.huggingface.co",
		EnvKey:      "HUGGING_FACE_HUB_TOKEN",
		NamePrefix:  "huggingface",
		Description: "HuggingFace models",
		Category:    CategoryProvider,
	},
	"qwen": {
		BaseURL:     "https://dashscope.ali-bab.com/api/v1",
		EnvKey:      "DASHSCOPE_API_KEY",
		NamePrefix:  "dashscope",
		Description: "Qwen via DashScope",
		Category:    CategoryProvider,
	},
	"custom_openai": {
		EnvKey:      "OPENAI_API_KEY",
		NamePrefix:  "openai",
		Description: "Custom OpenAI-compatible endpoint (set OPENAI_BASE_URL)",
		Category:    CategoryCustom,
	},
}

// GetEnvVar returns the value of the environment variable for the given key,
// or the provided default when the variable is unset or empty.
func GetEnvVar(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return defaultValue
}

// DetectDefaultEndpoint returns the built-in default base URL for a provider,
// or an empty string when none is defined.
func DetectDefaultEndpoint(providerName string) string {
	config, ok := DEFAULT_PROVIDERS[strings.ToLower(providerName)]
	if !ok {
		return ""
	}
	return config.BaseURL
}

// GetProvidersByCategory returns the subset of [DEFAULT_PROVIDERS] whose
// category matches the given value.
func GetProvidersByCategory(category ProviderCategory) map[string]ProviderConfig {
	result := make(map[string]ProviderConfig)
	for name, config := range DEFAULT_PROVIDERS {
		if config.Category == category {
			result[name] = config
		}
	}
	return result
}

// GetAllProviders returns a copy of the built-in provider registry.
func GetAllProviders() map[string]ProviderConfig {
	result := make(map[string]ProviderConfig, len(DEFAULT_PROVIDERS))
	for name, config := range DEFAULT_PROVIDERS {
		result[name] = config
	}
	return result
}
