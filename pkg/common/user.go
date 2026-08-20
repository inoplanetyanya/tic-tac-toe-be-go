package common

import (
	"fmt"
	"slices"
)

type User struct {
	Id       int      `json:"-"`
	Email    string   `json:"email"`
	Username *string  `json:"username"`
	Roles    []string `json:"roles"`
}

type UserToCreate struct {
	Email    string  `json:"email"`
	Username *string `json:"username"`
	Password string  `json:"password"`
}

func (u User) String() string {
	var username string
	if u.Username != nil {
		username = *u.Username
	} else {
		username = "NULL"
	}

	return fmt.Sprintf("User{Id: %d, Email: %s, Username: %s, Roles: %v}", u.Id, u.Email, username, u.Roles)
}

func CompareRolesIgnoreOrder(u1, u2 User) bool {
	r1 := slices.Clone(u1.Roles)
	r2 := slices.Clone(u2.Roles)

	slices.Sort(r1)
	slices.Sort(r2)

	return slices.Equal(r1, r2)
}

func (u User) Equals(other User) bool {
	return u.Id == other.Id && u.Email == other.Email && u.Username == other.Username && CompareRolesIgnoreOrder(u, other)
}
