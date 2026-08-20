package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"tic-tac-toe/pkg/service"
)

type UnknownPathsHandler struct {
	services *service.Service
}

func NewUnknownPathsHandler(services *service.Service) *UnknownPathsHandler {
	return &UnknownPathsHandler{services}
}

func (h *UnknownPathsHandler) InitRoutes(router *http.ServeMux) {
	router.HandleFunc("/", h.writeNotFoud)
}

func (h *UnknownPathsHandler) writeNotFoud(w http.ResponseWriter, r *http.Request) {
	log.Printf("Unknown path requested: %s", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	response := map[string]string{
		"error":   "Not Found",
		"message": "The requested path does not exist",
	}
	json.NewEncoder(w).Encode(response)
}
