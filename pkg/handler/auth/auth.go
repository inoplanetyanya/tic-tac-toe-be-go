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
	router.HandleFunc(routes.Auth.SignUp.Path, h.signUp)
	router.HandleFunc(routes.Auth.SignIn.Path, h.SignIn)
}

type ResponseSuccess struct {
	UserID   int     `json:"user_id"`
	Username *string `json:"username"`
	Message  string  `json:"message"`
	Success  bool    `json:"success"`
	Token    string  `json:"access"`
}

func newResponseSuccess(user common.User, message string) ResponseSuccess {
	res := ResponseSuccess{
		Success:  true,
		Message:  message,
		UserID:   user.Id,
		Username: user.Username,
	}
	return res
}

func newResponseSuccessWithToken(user common.User, message, token string) ResponseSuccess {
	res := ResponseSuccess{
		Success:  true,
		Message:  message,
		UserID:   user.Id,
		Username: user.Username,
		Token:    token,
	}
	return res
}

type ResponseError struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func newResponseError(message string) ResponseError {
	res := ResponseError{}
	res.Success = false
	res.Message = message
	return res
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string, logMessage string) {
	log.Println(logMessage)
	w.WriteHeader(statusCode)
	response := newResponseError(message)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatal("[writeErrorResponse] Failed to encode response:", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func logStartEnd(handlerName string) func() {
	log.Println("[" + handlerName + "] start")
	return func() {
		log.Println("[" + handlerName + "] end\n")
	}
}
