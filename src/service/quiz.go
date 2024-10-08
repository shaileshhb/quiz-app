package service

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/shaileshhb/quiz/src/db"
	"github.com/shaileshhb/quiz/src/db/models"
)

// QuizService will consist of service methods that would be implemented by quizService
type QuizService interface {
	Create(quiz *models.Quiz) error
	GetQuiz(quizID uuid.UUID) (*models.Quiz, error)
}

// quizService will contain reference to db.
type quizService struct {
	db *db.Database
}

// NewQuizService will create new instance of quizService
func NewQuizService(db *db.Database) QuizService {
	return &quizService{db: db}
}

// Create will create new quiz in database.
func (service *quizService) Create(quiz *models.Quiz) error {
	err := service.checkTitleExist(quiz.Title)
	if err != nil {
		return err
	}

	if quiz.MaxTime == 0 {
		quiz.MaxTime = 2
	}

	service.assignIDs(quiz)

	service.db.Quiz = append(service.db.Quiz, *quiz)
	return nil
}

// GetQuiz will get quiz by ID from database.
func (service *quizService) GetQuiz(quizID uuid.UUID) (*models.Quiz, error) {
	currentQuiz := models.Quiz{}
	isQuizFound := false

	for _, q := range service.db.Quiz {
		if q.ID == quizID {
			isQuizFound = true

			currentQuiz = copyQuiz(q)
			break
		}
	}

	if !isQuizFound {
		return nil, errors.New("quiz not found")
	}

	return &currentQuiz, nil
}

// checkTitleExist will check if quiz with same title already exists in database.
func (service *quizService) checkTitleExist(title string) error {
	for _, quiz := range service.db.Quiz {
		if strings.EqualFold(strings.ToLower(quiz.Title), strings.ToLower(title)) {
			return errors.New("quiz with same title already exists")
		}
	}

	return nil
}

func (service *quizService) assignIDs(quiz *models.Quiz) {
	quiz.ID = uuid.New()

	for i := range quiz.Questions {
		quiz.Questions[i].ID = uuid.New()
		quiz.Questions[i].QuizID = quiz.ID

		for j := range quiz.Questions[i].Options {
			quiz.Questions[i].Options[j].ID = uuid.New()
			quiz.Questions[i].Options[j].QuestionID = quiz.Questions[i].ID
		}
	}
}

func copyQuiz(q models.Quiz) models.Quiz {
	var questions []models.Question

	for _, question := range q.Questions {
		questions = append(questions, copyQuestion(question))
	}

	return models.Quiz{
		ID:        q.ID,
		Title:     q.Title,
		MaxTime:   q.MaxTime,
		Questions: questions,
	}
}

func copyQuestion(q models.Question) models.Question {
	var options []models.Option

	for _, option := range q.Options {
		options = append(options, copyOption(option))
	}

	return models.Question{
		ID:      q.ID,
		QuizID:  q.QuizID,
		Text:    q.Text,
		Options: options,
	}
}

func copyOption(o models.Option) models.Option {
	return models.Option{
		ID:         o.ID,
		QuestionID: o.QuestionID,
		Answer:     o.Answer,
		IsCorrect:  nil,
	}
}
