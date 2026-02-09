package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/berkkaradalan/stackflow/config"
	"github.com/berkkaradalan/stackflow/models"
	"github.com/berkkaradalan/stackflow/repository/postgres"
)

type AgentService struct {
	agentRepo *repository.AgentRepository
}

func NewAgentService(agentRepo *repository.AgentRepository) *AgentService {
	return &AgentService{
		agentRepo: agentRepo,
	}
}

func (s *AgentService) CreateAgent(ctx context.Context, req *models.CreateAgentRequest, userID int) (*models.Agent, error) {
	agent := &models.Agent{
		Name:        req.Name,
		Description: req.Description,
		ProjectID:   req.ProjectID,
		CreatedBy:   userID,
		Role:        req.Role,
		Level:       req.Level,
		Provider:    req.Provider,
		Model:       req.Model,
		APIKey:      req.APIKey, // TODO: Encrypt this in production
		Config:      req.Config,
		Status:      "idle",
		IsActive:    true,
		TotalTokensUsed: 0,
		TotalCost:       0.0,
		TotalRequests:   0,
	}

	// Set default config values if not provided
	if agent.Config.Temperature == 0 {
		agent.Config.Temperature = 0.7
	}
	if agent.Config.MaxTokens == 0 {
		agent.Config.MaxTokens = 2000
	}
	if agent.Config.TopP == 0 {
		agent.Config.TopP = 1.0
	}

	err := s.agentRepo.Create(ctx, agent)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	return agent, nil
}

func (s *AgentService) GetAllAgents(ctx context.Context) (*models.AgentListResponse, error) {
	agents, err := s.agentRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get agents: %w", err)
	}

	// Remove API keys from response
	for i := range agents {
		agents[i].APIKey = ""
	}

	return &models.AgentListResponse{
		Agents:     agents,
		TotalCount: len(agents),
	}, nil
}

func (s *AgentService) GetAgentByID(ctx context.Context, id int) (*models.Agent, error) {
	agent, err := s.agentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}

	return agent, nil
}

func (s *AgentService) GetAgentsByProjectID(ctx context.Context, projectID int) (*models.AgentListResponse, error) {
	agents, err := s.agentRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agents by project: %w", err)
	}

	// Remove API keys from response
	for i := range agents {
		agents[i].APIKey = ""
	}

	return &models.AgentListResponse{
		Agents:     agents,
		TotalCount: len(agents),
	}, nil
}

func (s *AgentService) UpdateAgent(ctx context.Context, id int, req *models.UpdateAgentRequest) (*models.Agent, error) {
	// First check if agent exists
	_, err := s.agentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if req.Role != nil {
		updates["role"] = *req.Role
	}

	if req.Level != nil {
		updates["level"] = *req.Level
	}

	if req.Provider != nil {
		updates["provider"] = *req.Provider
	}

	if req.Model != nil {
		updates["model"] = *req.Model
	}

	if req.APIKey != nil {
		updates["api_key"] = *req.APIKey // TODO: Encrypt this in production
	}

	if req.Config != nil {
		updates["config"] = *req.Config
	}

	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	// Perform partial update
	updatedAgent, err := s.agentRepo.UpdatePartial(ctx, id, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update agent: %w", err)
	}

	return updatedAgent, nil
}

func (s *AgentService) DeleteAgent(ctx context.Context, id int) error {
	// Check if agent exists
	_, err := s.agentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}

	err = s.agentRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete agent: %w", err)
	}

	return nil
}

func (s *AgentService) GetAgentStatus(ctx context.Context, id int) (*models.AgentStatusResponse, error) {
	status, err := s.agentRepo.GetStatus(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent status: %w", err)
	}

	return status, nil
}

func (s *AgentService) GetAgentWorkload(ctx context.Context, id int) (*models.AgentWorkloadResponse, error) {
	workload, err := s.agentRepo.GetWorkload(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent workload: %w", err)
	}

	return workload, nil
}

func (s *AgentService) GetAgentPerformance(ctx context.Context, id int) (*models.AgentPerformanceResponse, error) {
	performance, err := s.agentRepo.GetPerformance(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent performance: %w", err)
	}

	return performance, nil
}

