package models

import (
	"golang.org/x/net/websocket"
)

type Game struct {
	Rooms    map[string]*Room
	MaxRooms int
	Register chan *websocket.Conn
}
