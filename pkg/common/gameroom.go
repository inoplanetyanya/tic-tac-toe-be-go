package common

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type GameRoom struct {
	Players   map[*websocket.Conn]User
	Mu        sync.Mutex
	RoomUUID  string
	BoardSize int
	Started   bool
	Ended     bool
}

func NewGameRoom() *GameRoom {
	return &GameRoom{
		Players:   make(map[*websocket.Conn]User),
		Mu:        sync.Mutex{},
		RoomUUID:  uuid.NewString(),
		BoardSize: 3,
		Started:   false,
		Ended:     false,
	}
}

func (g *GameRoom) AddPlayer(user User, conn *websocket.Conn) {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	g.Players[conn] = user

	for _, p := range g.Players {
		if p == g.Players[conn] {
			log.Printf("[gr] user %v already in room %s\n", p, g.RoomUUID)
		}
		return
	}

	log.Printf("[gr] user %v added to room %s\n", user, g.RoomUUID)
}

func (g *GameRoom) RemovePlayer(conn *websocket.Conn) {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	user := g.Players[conn]
	delete(g.Players, conn)

	if len(g.Players) > 0 {
		chatMsg := ChatMessage{
			From:    "server",
			Message: fmt.Sprintf("user %s removed from game room", user.Username),
			Type:    "chat",
		}

		jsonMsg, err := json.Marshal(chatMsg)
		if err != nil {
			log.Printf("Failed to marshal chat message: %v", err)
			return
		}

		for conn := range g.Players {
			if _, err := conn.Write(jsonMsg); err != nil {
				log.Printf("Failed to send message to player: %v", err)
			}
		}
	}

	log.Printf("[gr] user %v removed from room %s\n", user, g.RoomUUID)
	if len(g.Players) < 2 {
		go g.EndGame()
	}
}

func (g *GameRoom) StartGame() {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	g.Started = true

	chatMsg := ChatMessage{
		From:    "server",
		Message: string("game start"),
		Type:    "chat",
	}

	jsonMsg, err := json.Marshal(chatMsg)
	if err != nil {
		log.Printf("Failed to marshal chat message: %v", err)
		return
	}

	for conn := range g.Players {
		if _, err := conn.Write(jsonMsg); err != nil {
			log.Printf("Failed to send message to player: %v", err)
		}
	}
}

func (g *GameRoom) EndGame() {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	chatMsg := ChatMessage{
		From:    "server",
		Message: string("game end"),
		Type:    "chat",
	}

	jsonMsg, err := json.Marshal(chatMsg)
	if err != nil {
		log.Printf("Failed to marshal chat message: %v", err)
		return
	}

	for conn := range g.Players {
		if _, err := conn.Write(jsonMsg); err != nil {
			log.Printf("Failed to send message to player: %v", err)
		}
	}
}

type ChatMessage struct {
	From    string `json:"from"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (g *GameRoom) Chat(from *websocket.Conn, msg string) {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	fromUser := g.Players[from]

	chatMsg := ChatMessage{
		From:    fromUser.Username,
		Message: string(msg),
		Type:    "chat",
	}

	jsonMsg, err := json.Marshal(chatMsg)
	if err != nil {
		log.Printf("Failed to marshal chat message: %v", err)
		return
	}

	for conn := range g.Players {
		if _, err := conn.Write(jsonMsg); err != nil {
			log.Printf("Failed to send message to player: %v", err)
		}
	}
}
