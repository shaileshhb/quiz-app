package models

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shaileshhb/quiz/src/utils"
)

// Option will contain options for specific question.
type Option struct {
	ID         uuid.UUID `json:"id"`
	QuestionID uuid.UUID `json:"questionID"`
	Answer     string    `json:"answer"`
	IsCorrect  *bool     `json:"isCorrect,omitempty"`
}

// Validate will validate if all fields for a option are correctly specified.
func (o *Option) Validate() error {
	if len(o.Answer) == 0 {
		return errors.New("answer must be specified")
	}

	if len(o.Answer) > 200 {
		return errors.New("answer should not exceed 200 characters")
	}

	isValid, err := utils.ValidateString(o.Answer, `^[a-zA-Z0-9@$!%*/&\s]+$`)
	if err != nil {
		return err
	}

	if !isValid {
		return errors.New("answer contains invalid characters")
	}

	if o.IsCorrect == nil {
		return errors.New("whether answer is correct or not must be specified")
	}

	return nil
}
