package entity

import (
	"time"

	"github.com/google/uuid"
)

type PlaylistMeta struct {
	ID          uuid.UUID
	OwnerID     uuid.UUID
	Name        string
	Description string
	IsPrivate   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Rating      int
}

type PlaylistAccessMeta struct {
	OwnerID   uuid.UUID
	IsPrivate bool
}
