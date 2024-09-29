package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	model "github.com/shaileshhb/quiz/src/db/models"
)

// TestValidateEmptyText will test for empty question text
func TestValidateEmptyText(t *testing.T) {
	question := model.Question{
		Text:    "",
		QuizID:  uuid.New(),
		Options: []model.Option{},
	}

	err := question.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, "text must be specified", err.Error())
}

// TestValidateTextRegex will test for valid characters in text
func TestValidateTextRegex(t *testing.T) {
	question := model.Question{
		Text:    "# This is question text",
		QuizID:  uuid.New(),
		Options: []model.Option{},
	}

	err := question.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, "question text contains invalid characters", err.Error())
}

// TestValidateOptions will test for valid number of options
func TestValidateOptions(t *testing.T) {
	question := model.Question{
		Text:    "This is question text",
		QuizID:  uuid.New(),
		Options: []model.Option{},
	}

	err := question.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, "question should have exactly 4 options", err.Error())
}

// TestValidateOptionsMax will test for valid number of options
func TestValidateOptionsMax(t *testing.T) {
	trueValue := true
	falseValue := false

	question := model.Question{
		Text:   "This is question text",
		QuizID: uuid.New(),
		Options: []model.Option{
			{
				Answer:    "This is answer 1",
				IsCorrect: &trueValue,
			}, {
				Answer:    "This is answer 2",
				IsCorrect: &falseValue,
			}, {
				Answer:    "This is answer 3",
				IsCorrect: &falseValue,
			}, {
				Answer:    "This is answer 4",
				IsCorrect: &falseValue,
			}, {
				Answer:    "This is answer 5",
				IsCorrect: &falseValue,
			},
		},
	}

	err := question.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, "question should have exactly 4 options", err.Error())
}

// TestValidateIsCorrect will test for one correct option exist
func TestValidateIsCorrect(t *testing.T) {
	falseValue := false

	question := model.Question{
		Text:   "This is question text",
		QuizID: uuid.New(),
		Options: []model.Option{
			{
				Answer:    "This is answer 1",
				IsCorrect: &falseValue,
			}, {
				Answer:    "This is answer 2",
				IsCorrect: &falseValue,
			}, {
				Answer:    "This is answer 3",
				IsCorrect: &falseValue,
			}, {
				Answer:    "This is answer 4",
				IsCorrect: &falseValue,
			},
		},
	}

	err := question.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, "atleast one correct option must be present", err.Error())
}
