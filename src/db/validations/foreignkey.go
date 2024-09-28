package validations

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shaileshhb/quiz/src/db"
)

// DoesUserIDExist will check if userID exist in the database, if not then return an error
func DoesUserIDExist(database *db.Database, userID uuid.UUID) error {
	for _, user := range database.Users {
		if user.ID == userID {
			return nil
		}
	}

	return errors.New("user not found")
}

// DoesQuizIDExist will check if quizID exist in the database, if not then return an error
func DoesQuizIDExist(database *db.Database, quizID uuid.UUID) error {
	for _, quiz := range database.Quiz {
		if quiz.ID == quizID {
			return nil
		}
	}

	return errors.New("quiz not found")
}
