package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"tic-tac-toe/pkg/common"
)

type SingUpRequest struct {
	Email           string  `json:"email"`
	Username        *string `json:"username"`
	Password        string  `json:"password"`
	PasswordConfirm string  `json:"passwordConfirm"`
}

func (h *HandlerAuth) Register(w http.ResponseWriter, r *http.Request) {
	defer logStartEnd("register")()

	var body SingUpRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error(), "[register] "+err.Error())
		return
	}

	log.Println("[register] json decode success")

	writeResponseWithMessage := func(message string) {
		writeErrorResponse(w, http.StatusBadRequest, message, "[register] "+message)
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

	log.Println("[register] payload is correct")

	existUser, err := h.services.FindUserByIdentity(body.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			writeResponseWithMessage(err.Error())
			return
		}
	}

	if existUser.Id != 0 {
		writeResponseWithMessage("user already exist")
		return
	}

	log.Println("[register] no conflict")

	user, err := h.services.Auth.CreateUser(common.UserToCreate{
		Email:    body.Email,
		Username: body.Username,
		Password: body.Password,
	})

	if err != nil {
		writeResponseWithMessage(err.Error())
		return
	}

	log.Println("[register] user created")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		writeResponseWithMessage(err.Error())
		return
	}
}
