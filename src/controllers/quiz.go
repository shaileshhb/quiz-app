package controllers

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/shaileshhb/quiz/src/db"
	"github.com/shaileshhb/quiz/src/db/models"
	"github.com/shaileshhb/quiz/src/utils"
)

// QuizController will consist of controllers methods that would be implemented by quizController
type QuizController interface {
	Create(quiz *models.Quiz) error
	GetQuiz(quizID uuid.UUID) (*models.Quiz, error)
}

// quizController will contain reference to db.
type quizController struct {
	db *db.Database
}

// NewQuizController will create new instance of quizController
func NewQuizController(db *db.Database) QuizController {
	return &quizController{db: db}
}

// Create will create new quiz in database.
func (controller *quizController) Create(quiz *models.Quiz) error {
	quiz.Title = strings.TrimSpace(quiz.Title)

	err := controller.checkTitleExist(quiz.Title)
	if err != nil {
		return err
	}

	controller.assignIDs(quiz)

	controller.db.Quiz = append(controller.db.Quiz, *quiz)
	return nil
}

// GetQuiz will get quiz by ID from database.
func (controller *quizController) GetQuiz(quizID uuid.UUID) (*models.Quiz, error) {
	currentQuiz := models.Quiz{}
	isQuizFound := false

	for _, q := range controller.db.Quiz {
		if q.ID == quizID {
			isQuizFound = true
			cq, _ := json.Marshal(q)
			json.Unmarshal(cq, &currentQuiz)
			currentQuiz = utils.CopyQuiz(q)
			break
		}
	}

	if !isQuizFound {
		return nil, errors.New("quiz not found")
	}

	return &currentQuiz, nil
}

// checkTitleExist will check if quiz with same title already exists in database.
func (controller *quizController) checkTitleExist(title string) error {
	for _, quiz := range controller.db.Quiz {
		if strings.EqualFold(strings.ToLower(quiz.Title), strings.ToLower(title)) {
			return errors.New("quiz with same title already exists")
		}
	}

	return nil
}

func (controller *quizController) assignIDs(quiz *models.Quiz) {
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
