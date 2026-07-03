export interface ProviderConfig {
    base_url: string | null;
    env_key: string | null;
    name_prefix: string;
    description: string;
    category: 'provider' | 'gateway' | 'local' | 'custom' | 'alias';
}

export interface ModelOptions {
    modelName: string;
    provider?: string;
    apiKey?: string;
    apiBase?: string;
    temperature?: number;
    maxTokens?: number;
    topP?: number;
    frequencyPenalty?: number;
    presencePenalty?: number;
    stop?: string[];
    timeout?: number;
    stream?: boolean;
    [key: string]: unknown;
}

export interface LiteLlmOptions {
    model: string;
    apiKey?: string;
    apiBase?: string;
    temperature?: number;
    maxTokens?: number;
    topP?: number;
    frequencyPenalty?: number;
    presencePenalty?: number;
    stop?: string[];
    timeout?: number;
    stream?: boolean;
    [key: string]: unknown;
}

export interface ProviderRegistry {
    [key: string]: ProviderConfig;
}

export interface EnvConfig {
    apiKey?: string;
    base_url?: string;
}