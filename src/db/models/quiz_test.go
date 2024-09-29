package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidateEmptyTitle will test for empty quiz title
func TestValidateEmptyTitle(t *testing.T) {
	quiz := Quiz{
		Title:     "",
		MaxTime:   2,
		Questions: []Question{},
	}

	err := quiz.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, "title must be specified", err.Error())
}

// TestValidateTitleLessThanFive will test for quiz title less than five characters
func TestValidateTitleLessThanFive(t *testing.T) {
	quiz := Quiz{
		Title:     "one",
		MaxTime:   2,
		Questions: []Question{},
	}

	err := quiz.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, "title must be between 5 and 50 characters", err.Error())
}

// TestValidateTitleGreaterThanFifty will test for quiz title more than twenty characters
func TestValidateTitleGreaterThanFifty(t *testing.T) {
	quiz := Quiz{
		Title:     "thisisaverylongtitlethatexceeds50charactersneedtoaddmorecharacters",
		MaxTime:   2,
		Questions: []Question{},
	}

	err := quiz.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, "title must be between 5 and 50 characters", err.Error())
}

// TestValidateTitleRegex will test if title has valid characters
func TestValidateTitleRegex(t *testing.T) {
	quiz := Quiz{
		Title:     "# invalid regex (title)",
		MaxTime:   2,
		Questions: []Question{},
	}

	err := quiz.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, "title contains invalid characters", err.Error())
}

// TestValidateTotalQuestions will test if atleast one question is specified
func TestValidateTotalQuestions(t *testing.T) {
	quiz := Quiz{
		Title:     "Sample quiz",
		MaxTime:   2,
		Questions: []Question{},
	}

	err := quiz.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, "at least one question is required", err.Error())
}
