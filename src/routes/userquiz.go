package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shaileshhb/quiz/src/controllers"
	"github.com/shaileshhb/quiz/src/db/models"
	"github.com/shaileshhb/quiz/src/security"
)

// userQuizRoute contains reference to quiz controller and logger
type userQuizRoute struct {
	con controllers.UserQuizController
	log zerolog.Logger
}

// NewUserQuizRoute will create new instance of userQuizRoute.
func NewUserQuizRoute(con controllers.UserQuizController, log zerolog.Logger) *userQuizRoute {
	return &userQuizRoute{
		con: con,
		log: log,
	}
}

// RegisterRoute registers all endpoints to router.
func (u *userQuizRoute) RegisterRoute(router fiber.Router) {
	router.Post("/user/quiz/:quizID/start", security.MandatoryAuthMiddleware, u.startQuiz)
	router.Post("/user/quiz/:quizID/attempt/:attemptID", security.MandatoryAuthMiddleware, u.submitAnswer)
	u.log.Info().Msg("User quiz routes registered")
}

// startQuiz will start a quiz for a user.
func (qr *userQuizRoute) startQuiz(c *fiber.Ctx) error {
	userQuiz := models.UserQuizAttempts{}

	quizID, err := uuid.Parse(c.Params("quizID"))
	if err != nil {
		qr.log.Error().Err(err).Msg("")
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
		qr.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = qr.con.StartQuiz(&userQuiz)
	if err != nil {
		qr.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(userQuiz)
}

// submitAnswer will submit user's answer for a given question.
func (ur *userQuizRoute) submitAnswer(c *fiber.Ctx) error {
	userResponse := models.UserResponse{}

	err := c.BodyParser(&userResponse)
	if err != nil {
		ur.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userResponse.QuizID, err = uuid.Parse(c.Params("quizID"))
	if err != nil {
		ur.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userResponse.UserQuizAttemptID, err = uuid.Parse(c.Params("attemptID"))
	if err != nil {
		ur.log.Error().Err(err).Msg("")
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
		ur.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = ur.con.SubmitAnswer(&userResponse)
	if err != nil {
		ur.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(userResponse)
}
