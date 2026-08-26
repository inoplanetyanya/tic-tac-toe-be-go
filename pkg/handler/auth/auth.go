package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"tic-tac-toe/pkg/common"
	"tic-tac-toe/pkg/handler/routes"
	"tic-tac-toe/pkg/service"
)

type HandlerAuth struct {
	services *service.Service
}

func NewHandler(services *service.Service) *HandlerAuth {
	return &HandlerAuth{services: services}
}

func (h *HandlerAuth) InitRoutes(router *http.ServeMux) {
	router.HandleFunc(routes.Auth.Register.Path, h.Register)
	router.HandleFunc(routes.Auth.Login.Path, h.Login)
}

type ResponseSuccess struct {
	User  common.User `json:"user"`
	Token string      `json:"token_access"`
}

func NewResponseSuccess(user common.User) ResponseSuccess {
	res := ResponseSuccess{
		User: user,
	}
	return res
}

func NewResponseSuccessWithToken(user common.User, token string) ResponseSuccess {
	res := ResponseSuccess{
		User:  user,
		Token: token,
	}
	return res
}

type ResponseError struct {
	Error string `json:"error"`
}

func NewResponseError(err string) ResponseError {
	res := ResponseError{
		Error: err,
	}
	return res
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string, logMessage string) {
	log.Println(logMessage)
	w.WriteHeader(statusCode)
	response := NewResponseError(message)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatal("[writeErrorResponse] Failed to encode response:", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func LogStartEnd(handlerName string) func() {
	log.Println("[" + handlerName + "] start")
	return func() {
		log.Println("[" + handlerName + "] end\n")
	}
}
