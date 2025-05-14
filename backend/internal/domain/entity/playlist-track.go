package entity

import "github.com/google/uuid"

type PlaylistTrack struct {
	PlaylistID uuid.UUID
	TrackID    uuid.UUID
	Position   int
}
