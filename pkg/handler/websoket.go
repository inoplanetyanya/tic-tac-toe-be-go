package handler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"tic-tac-toe/pkg/common"
	"tic-tac-toe/pkg/service"
	"time"

	"golang.org/x/net/websocket"
)

type WebSocketHandler struct {
	service *service.Service
}

func NewWebSocketHandler(service *service.Service) *WebSocketHandler {
	return &WebSocketHandler{service: service}
}

func (h *WebSocketHandler) InitRoutes(router *http.ServeMux) {
	router.HandleFunc("/ws", websocket.Handler(h.hws).ServeHTTP)
}

func getUsername(ptr *string) string {
	if ptr == nil {
		return "Anonymous"
	}

	return *ptr
}

func (h *WebSocketHandler) hws(ws *websocket.Conn) {
	log.Println("[ws] New client connection: ", ws.RemoteAddr())

	user, err := h.handleAuth(ws)
	if err != nil {
		log.Printf("[ws] Auth failed for %v: %v", ws.RemoteAddr(), err)
		_ = ws.Close()
		return
	}

	username := getUsername(user.Username)
	log.Printf("[ws] User '%s' (ID: %d) successfully authenticated", username, user.Id)

	h.readLoop(ws, user)
}

func (h *WebSocketHandler) handleAuth(ws *websocket.Conn) (common.User, error) {
	_ = ws.SetReadDeadline(time.Now().Add(5 * time.Second))

	var msg string
	if err := websocket.Message.Receive(ws, &msg); err != nil {
		return common.User{}, fmt.Errorf("failed to receive auth message: %w", err)
	}

	_ = ws.SetReadDeadline(time.Time{})

	sm := strings.Split(msg, " ")
	if len(sm) < 2 || sm[0] != "/auth" {
		return common.User{}, errors.New("first message must be /auth <token>")
	}

	token := sm[1]
	user, err := h.service.Auth.GetUserFromToken(token)
	if err != nil {
		return common.User{}, fmt.Errorf("invalid token: %w", err)
	}

	return user, nil
}

func (h *WebSocketHandler) removePlayer(ws *websocket.Conn) {
	room := h.service.Game.GameRoomList().MapByConn[ws]
	if room == nil {
		log.Println("[ws] Player was not in any game room, skip room cleanup")
		return
	}
	room.RemovePlayer(ws)
}

func (h *WebSocketHandler) readLoop(ws *websocket.Conn, user common.User) {
	p := common.Player{
		User: user,
		Conn: ws,
	}
	g := h.service.Game

	defer func() {
		username := getUsername(user.Username)

		log.Printf("[ws] Connection closed for user %s", username)
		h.removePlayer(ws)
		_, _ = g.RemovePlayerFromQueue(p)
		_ = ws.Close()
	}()

	for {
		var msg string
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			if errors.Is(err, io.EOF) || strings.Contains(err.Error(), "connection reset by peer") {
				log.Println("[ws] Connection closed or reset by peer")
			} else {
				log.Printf("[ws] Read error: %v", err)
			}
			break
		}

		sm := strings.Split(msg, " ")
		command := sm[0]

		if command == "/chat" {
			if len(sm) < 2 {
				continue
			}
			gr := h.service.Game.GameRoomList().MapByConn[ws]
			if gr == nil {
				continue
			}
			gr.Chat(ws, strings.Join(sm[1:], " "))
			continue
		}

		username := getUsername(user.Username)
		log.Printf("[ws] Message from user %s: %s\n", username, msg)

		if command == "/connect" {
			if err := g.AddPlayerToQueue(p); err != nil {
				log.Println("[ws] AddPlayerToQueue error:", err)
			}
		}

		if command == "/disconnect" {
			rp, err := g.RemovePlayerFromQueue(p)
			if err != nil {
				log.Println("[ws] RemovePlayerFromQueue error:", err)
				continue
			}
			log.Printf("[ws] player removed from queue %v\n", rp)
		}
	}
}
