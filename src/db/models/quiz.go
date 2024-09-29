package models

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/shaileshhb/quiz/src/utils"
)

// Quiz will contain details related to quiz
type Quiz struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	MaxTime   uint64     `json:"maxTime"` // this will store time in minutes. Default value is 2 minutes
	Questions []Question `json:"questions"`
}

// Validate will validate if all fields of quiz are valid.
func (q *Quiz) Validate() error {
	q.Title = strings.TrimSpace(q.Title)

	if len(q.Title) == 0 {
		return errors.New("title must be specified")
	}

	if len(q.Title) < 5 || len(q.Title) > 50 {
		return errors.New("title must be between 5 and 50 characters")
	}

	isValid, err := utils.ValidateString(q.Title, `^[a-zA-Z0-9@$()!%*/?&\s]+$`)
	if err != nil {
		return err
	}

	if !isValid {
		return errors.New("title contains invalid characters")
	}

	if len(q.Questions) == 0 {
		return errors.New("at least one question is required")
	}

	for _, question := range q.Questions {
		err := question.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
