package db

import (
	"github.com/shaileshhb/quiz/src/db/models"
)

// Database will mimic a database
type Database struct {
	Quiz []models.Quiz
}

// NewDatabase will initialize a new database instance
func NewDatabase() *Database {
	return &Database{
		Quiz: []models.Quiz{},
	}
}
