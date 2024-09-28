package controllers

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/shaileshhb/quiz/src/db"
	"github.com/shaileshhb/quiz/src/db/models"
	"github.com/shaileshhb/quiz/src/security"
)

type UserController interface {
	Register(*models.User) (*models.LoginResponse, error)
	Login(*models.Login) (*models.LoginResponse, error)
}

type userController struct {
	db *db.Database
}

func NewUserController(db *db.Database) UserController {
	return &userController{
		db: db,
	}
}

// Register will register new user in the system.
func (controller *userController) Register(user *models.User) (*models.LoginResponse, error) {
	user.Name = strings.TrimSpace(user.Name)
	user.Username = strings.TrimSpace(user.Username)
	user.Password = strings.TrimSpace(user.Password)

	err := controller.checkDuplicateExist(user.Username)
	if err != nil {
		return nil, err
	}

	password, err := security.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = string(password)
	user.ID = uuid.New()

	controller.db.Users = append(controller.db.Users, *user)

	loginResponse := models.LoginResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
	}

	loginResponse.Token, err = security.GenerateJWT(user)
	if err != nil {
		return nil, err
	}

	return &loginResponse, nil
}

// Login user.
func (controller *userController) Login(login *models.Login) (*models.LoginResponse, error) {
	login.Username = strings.TrimSpace(login.Username)
	login.Password = strings.TrimSpace(login.Password)

	user, err := controller.getUserByUsername(login.Username)
	if err != nil {
		return nil, err
	}

	err = security.ComparePassword(user.Password, login.Password)
	if err != nil {
		return nil, errors.New("username or password is incorrect")
	}

	loginResponse := models.LoginResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
	}

	loginResponse.Token, err = security.GenerateJWT(user)
	if err != nil {
		return nil, err
	}

	return &loginResponse, nil
}

// checkDuplicateExist will check if same username already exists in database.
func (controller *userController) checkDuplicateExist(username string) error {
	for _, user := range controller.db.Users {
		if strings.EqualFold(strings.ToLower(user.Username), strings.ToLower(username)) {
			return errors.New("same username already exists")
		}
	}

	return nil
}

// getUserByUsername will fetch user by username. If not found, it will return error
func (controller *userController) getUserByUsername(username string) (*models.User, error) {
	for _, user := range controller.db.Users {
		if strings.EqualFold(strings.ToLower(user.Username), strings.ToLower(username)) {
			return &user, nil
		}
	}

	return nil, errors.New("username or password is incorrect")
}
