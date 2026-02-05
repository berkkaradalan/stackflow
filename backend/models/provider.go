package models

// ProviderConfig represents a provider configuration
type ProviderConfig struct {
	Name               string        `json:"name"`
	DisplayName        string        `json:"display_name"`
	BaseURL            string        `json:"base_url"`
	HealthCheckPath    string        `json:"health_check_path"`
	ChatCompletionPath string        `json:"chat_completion_path"`
	RequiresAPIKey     bool          `json:"requires_api_key"`
	Models             []ModelConfig `json:"models"`
}

// ModelConfig represents a model configuration
type ModelConfig struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	MaxTokens          int     `json:"max_tokens"`
	SupportsStreaming  bool    `json:"supports_streaming"`
	SupportsVision     bool    `json:"supports_vision"`
	InputPricePerMToken  float64 `json:"input_price_per_m_token"`
	OutputPricePerMToken float64 `json:"output_price_per_m_token"`
}

// ProviderListResponse is the response model for listing providers
type ProviderListResponse struct {
	Providers []ProviderConfig `json:"providers"`
}

// ModelListResponse is the response model for listing models
type ModelListResponse struct {
	Provider string        `json:"provider"`
	Models   []ModelConfig `json:"models"`
}
