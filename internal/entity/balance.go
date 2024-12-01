package entity

type Balance struct {
	CurrentBalance   float64 `json:"current"`
	WithdrawnBalance float64 `json:"withdrawn"`
}
