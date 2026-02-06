package repository

import (
	"context"
	"fmt"

	"github.com/berkkaradalan/stackflow/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	pool *pgxpool.Pool
}

func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{
		pool: pool,
	}
}

func (r *ProjectRepository) Create(ctx context.Context, project *models.Project) error {
	query := `INSERT INTO projects (name, description, status, created_by)
	          VALUES ($1, $2, $3, $4)
	          RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		project.Name, project.Description, project.Status, project.CreatedBy,
	).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)
}

func (r *ProjectRepository) GetByID(ctx context.Context, id int) (*models.Project, error) {
	query := `SELECT id, name, description, status, created_by, created_at, updated_at
	          FROM projects WHERE id = $1`

	var project models.Project
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&project.ID, &project.Name, &project.Description, &project.Status,
		&project.CreatedBy, &project.CreatedAt, &project.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) GetAll(ctx context.Context) ([]models.Project, error) {
	query := `SELECT id, name, description, status, created_by, created_at, updated_at
	          FROM projects
	          ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		err := rows.Scan(
			&project.ID, &project.Name, &project.Description, &project.Status,
			&project.CreatedBy, &project.CreatedAt, &project.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ProjectRepository) Update(ctx context.Context, project *models.Project) error {
	query := `UPDATE projects
	          SET name = $1, description = $2, status = $3, updated_at = NOW()
	          WHERE id = $4
	          RETURNING updated_at`

	return r.pool.QueryRow(ctx, query,
		project.Name, project.Description, project.Status, project.ID,
	).Scan(&project.UpdatedAt)
}

func (r *ProjectRepository) UpdatePartial(ctx context.Context, id int, updates map[string]interface{}) (*models.Project, error) {
	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}

	query := "UPDATE projects SET "
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

func (r *ProjectRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM projects WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *ProjectRepository) GetStats(ctx context.Context, projectID int) (*models.ProjectStats, error) {
	stats := &models.ProjectStats{}

	// Get task counts
	taskQuery := `SELECT
		COUNT(*) as total,
		COUNT(*) FILTER (WHERE status IN ('done', 'closed')) as completed,
		COUNT(*) FILTER (WHERE status IN ('open', 'in_progress')) as pending
	FROM tasks WHERE project_id = $1`

	err := r.pool.QueryRow(ctx, taskQuery, projectID).Scan(
		&stats.TotalTasks,
		&stats.CompletedTasks,
		&stats.PendingTasks,
	)
	if err != nil {
		// If tasks table doesn't exist yet, use zeros
		stats.TotalTasks = 0
		stats.CompletedTasks = 0
		stats.PendingTasks = 0
	}

	// Get agent count
	agentQuery := `SELECT COUNT(*) FROM agents WHERE project_id = $1`
	err = r.pool.QueryRow(ctx, agentQuery, projectID).Scan(&stats.TotalAgents)
	if err != nil {
		stats.TotalAgents = 0
	}

	// Workflows not implemented yet
	stats.TotalWorkflows = 0

	return stats, nil
}
