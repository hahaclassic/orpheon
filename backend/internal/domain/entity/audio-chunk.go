package entity

import "github.com/google/uuid"

type AudioChunk struct {
	Data    []byte
	TrackID uuid.UUID
	Start   int64
	End     int64
}
