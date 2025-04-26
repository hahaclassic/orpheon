package entity

import (
	"time"

	"github.com/google/uuid"
)

type AlbumMeta struct {
	ID          uuid.UUID
	Title       string
	Label       string
	ReleaseDate time.Time
}
