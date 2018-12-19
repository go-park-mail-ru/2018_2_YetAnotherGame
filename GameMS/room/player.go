package room

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Position struct {
	X string
	Y string
}

type PlayerData struct {
	ID       string
	Username string
	Score       string
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

func (p *Player) Listen() {
	for {
		m := &IncomingMessage{}

		err := p.Conn.ReadJSON(m)
		if websocket.IsUnexpectedCloseError(err) {
			log.Printf("player %s was disconnected", p.ID)
			p.Room.Unregister <- p
			//p.Conn.Close()
			return
		}
		m.Player = p
		p.Room.Message <- m
	}
}

func (p *Player) Send(s *State) error {
	err := p.Conn.WriteJSON(s)
	if err != nil {
		fmt.Println("cant send msg")
		return err
	}
	return nil
}
