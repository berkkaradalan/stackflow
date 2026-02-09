"""
Provider Registry for LLM configurations and LangChain client factory.
Converted from Go provider registry system.
"""

from typing import List, Optional, Union
from pydantic import BaseModel, Field


class ModelConfig(BaseModel):
    """Represents a model configuration."""
    id: str = Field(..., alias="id")
    name: str = Field(..., alias="name")
    description: str = Field(..., alias="description")
    max_tokens: int = Field(..., alias="max_tokens")
    supports_streaming: bool = Field(..., alias="supports_streaming")
    supports_vision: bool = Field(..., alias="supports_vision")
    input_price_per_m_token: float = Field(..., alias="input_price_per_m_token")
    output_price_per_m_token: float = Field(..., alias="output_price_per_m_token")

    class Config:
        populate_by_name = True


class ProviderConfig(BaseModel):
    """Represents a provider configuration."""
    name: str = Field(..., alias="name")
    display_name: str = Field(..., alias="display_name")
    base_url: str = Field(..., alias="base_url")
    health_check_path: str = Field(..., alias="health_check_path")
    chat_completion_path: str = Field(..., alias="chat_completion_path")
    requires_api_key: bool = Field(..., alias="requires_api_key")
    models: List[ModelConfig] = Field(..., alias="models")

    class Config:
        populate_by_name = True


class ProviderListResponse(BaseModel):
    """Response model for listing providers."""
    providers: List[ProviderConfig]


class ModelListResponse(BaseModel):
    """Response model for listing models."""
    provider: str
    models: List[ModelConfig]


