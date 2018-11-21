package room

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"
	"time"
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

	Score string `json:"score"`
}

type State struct {
	Players []PlayerData
}

func (r *Room) Run() {
	r.Ticker = time.NewTicker(time.Millisecond*100)
	go r.RunBroadcast()
	for {
		<-r.Ticker.C
		log.Printf("room tick %s players %v", r.ID, len(r.Players))
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
			p.Send(s)
		}
	}
}
type message struct {
	// the json tag means this will serialize as a lowercased field
	Message string `json:"message"`
}
func (r *Room) ListenToPlayers() {
	for {
		m := <-r.Message
		log.Printf("rmessage %s %v", m.Player.ID, string(m.Payload))
		switch m.Type {
		case "newPlayer":
			np := &NewPlayer{}
			json.Unmarshal(m.Payload, np)
			m.Player.Data.Username = np.Username
		case "Score":
			log.Printf("rmessage %s %v", m.Player.ID, string(m.Payload))
			ns := &NewScore{}
			json.Unmarshal(m.Payload, ns)
			m.Player.Data.Score = ns.Score
			//for k, v :=  range  r.Players {
				//if k!=m.Player.ID{
					//m2 := message{"Thanks for the message!"}
					//v.Conn.WriteJSON(m2)
				//}
			//}
		}
	}
}

func New() *Room {
	id2, _ := exec.Command("uuidgen").Output()

	stringID := string(id2[:])
	id := strings.Trim(stringID, "\n")
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
