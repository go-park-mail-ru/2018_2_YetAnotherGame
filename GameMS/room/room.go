package room

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-park-mail-ru/2018_2_YetAnotherGame/GameMS/Collision"
	"github.com/google/uuid"
)

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
type NewScore struct {
	User        string `json:"user"`
	Score       string `json:"score"`
	X           string `json:"x"`
	Y           string `json:"y"`
	CollisionX  string `json:"xblock"`
	CollisionX2 string `json:"x2block"`
	CollisionY  string `json:"yblock"`
}

type State struct {
	Players []PlayerData
	Message *Message
}

func (r *Room) Run() {
	r.Ticker = time.NewTicker(time.Millisecond * 100)
	go r.RunBroadcast()
	for {
		<-r.Ticker.C
		players := []PlayerData{}
		for _, p := range r.Players {
			players = append(players, p.Data)
		}
		state := State{
			Players: players,
		}
		r.Broadcast <- &state
	}
}

func (r *Room) RunBroadcast() {
	for {
		s := <-r.Broadcast
		for _, p := range r.Players {
			err := p.Send(s)
			if err != nil {
				return
			}
		}
	}
}

type Message struct {
	Author  string `json:"author"`
	Message string `json:"message"`
}

//type ColDetect struct {
//	Collision string `json:"collision"`
//}
func (r *Room) ListenToPlayers() {
	for {
		m := <-r.Message
		//log.Printf("rmessage %s %v", m.Player.ID, string(m.Payload))
		switch m.Type {
		case "newPlayer":
			np := &NewPlayer{}
			err := json.Unmarshal(m.Payload, np)
			if err != nil {
				log.Println(err)
			}
			m.Player.Data.Username = np.Username
		case "Info":
			//	log.Printf("rmessage %s %v", m.Player.ID, string(m.Payload))
			ns := &NewScore{}
			err := json.Unmarshal(m.Payload, ns)
			if err != nil {
				log.Println(err)
			}
			m.Player.Data.Score = ns.Score
			m.Player.Data.Position.X = ns.X
			m.Player.Data.Position.Y = ns.Y
			ok := Collision.Collision(ns.X, ns.Y, ns.CollisionX, ns.CollisionY, ns.CollisionX2)
			if ok {
				//name:=[]PlayerData{}
				col := Message{Message: "Collision", Author: ns.User}
				//us:=PlayerData{Username:ns.User}
				//name=append(name,us)
				fmt.Println("coll")
				state := State{Message: &col}
				r.Broadcast <- &state
			}
		case "Chat":
			//log.Printf("rmessage %s %v", m.Player.ID, string(m.Payload))
			msg := &Message{}
			err := json.Unmarshal(m.Payload, msg)
			if err != nil {
				log.Println(err)
			}
			state := State{
				Message: msg,
			}
			r.Broadcast <- &state

		}
	}
}

func New() *Room {
	id := uuid.New().String()
	return &Room{
		ID:         id,
		MaxPlayers: 2,
		Players:    make(map[string]*Player),
		Register:   make(chan *Player),
		Unregister: make(chan *Player),
		Broadcast:  make(chan *State),
		Message:    make(chan *IncomingMessage),
	}
}
