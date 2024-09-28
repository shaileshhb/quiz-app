package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shaileshhb/quiz/src/controllers"
	"github.com/shaileshhb/quiz/src/db/models"
)

// quizRoute contains reference to quiz controller and logger
type quizRoute struct {
	con controllers.QuizController
	log zerolog.Logger
}

func NewQuizRoute(con controllers.QuizController,
	log zerolog.Logger) *quizRoute {
	return &quizRoute{
		con: con,
		log: log,
	}
}

func (qr *quizRoute) RegisterRoute(router fiber.Router) {
	router.Post("/quiz", qr.CreateQuiz)
	router.Get("/quiz/:quizID", qr.GetQuiz)
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

	err = quiz.Validate()
	if err != nil {
		qr.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = qr.con.Create(&quiz)
	if err != nil {
		qr.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(map[string]uuid.UUID{
		"quizID": quiz.ID,
	})
}

// CreateQuiz will create new quiz.
func (qr *quizRoute) GetQuiz(c *fiber.Ctx) error {
	quizID, err := uuid.Parse(c.Params("quizID"))
	if err != nil {
		qr.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	quiz, err := qr.con.GetQuiz(quizID)
	if err != nil {
		qr.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(quiz)
}
