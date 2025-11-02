package entity

type AccessLevel int

const (
	UnauthorizedLvl AccessLevel = iota
	UserLvl
	AdminLvl
)

func (a AccessLevel) String() string {
	return []string{"Unauthorized Lvl", "User Lvl", "Admin Lvl"}[a]
}
