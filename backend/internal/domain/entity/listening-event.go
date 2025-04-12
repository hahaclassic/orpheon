package entity

import "github.com/google/uuid"

type ListeningEvent struct {
	TrackID uuid.UUID
	UserID  uuid.UUID
	Ranges  []*Range // [[2, 39], [55, 141]] - listened from 2 to 39 seconds, then 55 to 141
}
