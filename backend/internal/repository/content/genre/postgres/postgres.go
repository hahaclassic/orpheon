package genre_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type GenreRepository struct {
	pool *pgxpool.Pool
}

func NewGenreRepository(pool *pgxpool.Pool) *GenreRepository {
	return &GenreRepository{pool: pool}
}

func (r *GenreRepository) Create(ctx context.Context, genre *entity.Genre) error {
	query := `INSERT INTO genres (id, title) VALUES ($1, $2)`
	_, err := r.pool.Exec(ctx, query, genre.ID, genre.Title)
	return err
}

func (r *GenreRepository) Get(ctx context.Context, genreID uuid.UUID) (*entity.Genre, error) {
	query := `SELECT id, title FROM genres WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, genreID)

	var g entity.Genre
	err := row.Scan(&g.ID, &g.Title)
	if err != nil {
		return nil, fmt.Errorf("genre not found: %w", err)
	}
	return &g, nil
}

func (r *GenreRepository) Update(ctx context.Context, genre *entity.Genre) error {
	query := `UPDATE genres SET description = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, genre.Title, genre.ID)
	return err
}

func (r *GenreRepository) Delete(ctx context.Context, genreID uuid.UUID) error {
	query := `DELETE FROM genres WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, genreID)
	return err
}
