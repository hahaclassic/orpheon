package deleteservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

// TODO: нет транзакции. Нет восстанавливающих операций. Надо продумать,
// как сделать так, чтобы сохранялась целостность данных.

type metaDeletionService interface {
	Delete(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) error
}

type trackDeletionService interface {
	Delete(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) error
}

type favoritesDeletionService interface {
	Delete(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) error
}

type PlaylistDeleter struct {
	meta      metaDeletionService
	tracks    trackDeletionService
	favorites favoritesDeletionService
}

func New(meta metaDeletionService, track trackDeletionService,
	favorites favoritesDeletionService) *PlaylistDeleter {
	return &PlaylistDeleter{
		meta:      meta,
		tracks:    track,
		favorites: favorites,
	}
}

func (p *PlaylistDeleter) Delete(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) error {
	err := p.favorites.Delete(ctx, claims, playlistID)
	if err != nil {
		return err
	}

	err = p.tracks.Delete(ctx, claims, playlistID)
	if err != nil {
		return err
	}

	err = p.meta.Delete(ctx, claims, playlistID)
	if err != nil {

		return err
	}

	return nil
}
