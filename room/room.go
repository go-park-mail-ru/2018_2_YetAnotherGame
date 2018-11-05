package room

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"
	"time"
)

type Game struct {
	Rooms    map[string]*Room
	MaxRooms int
}

type Room struct {
	MaxPlayers int
	Ticker *time.Ticker
	ID         string
	Players    map[string]*Player
	Register   chan *Player
	Unregister chan *Player
	Broadcast  chan *State
	Commands []*Command
	Message chan *IncomingMessage
}
type Command struct{
	Player *Player
	Command string
}


func (r *Room) Run() {
	r.Ticker = time.NewTicker(time.Second)
go r.RunBroadcast()
	for{
		 <-r.Ticker.C
		log.Printf("room tick %s players %v", r.ID, len(r.Players))
		 players:=[]PlayerData{}
for _,p:=range r.Players{
	players=append(players,p.Data)
}
state:=State{
	Players: players,
}
r.Broadcast<-&state
}}



func (r *Room) RunBroadcast(){
	for{
		s:=<-r.Broadcast
			for _, p:= range r.Players{
				p.Send(s)
			}
	}
}


type NewPlayer struct {
	Username string `json:"username"`
}

func (r *Room) ListenToPlayers(){
	for{
		m:=<-r.Message
		log.Printf("rmessage %s %v", m.Player.ID, string(m.Payload))
		switch m.Type {
		case "newPlayer":
			np:=&NewPlayer{}
			json.Unmarshal(m.Payload,np)
			m.Player.Data.Username=np.Username

		}
	}
}



func New() *Room {
	id2, _ := exec.Command("uuidgen").Output()

	stringID := string(id2[:])
	id:= strings.Trim(stringID, "\n")
	return &Room{
	ID: id,
	MaxPlayers:2,
		Players: make(map[string]*Player),
		Register: make(chan *Player),
		Unregister: make(chan *Player),
		Broadcast: make(chan *State),
		Message: make(chan *IncomingMessage),
	}
}



type State struct {
	Players []PlayerData
}
