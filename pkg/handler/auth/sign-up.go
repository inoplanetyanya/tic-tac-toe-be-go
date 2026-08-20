package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"tic-tac-toe/pkg/common"
)

type SingUpResponse struct{}

type SingUpRequest struct {
	Email           string  `json:"email"`
	Username        *string `json:"username"`
	Password        string  `json:"password"`
	PasswordConfirm string  `json:"passwordConfirm"`
}

func (h *HandlerAuth) signUp(w http.ResponseWriter, r *http.Request) {
	defer logStartEnd("signup")()

	var body SingUpRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error(), "[signup] "+err.Error())
		return
	}

	log.Println("[signup] json decode success")

	writeResponseWithMessage := func(message string) {
		writeErrorResponse(w, http.StatusBadRequest, message, "[signup] "+message)
	}

	if body.Email == "" {
		writeResponseWithMessage("Email is required")
		return
	}

	_, err = mail.ParseAddress(body.Email)
	if err != nil {
		writeResponseWithMessage("Email is invalid")
		return
	}

	if body.Email == "" {
		writeResponseWithMessage("Email is required")
		return
	}

	if body.Password == "" {
		writeResponseWithMessage("Password is required")
		return
	}

	if body.PasswordConfirm == "" {
		writeResponseWithMessage("PasswordConfirm is required")
		return
	}

	if body.Password != body.PasswordConfirm {
		writeResponseWithMessage("Password and PasswordConfirm are not equal")
		return
	}

	log.Println("[signup] payload is correct")

	existUser, err := h.services.FindUserByIdentity(body.Email)
	if err != nil {
		writeResponseWithMessage(err.Error())
		return
	}

	if existUser.Id != 0 {
		writeResponseWithMessage("user already exist")
		return
	}

	log.Println("[signup] no conflict")

	userID, err := h.services.Auth.CreateUser(common.UserToCreate{
		Email:    body.Email,
		Username: body.Username,
		Password: body.Password,
	})

	if err != nil {
		writeResponseWithMessage(err.Error())
		return
	}

	log.Println("[signup] user created")

	response := newResponseSuccess(
		common.User{
			Id:       userID,
			Username: body.Username,
			Email:    body.Email,
			Roles:    []string{"User"},
		},
		"User registered successfully",
	)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		writeResponseWithMessage(err.Error())
		return
	}
}
