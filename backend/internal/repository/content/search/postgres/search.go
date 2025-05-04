package search_postgres

// import (
// 	"context"
// 	"fmt"

// 	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// type SearchRepositoryPgx struct {
// 	db *pgxpool.Pool
// }

// // Новый конструктор для SearchRepositoryPgx
// func NewSearchRepositoryPgx(db *pgxpool.Pool) *SearchRepositoryPgx {
// 	return &SearchRepositoryPgx{
// 		db: db,
// 	}
// }

// // SearchTracks реализует поиск треков
// func (r *SearchRepositoryPgx) SearchTracks(ctx context.Context, req *entity.SearchRequest) ([]*entity.TrackMeta, error) {
// 	query := `
// 		SELECT t.id, t.name, t.explicit, t.duration, t.stream_count
// 		FROM tracks t
// 		WHERE LOWER(t.name) LIKE LOWER($1)
// 	`
// 	args := []interface{}{fmt.Sprintf("%%%s%%", req.Query)}

// 	// Добавляем фильтры по жанру и стране, если они заданы
// 	if req.Filters.Genre != "" {
// 		query += " AND t.genre = $2"
// 		args = append(args, req.Filters.Genre)
// 	}
// 	if req.Filters.Country != "" {
// 		query += " AND t.country = $3"
// 		args = append(args, req.Filters.Country)
// 	}

// 	rows, err := r.db.Query(ctx, query, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to search tracks: %w", err)
// 	}
// 	defer rows.Close()

// 	var tracks []*entity.TrackMeta
// 	for rows.Next() {
// 		var track entity.TrackMeta
// 		if err := rows.Scan(&track.ID, &track.Name, &track.Explicit, &track.Duration, &track.StreamCount); err != nil {
// 			return nil, fmt.Errorf("failed to scan track: %w", err)
// 		}
// 		tracks = append(tracks, &track)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("failed to read rows: %w", err)
// 	}

// 	return tracks, nil
// }

// // SearchAlbums реализует поиск альбомов
// func (r *SearchRepositoryPgx) SearchAlbums(ctx context.Context, req *entity.SearchRequest) ([]*entity.AlbumMeta, error) {
// 	query := `
// 		SELECT a.id, a.name, a.artist, a.release_date, a.genre
// 		FROM albums a
// 		WHERE LOWER(a.name) LIKE LOWER($1)
// 	`
// 	args := []interface{}{fmt.Sprintf("%%%s%%", req.Query)}

// 	// Добавляем фильтры по жанру и стране, если они заданы
// 	if req.Filters.Genre != "" {
// 		query += " AND a.genre = $2"
// 		args = append(args, req.Filters.Genre)
// 	}
// 	if req.Filters.Country != "" {
// 		query += " AND a.country = $3"
// 		args = append(args, req.Filters.Country)
// 	}

// 	rows, err := r.db.Query(ctx, query, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to search albums: %w", err)
// 	}
// 	defer rows.Close()

// 	var albums []*entity.AlbumMeta
// 	for rows.Next() {
// 		var album entity.AlbumMeta
// 		if err := rows.Scan(&album.ID, &album.T, &album.Artist, &album.ReleaseDate, &album.Genre); err != nil {
// 			return nil, fmt.Errorf("failed to scan album: %w", err)
// 		}
// 		albums = append(albums, &album)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("failed to read rows: %w", err)
// 	}

// 	return albums, nil
// }

// // SearchArtists реализует поиск исполнителей
// func (r *SearchRepositoryPgx) SearchArtists(ctx context.Context, req *entity.SearchRequest) ([]*entity.ArtistMeta, error) {
// 	query := `
// 		SELECT a.id, a.name, a.genre, a.country
// 		FROM artists a
// 		WHERE LOWER(a.name) LIKE LOWER($1)
// 	`
// 	args := []interface{}{fmt.Sprintf("%%%s%%", req.Query)}

// 	// Добавляем фильтры по жанру и стране, если они заданы
// 	if req.Genre != "" {
// 		query += " AND a.genre = $2"
// 		args = append(args, req.Genre)
// 	}
// 	if req.Country != "" {
// 		query += " AND a.country = $3"
// 		args = append(args, req.Country)
// 	}

// 	rows, err := r.db.Query(ctx, query, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to search artists: %w", err)
// 	}
// 	defer rows.Close()

// 	var artists []*entity.ArtistMeta
// 	for rows.Next() {
// 		var artist entity.ArtistMeta
// 		if err := rows.Scan(&artist.ID, &artist.Name, &artist.Genre, &artist.Country); err != nil {
// 			return nil, fmt.Errorf("failed to scan artist: %w", err)
// 		}
// 		artists = append(artists, &artist)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("failed to read rows: %w", err)
// 	}

// 	return artists, nil
// }

// // SearchPlaylists реализует поиск плейлистов
// func (r *SearchRepositoryPgx) SearchPlaylists(ctx context.Context, req *entity.SearchRequest) ([]*entity.PlaylistMeta, error) {
// 	query := `
// 		SELECT p.id, p.name, p.owner_id, p.is_private, p.created_at
// 		FROM playlists p
// 		WHERE LOWER(p.name) LIKE LOWER($1)
// 	`
// 	args := []interface{}{fmt.Sprintf("%%%s%%", req.Query)}

// 	// Добавляем фильтры по жанру и стране, если они заданы
// 	if req.Genre != "" {
// 		query += " AND p.genre = $2"
// 		args = append(args, req.Genre)
// 	}
// 	if req.Country != "" {
// 		query += " AND p.country = $3"
// 		args = append(args, req.Country)
// 	}

// 	rows, err := r.db.Query(ctx, query, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to search playlists: %w", err)
// 	}
// 	defer rows.Close()

// 	var playlists []*entity.PlaylistMeta
// 	for rows.Next() {
// 		var playlist entity.PlaylistMeta
// 		if err := rows.Scan(&playlist.ID, &playlist.Name, &playlist.OwnerID, &playlist.IsPrivate, &playlist.CreatedAt); err != nil {
// 			return nil, fmt.Errorf("failed to scan playlist: %w", err)
// 		}
// 		playlists = append(playlists, &playlist)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("failed to read rows: %w", err)
// 	}

// 	return playlists, nil
// }
