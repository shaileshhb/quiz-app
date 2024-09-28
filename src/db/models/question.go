package models

import (
	"errors"

	"github.com/google/uuid"
)

// Question will contain question details for a quiz
type Question struct {
	ID      uuid.UUID `json:"id"`
	Text    string    `json:"text"`
	QuizID  uuid.UUID `json:"quizID"`
	Options []Option  `json:"options"`
}

// Validate will validate if all fields for a question are correctly specified.
func (q *Question) Validate() error {
	if len(q.Text) > 200 {
		return errors.New("text should not exceed 200 characters")
	}

	if len(q.Options) != 4 {
		return errors.New("question should have exactly 4 options")
	}

	isCorrectPresent := false

	for _, option := range q.Options {
		err := option.Validate()
		if err != nil {
			return err
		}

		if option.IsCorrect != nil && *option.IsCorrect {
			isCorrectPresent = true
		}
	}

	if !isCorrectPresent {
		return errors.New("atleast one correct option must be present")
	}

	return nil
}
