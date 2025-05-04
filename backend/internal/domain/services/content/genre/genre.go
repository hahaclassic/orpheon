package genre

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type GenreRepository interface {
	Create(ctx context.Context, genre *entity.Genre) error
	Get(ctx context.Context, genreID uuid.UUID) (*entity.Genre, error)
	Update(ctx context.Context, genre *entity.Genre) error
	Delete(ctx context.Context, genreID uuid.UUID) error
}

var (
	ErrCreateGenre = errors.New("failed to create genre")
	ErrGetGenre    = errors.New("failed to get genre")
	ErrUpdateGenre = errors.New("failed to update genre")
	ErrDeleteGenre = errors.New("failed to delete genre")

	ErrForbidden      = errors.New("permission denied error")
	ErrInvalidGenreID = errors.New("invalid genre ID")
)

type GenreService struct {
	repo GenreRepository
}

func NewGenreService(repo GenreRepository) *GenreService {
	return &GenreService{repo: repo}
}

func (s *GenreService) CreateGenre(ctx context.Context, claims *entity.Claims, genre *entity.Genre) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(ErrCreateGenre, err)
	}()

	if claims.AccessLvl != entity.Admin {
		return ErrForbidden
	}

	if genre.ID == uuid.Nil {
		return ErrInvalidGenreID
	}

	return s.repo.Create(ctx, genre)
}

func (s *GenreService) GetGenre(ctx context.Context, genreID uuid.UUID) (_ *entity.Genre, err error) {
	defer func() {
		err = errwrap.WrapIfErr(ErrGetGenre, err)
	}()

	if genreID == uuid.Nil {
		return nil, ErrInvalidGenreID
	}

	return s.repo.Get(ctx, genreID)
}

func (s *GenreService) UpdateGenre(ctx context.Context, claims *entity.Claims, genre *entity.Genre) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(ErrUpdateGenre, err)
	}()

	if claims.AccessLvl != entity.Admin {
		return ErrForbidden
	}

	if genre.ID == uuid.Nil {
		return ErrInvalidGenreID
	}

	return s.repo.Update(ctx, genre)
}

func (s *GenreService) DeleteGenre(ctx context.Context, claims *entity.Claims, genreID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(ErrDeleteGenre, err)
	}()

	if claims.AccessLvl != entity.Admin {
		return ErrForbidden
	}

	if genreID == uuid.Nil {
		return ErrInvalidGenreID
	}

	return s.repo.Delete(ctx, genreID)
}
