package entities

type PlaylistAccessLvl int

const (
	PlaylistUndefinedLvl PlaylistAccessLvl = iota
	PlaylistViewer
	PlaylistOwnerLvl
)
