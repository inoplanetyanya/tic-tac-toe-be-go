package handler

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"tic-tac-toe/pkg/common"
	"tic-tac-toe/pkg/service"

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

func (h *WebSocketHandler) hws(ws *websocket.Conn) {
	log.Println("[ws] New client connection: ", ws.RemoteAddr())
	h.readLoop(ws)

}

func (h *WebSocketHandler) removePlayer(ws *websocket.Conn) {
	h.service.Game.GameRoomList().MapByConn[ws].RemovePlayer(ws)
}

func (h *WebSocketHandler) readLoop(ws *websocket.Conn) {
	for {
		var msg string
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("[ws] Connection closed by client")
				h.removePlayer(ws)
			} else if strings.Contains(err.Error(), "connection reset by peer") {
				log.Println("[ws] Connection reset by peer")
				h.removePlayer(ws)
			} else {
				log.Printf("[ws] Read error: %v", err)
				h.removePlayer(ws)
			}
			break
		}

		sm := strings.Split(msg, " ")
		if len(sm) < 2 {
			log.Println("[ws] Invalid message format:", msg)
			continue
		}

		if sm[0] == "/chat" {
			gr := h.service.Game.GameRoomList().MapByConn[ws]
			if gr == nil {
				continue
			}

			gr.Chat(ws, strings.Join(sm[1:], " "))
			continue
		}

		userId, err := h.service.Auth.ParseToken(sm[1])
		if err != nil {
			log.Println("[ws] Token parse error:", err)
			continue
		}
		log.Printf("[ws] message from user with id %d: %s\n", userId, msg)

		if sm[0] == "/connect" {
			user, err := h.service.Auth.GetUserFromToken(sm[1])
			if err != nil {
				log.Println("[ws] GetUserByToken error:", err)
				continue
			}

			p := common.Player{
				User: user,
				Conn: ws,
			}

			g := h.service.Game

			if err := g.AddPlayerToQueue(p); err != nil {
				log.Println("[ws] AddPlayerToQueue error:", err)
			}
		}

		if sm[0] == "/disconnect" {
			user, err := h.service.Auth.GetUserFromToken(sm[1])
			if err != nil {
				log.Println("[ws] GetUserByToken error:", err)
				continue
			}

			p := common.Player{
				User: user,
				Conn: ws,
			}

			g := h.service.Game

			rp, err := g.RemovePlayerFromQueue(p)
			if err != nil {
				log.Println("[ws] RemovePlayerFromQueue error:", err)
				continue
			} else {
				log.Printf("[ws] player removed from queue %v\n", rp)
			}
		}
	}
}
