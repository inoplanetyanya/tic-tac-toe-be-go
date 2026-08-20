package repository

import (
	"database/sql"
	"tic-tac-toe/pkg/common"
)

type AuthRepository interface {
	CreateUser(user common.UserToCreate) (int, error)
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
