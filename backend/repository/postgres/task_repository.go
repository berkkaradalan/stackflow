package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/berkkaradalan/stackflow/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		pool: pool,
	}
}

// Create creates a new task
func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	tagsJSON, err := json.Marshal(task.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `INSERT INTO tasks (project_id, title, description, status, priority, assigned_agent_id, reviewer_id, created_by, creator_type, tags)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	          RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		task.ProjectID, task.Title, task.Description, task.Status, task.Priority,
		task.AssignedAgentID, task.ReviewerID, task.CreatedBy, task.CreatorType, tagsJSON,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
}

// GetByID retrieves a task by ID
func (r *TaskRepository) GetByID(ctx context.Context, id int) (*models.Task, error) {
	query := `SELECT id, project_id, title, description, status, priority, assigned_agent_id, reviewer_id, created_by, creator_type, tags, created_at, updated_at
	          FROM tasks WHERE id = $1`

	var task models.Task
	var tagsJSON []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&task.ID, &task.ProjectID, &task.Title, &task.Description, &task.Status, &task.Priority,
		&task.AssignedAgentID, &task.ReviewerID, &task.CreatedBy, &task.CreatorType, &tagsJSON,
		&task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(tagsJSON, &task.Tags); err != nil {
		task.Tags = []string{}
	}

	return &task, nil
}

// GetByIDWithDetails retrieves a task by ID with related entity names
func (r *TaskRepository) GetByIDWithDetails(ctx context.Context, id int) (*models.TaskWithDetails, error) {
	query := `SELECT
		t.id, t.project_id, t.title, t.description, t.status, t.priority,
		t.assigned_agent_id, t.reviewer_id, t.created_by, t.creator_type, t.tags,
		t.created_at, t.updated_at,
		a.name as agent_name,
		u.username as reviewer_name,
		p.name as project_name,
		CASE
			WHEN t.creator_type = 'user' THEN (SELECT username FROM users WHERE id = t.created_by)
			ELSE (SELECT name FROM agents WHERE id = t.created_by)
		END as creator_name
	FROM tasks t
	LEFT JOIN agents a ON t.assigned_agent_id = a.id
	LEFT JOIN users u ON t.reviewer_id = u.id
	LEFT JOIN projects p ON t.project_id = p.id
	WHERE t.id = $1`

	var task models.TaskWithDetails
	var tagsJSON []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&task.ID, &task.ProjectID, &task.Title, &task.Description, &task.Status, &task.Priority,
		&task.AssignedAgentID, &task.ReviewerID, &task.CreatedBy, &task.CreatorType, &tagsJSON,
		&task.CreatedAt, &task.UpdatedAt,
		&task.AssignedAgentName, &task.ReviewerName, &task.ProjectName, &task.CreatorName,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(tagsJSON, &task.Tags); err != nil {
		task.Tags = []string{}
	}

	return &task, nil
}

// GetAll retrieves all tasks with optional filters
func (r *TaskRepository) GetAll(ctx context.Context, filters *models.TaskFilters) ([]models.TaskWithDetails, error) {
	query := `SELECT
		t.id, t.project_id, t.title, t.description, t.status, t.priority,
		t.assigned_agent_id, t.reviewer_id, t.created_by, t.creator_type, t.tags,
		t.created_at, t.updated_at,
		a.name as agent_name,
		u.username as reviewer_name,
		p.name as project_name,
		CASE
			WHEN t.creator_type = 'user' THEN (SELECT username FROM users WHERE id = t.created_by)
			ELSE (SELECT name FROM agents WHERE id = t.created_by)
		END as creator_name
	FROM tasks t
	LEFT JOIN agents a ON t.assigned_agent_id = a.id
	LEFT JOIN users u ON t.reviewer_id = u.id
	LEFT JOIN projects p ON t.project_id = p.id
	WHERE 1=1`

	args := []interface{}{}
	argPos := 1

	if filters != nil {
		if filters.ProjectID != nil {
			query += fmt.Sprintf(" AND t.project_id = $%d", argPos)
			args = append(args, *filters.ProjectID)
			argPos++
		}
		if filters.Status != nil {
			query += fmt.Sprintf(" AND t.status = $%d", argPos)
			args = append(args, *filters.Status)
			argPos++
		}
		if filters.Priority != nil {
			query += fmt.Sprintf(" AND t.priority = $%d", argPos)
			args = append(args, *filters.Priority)
			argPos++
		}
		if filters.AssignedAgentID != nil {
			query += fmt.Sprintf(" AND t.assigned_agent_id = $%d", argPos)
			args = append(args, *filters.AssignedAgentID)
			argPos++
		}
		if filters.ReviewerID != nil {
			query += fmt.Sprintf(" AND t.reviewer_id = $%d", argPos)
			args = append(args, *filters.ReviewerID)
			argPos++
		}
	}

	query += " ORDER BY t.created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.TaskWithDetails
	for rows.Next() {
		var task models.TaskWithDetails
		var tagsJSON []byte
		err := rows.Scan(
			&task.ID, &task.ProjectID, &task.Title, &task.Description, &task.Status, &task.Priority,
			&task.AssignedAgentID, &task.ReviewerID, &task.CreatedBy, &task.CreatorType, &tagsJSON,
			&task.CreatedAt, &task.UpdatedAt,
			&task.AssignedAgentName, &task.ReviewerName, &task.ProjectName, &task.CreatorName,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(tagsJSON, &task.Tags); err != nil {
			task.Tags = []string{}
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetByProjectID retrieves all tasks for a project
func (r *TaskRepository) GetByProjectID(ctx context.Context, projectID int) ([]models.TaskWithDetails, error) {
	query := `SELECT
		t.id, t.project_id, t.title, t.description, t.status, t.priority,
		t.assigned_agent_id, t.reviewer_id, t.created_by, t.creator_type, t.tags,
		t.created_at, t.updated_at,
		a.name as agent_name,
		u.username as reviewer_name,
		p.name as project_name,
		CASE
			WHEN t.creator_type = 'user' THEN (SELECT username FROM users WHERE id = t.created_by)
			ELSE (SELECT name FROM agents WHERE id = t.created_by)
		END as creator_name
	FROM tasks t
	LEFT JOIN agents a ON t.assigned_agent_id = a.id
	LEFT JOIN users u ON t.reviewer_id = u.id
	LEFT JOIN projects p ON t.project_id = p.id
	WHERE t.project_id = $1
	ORDER BY t.created_at DESC`

	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.TaskWithDetails
	for rows.Next() {
		var task models.TaskWithDetails
		var tagsJSON []byte
		err := rows.Scan(
			&task.ID, &task.ProjectID, &task.Title, &task.Description, &task.Status, &task.Priority,
			&task.AssignedAgentID, &task.ReviewerID, &task.CreatedBy, &task.CreatorType, &tagsJSON,
			&task.CreatedAt, &task.UpdatedAt,
			&task.AssignedAgentName, &task.ReviewerName, &task.ProjectName, &task.CreatorName,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(tagsJSON, &task.Tags); err != nil {
			task.Tags = []string{}
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Update updates a task
func (r *TaskRepository) Update(ctx context.Context, task *models.Task) error {
	tagsJSON, err := json.Marshal(task.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `UPDATE tasks
	          SET title = $1, description = $2, status = $3, priority = $4,
	              assigned_agent_id = $5, reviewer_id = $6, tags = $7, updated_at = NOW()
	          WHERE id = $8
	          RETURNING updated_at`

	return r.pool.QueryRow(ctx, query,
		task.Title, task.Description, task.Status, task.Priority,
		task.AssignedAgentID, task.ReviewerID, tagsJSON, task.ID,
	).Scan(&task.UpdatedAt)
}

// UpdatePartial partially updates a task
func (r *TaskRepository) UpdatePartial(ctx context.Context, id int, updates map[string]interface{}) (*models.Task, error) {
	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}

	query := "UPDATE tasks SET "
	args := make([]interface{}, 0, len(updates)+1)
	argPos := 1

	first := true
	for key, value := range updates {
		if !first {
			query += ", "
		}
		if key == "tags" {
			tagsJSON, err := json.Marshal(value)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal tags: %w", err)
			}
			value = tagsJSON
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

// UpdateStatus updates task status
func (r *TaskRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE tasks SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, status, id)
	return err
}

// AssignAgent assigns an agent to a task
func (r *TaskRepository) AssignAgent(ctx context.Context, taskID int, agentID int) error {
	query := `UPDATE tasks SET assigned_agent_id = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, agentID, taskID)
	return err
}

