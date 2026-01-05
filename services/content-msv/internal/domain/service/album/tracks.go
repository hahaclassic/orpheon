package album

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/pkg/errwrap"
	"github.com/hahaclassic/orpheon/services/content-msv/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/services/content-msv/internal/domain/usecase/album"
)

type AlbumTrackRepository interface {
	GetAllTracks(ctx context.Context, albumID uuid.UUID) ([]*entity.TrackMeta, error)
}

type AlbumTrackService struct {
	albumTrackRepository AlbumTrackRepository
}

func NewAlbumTrackService(albumTrackRepository AlbumTrackRepository) *AlbumTrackService {
	return &AlbumTrackService{albumTrackRepository: albumTrackRepository}
}

func (s *AlbumTrackService) GetAllTracks(ctx context.Context, albumID uuid.UUID) (tracks []*entity.TrackMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetAllTracks, err)
	}()

	return s.albumTrackRepository.GetAllTracks(ctx, albumID)
}
