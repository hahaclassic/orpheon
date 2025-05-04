package license_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type LicenseRepository struct {
	pool *pgxpool.Pool
}

func NewLicenseRepository(pool *pgxpool.Pool) *LicenseRepository {
	return &LicenseRepository{pool: pool}
}

func (r *LicenseRepository) Create(ctx context.Context, license *entity.License) error {
	query := `INSERT INTO licenses (id, title, description) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, query, license.ID, license.Title, license.Description)
	return err
}

func (r *LicenseRepository) Get(ctx context.Context, licenseID uuid.UUID) (*entity.License, error) {
	query := `SELECT id, title, description FROM licenses WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, licenseID)

	var l entity.License
	err := row.Scan(&l.ID, &l.Title, &l.Description)
	if err != nil {
		return nil, fmt.Errorf("license not found: %w", err)
	}
	return &l, nil
}

func (r *LicenseRepository) Update(ctx context.Context, license *entity.License) error {
	query := `UPDATE licenses SET description = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, license.Description, license.ID)
	return err
}

func (r *LicenseRepository) Delete(ctx context.Context, licenseID uuid.UUID) error {
	query := `DELETE FROM licenses WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, licenseID)
	return err
}
