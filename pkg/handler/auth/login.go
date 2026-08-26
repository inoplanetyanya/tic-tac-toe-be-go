package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"tic-tac-toe/pkg/service"
)

type SingInResponse struct {
	Id       int      `json:"id"`
	Email    string   `json:"email"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

type SignInRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

func validateRequestBody(req SignInRequest) error {
	if req.Password == "" {
		return errors.New("password is required")
	}

	if req.Identity == "" {
		return errors.New("either username or email must be provided")
	}

	return nil
}

func (h *HandlerAuth) Login(w http.ResponseWriter, r *http.Request) {
	defer LogStartEnd("login")()

	var body SignInRequest

	writeBadRequestWithMessage := func(message string) {
		WriteErrorResponse(w, http.StatusBadRequest, message, "[login] "+message)
	}

	writeUnauthorizedWithMessage := func(message string) {
		WriteErrorResponse(w, http.StatusUnauthorized, message, "[login] "+message)
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		writeBadRequestWithMessage(err.Error())
		return
	}

	log.Println("[login] json decode success")

	// meaningful?
	err = validateRequestBody(body)
	if err != nil {
		writeBadRequestWithMessage("Fields 'identity'(username or email) and 'password' are required")
		return
	}

	log.Println("[login] payload is correct")

	writeInternalServerError := func(err error) {
		message := "internal server error"
		logMessage := fmt.Sprintf("[login][ERROR] login failed due to system error: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, message, logMessage)
	}

	user, token, err := h.services.GenerateToken(body.Identity, body.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			writeUnauthorizedWithMessage(err.Error())
			return
		}

		writeInternalServerError(err)

		return
	}

	response := NewResponseSuccessWithToken(user, token)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		writeInternalServerError(err)
	}
}
