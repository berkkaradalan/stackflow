package models

import "time"

// Agent represents an AI agent in the system
type Agent struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ProjectID   int        `json:"project_id"`
	CreatedBy   int        `json:"created_by"`
	Role        string     `json:"role"`
	Level       string     `json:"level"`
	Provider    string     `json:"provider"`
	Model       string     `json:"model"`
	APIKey      string     `json:"-"`
	Config      AgentConfig `json:"config"`
	Status      string     `json:"status"`
	IsActive    bool       `json:"is_active"`
	LastActiveAt *time.Time `json:"last_active_at,omitempty"`
	TotalTokensUsed int64   `json:"total_tokens_used"`
	TotalCost       float64 `json:"total_cost"`
	TotalRequests   int64   `json:"total_requests"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// AgentConfig holds the configuration for an agent
type AgentConfig struct {
	Temperature      float64 `json:"temperature"`
	MaxTokens        int     `json:"max_tokens"`
	TopP             float64 `json:"top_p"`
	FrequencyPenalty float64 `json:"frequency_penalty"`
	PresencePenalty  float64 `json:"presence_penalty"`
}

// CreateAgentRequest is the request model for creating an agent
type CreateAgentRequest struct {
	Name        string      `json:"name" binding:"required,min=3,max=100"`
	Description string      `json:"description" binding:"omitempty,max=500"`
	ProjectID   int         `json:"project_id" binding:"required"`
	Role        string      `json:"role" binding:"required,oneof=backend_developer frontend_developer fullstack_developer tester devops project_manager"`
	Level       string      `json:"level" binding:"required,oneof=junior mid senior"`
	Provider    string      `json:"provider" binding:"required,oneof=openrouter openai anthropic gemini groq together glm claude kimi"`
	Model       string      `json:"model" binding:"required"`
	APIKey      string      `json:"api_key" binding:"required"`
	Config      AgentConfig `json:"config"`
}

// UpdateAgentRequest is the request model for updating an agent
type UpdateAgentRequest struct {
	Name        *string      `json:"name" binding:"omitempty,min=3,max=100"`
	Description *string      `json:"description" binding:"omitempty,max=500"`
	Role        *string      `json:"role" binding:"omitempty,oneof=backend_developer frontend_developer fullstack_developer tester devops project_manager"`
	Level       *string      `json:"level" binding:"omitempty,oneof=junior mid senior"`
	Provider    *string      `json:"provider" binding:"omitempty,oneof=openrouter openai anthropic gemini groq together glm claude kimi"`
	Model       *string      `json:"model" binding:"omitempty"`
	APIKey      *string      `json:"api_key" binding:"omitempty"`
	Config      *AgentConfig `json:"config" binding:"omitempty"`
	Status      *string      `json:"status" binding:"omitempty,oneof=idle active busy error disabled initializing"`
	IsActive    *bool        `json:"is_active"`
}

// AgentListResponse is the response model for listing agents
type AgentListResponse struct {
	Agents     []Agent `json:"agents"`
	TotalCount int     `json:"total_count"`
}

// AgentStatusResponse is the response model for agent status
type AgentStatusResponse struct {
	Status       string     `json:"status"`
	IsActive     bool       `json:"is_active"`
	LastActiveAt *time.Time `json:"last_active_at,omitempty"`
}

// AgentWorkloadResponse is the response model for agent workload (total usage metrics)
type AgentWorkloadResponse struct {
	TotalRequests   int64      `json:"total_requests"`
	TotalTokensUsed int64      `json:"total_tokens_used"`
	TotalCost       float64    `json:"total_cost"`
	LastActiveAt    *time.Time `json:"last_active_at,omitempty"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Description     string     `json:"description"`
}

// AgentPerformanceResponse is the response model for agent performance metrics (averages & efficiency)
type AgentPerformanceResponse struct {
	TotalRequests     int64      `json:"total_requests"`
	TotalTokensUsed   int64      `json:"total_tokens_used"`
	TotalCost         float64    `json:"total_cost"`
	AverageTokensPerRequest float64 `json:"average_tokens_per_request"`
	AverageCostPerRequest   float64 `json:"average_cost_per_request"`
	LastActiveAt    *time.Time `json:"last_active_at,omitempty"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Description     string     `json:"description"`
}

// AgentHealthResponse is the response model for agent health check
type AgentHealthResponse struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Status       string     `json:"status"`
	IsActive     bool       `json:"is_active"`
	LastActiveAt *time.Time `json:"last_active_at,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Healthy      bool       `json:"healthy"`
	Message      string     `json:"message"`
}
