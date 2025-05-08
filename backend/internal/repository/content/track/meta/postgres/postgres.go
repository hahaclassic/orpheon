package track_meta_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type TrackMetaRepository struct {
	pool *pgxpool.Pool
}

func NewTrackMetaRepository(pool *pgxpool.Pool) *TrackMetaRepository {
	return &TrackMetaRepository{pool: pool}
}

func (r *TrackMetaRepository) GetByID(ctx context.Context, trackID uuid.UUID) (*entity.TrackMeta, error) {
	query := `
		SELECT id, genre_id, name, duration, explicit, license_id, album_id, track_number, total_streams
		FROM tracks
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, trackID)

	var track entity.TrackMeta
	err := row.Scan(
		&track.ID,
		&track.GenreID,
		&track.Name,
		&track.Duration,
		&track.Explicit,
		&track.LicenseID,
		&track.AlbumID,
		&track.TrackNumber,
		&track.TotalStreams,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get track: %w", err)
	}

	return &track, nil
}

func (r *TrackMetaRepository) Create(ctx context.Context, track *entity.TrackMeta) error {
	query := `
		INSERT INTO tracks (id, genre_id, name, duration, explicit, license_id, album_id, track_number, total_streams)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.pool.Exec(ctx, query,
		track.ID,
		track.GenreID,
		track.Name,
		track.Duration,
		track.Explicit,
		track.LicenseID,
		track.AlbumID,
		track.TrackNumber,
		track.TotalStreams,
	)

	if err != nil {
		return fmt.Errorf("failed to create track: %w", err)
	}

	return nil
}

func (r *TrackMetaRepository) Update(ctx context.Context, track *entity.TrackMeta) error {
	query := `
		UPDATE tracks
		SET genre_id = $2,
			name = $3,
			duration = $4,
			explicit = $5,
			license_id = $6,
			album_id = $7,
			track_number = $8,
			total_streams = $9
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		track.ID,
		track.GenreID,
		track.Name,
		track.Duration,
		track.Explicit,
		track.LicenseID,
		track.AlbumID,
		track.TrackNumber,
		track.TotalStreams,
	)

	if err != nil {
		return fmt.Errorf("failed to update track: %w", err)
	}

	return nil
}

func (r *TrackMetaRepository) Delete(ctx context.Context, trackID uuid.UUID) error {
	query := `
		DELETE FROM tracks
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, trackID)
	if err != nil {
		return fmt.Errorf("failed to delete track: %w", err)
	}

	return nil
}
