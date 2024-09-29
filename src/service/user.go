package service

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/shaileshhb/quiz/src/db"
	"github.com/shaileshhb/quiz/src/db/models"
	"github.com/shaileshhb/quiz/src/security"
)

// UserService will consist of service methods that would be implemented by userService
type UserService interface {
	Register(*models.User) (*models.LoginResponse, error)
	Login(*models.Login) (*models.LoginResponse, error)
}

// userService will contain reference to db.
type userService struct {
	db *db.Database
}

// NewUserService will create new instance of userService
func NewUserService(db *db.Database) UserService {
	return &userService{
		db: db,
	}
}

// Register will register new user in the system.
func (service *userService) Register(user *models.User) (*models.LoginResponse, error) {
	user.Name = strings.TrimSpace(user.Name)
	user.Username = strings.TrimSpace(user.Username)
	user.Password = strings.TrimSpace(user.Password)

	err := service.checkDuplicateExist(user.Username)
	if err != nil {
		return nil, err
	}

	password, err := security.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = string(password)
	user.ID = uuid.New()

	service.db.Users = append(service.db.Users, *user)

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
func (service *userService) Login(login *models.Login) (*models.LoginResponse, error) {
	login.Username = strings.TrimSpace(login.Username)
	login.Password = strings.TrimSpace(login.Password)

	user, err := service.getUserByUsername(login.Username)
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
func (service *userService) checkDuplicateExist(username string) error {
	for _, user := range service.db.Users {
		if strings.EqualFold(strings.ToLower(user.Username), strings.ToLower(username)) {
			return errors.New("same username already exists")
		}
	}

	return nil
}

// getUserByUsername will fetch user by username. If not found, it will return error
func (service *userService) getUserByUsername(username string) (*models.User, error) {
	for _, user := range service.db.Users {
		if strings.EqualFold(strings.ToLower(user.Username), strings.ToLower(username)) {
			return &user, nil
		}
	}

	return nil, errors.New("username or password is incorrect")
}
