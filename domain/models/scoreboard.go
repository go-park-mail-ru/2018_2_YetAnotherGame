package models

type Scoreboard struct {
	Users   []User
	CanNext bool `json:"CanNext"`
}
