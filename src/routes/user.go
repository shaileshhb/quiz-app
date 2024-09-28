package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/shaileshhb/quiz/src/controllers"
	"github.com/shaileshhb/quiz/src/db/models"
)

// userRoute contains reference to user controller and logger
type userRoute struct {
	con controllers.UserController
	log zerolog.Logger
}

// NewUserRoute will create new instance of userRoute.
func NewUserRoute(con controllers.UserController, log zerolog.Logger) *userRoute {
	return &userRoute{
		con: con,
		log: log,
	}
}

// RegisterRoute registers all endpoints to router.
func (ur *userRoute) RegisterRoute(router fiber.Router) {
	router.Post("/register", ur.register)
	router.Post("/login", ur.login)
	ur.log.Info().Msg("User routes registered")
}

// register will add user.
func (ur *userRoute) register(c *fiber.Ctx) error {
	user := &models.User{}

	err := c.BodyParser(user)
	if err != nil {
		ur.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	loginResponse, err := ur.con.Register(user)
	if err != nil {
		ur.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(loginResponse)
}

// login will check user details and set the cookie
func (ur *userRoute) login(c *fiber.Ctx) error {
	login := &models.Login{}

	err := c.BodyParser(login)
	if err != nil {
		ur.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	loginResponse, err := ur.con.Login(login)
	if err != nil {
		ur.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(loginResponse)
}
