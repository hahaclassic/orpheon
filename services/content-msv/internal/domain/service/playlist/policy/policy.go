package policy

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/pkg/commonerr"
	"github.com/hahaclassic/orpheon/pkg/errwrap"
	"github.com/hahaclassic/orpheon/services/content-msv/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/services/content-msv/internal/domain/usecase/playlist"
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

	if !meta.IsPrivate || (claims != nil && claims.UserID == meta.OwnerID) {
		return nil // ok
	}

	return commonerr.ErrForbidden
}

func (p *PlaylistPolicyService) CanEdit(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrCanEdit, err)
	}()

	meta, err := p.accessRepo.GetAccessMeta(ctx, playlistID)
	if err != nil {
		return err
	}

	if claims != nil && claims.UserID == meta.OwnerID {
		return nil // ok
	}

	return commonerr.ErrForbidden
}

func (p *PlaylistPolicyService) CanDelete(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrCanEdit, err)
	}()

	meta, err := p.accessRepo.GetAccessMeta(ctx, playlistID)
	if err != nil {
		return err
	}

	if (claims != nil && claims.UserID == meta.OwnerID) ||
		(claims != nil && claims.AccessLvl == entity.AdminLvl && !meta.IsPrivate) {
		return nil // ok
	}

	return commonerr.ErrForbidden
}
