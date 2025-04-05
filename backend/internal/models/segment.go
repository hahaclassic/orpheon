package models

import (
	"github.com/google/uuid"
)

// meta information about
type Segment struct {
	ID          uuid.UUID
	TrackID     uuid.UUID
	ListenCount int64
	Start       int
	Duration    int
}
