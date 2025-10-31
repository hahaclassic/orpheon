package creds_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/services/auth-msv/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CredentialsRepository struct {
	pool *pgxpool.Pool
}

func NewCredentialsRepository(pool *pgxpool.Pool) *CredentialsRepository {
	return &CredentialsRepository{pool: pool}
}

func (r *CredentialsRepository) SaveCredentials(ctx context.Context, userID uuid.UUID, credentials *entity.UserCredentials) error {
	const query = `
		INSERT INTO credentials (user_id, login, password)
		VALUES ($1, $2, $3)
		ON CONFLICT (login) DO NOTHING
	`

	_, err := r.pool.Exec(ctx, query, userID, credentials.Login, credentials.Password)
	if err != nil {
		return fmt.Errorf("save credentials: %w", err)
	}

	return nil
}

func (r *CredentialsRepository) GetPasswordByLogin(ctx context.Context, login string) (string, error) {
	const query = `
		SELECT password
		FROM credentials
		WHERE login = $1
	`

	var password string
	err := r.pool.QueryRow(ctx, query, login).Scan(&password)
	if err != nil {
		return "", fmt.Errorf("get password by login: %w", err)
	}

	return password, nil
}

func (r *CredentialsRepository) GetPasswordByID(ctx context.Context, userID uuid.UUID) (string, error) {
	const query = `
		SELECT password
		FROM credentials
		WHERE user_id = $1
	`

	var password string
	err := r.pool.QueryRow(ctx, query, userID).Scan(&password)
	if err != nil {
		return "", fmt.Errorf("get password by id: %w", err)
	}

	return password, nil
}

func (r *CredentialsRepository) GetClaimsByLogin(ctx context.Context, login string) (*entity.Claims, error) {
	const query = `
		SELECT u.id, u.access_level FROM users u 
		JOIN credentials c ON u.id = c.user_id
		WHERE c.login = $1
	`

	var claims entity.Claims
	err := r.pool.QueryRow(ctx, query, login).Scan(&claims.UserID, &claims.AccessLvl)
	if err != nil {
		return nil, fmt.Errorf("get claims by login: %w", err)
	}

	return &claims, nil
}

func (r *CredentialsRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	const query = `
		UPDATE credentials
		SET password = $1
		WHERE user_id = $2
	`

	ct, err := r.pool.Exec(ctx, query, newPassword, userID)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("user %s not found", userID)
	}

	return nil
}
