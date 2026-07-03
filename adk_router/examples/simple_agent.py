"""
Example: Using adk_router with OpenRouter to create ADK agents.

This module demonstrates how to use OpenRouter as an LLM provider
through the adk_router module.
"""
import os
from typing import Optional

from adk_router import get_model, OpenRouterProvider


def create_openrouter_model(
    model_name: str,
    api_key: Optional[str] = None,
    api_base: Optional[str] = None,
) -> "LiteLlm":
    provider = OpenRouterProvider(api_key=api_key, api_base=api_base)
    return provider.create_model(model_name)


def create_agent_from_env(model_name: str):
    from google.adk.agents.lm_agent import LlmAgent

    provider = OpenRouterProvider.from_env()
    model = provider.create_model(model_name)

    agent = LlmAgent(
        model=model,
        name="openrouter_agent",
        instruction="You are a helpful assistant powered by OpenRouter.",
    )
    return agent


def create_custom_provider_model(
    model_name: str,
    api_key: str,
    base_url: str,
):
    return get_model(
        model_name=model_name,
        provider="openai",
        api_key=api_key,
        api_base=base_url,
    )


if __name__ == "__main__":
    from google.adk.agents.lm_agent import LlmAgent

    api_key = os.getenv("OPENROUTER_API_KEY")
    if api_key:
        model = create_openrouter_model(
            model_name="xiaomi/mimo-v2-flash:free",
            api_key=api_key,
        )
        agent = LlmAgent(
            model=model,
            name="openrouter_agent",
            instruction="You are a helpful assistant.",
        )
        print("✅ Agent with OpenRouter model created successfully!")
    else:
        print("⚠️ Set OPENROUTER_API_KEY environment variable to test")