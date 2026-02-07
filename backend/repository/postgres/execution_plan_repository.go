package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/berkkaradalan/stackflow/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExecutionPlanRepository struct {
	pool *pgxpool.Pool
}

func NewExecutionPlanRepository(pool *pgxpool.Pool) *ExecutionPlanRepository {
	return &ExecutionPlanRepository{
		pool: pool,
	}
}

// --- Execution Plans ---

// CreatePlan creates a new execution plan
func (r *ExecutionPlanRepository) CreatePlan(ctx context.Context, plan *models.ExecutionPlan) error {
	planDataJSON, err := json.Marshal(plan.PlanData)
	if err != nil {
		return fmt.Errorf("failed to marshal plan_data: %w", err)
	}

	query := `INSERT INTO execution_plans (project_id, created_by, creator_type, plan_data, status)
	          VALUES ($1, $2, $3, $4, $5)
	          RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		plan.ProjectID, plan.CreatedBy, plan.CreatorType, planDataJSON, plan.Status,
	).Scan(&plan.ID, &plan.CreatedAt, &plan.UpdatedAt)
}

// GetPlanByID retrieves an execution plan by ID
func (r *ExecutionPlanRepository) GetPlanByID(ctx context.Context, id int) (*models.ExecutionPlan, error) {
	query := `SELECT id, project_id, created_by, creator_type, plan_data, status, created_at, updated_at
	          FROM execution_plans WHERE id = $1`

	var plan models.ExecutionPlan
	var planDataJSON []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&plan.ID, &plan.ProjectID, &plan.CreatedBy, &plan.CreatorType,
		&planDataJSON, &plan.Status, &plan.CreatedAt, &plan.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(planDataJSON, &plan.PlanData); err != nil {
		plan.PlanData = models.PlanData{}
	}

	return &plan, nil
}

// GetPlanByIDWithDetails retrieves a plan with project and creator names
func (r *ExecutionPlanRepository) GetPlanByIDWithDetails(ctx context.Context, id int) (*models.ExecutionPlanWithDetails, error) {
	query := `SELECT
		ep.id, ep.project_id, ep.created_by, ep.creator_type, ep.plan_data, ep.status,
		ep.created_at, ep.updated_at,
		p.name as project_name,
		CASE
			WHEN ep.creator_type = 'user' THEN (SELECT username FROM users WHERE id = ep.created_by)
			ELSE (SELECT name FROM agents WHERE id = ep.created_by)
		END as creator_name
	FROM execution_plans ep
	LEFT JOIN projects p ON ep.project_id = p.id
	WHERE ep.id = $1`

	var plan models.ExecutionPlanWithDetails
	var planDataJSON []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&plan.ID, &plan.ProjectID, &plan.CreatedBy, &plan.CreatorType,
		&planDataJSON, &plan.Status, &plan.CreatedAt, &plan.UpdatedAt,
		&plan.ProjectName, &plan.CreatorName,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(planDataJSON, &plan.PlanData); err != nil {
		plan.PlanData = models.PlanData{}
	}

	return &plan, nil
}

// GetActivePlanByProjectID retrieves the active execution plan for a project
func (r *ExecutionPlanRepository) GetActivePlanByProjectID(ctx context.Context, projectID int) (*models.ExecutionPlanWithDetails, error) {
	query := `SELECT
		ep.id, ep.project_id, ep.created_by, ep.creator_type, ep.plan_data, ep.status,
		ep.created_at, ep.updated_at,
		p.name as project_name,
		CASE
			WHEN ep.creator_type = 'user' THEN (SELECT username FROM users WHERE id = ep.created_by)
			ELSE (SELECT name FROM agents WHERE id = ep.created_by)
		END as creator_name
	FROM execution_plans ep
	LEFT JOIN projects p ON ep.project_id = p.id
	WHERE ep.project_id = $1 AND ep.status = 'active'
	ORDER BY ep.created_at DESC
	LIMIT 1`

	var plan models.ExecutionPlanWithDetails
	var planDataJSON []byte
	err := r.pool.QueryRow(ctx, query, projectID).Scan(
		&plan.ID, &plan.ProjectID, &plan.CreatedBy, &plan.CreatorType,
		&planDataJSON, &plan.Status, &plan.CreatedAt, &plan.UpdatedAt,
		&plan.ProjectName, &plan.CreatorName,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(planDataJSON, &plan.PlanData); err != nil {
		plan.PlanData = models.PlanData{}
	}

	return &plan, nil
}

// GetPlansByProjectID retrieves all plans for a project
func (r *ExecutionPlanRepository) GetPlansByProjectID(ctx context.Context, projectID int) ([]models.ExecutionPlanWithDetails, error) {
	query := `SELECT
		ep.id, ep.project_id, ep.created_by, ep.creator_type, ep.plan_data, ep.status,
		ep.created_at, ep.updated_at,
		p.name as project_name,
		CASE
			WHEN ep.creator_type = 'user' THEN (SELECT username FROM users WHERE id = ep.created_by)
			ELSE (SELECT name FROM agents WHERE id = ep.created_by)
		END as creator_name
	FROM execution_plans ep
	LEFT JOIN projects p ON ep.project_id = p.id
	WHERE ep.project_id = $1
	ORDER BY ep.created_at DESC`

	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []models.ExecutionPlanWithDetails
	for rows.Next() {
		var plan models.ExecutionPlanWithDetails
		var planDataJSON []byte
		err := rows.Scan(
			&plan.ID, &plan.ProjectID, &plan.CreatedBy, &plan.CreatorType,
			&planDataJSON, &plan.Status, &plan.CreatedAt, &plan.UpdatedAt,
			&plan.ProjectName, &plan.CreatorName,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(planDataJSON, &plan.PlanData); err != nil {
			plan.PlanData = models.PlanData{}
		}

		plans = append(plans, plan)
	}

	return plans, rows.Err()
}

// UpdatePlan updates an execution plan
func (r *ExecutionPlanRepository) UpdatePlan(ctx context.Context, id int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	query := "UPDATE execution_plans SET "
	args := make([]interface{}, 0, len(updates)+1)
	argPos := 1

	first := true
	for key, value := range updates {
		if !first {
			query += ", "
		}
		if key == "plan_data" {
			dataJSON, err := json.Marshal(value)
			if err != nil {
				return fmt.Errorf("failed to marshal plan_data: %w", err)
			}
			value = dataJSON
		}
		query += fmt.Sprintf("%s = $%d", key, argPos)
		args = append(args, value)
		argPos++
		first = false
	}

	query += fmt.Sprintf(", updated_at = NOW() WHERE id = $%d", argPos)
	args = append(args, id)

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// DeletePlan deletes an execution plan
func (r *ExecutionPlanRepository) DeletePlan(ctx context.Context, id int) error {
	query := `DELETE FROM execution_plans WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// --- Agent Assignments ---

// CreateAssignment creates a new agent assignment
func (r *ExecutionPlanRepository) CreateAssignment(ctx context.Context, assignment *models.AgentAssignment) error {
	query := `INSERT INTO agent_assignments (plan_id, agent_id, task_id, status)
	          VALUES ($1, $2, $3, $4)
	          RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		assignment.PlanID, assignment.AgentID, assignment.TaskID, assignment.Status,
	).Scan(&assignment.ID, &assignment.CreatedAt, &assignment.UpdatedAt)
}

// GetNextAssignmentForAgent finds the next pending assignment for an agent
func (r *ExecutionPlanRepository) GetNextAssignmentForAgent(ctx context.Context, agentID int) (*models.AgentAssignmentWithDetails, error) {
	query := `SELECT
		aa.id, aa.plan_id, aa.agent_id, aa.task_id, aa.status,
		aa.started_at, aa.completed_at, aa.report_data, aa.created_at, aa.updated_at,
		ag.name as agent_name,
		t.title as task_title
	FROM agent_assignments aa
	LEFT JOIN agents ag ON aa.agent_id = ag.id
	LEFT JOIN tasks t ON aa.task_id = t.id
	WHERE aa.agent_id = $1 AND aa.status = 'pending'
	ORDER BY aa.created_at ASC
	LIMIT 1`

	var assignment models.AgentAssignmentWithDetails
	var reportDataJSON []byte
	err := r.pool.QueryRow(ctx, query, agentID).Scan(
		&assignment.ID, &assignment.PlanID, &assignment.AgentID, &assignment.TaskID,
		&assignment.Status, &assignment.StartedAt, &assignment.CompletedAt,
		&reportDataJSON, &assignment.CreatedAt, &assignment.UpdatedAt,
		&assignment.AgentName, &assignment.TaskTitle,
	)
	if err != nil {
		return nil, err
	}

	if reportDataJSON != nil {
		_ = json.Unmarshal(reportDataJSON, &assignment.ReportData)
	}

	return &assignment, nil
}

// GetAssignmentsByPlanID retrieves all assignments for a plan
func (r *ExecutionPlanRepository) GetAssignmentsByPlanID(ctx context.Context, planID int) ([]models.AgentAssignmentWithDetails, error) {
	query := `SELECT
		aa.id, aa.plan_id, aa.agent_id, aa.task_id, aa.status,
		aa.started_at, aa.completed_at, aa.report_data, aa.created_at, aa.updated_at,
		ag.name as agent_name,
		t.title as task_title
	FROM agent_assignments aa
	LEFT JOIN agents ag ON aa.agent_id = ag.id
	LEFT JOIN tasks t ON aa.task_id = t.id
	WHERE aa.plan_id = $1
	ORDER BY aa.created_at ASC`

	rows, err := r.pool.Query(ctx, query, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []models.AgentAssignmentWithDetails
	for rows.Next() {
		var assignment models.AgentAssignmentWithDetails
		var reportDataJSON []byte
		err := rows.Scan(
			&assignment.ID, &assignment.PlanID, &assignment.AgentID, &assignment.TaskID,
			&assignment.Status, &assignment.StartedAt, &assignment.CompletedAt,
			&reportDataJSON, &assignment.CreatedAt, &assignment.UpdatedAt,
			&assignment.AgentName, &assignment.TaskTitle,
		)
		if err != nil {
			return nil, err
		}

		if reportDataJSON != nil {
			_ = json.Unmarshal(reportDataJSON, &assignment.ReportData)
		}

		assignments = append(assignments, assignment)
	}

	return assignments, rows.Err()
}

// CompleteAssignment marks an assignment as completed with report data
func (r *ExecutionPlanRepository) CompleteAssignment(ctx context.Context, assignmentID int, reportData any) error {
	reportJSON, err := json.Marshal(reportData)
	if err != nil {
		return fmt.Errorf("failed to marshal report_data: %w", err)
	}

	query := `UPDATE agent_assignments
	          SET status = 'completed', completed_at = NOW(), report_data = $1, updated_at = NOW()
	          WHERE id = $2`

	_, err = r.pool.Exec(ctx, query, reportJSON, assignmentID)
	return err
}

// StartAssignment marks an assignment as in_progress
func (r *ExecutionPlanRepository) StartAssignment(ctx context.Context, assignmentID int) error {
	query := `UPDATE agent_assignments
	          SET status = 'in_progress', started_at = NOW(), updated_at = NOW()
	          WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, assignmentID)
	return err
}

