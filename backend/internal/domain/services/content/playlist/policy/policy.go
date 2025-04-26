package policy

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type PlaylistAccessRepository interface {
	GetAccessMeta(ctx context.Context, playlistID uuid.UUID) (*entity.PlaylistAccessMeta, error)
}

type PlaylistPolicyService struct {
	accessRepo PlaylistAccessRepository
}

func New(accessRepo PlaylistAccessRepository) *PlaylistPolicyService {
	return &PlaylistPolicyService{
		accessRepo: accessRepo,
	}
}

func (p *PlaylistPolicyService) CanView(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrCanView, err)
	}()

	meta, err := p.accessRepo.GetAccessMeta(ctx, playlistID)
	if err != nil {
		return err
	}

	if !meta.IsPrivate || claims.UserID == meta.OwnerID {
		return nil // ok
	}

	return usecase.ErrForbidden
}

func (p *PlaylistPolicyService) CanEdit(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrCanEdit, err)
	}()

	meta, err := p.accessRepo.GetAccessMeta(ctx, playlistID)
	if err != nil {
		return err
	}

	if claims.UserID == meta.OwnerID {
		return nil // ok
	}

	return usecase.ErrForbidden
}

func (p *PlaylistPolicyService) CanDelete(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrCanEdit, err)
	}()

	meta, err := p.accessRepo.GetAccessMeta(ctx, playlistID)
	if err != nil {
		return err
	}

	if claims.UserID == meta.OwnerID ||
		(claims.AccessLvl == entity.Admin && !meta.IsPrivate) {
		return nil // ok
	}

	return usecase.ErrForbidden
}
