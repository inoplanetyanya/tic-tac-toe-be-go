package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"tic-tac-toe/pkg/common"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user common.UserToCreate) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (email, username, password_hash, roles) VALUES ($1, $2, $3, $4) RETURNING id", usersTable)

	defaultRoles := []string{"user"}

	row := r.db.QueryRow(query, user.Email, user.Username, user.Password, defaultRoles)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUserByUsernameAndPassword(username, password_hash string) (common.User, error) {
	var user common.User

	query := fmt.Sprintf("SELECT id, username, email, roles FROM %s WHERE username=$1 AND password_hash=$2", usersTable)

	row := r.db.QueryRow(query, username, password_hash)

	if err := row.Scan(&user.Id, &user.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("row scan error: %w", err)
	}

	return user, nil
}

func (r *AuthPostgres) GetUserByEmailAndPassword(email, password_hash string) (common.User, error) {
	var user common.User

	query := fmt.Sprintf("SELECT id, username, email, roles  FROM %s WHERE email=$1 AND password_hash=$2", usersTable)

	row := r.db.QueryRow(query, email, password_hash)

	if err := row.Scan(&user.Id, &user.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("row scan error: %w", err)
	}

	return user, nil
}

func (r *AuthPostgres) FindUserByIdentity(identity string) (common.User, error) {
	var user common.User

	query := fmt.Sprintf("SELECT id, email, username, roles FROM %s WHERE username=$1 OR email=$1", usersTable)

	row := r.db.QueryRow(query, identity)

	if err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Roles); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("row scan error: %w", err)
	}

	return user, nil
}
