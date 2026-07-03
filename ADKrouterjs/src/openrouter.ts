import { LiteLLMProvider } from './lite-llm-provider';
import { LiteLlmOptions } from './types';
import { DEFAULT_PROVIDERS } from './config';

const OPENROUTER_DEFAULT_BASE_URL = "https://openrouter.ai/api/v1";

export class OpenRouterProvider {
    static readonly DEFAULT_BASE_URL = OPENROUTER_DEFAULT_BASE_URL;
    
    private apiKey: string | undefined;
    private apiBase: string;
    private provider: LiteLLMProvider;

    constructor(apiKey?: string, apiBase?: string) {
        this.apiKey = apiKey;
        this.apiBase = apiBase ?? OPENROUTER_DEFAULT_BASE_URL;
        this.provider = LiteLLMProvider.getInstance();
    }

    public async createModel(
        modelName: string,
        ...kwargs: LiteLlmOptions[]
    ): Promise<LiteLlmOptions> {
        const options = kwargs[0] ?? {};
        
        return this.provider.createModel(
            modelName,
            'openrouter',
            this.apiKey ?? options.apiKey,
            this.apiBase ?? options.apiBase,
            options
        );
    }

    public static fromEnv(): OpenRouterProvider {
        const apiKey = typeof process !== 'undefined' ? process.env.OPENROUTER_API_KEY : undefined;
        const baseUrl = typeof process !== 'undefined' 
            ? process.env.OPENROUTER_BASE_URL ?? OPENROUTER_DEFAULT_BASE_URL
            : OPENROUTER_DEFAULT_BASE_URL;
            
        return new OpenRouterProvider(apiKey, baseUrl);
    }
}