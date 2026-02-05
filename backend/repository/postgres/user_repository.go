package repository

import (
	"context"
	"fmt"

	"github.com/berkkaradalan/stackflow/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, avatar_url, role, is_active, created_at, updated_at 
	          FROM users WHERE email = $1`

	var user models.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.AvatarUrl, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash, avatar_url, role) 
	          VALUES ($1, $2, $3, $4, $5) 
	          RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		user.Username, user.Email, user.PasswordHash, user.AvatarUrl, user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, avatar_url, role, is_active, created_at, updated_at 
	          FROM users WHERE id = $1`
	
	var user models.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.AvatarUrl, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `UPDATE users
	          SET username = $1, email = $2, avatar_url = $3, password_hash = $4, updated_at = NOW()
	          WHERE id = $5
	          RETURNING updated_at`

	return r.pool.QueryRow(ctx, query,
		user.Username, user.Email, user.AvatarUrl, user.PasswordHash, user.ID,
	).Scan(&user.UpdatedAt)
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, username, email, password_hash, avatar_url, role, is_active, created_at, updated_at
	          FROM users
	          ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.AvatarUrl, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *UserRepository) UpdatePartial(ctx context.Context, id int, updates map[string]interface{}) (*models.User, error) {
	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}

	query := "UPDATE users SET "
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