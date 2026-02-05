package service

import (
	"context"
	"fmt"

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

	// Remove API key from response
	agent.APIKey = ""

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

	// Remove API key from response
	updatedAgent.APIKey = ""

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
