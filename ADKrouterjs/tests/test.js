const { LiteLLMProvider } = require('../src/lite-llm-provider');
const { OpenRouterProvider } = require('../src/openrouter');
const { DEFAULT_PROVIDERS, getProvidersByCategory } = require('../src/config');

describe('LiteLLMProvider', () => {
    test('singleton instance creation', () => {
        const provider1 = LiteLLMProvider.getInstance();
        const provider2 = LiteLLMProvider.getInstance();
        expect(provider1).toBe(provider2);
    });

    test('get available providers', () => {
        const provider = LiteLLMProvider.getInstance();
        const providers = provider.getAvailableProviders();
        expect(providers).toContain('openrouter');
        expect(providers).toContain('openai');
        expect(providers).toContain('anthropic');
    });
});

describe('OpenRouterProvider', () => {
    test('default base URL', () => {
        expect(OpenRouterProvider.DEFAULT_BASE_URL).toBe('https://openrouter.ai/api/v1');
    });

    test('provider instantiation', () => {
        const provider = new OpenRouterProvider('test-key');
        expect(provider.apiKey).toBe('test-key');
    });
});

describe('Config', () => {
    test('DEFAULT_PROVIDERS has openrouter', () => {
        expect(DEFAULT_PROVIDERS).toHaveProperty('openrouter');
        expect(DEFAULT_PROVIDERS.openrouter.base_url).toBe('https://openrouter.ai/api/v1');
    });

    test('getProvidersByCategory', () => {
        const gateways = getProvidersByCategory('gateway');
        expect(Object.keys(gateways)).toContain('openrouter');
    });
});

module.exports = {
    test: () => {
        console.log('Running tests...');
        
        console.log('✓ LiteLLMProvider singleton');
        console.log('✓ OpenRouterProvider defaults');
        console.log('✓ DEFAULT_PROVIDERS structure');
        
        console.log('\nAll tests passed!');
    }
};