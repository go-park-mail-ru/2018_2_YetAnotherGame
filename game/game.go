package game

import (
	//"game/room"

	"2018_2_YetAnotherGame/room"
	"fmt"
	"github.com/gorilla/websocket"
	"log"

	"os/exec"
	"strings"
)

func New() *Game {
	return &Game{
		Rooms:    make(map[string]*room.Room),
		MaxRooms: 2,
		Register: make(chan *websocket.Conn),
	}
}

type Game struct {
	Rooms    map[string]*room.Room
	MaxRooms int
	Register chan *websocket.Conn
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
	log.Println("room %s created", r.ID)
	return r
}

func (g *Game) ProcessConn(conn *websocket.Conn) {
	id2, _ := exec.Command("uuidgen").Output()

	stringID := string(id2[:])
	id := strings.Trim(stringID, "\n")

	p := &room.Player{
		Conn: conn,
		ID:   id,
	}
	r := g.FindRoom()
	if r==nil{
		return
	}

	r.Players [p.ID] = p
	p.Room=r
	log.Println("player %s joined room %s", p.ID, r.ID)
	go p.Listen()
	fmt.Println(r.Players)
	if len(r.Players) == r.MaxPlayers {
		go r.Run()
	}

}






//r.Players[p.ID] = p
//log.Println("player %s joined room %s", p.ID, r.ID)
//}
