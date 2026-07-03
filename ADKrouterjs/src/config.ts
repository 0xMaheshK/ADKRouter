import { ProviderConfig, ProviderRegistry } from './types';

export const DEFAULT_PROVIDERS: ProviderRegistry = {
    openrouter: {
        base_url: "https://openrouter.ai/api/v1",
        env_key: "OPENROUTER_API_KEY",
        name_prefix: "openrouter",
        description: "OpenRouter - Unified API gateway for 100+ models including free options",
        category: "gateway",
    },
    openai: {
        base_url: "https://api.openai.com/v1",
        env_key: "OPENAI_API_KEY",
        name_prefix: "openai",
        description: "OpenAI GPT models",
        category: "provider",
    },
    anthropic: {
        base_url: "https://api.anthropic.com",
        env_key: "ANTHROPIC_API_KEY",
        name_prefix: "anthropic",
        description: "Anthropic Claude models",
        category: "provider",
    },
    azure: {
        base_url: null,
        env_key: "AZURE_OPENAI_API_KEY",
        name_prefix: "azure",
        description: "Azure OpenAI",
        category: "provider",
    },
    azure_ai: {
        base_url: null,
        env_key: "AZURE_OPENAI_API_KEY",
        name_prefix: "azure_ai",
        description: "Azure AI Studio",
        category: "provider",
    },
    gemini: {
        base_url: null,
        env_key: "GOOGLE_API_KEY",
        name_prefix: "gemini",
        description: "Google Gemini",
        category: "provider",
    },
    vertex: {
        base_url: null,
        env_key: "GOOGLE_APPLICATION_CREDENTIALS",
        name_prefix: "vertex",
        description: "Google Vertex AI",
        category: "provider",
    },
    bedrock: {
        base_url: null,
        env_key: "AWS_ACCESS_KEY_ID",
        name_prefix: "bedrock",
        description: "AWS Bedrock",
        category: "provider",
    },
    groq: {
        base_url: "https://api.groq.com/v1",
        env_key: "GROQ_API_KEY",
        name_prefix: "groq",
        description: "Groq Llama models",
        category: "provider",
    },
    mistral: {
        base_url: "https://api.mistral.ai/v1",
        env_key: "MISTRAL_API_KEY",
        name_prefix: "mistral",
        description: "Mistral AI models",
        category: "provider",
    },
    xai: {
        base_url: "https://api.x.ai/v1",
        env_key: "XAI_API_KEY",
        name_prefix: "xai",
        description: "xAI Grok models",
        category: "provider",
    },
    deepseek: {
        base_url: "https://api.deepseek.com/v1",
        env_key: "DEEPSEEK_API_KEY",
        name_prefix: "deepseek",
        description: "DeepSeek models",
        category: "provider",
    },
    cohere: {
        base_url: "https://api.cohere.com/v1",
        env_key: "COHERE_API_KEY",
        name_prefix: "cohere",
        description: "Cohere models",
        category: "provider",
    },
    nvidia_nim: {
        base_url: null,
        env_key: "NVIDIA_NIM_API_KEY",
        name_prefix: "nvidia_nim",
        description: "NVIDIA NIM models",
        category: "provider",
    },
    nebius: {
        base_url: "https://api.nebius.ai/v1",
        env_key: "NEBIUS_API_KEY",
        name_prefix: "nebius",
        description: "Nebius AI Studio",
        category: "provider",
    },
    together: {
        base_url: "https://api.together.ai/v1",
        env_key: "TOGETHER_API_KEY",
        name_prefix: "together",
        description: "Together AI models",
        category: "provider",
    },
    moonshot: {
        base_url: "https://api.moonshot.ai/v1",
        env_key: "MOONSHOT_API_KEY",
        name_prefix: "moonshot",
        description: "Moonshot AI Kimi models",
        category: "provider",
    },
    ollama: {
        base_url: "http://localhost:11434",
        env_key: null,
        name_prefix: "ollama",
        description: "Ollama local models",
        category: "local",
    },
    huggingface: {
        base_url: "https://api-inference.huggingface.co",
        env_key: "HUGGING_FACE_HUB_TOKEN",
        name_prefix: "huggingface",
        description: "HuggingFace models",
        category: "provider",
    },
    replicate: {
        base_url: "https://api.replicate.com/v1",
        env_key: "REPLICATE_API_TOKEN",
        name_prefix: "replicate",
        description: "Replicate models",
        category: "provider",
    },
    qwen: {
        base_url: "https://dashscope.ali-bab.com/api/v1",
        env_key: "DASHSCOPE_API_KEY",
        name_prefix: "dashscope",
        description: "Qwen via DashScope",
        category: "provider",
    },
    fireworks: {
        base_url: "https://api.fireworks.ai/inference/v1",
        env_key: "FIREWORKS_API_KEY",
        name_prefix: "fireworks_ai",
        description: "Fireworks AI",
        category: "provider",
    },
    minix: {
        base_url: "https://api.minimax.chat/v1",
        env_key: "MINIMAX_API_KEY",
        name_prefix: "minimax",
        description: "Minimax models",
        category: "provider",
    },
};

export function loadEnvConfig(envPath?: string): Record<string, string | undefined> {
    if (typeof process === 'undefined') return {};
    
    try {
        // Dynamic import for dotenv
        // In Node.js environment, dotenv is used
        const result: Record<string, string | undefined> = {};
        
        const envKey = process.env.OPENROUTER_API_KEY;
        const baseUrl = process.env.OPENROUTER_BASE_URL || DEFAULT_PROVIDERS.openrouter.base_url;
        
        if (envKey) result.OPENROUTER_API_KEY = envKey;
        if (baseUrl) result.OPENROUTER_BASE_URL = baseUrl;
        
        return result;
    } catch {
        return {};
    }
}

export function getEnvVar(key: string, defaultValue?: string): string | undefined {
    if (typeof process === 'undefined') return defaultValue;
    return process.env[key] ?? defaultValue;
}

export function detectDefaultEndpoint(providerName: string): string | undefined {
    const config = DEFAULT_PROVIDERS[providerName.toLowerCase()];
    return config?.base_url ?? undefined;
}

export function getProvidersByCategory(category: string): ProviderRegistry {
    const result: ProviderRegistry = {};
    for (const [name, config] of Object.entries(DEFAULT_PROVIDERS)) {
        if (config.category === category) {
            result[name] = config;
        }
    }
    return result;
}

export function getAllProviders(): ProviderRegistry {
    return { ...DEFAULT_PROVIDERS };
}