def get_provider_registry() -> List[ProviderConfig]:
    """Returns all available providers with their configurations."""
    return [
        ProviderConfig(
            name="zai",
            display_name="Z.AI",
            base_url="https://api.z.ai/api/paas/v4",
            health_check_path="/models",
            chat_completion_path="/chat/completions",
            requires_api_key=True,
            models=[
                ModelConfig(
                    id="glm-4.7",
                    name="GLM-4.7",
                    description="Latest GLM model (Jan 2025)",
                    max_tokens=128000,
                    supports_streaming=True,
                    supports_vision=True,
                    input_price_per_m_token=1.0,
                    output_price_per_m_token=1.0,
                ),
                ModelConfig(
                    id="glm-4.7-flash",
                    name="GLM-4.7 Flash",
                    description="Fast GLM-4.7 variant",
                    max_tokens=128000,
                    supports_streaming=True,
                    supports_vision=False,
                    input_price_per_m_token=0.5,
                    output_price_per_m_token=0.5,
                ),
                ModelConfig(
                    id="glm-4.6",
                    name="GLM-4.6",
                    description="GLM model from Dec 2024",
                    max_tokens=128000,
                    supports_streaming=True,
                    supports_vision=True,
                    input_price_per_m_token=1.0,
                    output_price_per_m_token=1.0,
                ),
                ModelConfig(
                    id="glm-4.5",
                    name="GLM-4.5",
                    description="Stable GLM-4.5 model",
                    max_tokens=128000,
                    supports_streaming=True,
                    supports_vision=True,
                    input_price_per_m_token=1.0,
                    output_price_per_m_token=1.0,
                ),
                ModelConfig(
                    id="glm-4.5-air",
                    name="GLM-4.5 Air",
                    description="Lightweight GLM-4.5 variant",
                    max_tokens=128000,
                    supports_streaming=True,
                    supports_vision=False,
                    input_price_per_m_token=0.5,
                    output_price_per_m_token=0.5,
                ),
            ],
        ),
        ProviderConfig(
            name="anthropic",
            display_name="Anthropic",
            base_url="https://api.anthropic.com/v1",
            health_check_path="/messages",
            chat_completion_path="/messages",
            requires_api_key=True,
            models=[
                ModelConfig(
                    id="claude-opus-4-5-20251101",
                    name="Claude Opus 4.5",
                    description="Most capable Claude model",
                    max_tokens=200000,
                    supports_streaming=True,
                    supports_vision=True,
                    input_price_per_m_token=15.0,
                    output_price_per_m_token=75.0,
                ),
                ModelConfig(
                    id="claude-sonnet-4-5-20250929",
                    name="Claude Sonnet 4.5",
                    description="Balanced performance and speed",
                    max_tokens=200000,
                    supports_streaming=True,
                    supports_vision=True,
                    input_price_per_m_token=3.0,
                    output_price_per_m_token=15.0,
                ),
                ModelConfig(
                    id="claude-haiku-4-5-20251001",
                    name="Claude Haiku 4.5",
                    description="Fastest Claude model",
                    max_tokens=200000,
                    supports_streaming=True,
                    supports_vision=False,
                    input_price_per_m_token=0.8,
                    output_price_per_m_token=4.0,
                ),
            ],
        ),
        ProviderConfig(
            name="gemini",
            display_name="Google Gemini",
            base_url="https://generativelanguage.googleapis.com/v1beta",
            health_check_path="/models",
            chat_completion_path="/models/{model}:generateContent",
            requires_api_key=True,
            models=[
                ModelConfig(
                    id="gemini-2.0-flash-exp",
                    name="Gemini 2.0 Flash (Experimental)",
                    description="Latest experimental Gemini model",
                    max_tokens=1000000,
                    supports_streaming=True,
                    supports_vision=True,
                    input_price_per_m_token=0.0,
                    output_price_per_m_token=0.0,
                ),
                ModelConfig(
                    id="gemini-1.5-pro",
                    name="Gemini 1.5 Pro",
                    description="Most capable Gemini 1.5 model",
                    max_tokens=2000000,
                    supports_streaming=True,
                    supports_vision=True,
                    input_price_per_m_token=1.25,
                    output_price_per_m_token=5.0,
                ),
                ModelConfig(
                    id="gemini-1.5-flash",
                    name="Gemini 1.5 Flash",
                    description="Fast and efficient Gemini model",
                    max_tokens=1000000,
                    supports_streaming=True,
                    supports_vision=True,
                    input_price_per_m_token=0.075,
                    output_price_per_m_token=0.3,
                ),
            ],
        ),
        ProviderConfig(
            name="kimi",
            display_name="Kimi (Moonshot AI)",
            base_url="https://api.moonshot.cn/v1",
            health_check_path="/models",
            chat_completion_path="/chat/completions",
            requires_api_key=True,
            models=[
                ModelConfig(
                    id="moonshot-v1-128k",
                    name="Moonshot v1 128K",
                    description="Moonshot model with 128K context",
                    max_tokens=128000,
                    supports_streaming=True,
                    supports_vision=False,
                    input_price_per_m_token=12.0,
                    output_price_per_m_token=12.0,
                ),
                ModelConfig(
                    id="moonshot-v1-32k",
                    name="Moonshot v1 32K",
                    description="Moonshot model with 32K context",
                    max_tokens=32000,
                    supports_streaming=True,
                    supports_vision=False,
                    input_price_per_m_token=24.0,
                    output_price_per_m_token=24.0,
                ),
            ],
        ),
        ProviderConfig(
            name="openrouter",
            display_name="OpenRouter",
            base_url="https://openrouter.ai/api/v1",
            health_check_path="/models",
            chat_completion_path="/chat/completions",
            requires_api_key=True,
            models=[
                ModelConfig(
                    id="anthropic/claude-opus-4-5",
                    name="Claude Opus 4.5 (via OpenRouter)",
                    description="Access Claude Opus 4.5 through OpenRouter",
                    max_tokens=200000,
                    supports_streaming=True,
                    supports_vision=True,
                    input_price_per_m_token=15.0,
                    output_price_per_m_token=75.0,
                ),
                ModelConfig(
                    id="anthropic/claude-sonnet-4-5",
                    name="Claude Sonnet 4.5 (via OpenRouter)",
                    description="Access Claude Sonnet 4.5 through OpenRouter",
                    max_tokens=200000,
                    supports_streaming=True,
                    supports_vision=True,
                    input_price_per_m_token=3.0,
                    output_price_per_m_token=15.0,
                ),
                ModelConfig(
                    id="openai/gpt-4-turbo",
                    name="GPT-4 Turbo (via OpenRouter)",
                    description="Access GPT-4 Turbo through OpenRouter",
                    max_tokens=128000,
                    supports_streaming=True,
                    supports_vision=True,
                    input_price_per_m_token=10.0,
                    output_price_per_m_token=30.0,
                ),
                ModelConfig(
                    id="google/gemini-2.0-flash-exp:free",
                    name="Gemini 2.0 Flash (Free via OpenRouter)",
                    description="Free access to Gemini 2.0 Flash",
                    max_tokens=1000000,
                    supports_streaming=True,
                    supports_vision=True,
                    input_price_per_m_token=0.0,
                    output_price_per_m_token=0.0,
                ),
            ],
        ),
    ]


