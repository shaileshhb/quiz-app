package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/shaileshhb/quiz/src/db/models"
)

// quizRoute contains reference to quiz controller and logger
type quizRoute struct {
	log zerolog.Logger
}

func NewQuizRoute(log zerolog.Logger) *quizRoute {
	return &quizRoute{
		log: log,
	}
}

func (qr *quizRoute) RegisterRoute(router fiber.Router) {
	router.Post("/register", qr.CreateQuiz)
	qr.log.Info().Msg("Quiz routes registered")
}

// CreateQuiz will create new quiz.
func (qr *quizRoute) CreateQuiz(c *fiber.Ctx) error {
	quiz := models.Quiz{}

	err := c.BodyParser(&quiz)
	if err != nil {
		qr.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Add logic to create quiz in the database

	return c.Status(http.StatusCreated).JSON(quiz)
}
