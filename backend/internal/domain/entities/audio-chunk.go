package entities

import "github.com/google/uuid"

type AudioChunk struct {
	Data    []byte
	TrackID uuid.UUID
	Start   uint64
	End     uint64
}
