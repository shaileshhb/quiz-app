package db

import (
	"github.com/shaileshhb/quiz/src/db/models"
)

var Quiz []models.Quiz

func init() {
	Quiz = make([]models.Quiz, 0)
}