// GetAssignmentByID retrieves an assignment by ID
func (r *ExecutionPlanRepository) GetAssignmentByID(ctx context.Context, id int) (*models.AgentAssignment, error) {
	query := `SELECT id, plan_id, agent_id, task_id, status, started_at, completed_at, report_data, created_at, updated_at
	          FROM agent_assignments WHERE id = $1`

	var assignment models.AgentAssignment
	var reportDataJSON []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&assignment.ID, &assignment.PlanID, &assignment.AgentID, &assignment.TaskID,
		&assignment.Status, &assignment.StartedAt, &assignment.CompletedAt,
		&reportDataJSON, &assignment.CreatedAt, &assignment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if reportDataJSON != nil {
		_ = json.Unmarshal(reportDataJSON, &assignment.ReportData)
	}

	return &assignment, nil
}

// GetAssignmentByAgentAndTask finds assignment by agent and task
func (r *ExecutionPlanRepository) GetAssignmentByAgentAndTask(ctx context.Context, agentID int, taskID int) (*models.AgentAssignment, error) {
	query := `SELECT id, plan_id, agent_id, task_id, status, started_at, completed_at, report_data, created_at, updated_at
	          FROM agent_assignments
	          WHERE agent_id = $1 AND task_id = $2 AND status IN ('pending', 'in_progress')
	          ORDER BY created_at DESC LIMIT 1`

	var assignment models.AgentAssignment
	var reportDataJSON []byte
	err := r.pool.QueryRow(ctx, query, agentID, taskID).Scan(
		&assignment.ID, &assignment.PlanID, &assignment.AgentID, &assignment.TaskID,
		&assignment.Status, &assignment.StartedAt, &assignment.CompletedAt,
		&reportDataJSON, &assignment.CreatedAt, &assignment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if reportDataJSON != nil {
		_ = json.Unmarshal(reportDataJSON, &assignment.ReportData)
	}

	return &assignment, nil
}

// --- Execution Reports ---

// CreateReport creates a new execution report
func (r *ExecutionPlanRepository) CreateReport(ctx context.Context, report *models.ExecutionReport) error {
	reportJSON, err := json.Marshal(report.ReportData)
	if err != nil {
		return fmt.Errorf("failed to marshal report_data: %w", err)
	}

	query := `INSERT INTO execution_reports (project_id, report_type, generated_by, generator_type, report_data)
	          VALUES ($1, $2, $3, $4, $5)
	          RETURNING id, created_at`

	return r.pool.QueryRow(ctx, query,
		report.ProjectID, report.ReportType, report.GeneratedBy, report.GeneratorType, reportJSON,
	).Scan(&report.ID, &report.CreatedAt)
}

// GetReportsByProjectID retrieves all reports for a project
func (r *ExecutionPlanRepository) GetReportsByProjectID(ctx context.Context, projectID int) ([]models.ExecutionReportWithDetails, error) {
	query := `SELECT
		er.id, er.project_id, er.report_type, er.generated_by, er.generator_type, er.report_data, er.created_at,
		p.name as project_name,
		CASE
			WHEN er.generator_type = 'user' THEN (SELECT username FROM users WHERE id = er.generated_by)
			ELSE (SELECT name FROM agents WHERE id = er.generated_by)
		END as generator_name
	FROM execution_reports er
	LEFT JOIN projects p ON er.project_id = p.id
	WHERE er.project_id = $1
	ORDER BY er.created_at DESC`

	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.ExecutionReportWithDetails
	for rows.Next() {
		var report models.ExecutionReportWithDetails
		var reportDataJSON []byte
		err := rows.Scan(
			&report.ID, &report.ProjectID, &report.ReportType, &report.GeneratedBy,
			&report.GeneratorType, &reportDataJSON, &report.CreatedAt,
			&report.ProjectName, &report.GeneratorName,
		)
		if err != nil {
			return nil, err
		}

		if reportDataJSON != nil {
			_ = json.Unmarshal(reportDataJSON, &report.ReportData)
		}

		reports = append(reports, report)
	}

	return reports, rows.Err()
}

// GetReportsByType retrieves reports for a project filtered by type
func (r *ExecutionPlanRepository) GetReportsByType(ctx context.Context, projectID int, reportType string) ([]models.ExecutionReportWithDetails, error) {
	query := `SELECT
		er.id, er.project_id, er.report_type, er.generated_by, er.generator_type, er.report_data, er.created_at,
		p.name as project_name,
		CASE
			WHEN er.generator_type = 'user' THEN (SELECT username FROM users WHERE id = er.generated_by)
			ELSE (SELECT name FROM agents WHERE id = er.generated_by)
		END as generator_name
	FROM execution_reports er
	LEFT JOIN projects p ON er.project_id = p.id
	WHERE er.project_id = $1 AND er.report_type = $2
	ORDER BY er.created_at DESC`

	rows, err := r.pool.Query(ctx, query, projectID, reportType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.ExecutionReportWithDetails
	for rows.Next() {
		var report models.ExecutionReportWithDetails
		var reportDataJSON []byte
		err := rows.Scan(
			&report.ID, &report.ProjectID, &report.ReportType, &report.GeneratedBy,
			&report.GeneratorType, &reportDataJSON, &report.CreatedAt,
			&report.ProjectName, &report.GeneratorName,
		)
		if err != nil {
			return nil, err
		}

		if reportDataJSON != nil {
			_ = json.Unmarshal(reportDataJSON, &report.ReportData)
		}

		reports = append(reports, report)
	}

	return reports, rows.Err()
}
