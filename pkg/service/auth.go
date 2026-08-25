package service

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"tic-tac-toe/pkg/common"
	"tic-tac-toe/pkg/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getTokenTtlHours() time.Duration {
	ttlStr := os.Getenv("TOKEN_TTL_HOURS")
	if ttlStr == "" {
		log.Fatal("TOKEN_TTL_HOURS is not provided via evn")
	}

	ttlHours, err := strconv.Atoi(ttlStr)
	if err != nil {
		log.Fatalf("TOKEN_TTL_HOURS parsing error '%s': '%v'", ttlStr, err)
	}

	return time.Duration(ttlHours) * time.Hour
}

var (
	salt       = os.Getenv("TOKEN_SALT")
	signingKey = os.Getenv("TOKEN_SIGNING_KEY")
	tokenTTL   = getTokenTtlHours()
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId   int      `json:"user_id"`
	Username *string  `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}

type AuthService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user common.UserToCreate) (int, error) {
	if s.repo == nil {
		return 0, errors.New("repository is not initialized")
	}
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GetUserByUsernameAndPassword(username, password string) (common.User, error) {
	var user common.User

	if s.repo == nil {
		return user, errors.New("repository is not initialized")
	}

	password_hash := generatePasswordHash(password)
	user, err := s.repo.GetUserByUsernameAndPassword(username, password_hash)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *AuthService) GetUserByEmailAndPassword(email, password string) (common.User, error) {
	var user common.User

	if s.repo == nil {
		return user, errors.New("repository is not initialized")
	}

	password_hash := generatePasswordHash(password)
	user, err := s.repo.GetUserByEmailAndPassword(email, password_hash)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *AuthService) FindUserByIdentity(identity string) (common.User, error) {
	if s.repo == nil {
		return common.User{}, errors.New("repository is not initialized")
	}

	user, err := s.repo.FindUserByIdentity(identity)

	return user, err
}

func generatePasswordHash(password string) string {
	hash := sha256.New()

	hash.Write([]byte(password))
	hash.Write([]byte(salt))

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (s *AuthService) GenerateToken(identity, password string) (common.User, string, error) {
	user, err := s.repo.FindUserByIdentity(identity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("user '%s' does not exist", identity)
			return user, "", ErrInvalidCredentials
		}
		return user, "", err
	}

	hash := generatePasswordHash(password)
	user, err = s.repo.GetUserByEmailAndPassword(user.Email, hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[WARN] login attempt failed: incorrect password for user '%s'", identity)
			return user, "", ErrInvalidCredentials
		}
		return user, "", err
	}

	claims := jwt.MapClaims{
		"user_id":  user.Id,
		"username": user.Username,
		"exp":      time.Now().Add(tokenTTL).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return user, "", fmt.Errorf("failed to sign token: %w", err)
	}

	return user, signedToken, nil
}

// TODO rename(?)
func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		fmt.Println("parse token error: ", err)
		return -1, err
	}

	if !token.Valid {
		msg := "token is invalid"
		fmt.Println(msg)
		return -1, errors.New(msg)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		fmt.Println("invalid token claims")
		return -1, errors.New("invalid token claims")
	}

	return claims.UserId, nil
}

func (s *AuthService) GetUserFromToken(accessToken string) (common.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		fmt.Println("parse token error: ", err)
		return common.User{}, err
	}

	if !token.Valid {
		msg := "token is invalid"
		fmt.Println(msg)
		return common.User{}, errors.New(msg)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		fmt.Println("invalid token claims")
		return common.User{}, errors.New("invalid token claims")
	}

	user := common.User{
		Id:       claims.UserId,
		Username: claims.Username,
		Email:    claims.Email,
	}

	return user, nil
}
