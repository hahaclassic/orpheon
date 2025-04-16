package entity

import (
	"time"

	"github.com/google/uuid"
)

type Playlist struct {
	ID          uuid.UUID
	OwnerID     uuid.UUID
	Name        string
	Description string
	IsPrivate   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PlaylistAccessMeta struct {
	OwnerID   uuid.UUID
	IsPrivate bool
}
