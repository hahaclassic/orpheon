package search_postgres

import (
	"context"
	"fmt"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SearchRepository struct {
	db *pgxpool.Pool
}

func NewSearchRepository(db *pgxpool.Pool) *SearchRepository {
	return &SearchRepository{db: db}
}

func (r *SearchRepository) SearchTracks(ctx context.Context, req *entity.SearchRequest) ([]*entity.TrackMeta, error) {
	query := `
		SELECT t.id, t.genre_id, t.name, t.duration, t.explicit, t.license_id, t.album_id, t.track_number, t.total_streams
		FROM tracks t
		JOIN albums a ON t.album_id = a.id
		JOIN artists ar ON a.artist_id = ar.id
		WHERE LOWER(t.name) LIKE LOWER($1)
	`
	args := []interface{}{fmt.Sprintf("%%%s%%", req.Query)}
	argIdx := 2

	if req.Filters.Genre != "" {
		query += fmt.Sprintf(" AND t.genre_id = $%d", argIdx)
		args = append(args, req.Filters.Genre)
		argIdx++
	}
	if req.Filters.Country != "" {
		query += fmt.Sprintf(" AND ar.country = $%d", argIdx)
		args = append(args, req.Filters.Country)
		argIdx++
	}

	query += " ORDER BY t.name LIMIT $" + fmt.Sprint(argIdx) + " OFFSET $" + fmt.Sprint(argIdx+1)
	args = append(args, req.Limit, req.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search tracks: %w", err)
	}
	defer rows.Close()

	var tracks []*entity.TrackMeta
	for rows.Next() {
		var track entity.TrackMeta
		err := rows.Scan(&track.ID, &track.GenreID, &track.Name, &track.Duration, &track.Explicit, &track.LicenseID, &track.AlbumID, &track.TrackNumber, &track.TotalStreams)
		if err != nil {
			return nil, fmt.Errorf("failed to scan track: %w", err)
		}
		tracks = append(tracks, &track)
	}

	return tracks, rows.Err()
}

func (r *SearchRepository) SearchAlbums(ctx context.Context, req *entity.SearchRequest) ([]*entity.AlbumMeta, error) {
	query := `
		SELECT id, title, label, license_id, release_date
		FROM albums
		WHERE LOWER(title) LIKE LOWER($1)
	`
	args := []interface{}{fmt.Sprintf("%%%s%%", req.Query)}
	argIdx := 2

	// Жанра у albums напрямую нет, если нужно — делать JOIN с tracks

	query += " ORDER BY release_date DESC LIMIT $" + fmt.Sprint(argIdx) + " OFFSET $" + fmt.Sprint(argIdx+1)
	args = append(args, req.Limit, req.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search albums: %w", err)
	}
	defer rows.Close()

	var albums []*entity.AlbumMeta
	for rows.Next() {
		var album entity.AlbumMeta
		err := rows.Scan(&album.ID, &album.Title, &album.Label, &album.LicenseID, &album.ReleaseDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan album: %w", err)
		}
		albums = append(albums, &album)
	}

	return albums, rows.Err()
}

func (r *SearchRepository) SearchArtists(ctx context.Context, req *entity.SearchRequest) ([]*entity.ArtistMeta, error) {
	query := `
		SELECT id, name, country, description
		FROM artists
		WHERE LOWER(name) LIKE LOWER($1)
	`
	args := []interface{}{fmt.Sprintf("%%%s%%", req.Query)}
	argIdx := 2

	if req.Filters.Country != "" {
		query += fmt.Sprintf(" AND country = $%d", argIdx)
		args = append(args, req.Filters.Country)
		argIdx++
	}

	query += " ORDER BY name LIMIT $" + fmt.Sprint(argIdx) + " OFFSET $" + fmt.Sprint(argIdx+1)
	args = append(args, req.Limit, req.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search artists: %w", err)
	}
	defer rows.Close()

	var artists []*entity.ArtistMeta
	for rows.Next() {
		var artist entity.ArtistMeta
		err := rows.Scan(&artist.ID, &artist.Name, &artist.Country, &artist.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan artist: %w", err)
		}
		artists = append(artists, &artist)
	}

	return artists, rows.Err()
}

func (r *SearchRepository) SearchPlaylists(ctx context.Context, req *entity.SearchRequest) ([]*entity.PlaylistMeta, error) {
	query := `
		SELECT id, name, description, is_private, owner_id, created_at, updated_at
		FROM playlists
		WHERE LOWER(name) LIKE LOWER($1)
	`
	args := []interface{}{fmt.Sprintf("%%%s%%", req.Query)}
	argIdx := 2

	query += " ORDER BY created_at DESC LIMIT $" + fmt.Sprint(argIdx) + " OFFSET $" + fmt.Sprint(argIdx+1)
	args = append(args, req.Limit, req.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search playlists: %w", err)
	}
	defer rows.Close()

	var playlists []*entity.PlaylistMeta
	for rows.Next() {
		var playlist entity.PlaylistMeta
		err := rows.Scan(&playlist.ID, &playlist.Name, &playlist.Description, &playlist.IsPrivate, &playlist.OwnerID, &playlist.CreatedAt, &playlist.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan playlist: %w", err)
		}
		playlists = append(playlists, &playlist)
	}

	return playlists, rows.Err()
}
