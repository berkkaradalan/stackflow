package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/berkkaradalan/stackflow/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AgentRepository struct {
	pool *pgxpool.Pool
}

func NewAgentRepository(pool *pgxpool.Pool) *AgentRepository {
	return &AgentRepository{
		pool: pool,
	}
}

func (r *AgentRepository) Create(ctx context.Context, agent *models.Agent) error {
	configJSON, err := json.Marshal(agent.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	query := `INSERT INTO agents (name, description, project_id, created_by, role, level, provider, model, api_key, config, status, is_active)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	          RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		agent.Name, agent.Description, agent.ProjectID, agent.CreatedBy,
		agent.Role, agent.Level, agent.Provider, agent.Model, agent.APIKey,
		configJSON, agent.Status, agent.IsActive,
	).Scan(&agent.ID, &agent.CreatedAt, &agent.UpdatedAt)
}

func (r *AgentRepository) GetByID(ctx context.Context, id int) (*models.Agent, error) {
	query := `SELECT id, name, description, project_id, created_by, role, level, provider, model, api_key, config,
	          status, is_active, last_active_at, total_tokens_used, total_cost, total_requests, created_at, updated_at
	          FROM agents WHERE id = $1`

	var agent models.Agent
	var configJSON []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&agent.ID, &agent.Name, &agent.Description, &agent.ProjectID, &agent.CreatedBy,
		&agent.Role, &agent.Level, &agent.Provider, &agent.Model, &agent.APIKey,
		&configJSON, &agent.Status, &agent.IsActive, &agent.LastActiveAt,
		&agent.TotalTokensUsed, &agent.TotalCost, &agent.TotalRequests,
		&agent.CreatedAt, &agent.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(configJSON, &agent.Config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &agent, nil
}

func (r *AgentRepository) GetAll(ctx context.Context) ([]models.Agent, error) {
	query := `SELECT id, name, description, project_id, created_by, role, level, provider, model, api_key, config,
	          status, is_active, last_active_at, total_tokens_used, total_cost, total_requests, created_at, updated_at
	          FROM agents
	          ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []models.Agent
	for rows.Next() {
		var agent models.Agent
		var configJSON []byte

		err := rows.Scan(
			&agent.ID, &agent.Name, &agent.Description, &agent.ProjectID, &agent.CreatedBy,
			&agent.Role, &agent.Level, &agent.Provider, &agent.Model, &agent.APIKey,
			&configJSON, &agent.Status, &agent.IsActive, &agent.LastActiveAt,
			&agent.TotalTokensUsed, &agent.TotalCost, &agent.TotalRequests,
			&agent.CreatedAt, &agent.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(configJSON, &agent.Config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		agents = append(agents, agent)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return agents, nil
}

func (r *AgentRepository) GetByProjectID(ctx context.Context, projectID int) ([]models.Agent, error) {
	query := `SELECT id, name, description, project_id, created_by, role, level, provider, model, api_key, config,
	          status, is_active, last_active_at, total_tokens_used, total_cost, total_requests, created_at, updated_at
	          FROM agents
	          WHERE project_id = $1
	          ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []models.Agent
	for rows.Next() {
		var agent models.Agent
		var configJSON []byte

		err := rows.Scan(
			&agent.ID, &agent.Name, &agent.Description, &agent.ProjectID, &agent.CreatedBy,
			&agent.Role, &agent.Level, &agent.Provider, &agent.Model, &agent.APIKey,
			&configJSON, &agent.Status, &agent.IsActive, &agent.LastActiveAt,
			&agent.TotalTokensUsed, &agent.TotalCost, &agent.TotalRequests,
			&agent.CreatedAt, &agent.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(configJSON, &agent.Config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		agents = append(agents, agent)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return agents, nil
}

func (r *AgentRepository) UpdatePartial(ctx context.Context, id int, updates map[string]interface{}) (*models.Agent, error) {
	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}

	// Handle config field specially
	if config, ok := updates["config"]; ok {
		configJSON, err := json.Marshal(config)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal config: %w", err)
		}
		updates["config"] = configJSON
	}

	query := "UPDATE agents SET "
	args := make([]interface{}, 0, len(updates)+1)
	argPos := 1

	first := true
	for key, value := range updates {
		if !first {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", key, argPos)
		args = append(args, value)
		argPos++
		first = false
	}

	query += fmt.Sprintf(", updated_at = NOW() WHERE id = $%d", argPos)
	args = append(args, id)

	_, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *AgentRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM agents WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *AgentRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE agents SET status = $1, last_active_at = NOW(), updated_at = NOW() WHERE id = $2`

	_, err := r.pool.Exec(ctx, query, status, id)
	return err
}

func (r *AgentRepository) IncrementUsage(ctx context.Context, id int, tokensUsed int64, cost float64) error {
	query := `UPDATE agents
	          SET total_tokens_used = total_tokens_used + $1,
	              total_cost = total_cost + $2,
	              total_requests = total_requests + 1,
	              last_active_at = NOW(),
	              updated_at = NOW()
	          WHERE id = $3`

	_, err := r.pool.Exec(ctx, query, tokensUsed, cost, id)
	return err
}

func (r *AgentRepository) GetStatus(ctx context.Context, id int) (*models.AgentStatusResponse, error) {
	query := `SELECT status, is_active, last_active_at FROM agents WHERE id = $1`

	var status models.AgentStatusResponse
	err := r.pool.QueryRow(ctx, query, id).Scan(&status.Status, &status.IsActive, &status.LastActiveAt)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (r *AgentRepository) GetWorkload(ctx context.Context, id int) (*models.AgentWorkloadResponse, error) {
	query := `SELECT total_requests, total_tokens_used, total_cost, last_active_at, updated_at FROM agents WHERE id = $1`

	var workload models.AgentWorkloadResponse
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&workload.TotalRequests,
		&workload.TotalTokensUsed,
		&workload.TotalCost,
		&workload.LastActiveAt,
		&workload.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	workload.Description = fmt.Sprintf("Agent has processed %d requests, used %d tokens, and cost $%.4f in total",
		workload.TotalRequests, workload.TotalTokensUsed, workload.TotalCost)

	return &workload, nil
}

func (r *AgentRepository) GetPerformance(ctx context.Context, id int) (*models.AgentPerformanceResponse, error) {
	query := `SELECT total_requests, total_tokens_used, total_cost, last_active_at, updated_at FROM agents WHERE id = $1`

	var perf models.AgentPerformanceResponse
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&perf.TotalRequests,
		&perf.TotalTokensUsed,
		&perf.TotalCost,
		&perf.LastActiveAt,
		&perf.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Calculate averages
	if perf.TotalRequests > 0 {
		perf.AverageTokensPerRequest = float64(perf.TotalTokensUsed) / float64(perf.TotalRequests)
		perf.AverageCostPerRequest = perf.TotalCost / float64(perf.TotalRequests)
		perf.Description = fmt.Sprintf("Agent averages %.2f tokens and $%.4f per request",
			perf.AverageTokensPerRequest, perf.AverageCostPerRequest)
	} else {
		perf.Description = "No requests processed yet"
	}

	return &perf, nil
}

func (r *AgentRepository) HealthCheck(ctx context.Context, id int) (*models.AgentHealthResponse, error) {
	query := `SELECT id, name, status, is_active, last_active_at, updated_at FROM agents WHERE id = $1`

	var health models.AgentHealthResponse
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&health.ID,
		&health.Name,
		&health.Status,
		&health.IsActive,
		&health.LastActiveAt,
		&health.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Determine health status
	health.Healthy = health.IsActive && health.Status != "error" && health.Status != "disabled"

	if health.Healthy {
		health.Message = "Agent is healthy and operational"
	} else if health.Status == "error" {
		health.Message = "Agent encountered an error"
	} else if health.Status == "disabled" {
		health.Message = "Agent is disabled"
	} else if !health.IsActive {
		health.Message = "Agent is inactive"
	}

	// Update last_active_at and status if agent is being checked
	if health.Healthy && health.Status == "idle" {
		updateQuery := `UPDATE agents SET status = 'active', last_active_at = NOW(), updated_at = NOW() WHERE id = $1`
		_, _ = r.pool.Exec(ctx, updateQuery, id)
		health.Status = "active"
	}

	return &health, nil
}
