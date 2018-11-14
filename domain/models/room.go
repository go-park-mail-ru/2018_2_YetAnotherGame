package models

import "time"

type Room struct {
	MaxPlayers int
	Ticker     *time.Ticker
	ID         string
	Players    map[string]*Player
	Register   chan *Player
	Unregister chan *Player
	Broadcast  chan *State
	Commands   []*Command
	Message    chan *IncomingMessage
}
type Command struct {
	Player  *Player
	Command string
}

type NewPlayer struct {
	Username string `json:"username"`
}

type State struct {
	Players []PlayerData
}
