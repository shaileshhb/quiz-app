package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shaileshhb/quiz/src/db"
	"github.com/shaileshhb/quiz/src/db/models"
	"github.com/stretchr/testify/assert"
)

// TestStartQuizForInvalidUser will test start quiz for a user who does not exist.
func TestStartQuizForInvalidUser(t *testing.T) {
	database := db.NewDatabase()
	serv := NewUserQuizService(database)

	userQuiz := models.UserQuizAttempts{
		UserID: uuid.New(),
		QuizID: uuid.New(),
	}

	err := serv.StartQuiz(&userQuiz)

	assert.NotNil(t, err)
	assert.Equal(t, "user not found", err.Error())
}

// TestStartQuizForInvalidQuiz will test start quiz for a quiz which does not exist.
func TestStartQuizForInvalidQuiz(t *testing.T) {
	database := db.NewDatabase()
	serv := NewUserQuizService(database)

	userID, _ := uuid.Parse("bfc8ec19-124b-40a1-8936-12dace6fd162")

	userQuiz := models.UserQuizAttempts{
		UserID: userID,
		QuizID: uuid.New(),
	}

	err := serv.StartQuiz(&userQuiz)

	assert.NotNil(t, err)
	assert.Equal(t, "quiz not found", err.Error())
}

// TestStartQuizForAttemptedQuiz will test start quiz for a quiz for which user has already attempted.
func TestStartQuizForAttemptedQuiz(t *testing.T) {
	database := db.NewDatabase()
	serv := NewUserQuizService(database)

	quizID, _ := uuid.Parse("997f06f9-89d1-4f95-9300-09caee4d6b40")
	userID, _ := uuid.Parse("bfc8ec19-124b-40a1-8936-12dace6fd162")

	userQuiz := models.UserQuizAttempts{
		UserID: userID,
		QuizID: quizID,
	}

	_ = serv.StartQuiz(&userQuiz)
	err := serv.StartQuiz(&userQuiz)

	assert.NotNil(t, err)
	assert.Equal(t, "user has already attempted this quiz", err.Error())
}

// TestStartQuiz will test start quiz.
func TestStartQuiz(t *testing.T) {
	database := db.NewDatabase()
	serv := NewUserQuizService(database)

	quizID, _ := uuid.Parse("997f06f9-89d1-4f95-9300-09caee4d6b40")
	userID, _ := uuid.Parse("bfc8ec19-124b-40a1-8936-12dace6fd162")

	userQuiz := models.UserQuizAttempts{
		UserID: userID,
		QuizID: quizID,
	}

	totalQuizzes := len(database.UserQuizAttempts)
	err := serv.StartQuiz(&userQuiz)

	assert.Nil(t, err)
	assert.Equal(t, len(database.UserQuizAttempts), totalQuizzes+1)
}
