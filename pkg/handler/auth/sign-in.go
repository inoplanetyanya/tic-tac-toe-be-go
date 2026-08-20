package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"tic-tac-toe/pkg/common"
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

func (h *HandlerAuth) SignIn(w http.ResponseWriter, r *http.Request) {
	defer logStartEnd("signin")()

	var body SignInRequest

	writeResponseWithMessage := func(message string) {
		writeErrorResponse(w, http.StatusBadRequest, message, "[signin] "+message)
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		writeResponseWithMessage(err.Error())
		return
	}

	log.Println("[signin] json decode success")

	// meaningful?
	err = validateRequestBody(body)
	if err != nil {
		writeResponseWithMessage("Fields 'identity'(username or email) and 'password' are required")
		return
	}

	log.Println("[signin] payload is correct")

	var getUserMethod func(identity string, password string) (common.User, error)

	// TODO better way to check if identity is email or username?
	if strings.Contains(body.Identity, "@") {
		getUserMethod = h.services.GetUserByEmailAndPassword
	} else {
		getUserMethod = h.services.GetUserByUsernameAndPassword
	}

	user, err := getUserMethod(body.Identity, body.Password)

	if err != nil {
		writeResponseWithMessage(err.Error())
		return
	}

	token, err := h.services.GenerateToken(body.Identity, body.Password)
	if err != nil {
		writeResponseWithMessage(err.Error())
		return
	}

	response := newResponseSuccessWithToken(user, "successfully signed in", token)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		writeResponseWithMessage(err.Error())
		return
	}
}
