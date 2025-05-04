package entity

import (
	"time"

	"github.com/google/uuid"
)

type AlbumMeta struct {
	ID          uuid.UUID
	Title       string
	Label       string
	LicenseID   uuid.UUID
	ReleaseDate time.Time
}
