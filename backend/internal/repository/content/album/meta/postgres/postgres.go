package album_meta_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlbumRepository struct {
	pool *pgxpool.Pool
}

func NewAlbumRepository(pool *pgxpool.Pool) *AlbumRepository {
	return &AlbumRepository{pool: pool}
}

func (r *AlbumRepository) CreateAlbum(ctx context.Context, album *entity.AlbumMeta) error {
	query := `
		INSERT INTO albums (id, title, label, license_id, release_date)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, album.ID, album.Title, album.Label, album.LicenseID, album.ReleaseDate)
	if err != nil {
		return fmt.Errorf("create album: %w", err)
	}
	return nil
}

func (r *AlbumRepository) GetAlbum(ctx context.Context, id uuid.UUID) (*entity.AlbumMeta, error) {
	query := `
		SELECT id, title, label, license_id, release_date
		FROM albums
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)

	var album entity.AlbumMeta
	err := row.Scan(&album.ID, &album.Title, &album.Label, &album.LicenseID, &album.ReleaseDate)
	if err != nil {
		return nil, fmt.Errorf("get album: %w", err)
	}
	return &album, nil
}

func (r *AlbumRepository) UpdateAlbum(ctx context.Context, album *entity.AlbumMeta) error {
	query := `
		UPDATE albums
		SET title = $1, label = $2, release_date = $3
		WHERE id = $4
	`
	cmd, err := r.pool.Exec(ctx, query, album.Title, album.Label, album.ReleaseDate, album.ID)
	if err != nil {
		return fmt.Errorf("update album: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("update album: no rows affected")
	}
	return nil
}

func (r *AlbumRepository) DeleteAlbum(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM albums
		WHERE id = $1
	`
	cmd, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete album: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("delete album: no rows affected")
	}
	return nil
}
