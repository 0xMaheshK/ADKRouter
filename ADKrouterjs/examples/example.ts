import { getModel, OpenRouterProvider, LiteLLMProvider, registerProvider } from '../src/index';

async function main() {
    console.log('ADK Router JS - TypeScript Example\n');

    // Method 1: Using getModel helper
    console.log('Method 1: Using getModel');
    const model1 = getModel({
        modelName: "xiaomi/mimo-v2-flash:free",
        provider: "openrouter",
        apiKey: process.env.OPENROUTER_API_KEY,
        temperature: 0.7,
        maxTokens: 4096,
    });
    console.log('Model options:', JSON.stringify(model1, null, 2));

    // Method 2: Using OpenRouterProvider
    console.log('\nMethod 2: Using OpenRouterProvider');
    const provider = OpenRouterProvider.fromEnv();
    console.log('Provider base URL:', provider.apiBase);
    
    // Method 3: Using LiteLLMProvider singleton
    console.log('\nMethod 3: Using LiteLLMProvider singleton');
    const liteProvider = LiteLLMProvider.getInstance();
    console.log('Available providers:', liteProvider.getAvailableProviders());

    // Method 4: Custom provider
    console.log('\nMethod 4: Custom provider registration');
    registerProvider("my-custom-provider", {
        base_url: "https://custom.api.com/v1",
        env_key: "MY_CUSTOM_API_KEY",
        name_prefix: "custom",
        description: "My custom OpenAI-compatible API",
        category: "custom",
    });
    console.log('Custom provider registered!');

    // Query providers
    console.log('\nAvailable providers by category:');
    const gateways = getProvidersByCategory('gateway');
    console.log('Gateways:', Object.keys(gateways));
}

main().catch(console.error);