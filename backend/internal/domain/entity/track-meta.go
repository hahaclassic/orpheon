package entity

import "github.com/google/uuid"

type TrackMeta struct {
	ID           uuid.UUID
	GenreID      uuid.UUID
	Name         string
	Duration     int
	Explicit     bool
	LicenseID    uuid.UUID
	AlbumID      uuid.UUID
	TrackNumber  int // position in album
	TotalStreams int
}
