package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestValidateEmptyAnswer will test for empty answer in option
func TestValidateEmptyAnswer(t *testing.T) {
	falseValue := false

	option := Option{
		Answer:     "",
		QuestionID: uuid.New(),
		IsCorrect:  &falseValue,
	}

	err := option.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, "answer must be specified", err.Error())
}

// TestValidateAnswerRegex will test for valid characters in answer field
func TestValidateAnswerRegex(t *testing.T) {
	falseValue := false

	option := Option{
		Answer:     "# This is invalid answer",
		QuestionID: uuid.New(),
		IsCorrect:  &falseValue,
	}

	err := option.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, "answer contains invalid characters", err.Error())
}

// TestValidateOptionIsCorrect will test for valid value in iscorret field
func TestValidateOptionIsCorrect(t *testing.T) {

	option := Option{
		Answer:     "This is invalid answer",
		QuestionID: uuid.New(),
		IsCorrect:  nil,
	}

	err := option.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, "whether answer is correct or not must be specified", err.Error())
}
