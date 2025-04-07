package entities

import "github.com/google/uuid"

type CoverObjectType string

const (
	CoverAlbum    CoverObjectType = "album"
	CoverPlaylist CoverObjectType = "playlist"
)

type Cover struct {
	ObjectType CoverObjectType
	ObjectID   uuid.UUID
	Data       []byte
}
