package user_postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/services/user-msv/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, name, registration_date, access_level)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.pool.Exec(ctx, query,
		user.ID,
		user.Name,
		user.RegistrationDate,
		int(user.AccessLvl),
	)

	return err
}

func (r *UserRepository) GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	query := `
		SELECT id, name, registration_date, access_level
		FROM users
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, userID)

	var user entity.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.RegistrationDate,
		&user.AccessLvl,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET name = $1
		WHERE id = $2
	`
	cmdTag, err := r.pool.Exec(ctx, query,
		user.Name,
		user.ID,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	cmdTag, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}
