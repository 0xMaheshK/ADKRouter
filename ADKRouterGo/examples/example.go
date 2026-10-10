// Command example demonstrates the adkrouter API for creating model
// configurations and talking to OpenAI-compatible endpoints.
package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/adkrouter/adkroutergo"
)

func main() {
	fmt.Println("ADK Router Go - Example")
	fmt.Println(strings.Repeat("-", 40))

	// Method 1: Using the top-level GetModel helper.
	// The provider is used to route the request and infer defaults.
	model, err := adkrouter.GetModel(
		"xiaomi/mimo-v2-flash:free",
		adkrouter.WithProvider("openrouter"),
		adkrouter.WithAPIKey(os.Getenv("OPENROUTER_API_KEY")),
		adkrouter.WithTemperature(0.7),
		adkrouter.WithMaxTokens(4096),
	)
	if err != nil {
		fmt.Printf("GetModel error: %v\n", err)
	} else {
		fmt.Printf("Model key: %s\n", model.String())
		maybeChat(model)
	}

	// Method 2: OpenRouterProvider from environment variables.
	provider, err := adkrouter.NewOpenRouterProviderFromEnv()
	if err != nil {
		fmt.Printf("OpenRouterProvider: %v\n\n", err)
	} else {
		fmt.Printf("OpenRouter base URL: %s\n", provider.APIBase())
		openRouterModel, err := provider.CreateModel("xiaomi/mimo-v2-flash:free")
		if err != nil {
			fmt.Printf("CreateModel error: %v\n", err)
		} else {
			fmt.Printf("OpenRouter model key: %s\n", openRouterModel.String())
		}
	}

	// Method 3: Registering a custom provider.
	adkrouter.RegisterProvider("my-endpoint", adkrouter.ProviderConfig{
		BaseURL:     "https://my-api.example.com/v1",
		EnvKey:      "MY_API_KEY",
		NamePrefix:  "custom",
		Description: "My custom OpenAI-compatible API",
		Category:    adkrouter.CategoryCustom,
	})
	fmt.Println("Custom provider 'my-endpoint' registered.")

	// Query available providers.
	fmt.Printf("\nGateways: %v\n", sortedKeys(adkrouter.GetProvidersByCategory(adkrouter.CategoryGateway)))
	fmt.Printf("Local providers: %v\n", sortedKeys(adkrouter.GetProvidersByCategory(adkrouter.CategoryLocal)))
	fmt.Printf("Total providers: %d\n", len(adkrouter.GetAllProviders()))
}

// sortedKeys returns the sorted keys of a provider map.
func sortedKeys(m map[string]adkrouter.ProviderConfig) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// maybeChat performs a real chat completion when an API key is available.
func maybeChat(model *adkrouter.Model) {
	if os.Getenv("OPENROUTER_API_KEY") == "" {
		fmt.Println("(set OPENROUTER_API_KEY to run a live chat completion)")
		return
	}
	ctx := context.Background()
	resp, err := model.CreateChatCompletion(ctx, []adkrouter.ChatMessage{
		{Role: "user", Content: "Say hello in exactly three words."},
	})
	if err != nil {
		fmt.Printf("Chat completion error: %v\n", err)
		return
	}
	if len(resp.Choices) > 0 {
		fmt.Printf("Chat reply: %s\n", resp.Choices[0].Message.Content)
	}
}
