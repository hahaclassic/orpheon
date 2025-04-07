package entities

import "github.com/google/uuid"

type TrackMeta struct {
	ID       uuid.UUID
	Name     string
	Explicit bool
	Duration int
}
