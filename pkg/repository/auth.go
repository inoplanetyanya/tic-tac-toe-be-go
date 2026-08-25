package repository

import (
	"database/sql"
	"fmt"
	"tic-tac-toe/pkg/common"

	"github.com/lib/pq"
)

const userSelectFields = "id, email, username, roles"
const userInsertFields = "(email, username, password_hash, roles)"

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user common.UserToCreate) (int, error) {
	var id int
	query := fmt.Sprintf(
		"INSERT INTO %s %s VALUES ($1, $2, $3, $4) RETURNING id",
		usersTable,
		userInsertFields,
	)

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
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE username=$1 AND password_hash=$2",
		userSelectFields,
		usersTable,
	)

	queryRowArgs := []any{
		username,
		password_hash,
	}

	row := r.db.QueryRow(query, queryRowArgs...)

	return scanToUser(row)
}

func (r *AuthPostgres) GetUserByEmailAndPassword(email, password_hash string) (common.User, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE email=$1 AND password_hash=$2",
		userSelectFields,
		usersTable,
	)

	queryRowArgs := []any{
		email,
		password_hash,
	}

	row := r.db.QueryRow(query, queryRowArgs...)

	return scanToUser(row)
}

func (r *AuthPostgres) FindUserByIdentity(identity string) (common.User, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE username=$1 OR email=$1",
		userSelectFields,
		usersTable,
	)

	row := r.db.QueryRow(query, identity)

	return scanToUser(row)
}

func scanToUser(row *sql.Row) (common.User, error) {
	var user common.User

	scanArgs := []any{
		&user.Id,
		&user.Email,
		&user.Username,
		pq.Array(&user.Roles),
	}

	err := row.Scan(scanArgs...)

	return user, err
}
