package models

import (
	"errors"

	"github.com/google/uuid"
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
	if len(q.Title) > 50 {
		return errors.New("title should not exceed 50 characters")
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
