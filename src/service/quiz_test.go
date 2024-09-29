package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shaileshhb/quiz/src/db"
	"github.com/shaileshhb/quiz/src/db/models"
	"github.com/stretchr/testify/assert"
)

// TestDuplicateCreate will test for duplicate quiz title creation.
func TestDuplicateCreate(t *testing.T) {
	database := db.NewDatabase()
	quizService := NewQuizService(database)

	quizOne := models.Quiz{
		Title:   "Quiz Title 1",
		MaxTime: 10,
		Questions: []models.Question{
			{
				Text: "Question 1",
				Options: []models.Option{
					{
						Answer: "Answer 1",
					},
				},
			},
		},
	}

	quizTwo := models.Quiz{
		Title:   "Quiz Title 1",
		MaxTime: 10,
		Questions: []models.Question{
			{
				Text: "Question 1",
				Options: []models.Option{
					{
						Answer: "Answer 1",
					},
				},
			},
		},
	}

	quizService.Create(&quizOne)

	err := quizService.Create(&quizTwo)

	assert.NotNil(t, err)
	assert.Equal(t, "quiz with same title already exists", err.Error())
}

// TestDefaultQuizTime will test for default quiz time if not provided.
func TestDefaultQuizTime(t *testing.T) {
	database := db.NewDatabase()
	quizService := NewQuizService(database)

	quizOne := models.Quiz{
		Title: "Quiz Title 1",
		Questions: []models.Question{
			{
				Text: "Question 1",
				Options: []models.Option{
					{
						Answer: "Answer 1",
					},
				},
			},
		},
	}

	err := quizService.Create(&quizOne)

	assert.Nil(t, err)
	assert.Equal(t, uint64(2), quizOne.MaxTime)
}

// TestCreateAssignID will test for assigning unique IDs to quiz.
func TestCreateAssignID(t *testing.T) {
	database := db.NewDatabase()
	quizService := NewQuizService(database)

	quizOne := models.Quiz{
		Title: "Quiz Title 1",
		Questions: []models.Question{
			{
				Text: "Question 1",
				Options: []models.Option{
					{
						Answer: "Answer 1",
					},
				},
			},
		},
	}

	totalQuizzes := len(database.Quiz)

	err := quizService.Create(&quizOne)

	assert.Nil(t, err)
	assert.Equal(t, len(database.Quiz), totalQuizzes+1)
}

// TestGetQuizNotFound will test for not found quiz
func TestGetQuizNotFound(t *testing.T) {
	database := db.NewDatabase()
	quizService := NewQuizService(database)

	quizOne := models.Quiz{
		Title: "Quiz Title 1",
		Questions: []models.Question{
			{
				Text: "Question 1",
				Options: []models.Option{
					{
						Answer: "Answer 1",
					},
				},
			},
		},
	}

	quizID := uuid.New()
	_ = quizService.Create(&quizOne)
	_, err := quizService.GetQuiz(quizID)

	assert.NotNil(t, err)
	assert.Equal(t, "quiz not found", err.Error())
}

// TestGetQuiz will test for fetch quiz by quizID
func TestGetQuiz(t *testing.T) {
	database := db.NewDatabase()
	quizService := NewQuizService(database)

	quizID, _ := uuid.Parse("997f06f9-89d1-4f95-9300-09caee4d6b40")
	quiz, err := quizService.GetQuiz(quizID)

	assert.Nil(t, err)
	assert.Equal(t, quiz.ID, quizID)
}
