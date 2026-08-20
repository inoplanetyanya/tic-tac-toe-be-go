package handler

import (
	"net/http"
	"tic-tac-toe/pkg/handler/auth"
	"tic-tac-toe/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() http.Handler {
	router := http.NewServeMux()

	authHandler := auth.NewHandler(h.services)
	authHandler.InitRoutes(router)

	wsHandler := NewWebSocketHandler(h.services)
	wsHandler.InitRoutes(router)

	unknownPathHandler := NewUnknownPathsHandler(h.services)
	unknownPathHandler.InitRoutes(router)

	return enableCORS(router)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
