package models

import (
	"encoding/json"

	"golang.org/x/net/websocket"
)

type Position struct {
	X int
	Y int
}

type PlayerData struct {
	ID       string
	Username string
	HP       string
	Position Position
}

type Player struct {
	ID   string
	Room *Room
	Conn *websocket.Conn
	Data PlayerData
}

type IncomingMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	Player  *Player         `json:"-"`
}
