package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shaileshhb/quiz/src/db/models"
	"github.com/shaileshhb/quiz/src/security"
	serv "github.com/shaileshhb/quiz/src/service"
)

// userQuizController contains reference to user quiz serivce and logger
type userQuizController struct {
	service serv.UserQuizService
	log     zerolog.Logger
}

// NewUserQuizController will create new instance of userQuizRoute.
func NewUserQuizController(service serv.UserQuizService, log zerolog.Logger) *userQuizController {
	return &userQuizController{
		service: service,
		log:     log,
	}
}

// RegisterRoute registers all endpoints to router.
func (controller *userQuizController) RegisterRoute(router fiber.Router) {
	router.Post("/users/quizzes/:quizID/start", security.MandatoryAuthMiddleware, controller.startQuiz)
	router.Post("/users/quizzes/:quizID/attempts/:attemptID", security.MandatoryAuthMiddleware, controller.submitAnswer)
	router.Get("/users/quizzes/:quizID/results", security.MandatoryAuthMiddleware, controller.getUserQuizResults)
	controller.log.Info().Msg("User quiz routes registered")
}

// startQuiz will start a quiz for a user.
func (controller *userQuizController) startQuiz(c *fiber.Ctx) error {
	userQuiz := models.UserQuizAttempts{}

	quizID, err := uuid.Parse(c.Params("quizID"))
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userInterface := c.Locals("user")
	user := userInterface.(*models.User)
	userQuiz.UserID = user.ID
	userQuiz.QuizID = quizID

	err = userQuiz.Validate()
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = controller.service.StartQuiz(&userQuiz)
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(userQuiz)
}

// submitAnswer will submit user's answer for a given question.
func (controller *userQuizController) submitAnswer(c *fiber.Ctx) error {
	userResponse := models.UserResponse{}

	err := c.BodyParser(&userResponse)
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userResponse.QuizID, err = uuid.Parse(c.Params("quizID"))
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userResponse.UserQuizAttemptID, err = uuid.Parse(c.Params("attemptID"))
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userInterface := c.Locals("user")
	user := userInterface.(*models.User)
	userResponse.UserID = user.ID
	userResponse.IsCorrect = false

	err = userResponse.Validate()
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	correctOption, err := controller.service.SubmitAnswer(&userResponse)
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(map[string]interface{}{
		"isCorrect":     userResponse.IsCorrect,
		"correctOption": correctOption,
	})
}

// getUserQuizResults will return results for specific quiz for specified user.
func (controller *userQuizController) getUserQuizResults(c *fiber.Ctx) error {
	userInterface := c.Locals("user")
	user := userInterface.(*models.User)

	quizID, err := uuid.Parse(c.Params("quizID"))
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userQuiz, err := controller.service.GetUserQuizResults(user.ID, quizID)
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(userQuiz)
}
