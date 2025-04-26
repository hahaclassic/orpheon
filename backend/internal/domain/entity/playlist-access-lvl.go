package entity

type PlaylistAccessLvl int

const (
	PlaylistNoAccessLvl PlaylistAccessLvl = iota
	PlaylistViewerLvl
	PlaylistOwnerLvl
)