func (s *AgentService) UpdateAgentStatus(ctx context.Context, id int, status string) error {
	// Check if agent exists
	_, err := s.agentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}

	err = s.agentRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return fmt.Errorf("failed to update agent status: %w", err)
	}

	return nil
}

func (s *AgentService) IncrementAgentUsage(ctx context.Context, id int, tokensUsed int64, cost float64) error {
	err := s.agentRepo.IncrementUsage(ctx, id, tokensUsed, cost)
	if err != nil {
		return fmt.Errorf("failed to increment agent usage: %w", err)
	}

	return nil
}

func (s *AgentService) HealthCheck(ctx context.Context, id int) (*models.AgentHealthResponse, error) {
	health, err := s.agentRepo.HealthCheck(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to perform health check: %w", err)
	}

	return health, nil
}

// PerformRealHealthCheck performs a real API test to the provider
func (s *AgentService) PerformRealHealthCheck(ctx context.Context, id int) (*models.AgentHealthResponse, error) {
	// Get agent details
	agent, err := s.agentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}

	// Get provider configuration
	providerConfig := config.GetProviderByName(agent.Provider)
	if providerConfig == nil {
		return nil, fmt.Errorf("provider configuration not found for: %s", agent.Provider)
	}

	// Perform real API health check with a test message
	isHealthy, testResponse, errorMsg := s.testProviderAPIWithMessage(providerConfig, agent)

	// Update agent status based on test result
	var newStatus string
	var message string

	if isHealthy {
		newStatus = "active"
		message = "Agent is healthy and operational - API test successful"
		// Update status to active
		_ = s.agentRepo.UpdateStatus(ctx, id, newStatus)
	} else {
		newStatus = "error"
		message = fmt.Sprintf("Agent health check failed: %s", errorMsg)
		// Update status to error
		_ = s.agentRepo.UpdateStatus(ctx, id, newStatus)
	}

	// Return health response
	now := time.Now()
	return &models.AgentHealthResponse{
		ID:           agent.ID,
		Name:         agent.Name,
		Status:       newStatus,
		IsActive:     agent.IsActive,
		LastActiveAt: &now,
		UpdatedAt:    now,
		Healthy:      isHealthy,
		Message:      message,
		TestResponse: testResponse,
	}, nil
}

// testProviderAPIWithMessage sends a real test message to the AI and returns its response
func (s *AgentService) testProviderAPIWithMessage(providerConfig *models.ProviderConfig, agent *models.Agent) (bool, string, string) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// Construct chat completion URL
	chatURL := providerConfig.BaseURL + providerConfig.ChatCompletionPath

	// Prepare test message payload
	payload := map[string]interface{}{
		"model": agent.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": "This is a health check test. Please respond with a brief confirmation that you are operational.",
			},
		},
		"max_tokens":  50,
		"temperature": 0.7,
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return false, "", fmt.Sprintf("failed to marshal payload: %v", err)
	}

	// Create request
	req, err := http.NewRequest("POST", chatURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return false, "", fmt.Sprintf("failed to create request: %v", err)
	}

	// Add authorization header
	req.Header.Set("Authorization", "Bearer "+agent.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		return false, "", fmt.Sprintf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, "", fmt.Sprintf("API returned status code: %d", resp.StatusCode)
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, "", fmt.Sprintf("failed to decode response: %v", err)
	}

	// Extract AI response (OpenAI-compatible format)
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return true, content, ""
				}
			}
		}
	}

	return false, "", "unable to extract AI response from API"
}

// testProviderAPI tests if the provider API is accessible and working
func (s *AgentService) testProviderAPI(providerConfig *models.ProviderConfig, apiKey string) (bool, string) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Construct health check URL
	healthCheckURL := providerConfig.BaseURL + providerConfig.HealthCheckPath

	// Create request
	req, err := http.NewRequest("GET", healthCheckURL, nil)
	if err != nil {
		return false, fmt.Sprintf("failed to create request: %v", err)
	}

	// Add authorization header
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Sprintf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, ""
	}

	return false, fmt.Sprintf("API returned status code: %d", resp.StatusCode)
}
