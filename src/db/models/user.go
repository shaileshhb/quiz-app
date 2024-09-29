package models

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/shaileshhb/quiz/src/utils"
)

// User contains information related to user
type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

func (u *User) Validate() error {
	u.Name = strings.TrimSpace(u.Name)
	if len(u.Name) == 0 {
		return errors.New("name must be specified")
	}

	if len(u.Name) > 50 {
		return errors.New("name cannot be greater than 50 characters")
	}

	isValid, err := utils.ValidateString(u.Name, `^[a-zA-Z\s]+$`)
	if err != nil {
		return err
	}

	if !isValid {
		return errors.New("name contains invalid characters")
	}

	if len(strings.TrimSpace(u.Username)) == 0 {
		return errors.New("name must be specified")
	}

	if len(strings.TrimSpace(u.Username)) < 5 || len(strings.TrimSpace(u.Username)) > 20 {
		return errors.New("username must be between 5 and 20 characters")
	}

	isValid, err = utils.ValidateString(u.Username, `^[a-zA-Z0-9]+$`)
	if err != nil {
		return err
	}

	if !isValid {
		return errors.New("username contains invalid characters")
	}

	if len(strings.TrimSpace(u.Password)) == 0 {
		return errors.New("name must be specified")
	}

	if len(strings.TrimSpace(u.Password)) < 8 || len(strings.TrimSpace(u.Password)) > 20 {
		return errors.New("password must be between 8 and 20 characters")
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
