package repository

import (
	"context"

	"github.com/berkkaradalan/stackflow/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InviteTokenRepository struct {
	pool *pgxpool.Pool
}

func NewInviteTokenRepository(pool *pgxpool.Pool) *InviteTokenRepository {
	return &InviteTokenRepository{
		pool: pool,
	}
}

func (r *InviteTokenRepository) Create(ctx context.Context, invite *models.InviteToken) error {
	query := `INSERT INTO invite_tokens (email, username, role, token, expires_at)
	          VALUES ($1, $2, $3, $4, $5)
	          RETURNING id, created_at`

	return r.pool.QueryRow(ctx, query,
		invite.Email, invite.Username, invite.Role, invite.Token, invite.ExpiresAt,
	).Scan(&invite.ID, &invite.CreatedAt)
}

func (r *InviteTokenRepository) GetByToken(ctx context.Context, token string) (*models.InviteToken, error) {
	query := `SELECT id, email, username, role, token, expires_at, used_at, created_at
	          FROM invite_tokens
	          WHERE token = $1`

	var invite models.InviteToken
	err := r.pool.QueryRow(ctx, query, token).Scan(
		&invite.ID, &invite.Email, &invite.Username, &invite.Role,
		&invite.Token, &invite.ExpiresAt, &invite.UsedAt, &invite.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &invite, nil
}

func (r *InviteTokenRepository) MarkAsUsed(ctx context.Context, token string) error {
	query := `UPDATE invite_tokens
	          SET used_at = NOW()
	          WHERE token = $1`

	_, err := r.pool.Exec(ctx, query, token)
	return err
}

func (r *InviteTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM invite_tokens
	          WHERE expires_at < NOW()`

	_, err := r.pool.Exec(ctx, query)
	return err
}
