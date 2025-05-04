package assign_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ArtistAssignRepository struct {
	pool *pgxpool.Pool
}

func NewArtistAssignRepository(pool *pgxpool.Pool) *ArtistAssignRepository {
	return &ArtistAssignRepository{pool: pool}
}

func (r *ArtistAssignRepository) AssignArtistToTrack(ctx context.Context, artistID uuid.UUID, trackID uuid.UUID) error {
	query := `
		INSERT INTO artist_tracks (artist_id, track_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING;
	`

	_, err := r.pool.Exec(ctx, query, artistID, trackID)
	if err != nil {
		return fmt.Errorf("assign artist to track: %w", err)
	}

	return nil
}

func (r *ArtistAssignRepository) AssignArtistToAlbum(ctx context.Context, artistID uuid.UUID, albumID uuid.UUID) error {
	query := `
		INSERT INTO artist_albums (artist_id, album_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING;
	`

	_, err := r.pool.Exec(ctx, query, artistID, albumID)
	if err != nil {
		return fmt.Errorf("assign artist to album: %w", err)
	}

	return nil
}
