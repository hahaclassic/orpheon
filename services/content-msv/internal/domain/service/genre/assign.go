package genre

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/pkg/commonerr"
	"github.com/hahaclassic/orpheon/pkg/errwrap"
	"github.com/hahaclassic/orpheon/services/content-msv/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/services/content-msv/internal/domain/usecase/genre"
)

type GenreAssignRepository interface {
	AssignGenreToAlbum(ctx context.Context, genreID uuid.UUID, albumID uuid.UUID) error
	UnassignGenreFromAlbum(ctx context.Context, genreID uuid.UUID, albumID uuid.UUID) error
}

type GenreAssignService struct {
	repo GenreAssignRepository
}

func NewGenreAssignService(repo GenreAssignRepository) *GenreAssignService {
	return &GenreAssignService{repo: repo}
}

func (s *GenreAssignService) AssignGenreToAlbum(ctx context.Context, claims *entity.Claims, genreID uuid.UUID, albumID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrAssignGenreToAlbum, err)
	}()

	if claims == nil || claims.AccessLvl != entity.AdminLvl {
		return commonerr.ErrForbidden
	}

	return s.repo.AssignGenreToAlbum(ctx, genreID, albumID)
}

func (s *GenreAssignService) UnassignGenreFromAlbum(ctx context.Context, claims *entity.Claims, genreID uuid.UUID, albumID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUnassignGenreFromAlbum, err)
	}()

	if claims == nil || claims.AccessLvl != entity.AdminLvl {
		return commonerr.ErrForbidden
	}

	return s.repo.UnassignGenreFromAlbum(ctx, genreID, albumID)
}
