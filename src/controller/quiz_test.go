package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shaileshhb/quiz/src/db/models"
	"github.com/shaileshhb/quiz/src/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock for the service that handles quiz creation
type MockService struct {
	mock.Mock
}

var logger = log.InitializeLogger()

func (s *MockService) Create(quiz *models.Quiz) error {
	args := s.Called(quiz)
	return args.Error(0)
}

func (s *MockService) GetQuiz(quizID uuid.UUID) (*models.Quiz, error) {
	args := s.Called(quizID)
	return args.Get(0).(*models.Quiz), args.Error(1)
}

func TestCreateQuiz(t *testing.T) {
	app := fiber.New()

	mockService := new(MockService)

	quizController := NewQuizController(mockService, logger)

	app.Post("/quizzes", quizController.CreateQuiz)

	t.Run("Invalid Request Body", func(t *testing.T) {
		// Simulate an invalid body (non-JSON)
		req := httptest.NewRequest(http.MethodPost, "/quizzes", bytes.NewBuffer([]byte("invalid body")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		// Assert the response status is 400 BadRequest
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Validation Error", func(t *testing.T) {
		quiz := models.Quiz{}

		body, _ := json.Marshal(&quiz)

		req := httptest.NewRequest(http.MethodPost, "/quizzes", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		// Assert the response status is 400 BadRequest
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Service Error", func(t *testing.T) {
		trueValue := true
		falseValue := false

		quizOne := models.Quiz{
			Title:   "Quiz Title 1",
			MaxTime: 10,
			Questions: []models.Question{
				{
					Text: "Question 1",
					Options: []models.Option{
						{
							Answer:    "Answer 1",
							IsCorrect: &trueValue,
						}, {
							Answer:    "Answer 2",
							IsCorrect: &falseValue,
						}, {
							Answer:    "Answer 3",
							IsCorrect: &falseValue,
						}, {
							Answer:    "Answer 4",
							IsCorrect: &falseValue,
						},
					},
				},
			},
		}

		quizBytes, _ := json.Marshal(&quizOne)

		mockService.On("Create", mock.Anything).Return(errors.New("service error")).Once()

		req := httptest.NewRequest(http.MethodPost, "/quizzes", bytes.NewBuffer(quizBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		// Assert the response status is 400 BadRequest
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockService.AssertExpectations(t)
	})
}
