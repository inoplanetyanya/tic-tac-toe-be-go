package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"tic-tac-toe/pkg/common"

	"github.com/lib/pq"
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

	defaultRoles := pq.Array([]string{"user"})

	queryRowArgs := []any{
		user.Email,
		user.Username,
		user.Password,
		defaultRoles,
	}

	row := r.db.QueryRow(query, queryRowArgs...)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUserByUsernameAndPassword(username, password_hash string) (common.User, error) {
	var user common.User

	query := fmt.Sprintf("SELECT id, username, email, roles FROM %s WHERE username=$1 AND password_hash=$2", usersTable)

	queryRowArgs := []any{
		username,
		password_hash,
	}

	row := r.db.QueryRow(query, queryRowArgs...)

	scanArgs := []any{
		&user.Id,
		&user.Email,
		&user.Username,
		pq.Array(&user.Roles),
	}

	if err := row.Scan(scanArgs...); err != nil {
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

	queryRowArgs := []any{
		email,
		password_hash,
	}

	row := r.db.QueryRow(query, queryRowArgs...)

	scanArgs := []any{
		&user.Id,
		&user.Email,
		&user.Username,
		pq.Array(&user.Roles),
	}

	if err := row.Scan(scanArgs...); err != nil {
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

	scanArgs := []any{
		&user.Id,
		&user.Email,
		&user.Username,
		pq.Array(&user.Roles),
	}

	if err := row.Scan(scanArgs...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("row scan error: %w", err)
	}

	return user, nil
}
