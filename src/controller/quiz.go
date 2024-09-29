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

// quizController contains reference to quiz service and logger
type quizController struct {
	service serv.QuizService
	log     zerolog.Logger
}

// NewQuizController will create new instance of quizRoute.
func NewQuizController(service serv.QuizService, log zerolog.Logger) *quizController {
	return &quizController{
		service: service,
		log:     log,
	}
}

// RegisterRoute registers all endpoints to router.
func (controller *quizController) RegisterRoute(router fiber.Router) {
	router.Post("/quizzes", security.MandatoryAuthMiddleware, controller.CreateQuiz)
	router.Get("/quizzes/:quizID", security.MandatoryAuthMiddleware, controller.GetQuiz)

	controller.log.Info().Msg("Quiz routes registered")
}

// CreateQuiz will create new quiz.
func (controller *quizController) CreateQuiz(c *fiber.Ctx) error {
	quiz := models.Quiz{}

	err := c.BodyParser(&quiz)
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = quiz.Validate()
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = controller.service.Create(&quiz)
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(map[string]uuid.UUID{
		"quizID": quiz.ID,
	})
}

// CreateQuiz will create new quiz.
func (controller *quizController) GetQuiz(c *fiber.Ctx) error {
	quizID, err := uuid.Parse(c.Params("quizID"))
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	quiz, err := controller.service.GetQuiz(quizID)
	if err != nil {
		controller.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(quiz)
}
