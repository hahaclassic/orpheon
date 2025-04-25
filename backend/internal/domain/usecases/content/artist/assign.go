package artist

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrAssignArtistOnTrack = errors.New("failed to assign artist on track")
	ErrAssignArtistOnAlbum = errors.New("failed to assign artist on album")
)

type ArtistAssignService interface {
	AssignArtistToTrack(ctx context.Context, claims *entity.Claims, artistID uuid.UUID, trackID uuid.UUID) error
	AssignArtistToAlbum(ctx context.Context, claims *entity.Claims, artistID uuid.UUID, albumID uuid.UUID) error
}
