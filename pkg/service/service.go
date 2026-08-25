package service

import (
	"errors"
	"tic-tac-toe/pkg/common"
	"tic-tac-toe/pkg/repository"
)

type Auth interface {
	CreateUser(user common.UserToCreate) (int, error)
	GetUserByUsernameAndPassword(username, password string) (common.User, error)
	GetUserByEmailAndPassword(email, password string) (common.User, error)
	GetUserFromToken(token string) (common.User, error)
	FindUserByIdentity(identity string) (common.User, error)
	GenerateToken(username, password string) (common.User, string, error)
	ParseToken(token string) (int, error)
}

type Game interface {
	AddPlayerToQueue(player common.Player) error
	RemovePlayerFromQueue(player common.Player) (common.Player, error)
	PlayerInQueue(player common.Player) bool
	CreateGameRoom() *common.GameRoom
	GameRoomList() *common.GameList
}

type Service struct {
	Auth
	Game
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Auth: NewAuthService(repos.AuthRepository),
		Game: NewGameService(),
	}
}

var ErrInvalidCredentials = errors.New("invalid username or password")