// SetReviewer sets the reviewer for a task
func (r *TaskRepository) SetReviewer(ctx context.Context, taskID int, reviewerID int) error {
	query := `UPDATE tasks SET reviewer_id = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, reviewerID, taskID)
	return err
}

// Delete deletes a task
func (r *TaskRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// CreateActivity creates a new task activity
func (r *TaskRepository) CreateActivity(ctx context.Context, activity *models.TaskActivity) error {
	query := `INSERT INTO task_activities (task_id, actor_id, actor_type, action, old_value, new_value, message)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)
	          RETURNING id, created_at`

	return r.pool.QueryRow(ctx, query,
		activity.TaskID, activity.ActorID, activity.ActorType, activity.Action,
		activity.OldValue, activity.NewValue, activity.Message,
	).Scan(&activity.ID, &activity.CreatedAt)
}

// GetActivitiesByTaskID retrieves all activities for a task
func (r *TaskRepository) GetActivitiesByTaskID(ctx context.Context, taskID int) ([]models.TaskActivityWithDetails, error) {
	query := `SELECT
		ta.id, ta.task_id, ta.actor_id, ta.actor_type, ta.action, ta.old_value, ta.new_value, ta.message, ta.created_at,
		CASE
			WHEN ta.actor_type = 'user' THEN (SELECT username FROM users WHERE id = ta.actor_id)
			ELSE (SELECT name FROM agents WHERE id = ta.actor_id)
		END as actor_name
	FROM task_activities ta
	WHERE ta.task_id = $1
	ORDER BY ta.created_at DESC`

	rows, err := r.pool.Query(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.TaskActivityWithDetails
	for rows.Next() {
		var activity models.TaskActivityWithDetails
		err := rows.Scan(
			&activity.ID, &activity.TaskID, &activity.ActorID, &activity.ActorType,
			&activity.Action, &activity.OldValue, &activity.NewValue, &activity.Message,
			&activity.CreatedAt, &activity.ActorName,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return activities, nil
}

// GetTaskCountByProjectID returns the count of tasks for a project
func (r *TaskRepository) GetTaskCountByProjectID(ctx context.Context, projectID int) (int, error) {
	query := `SELECT COUNT(*) FROM tasks WHERE project_id = $1`
	var count int
	err := r.pool.QueryRow(ctx, query, projectID).Scan(&count)
	return count, err
}

// GetTaskCountByStatus returns task counts grouped by status for a project
func (r *TaskRepository) GetTaskCountByStatus(ctx context.Context, projectID int) (map[string]int, error) {
	query := `SELECT status, COUNT(*) FROM tasks WHERE project_id = $1 GROUP BY status`

	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		counts[status] = count
	}

	return counts, rows.Err()
}
