package dto

type AuthTokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
