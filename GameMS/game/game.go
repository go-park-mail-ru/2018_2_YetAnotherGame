package game

import (
	"github.com/go-park-mail-ru/2018_2_YetAnotherGame/GameMS/room"
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

type Metrics struct {
	Counter *prometheus.SummaryVec
}

type Game struct {
	Rooms       map[string]*room.Room
	MaxRooms    int
	Register    chan *websocket.Conn
	AmountRooms int
}

func New() *Game {
	return &Game{
		Rooms:       make(map[string]*room.Room),
		MaxRooms:    2,
		Register:    make(chan *websocket.Conn),
		AmountRooms: 0,
	}
}

func (g *Game) Run() {
	fmt.Println("fff")
	for {
		conn := <-g.Register
		g.ProcessConn(conn)
	}
}

func (g *Game) FindRoom() *room.Room {
	for _, r := range g.Rooms {
		if len(r.Players) < r.MaxPlayers {
			return r
		}
	}
	if len(g.Rooms) >= g.MaxRooms {
		return nil
	}
	r := room.New()
	go r.ListenToPlayers()
	g.Rooms[r.ID] = r
	g.AmountRooms += 1
	log.Println("room %s created", r.ID)
	return r
}

func (g *Game) ProcessConn(conn *websocket.Conn) {
	id := uuid.New().String()

	p := &room.Player{
		Conn: conn,
		ID:   id,
	}
	p.Data.ID = p.ID
	r := g.FindRoom()
	if r == nil {
		return
	}

	r.Players[p.ID] = p
	p.Room = r
	log.Println("player %s joined room %s", p.ID, r.ID)
	go p.Listen()
	fmt.Println(r.Players)
	if len(r.Players) == r.MaxPlayers {
		go r.Run()
	}

}
