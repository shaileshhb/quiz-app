package db

import (
	"github.com/google/uuid"
	"github.com/shaileshhb/quiz/src/db/models"
	"github.com/shaileshhb/quiz/src/security"
)

// Database will mimic a database
type Database struct {
	Quiz             []models.Quiz
	Users            []models.User
	UserQuizAttempts []models.UserQuizAttempts
}

// NewDatabase will initialize a new database instance
func NewDatabase() *Database {
	db := &Database{
		Quiz: []models.Quiz{},
	}

	createDummyQuiz(db)
	createDummyUsers(db)

	return db
}

func createDummyQuiz(db *Database) {
	quizID, _ := uuid.Parse("997f06f9-89d1-4f95-9300-09caee4d6b40")
	quiz := models.Quiz{
		ID:      quizID,
		Title:   "Sample Quiz",
		MaxTime: 1,
	}

	quiz.Questions = createDummyQuestions(quiz.ID)
	db.Quiz = append(db.Quiz, quiz)
}

func createDummyQuestions(quizID uuid.UUID) []models.Question {
	questions := make([]models.Question, 0)
	question := models.Question{
		ID:     uuid.New(),
		Text:   "What is the capital of France?",
		QuizID: quizID,
	}

	question.Options = createDummyOptions(question.ID)

	questions = append(questions, question)

	question = models.Question{
		ID:     uuid.New(),
		Text:   "What is the capital of India?",
		QuizID: quizID,
	}

	question.Options = createDummyOptions(question.ID)
	questions = append(questions, question)

	return questions
}

func createDummyOptions(questionID uuid.UUID) []models.Option {
	trueValue := true
	falseValue := false

	options := []models.Option{
		{
			ID:         uuid.New(),
			QuestionID: questionID,
			Answer:     "Paris",
			IsCorrect:  &trueValue,
		},
		{
			ID:         uuid.New(),
			QuestionID: questionID,
			Answer:     "Delhi",
			IsCorrect:  &falseValue,
		},
		{
			ID:         uuid.New(),
			QuestionID: questionID,
			Answer:     "London",
			IsCorrect:  &falseValue,
		},
		{
			ID:         uuid.New(),
			QuestionID: questionID,
			Answer:     "New York",
			IsCorrect:  &falseValue,
		},
	}

	return options
}

func createDummyUsers(db *Database) {
	id, _ := uuid.Parse("bfc8ec19-124b-40a1-8936-12dace6fd162")
	password, _ := security.HashPassword("userone")
	user := models.User{
		ID:       id,
		Name:     "User one",
		Username: "userone",
		Password: string(password),
	}

	db.Users = append(db.Users, user)

	password, _ = security.HashPassword("usertwo")
	user = models.User{
		ID:       uuid.New(),
		Name:     "User two",
		Username: "usertwo",
		Password: string(password),
	}

	db.Users = append(db.Users, user)
}
