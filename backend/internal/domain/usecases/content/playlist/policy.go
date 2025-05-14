package playlist

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrCanView          = errors.New("couldn't verify viewing permissions")
	ErrCanEdit          = errors.New("couldn't verify edition permissions")
	ErrCanDelete        = errors.New("couldn't verify deletion permissions")
	ErrCanUpdatePrivacy = errors.New("couldn't update permissions")
)

type PlaylistPolicyService interface {
	CanDelete(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
	CanEdit(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
	CanView(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error

	UpdatePrivacy(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID, isPrivate bool) error
	DeletePolicy(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
}
