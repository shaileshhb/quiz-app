package models

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

// User contains information related to user
type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

func (u *User) Validate() error {
	if len(strings.TrimSpace(u.Name)) == 0 {
		return errors.New("name must be specified")
	}

	if len(strings.TrimSpace(u.Username)) == 0 {
		return errors.New("name must be specified")
	}

	if len(strings.TrimSpace(u.Password)) == 0 {
		return errors.New("name must be specified")
	}
	return nil
}

// Login contains information related to user login
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse contains information related to user login response
type LoginResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Token    string    `json:"token"`
}
