package entity

import "strconv"

type AccessLevel int

const (
	UnauthorizedLvl AccessLevel = iota
	UserLvl
	AdminLvl
)

func (a AccessLevel) String() string {
	return strconv.Itoa(int(a))
}

func (a AccessLevel) Name() string {
	return []string{"Unauthorized Lvl", "User Lvl", "Admin Lvl"}[a]
}

func (a AccessLevel) IsValid() bool {
	return a >= UnauthorizedLvl && a <= AdminLvl
}
