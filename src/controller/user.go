package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/shaileshhb/quiz/src/db/models"
	serv "github.com/shaileshhb/quiz/src/service"
)

// userController contains reference to user service and logger
type userController struct {
	service serv.UserService
	log     zerolog.Logger
}

// NewUserController will create new instance of userRoute.
func NewUserController(service serv.UserService, log zerolog.Logger) *userController {
	return &userController{
		service: service,
		log:     log,
	}
}

// RegisterRoute registers all endpoints to router.
func (controller *userController) RegisterRoute(router fiber.Router) {
	router.Post("/register", controller.register)
	router.Post("/login", controller.login)
	controller.log.Info().Msg("User routes registered")
}

// register will add user.
func (controller *userController) register(c *fiber.Ctx) error {
	user := &models.User{}

	err := c.BodyParser(user)
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	loginResponse, err := controller.service.Register(user)
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(loginResponse)
}

// login will check user details and set the cookie
func (controller *userController) login(c *fiber.Ctx) error {
	login := &models.Login{}

	err := c.BodyParser(login)
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	loginResponse, err := controller.service.Login(login)
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(loginResponse)
}