def get_provider_by_name(name: str) -> Optional[ProviderConfig]:
    """Returns a specific provider configuration by name."""
    providers = get_provider_registry()
    for provider in providers:
        if provider.name == name:
            return provider
    return None


def get_model_config(provider_name: str, model_id: str) -> Optional[ModelConfig]:
    """Returns a specific model configuration by provider name and model ID."""
    provider = get_provider_by_name(provider_name)
    if provider:
        for model in provider.models:
            if model.id == model_id:
                return model
    return None


# LangChain imports for model factory
try:
    from langchain_openai import ChatOpenAI
    from langchain_anthropic import ChatAnthropic
    from langchain_google_genai import ChatGoogleGenerativeAI
    LANGCHAIN_AVAILABLE = True
except ImportError:
    LANGCHAIN_AVAILABLE = False


def create_model_client(
    provider_name: str,
    model_id: str,
    api_key: str,
    temperature: float = 0.7,
    **kwargs
) -> Union["ChatOpenAI", "ChatAnthropic", "ChatGoogleGenerativeAI"]:
    """
    Factory function that creates and returns the appropriate LangChain client instance.
    
    Args:
        provider_name: The provider name (zai, anthropic, gemini, kimi, openrouter)
        model_id: The model ID to use
        api_key: The API key for authentication
        temperature: Sampling temperature (default: 0.7)
        **kwargs: Additional arguments to pass to the LangChain client
    
    Returns:
        A LangChain chat model instance compatible with Deep Agents
    
    Raises:
        ImportError: If required LangChain packages are not installed
        ValueError: If provider or model is not found
    
    Example:
        >>> model = create_model_client("anthropic", "claude-opus-4-5-20251101", "sk-...")
        >>> agent = create_deep_agent(model=model, ...)
    """
    if not LANGCHAIN_AVAILABLE:
        raise ImportError(
            "LangChain packages are required. Install with: "
            "pip install langchain-openai langchain-anthropic langchain-google-genai"
        )
    
    provider = get_provider_by_name(provider_name)
    if not provider:
        raise ValueError(f"Provider '{provider_name}' not found in registry")
    
    model_config = get_model_config(provider_name, model_id)
    if not model_config:
        raise ValueError(f"Model '{model_id}' not found for provider '{provider_name}'")
    
    # OpenAI-compatible providers (zai, kimi, openrouter)
    if provider_name in ("zai", "kimi", "openrouter"):
        return ChatOpenAI(
            model=model_id,
            api_key=api_key,
            base_url=provider.base_url,
            temperature=temperature,
            max_tokens=model_config.max_tokens,
            **kwargs
        )
    
    # Anthropic
    elif provider_name == "anthropic":
        return ChatAnthropic(
            model=model_id,
            api_key=api_key,
            temperature=temperature,
            max_tokens=model_config.max_tokens,
            **kwargs
        )
    
    # Google Gemini
    elif provider_name == "gemini":
        return ChatGoogleGenerativeAI(
            model=model_id,
            api_key=api_key,
            temperature=temperature,
            max_output_tokens=model_config.max_tokens,
            **kwargs
        )
    
    else:
        raise ValueError(f"Unsupported provider: {provider_name}")
