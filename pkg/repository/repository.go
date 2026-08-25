package repository

import (
	"database/sql"
	"errors"
	"tic-tac-toe/pkg/common"
)

type AuthRepository interface {
	CreateUser(user common.UserToCreate) (common.User, error)
	GetUserByUsernameAndPassword(username, password string) (common.User, error)
	GetUserByEmailAndPassword(email, password string) (common.User, error)
	FindUserByIdentity(identity string) (common.User, error)
}

type Repository struct {
	AuthRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		AuthRepository: NewAuthPostgres(db),
	}
}

var ErrUserNotFound = errors.New("user not found")
