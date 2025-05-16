package entity

import (
	"github.com/google/uuid"
)

// meta information about segment
// !!! UPDATED: NO SEGMENT ID ONLY IDX
type Segment struct {
	TrackID      uuid.UUID
	Idx          int
	TotalStreams uint64
	Range        *Range
}

type SegmentsIdxs struct {
	TrackID uuid.UUID
	Idxs    []int
}
