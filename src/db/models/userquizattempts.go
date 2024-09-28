package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// UserQuizAttempts will contain details about a user and quiz they have given.
type UserQuizAttempts struct {
	ID            uuid.UUID      `json:"id"`
	UserID        uuid.UUID      `json:"userID"`
	QuizID        uuid.UUID      `json:"quizID"`
	StartedAt     *time.Time     `json:"startedAt"`
	EndedAt       *time.Time     `json:"endAt"`
	TotalScore    uint32         `json:"totalScore"`
	UserResponses []UserResponse `json:"userResponses"`
}

// Validate will check if valid userID and quizID are provided.
func (u *UserQuizAttempts) Validate() error {
	if u.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}

	if u.QuizID == uuid.Nil {
		return errors.New("quiz ID is required")
	}
	return nil
}
