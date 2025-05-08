package playlist_tracks_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaylistTracksRepository struct {
	pool *pgxpool.Pool
}

func NewPlaylistTracksRepository(pool *pgxpool.Pool) *PlaylistTracksRepository {
	return &PlaylistTracksRepository{pool: pool}
}

func (r *PlaylistTracksRepository) AddTrackToPlaylist(ctx context.Context, playlistID uuid.UUID, trackID uuid.UUID) error {
	const query = `
		INSERT INTO playlist_tracks (playlist_id, track_id, position)
		VALUES ($1, $2, COALESCE(
			(SELECT MAX(position) FROM playlist_tracks WHERE playlist_id = $1), -1
		) + 1)
		ON CONFLICT (playlist_id, track_id) DO NOTHING
	`

	_, err := r.pool.Exec(ctx, query, playlistID, trackID)
	if err != nil {
		return fmt.Errorf("add track to playlist: %w", err)
	}

	return nil
}

func (r *PlaylistTracksRepository) DeleteTrackFromPlaylist(ctx context.Context, playlistID uuid.UUID, trackID uuid.UUID) error {
	const query = `
		DELETE FROM playlist_tracks
		WHERE playlist_id = $1 AND track_id = $2
	`

	ct, err := r.pool.Exec(ctx, query, playlistID, trackID)
	if err != nil {
		return fmt.Errorf("delete track from playlist: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("track %s not found in playlist %s", trackID, playlistID)
	}

	return nil
}

func (r *PlaylistTracksRepository) DeleteAllTracksFromPlaylist(ctx context.Context, playlistID uuid.UUID) error {
	const query = `
		DELETE FROM playlist_tracks
		WHERE playlist_id = $1
	`

	ct, err := r.pool.Exec(ctx, query, playlistID)
	if err != nil {
		return fmt.Errorf("delete all tracks from playlist: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("no tracks found in playlist %s", playlistID)
	}

	return nil
}

func (r *PlaylistTracksRepository) GetAllPlaylistTracks(ctx context.Context, playlistID uuid.UUID) ([]*entity.TrackMeta, error) {
	const query = `
		SELECT 
			t.id, t.genre_id, t.name, t.duration, t.explicit,
			t.license_id, t.album_id, t.track_number, t.total_streams
		FROM playlist_tracks pt
		JOIN tracks t ON pt.track_id = t.id
		WHERE pt.playlist_id = $1
		ORDER BY pt.position ASC
	`

	rows, err := r.pool.Query(ctx, query, playlistID)
	if err != nil {
		return nil, fmt.Errorf("get all tracks from playlist: %w", err)
	}
	defer rows.Close()

	var tracks []*entity.TrackMeta
	for rows.Next() {
		var track entity.TrackMeta
		if err := rows.Scan(
			&track.ID,
			&track.GenreID,
			&track.Name,
			&track.Duration,
			&track.Explicit,
			&track.LicenseID,
			&track.AlbumID,
			&track.TrackNumber,
			&track.TotalStreams,
		); err != nil {
			return nil, fmt.Errorf("scan track: %w", err)
		}
		tracks = append(tracks, &track)
	}

	return tracks, nil
}
