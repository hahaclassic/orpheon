package search_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/search"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/search"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

func TestSearchTracks(t *testing.T) {
	req := &entity.SearchRequest{Query: "rock"}
	tracks := []*entity.TrackMeta{{ID: uuid.New(), Name: "Rock Anthem"}}

	tests := []struct {
		name    string
		mockRes []*entity.TrackMeta
		mockErr error
		wantErr error
	}{
		{"success", tracks, nil, nil},
		{"repo error", nil, errors.New("repo error"), usecase.ErrSearchTracks},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewSearchRepository(t)
			repo.On("SearchTracks", mock.Anything, req).Return(tt.mockRes, tt.mockErr)
			svc := search.NewSearchService(repo)

			res, err := svc.SearchTracks(context.Background(), req)
			assert.Equal(t, tt.mockRes, res)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestSearchAlbums(t *testing.T) {
	req := &entity.SearchRequest{Query: "album"}
	albums := []*entity.AlbumMeta{{ID: uuid.New(), Title: "The Album"}}

	tests := []struct {
		name    string
		mockRes []*entity.AlbumMeta
		mockErr error
		wantErr error
	}{
		{"success", albums, nil, nil},
		{"repo error", nil, errors.New("repo error"), usecase.ErrSearchAlbums},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewSearchRepository(t)
			repo.On("SearchAlbums", mock.Anything, req).Return(tt.mockRes, tt.mockErr)
			svc := search.NewSearchService(repo)

			res, err := svc.SearchAlbums(context.Background(), req)
			assert.Equal(t, tt.mockRes, res)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestSearchArtists(t *testing.T) {
	req := &entity.SearchRequest{Query: "artist"}
	artists := []*entity.ArtistMeta{{ID: uuid.New(), Name: "Famous Artist"}}

	tests := []struct {
		name    string
		mockRes []*entity.ArtistMeta
		mockErr error
		wantErr error
	}{
		{"success", artists, nil, nil},
		{"repo error", nil, errors.New("repo error"), usecase.ErrSearchArtists},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewSearchRepository(t)
			repo.On("SearchArtists", mock.Anything, req).Return(tt.mockRes, tt.mockErr)
			svc := search.NewSearchService(repo)

			res, err := svc.SearchArtists(context.Background(), req)
			assert.Equal(t, tt.mockRes, res)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestSearchPlaylists(t *testing.T) {
	userID := uuid.New()
	claims := &entity.Claims{UserID: userID}
	req := &entity.SearchRequest{Query: "playlist"}
	playlist1 := &entity.PlaylistMeta{ID: uuid.New(), IsPrivate: false}
	playlist2 := &entity.PlaylistMeta{ID: uuid.New(), IsPrivate: true, OwnerID: userID}
	playlist3 := &entity.PlaylistMeta{ID: uuid.New(), IsPrivate: true, OwnerID: uuid.New()}

	tests := []struct {
		name    string
		repoRes []*entity.PlaylistMeta
		repoErr error
		want    []*entity.PlaylistMeta
		wantErr error
	}{
		{
			"success with filtering",
			[]*entity.PlaylistMeta{playlist1, playlist2, playlist3},
			nil,
			[]*entity.PlaylistMeta{playlist1, playlist2},
			nil,
		},
		{
			"repo error",
			nil,
			errors.New("repo error"),
			nil,
			usecase.ErrSearchPlaylists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewSearchRepository(t)
			repo.On("SearchPlaylists", mock.Anything, req).Return(tt.repoRes, tt.repoErr)
			svc := search.NewSearchService(repo)

			res, err := svc.SearchPlaylists(context.Background(), claims, req)
			assert.Equal(t, tt.want, res)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
