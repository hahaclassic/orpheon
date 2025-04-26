package entity

type AccessLevel int

const (
	Unauthorized AccessLevel = iota
	User
	Admin
)
