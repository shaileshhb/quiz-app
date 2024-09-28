package models

import (
	"errors"

	"github.com/google/uuid"
)

// UserResponse is the structure for user response for specific quiz
type UserResponse struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"userID"`
	QuizID            uuid.UUID `json:"quizID"`
	UserQuizAttemptID uuid.UUID `json:"userQuizAttemptID"`
	QuestionID        uuid.UUID `json:"questionID"`
	SelectedOptionID  uuid.UUID `json:"selectedOptionID"`
	IsCorrect         bool      `json:"isCorrect"`
}

// Validate will validate if all fields for a user response are correctly specified.
func (u *UserResponse) Validate() error {
	if u.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}

	if u.QuizID == uuid.Nil {
		return errors.New("quiz ID is required")
	}

	if u.QuestionID == uuid.Nil {
		return errors.New("question ID is required")
	}

	if u.SelectedOptionID == uuid.Nil {
		return errors.New("selected option ID is required")
	}
	return nil
}
