package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
	defer LogStartEnd("register")()

	var body SingUpRequest

	writeInternalServerError := func(err error) {
		message := "internal server error"
		logMessage := fmt.Sprintf("[register][ERROR] register failed due to system error: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, message, logMessage)
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		writeInternalServerError(err)
		return
	}

	log.Println("[register] json decode success")

	writeBadRequestWithMessage := func(message string) {
		WriteErrorResponse(w, http.StatusBadRequest, message, "[register] "+message)
	}

	if body.Email == "" {
		writeBadRequestWithMessage("Email is required")
		return
	}

	_, err = mail.ParseAddress(body.Email)
	if err != nil {
		writeBadRequestWithMessage("Email is invalid")
		return
	}

	if body.Email == "" {
		writeBadRequestWithMessage("Email is required")
		return
	}

	if body.Password == "" {
		writeBadRequestWithMessage("Password is required")
		return
	}

	if body.PasswordConfirm == "" {
		writeBadRequestWithMessage("PasswordConfirm is required")
		return
	}

	if body.Password != body.PasswordConfirm {
		writeBadRequestWithMessage("Password and PasswordConfirm are not equal")
		return
	}

	log.Println("[register] payload is correct")

	existUser, err := h.services.FindUserByIdentity(body.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			writeBadRequestWithMessage(err.Error())
			return
		}
	}

	if existUser.Id != 0 {
		writeBadRequestWithMessage("user already exist")
		return
	}

	log.Println("[register] no conflict")

	user, err := h.services.Auth.CreateUser(common.UserToCreate{
		Email:    body.Email,
		Username: body.Username,
		Password: body.Password,
	})

	if err != nil {
		writeBadRequestWithMessage(err.Error())
		return
	}

	log.Println("[register] user created")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	response := NewResponseSuccess(user)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		writeInternalServerError(err)
	}
}
