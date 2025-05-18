package entity

import (
	"time"

	"github.com/google/uuid"
)

type AlbumMeta struct {
	ID          uuid.UUID `json:"id"`
	Title       string
	Label       string
	LicenseID   uuid.UUID `json:"license_id"`
	ReleaseDate time.Time `json:"release_date"`
}
