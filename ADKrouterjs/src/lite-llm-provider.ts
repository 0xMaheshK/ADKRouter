import { DEFAULT_PROVIDERS, getEnvVar, detectDefaultEndpoint } from './config';
import { ProviderConfig, LiteLlmOptions } from './types';
import { ConfigurationError, ProviderError } from './exceptions';

export class LiteLLMProvider {
    private static instance: LiteLLMProvider;
    private providers: Record<string, ProviderConfig>;

    private constructor() {
        this.providers = {};
        this.registerDefaultProviders();
    }

    public static getInstance(): LiteLLMProvider {
        if (!LiteLLMProvider.instance) {
            LiteLLMProvider.instance = new LiteLLMProvider();
        }
        return LiteLLMProvider.instance;
    }

    private registerDefaultProviders(): void {
        for (const [name, config] of Object.entries(DEFAULT_PROVIDERS)) {
            this.registerProvider(name, config);
        }
    }

    public registerProvider(name: string, config: ProviderConfig): void {
        this.providers[name.toLowerCase()] = config;
    }

    public createModel(
        modelName: string,
        provider?: string,
        apiKey?: string,
        apiBase?: string,
        ...kwargs: LiteLlmOptions[]
    ): LiteLlmOptions {
        const effectiveProvider = provider ?? this.inferProvider(modelName);
        const config = this.providers[effectiveProvider];

        if (!config) {
            throw new ProviderError(
                `Provider '${effectiveProvider}' is not registered. ` +
                `Available: ${Object.keys(this.providers).join(', ')}`
            );
        }

        let modelKey = modelName;
        if (effectiveProvider !== 'openai' && effectiveProvider !== 'anthropic') {
            modelKey = `${effectiveProvider}/${modelName}`;
        }

        const finalApiKey = apiKey ?? getEnvVar(config.env_key!);
        if (!finalApiKey && !config.env_key) {
            throw new ConfigurationError(
                `API key for provider '${effectiveProvider}' not provided. ` +
                `Set environment variable '${config.env_key}' or pass apiKey parameter.`
            );
        }

        const finalApiBase = apiBase ?? detectDefaultEndpoint(effectiveProvider);

        const options: LiteLlmOptions = {
            model: modelKey,
            ...(finalApiKey ? { apiKey: finalApiKey } : {}),
            ...(finalApiBase ? { apiBase: finalApiBase } : {}),
            ...kwargs[0],
        };

        return options;
    }

    private inferProvider(modelName: string): string {
        if (modelName.includes('/')) {
            return modelName.split('/')[0].toLowerCase();
        }
        return 'openai';
    }

    public getModelString(modelName: string, provider: string): string {
        if (modelName.includes('/')) {
            return modelName;
        }
        const prefix = this.providers[provider]?.name_prefix ?? provider;
        return `${prefix}/${modelName}`;
    }

    public getAvailableProviders(): string[] {
        return Object.keys(this.providers);
    }
}