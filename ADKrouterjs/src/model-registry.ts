import { LiteLLMProvider } from './lite-llm-provider';
import { LiteLlmOptions, ProviderConfig } from './types';
import { getAllProviders, getProvidersByCategory, DEFAULT_PROVIDERS } from './config';

const _registry: Record<string, ProviderConfig> = {};

export function registerProvider(providerName: string, config: ProviderConfig): void {
    _registry[providerName.toLowerCase()] = config;
}

export function getAvailableProviders(): string[] {
    return Object.keys(_registry).length > 0 ? Object.keys(_registry) : Object.keys(DEFAULT_PROVIDERS);
}

export function getModel(
    modelName: string,
    provider?: string,
    apiKey?: string,
    apiBase?: string,
    options?: Omit<LiteLlmOptions, 'model' | 'apiKey' | 'apiBase'>
): LiteLlmOptions {
    const liteProvider = LiteLLMProvider.getInstance();
    return liteProvider.createModel(
        modelName,
        provider,
        apiKey,
        apiBase,
        options
    );
}

export { getAllProviders, getProvidersByCategory, DEFAULT_PROVIDERS };