package deleter

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type AlbumMetaDeletionService interface {
	DeleteMeta(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) error
}

type AlbumCoverDeletionService interface {
	GetCover(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) (*entity.Cover, error)
	DeleteCover(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) error
	SaveCover(ctx context.Context, claims *entity.Claims, cover *entity.Cover) error
}

type AlbumTrackDeletionService interface {
	GetAllTracks(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) ([]uuid.UUID, error)
	DeleteAllTracks(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) error
	RestoreAllTracks(ctx context.Context, claims *entity.Claims, albumID uuid.UUID, trackIDs []uuid.UUID) error
}

type rollback func() error

type AlbumDeleter struct {
	meta   AlbumMetaDeletionService
	cover  AlbumCoverDeletionService
	tracks AlbumTrackDeletionService
}

type OptionFunc func(*AlbumDeleter)

func WithMetaDeletion(metaService AlbumMetaDeletionService) OptionFunc {
	return func(pd *AlbumDeleter) {
		pd.meta = metaService
	}
}

func WithTracksDeletion(trackService AlbumTrackDeletionService) OptionFunc {
	return func(pd *AlbumDeleter) {
		pd.tracks = trackService
	}
}

func WithFavoritesDeletion(favoriteService FavoritesDeletionService) OptionFunc {
	return func(pd *AlbumDeleter) {
		pd.favorites = favoriteService
	}
}

func NewAlbumDeleter(meta AlbumMetaDeletionService, cover AlbumCoverDeletionService, tracks AlbumTrackDeletionService) *AlbumDeleter {
	return &AlbumDeleter{meta: meta, cover: cover, tracks: tracks}
}

func (a *AlbumDeleter) DeleteAlbum(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) (err error) {
	var rollbacks []rollback
	defer func() {
		if err != nil {
			for i := len(rollbacks) - 1; i >= 0; i-- {
				_ = rollbacks[i]()
			}
		}
	}()

	if a.cover != nil {
		cover, err := a.cover.GetCover(ctx, claims, albumID)
		if err != nil {
			return err
		}
		err = a.cover.DeleteCover(ctx, claims, albumID)
		if err != nil {
			return err
		}
		rollbacks = append(rollbacks, func() error {
			return a.cover.SaveCover(ctx, claims, cover)
		})
	}

	if a.tracks != nil {
		trackIDs, err := a.tracks.GetAllTracks(ctx, claims, albumID)
		if err != nil {
			return err
		}
		err = a.tracks.DeleteAllTracks(ctx, claims, albumID)
		if err != nil {
			return err
		}
		rollbacks = append(rollbacks, func() error {
			return a.tracks.RestoreAllTracks(ctx, claims, albumID, trackIDs)
		})
	}

	if a.meta != nil {
		err = a.meta.DeleteMeta(ctx, claims, albumID)
		if err != nil {
			return err
		}
	}

	return nil
}
