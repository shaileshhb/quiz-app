package controllers

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/shaileshhb/quiz/src/db"
	"github.com/shaileshhb/quiz/src/db/models"
	"github.com/shaileshhb/quiz/src/db/validations"
)

// UserQuizController will consist of controllers methods that would be implemented by userQuizController
type UserQuizController interface {
	StartQuiz(userQuiz *models.UserQuizAttempts) error
	SubmitAnswer(userResponse *models.UserResponse) error
}

// userQuizController will contain reference to db.
type userQuizController struct {
	db *db.Database
}

// NewUserQuizController will create new instance of userQuizController
func NewUserQuizController(db *db.Database) UserQuizController {
	return &userQuizController{
		db: db,
	}
}

// StartQuiz will start a quiz for a user.
func (controller *userQuizController) StartQuiz(userQuiz *models.UserQuizAttempts) error {

	err := validations.DoesUserIDExist(controller.db, userQuiz.UserID)
	if err != nil {
		return err
	}

	err = validations.DoesQuizIDExist(controller.db, userQuiz.QuizID)
	if err != nil {
		return err
	}

	// check if user has attempted this quiz
	for _, attempt := range controller.db.UserQuizAttempts {
		if attempt.UserID == userQuiz.UserID && attempt.QuizID == userQuiz.QuizID {
			return errors.New("user has already attempted this quiz")
		}
	}

	userQuiz.TotalScore = 0
	userQuiz.ID = uuid.New()

	controller.db.UserQuizAttempts = append(controller.db.UserQuizAttempts, *userQuiz)

	return nil
}

// SubmitAnswer will submit user's answer for a given question.
func (controller *userQuizController) SubmitAnswer(userResponse *models.UserResponse) error {
	err := validations.DoesUserIDExist(controller.db, userResponse.UserID)
	if err != nil {
		return err
	}

	err = validations.DoesQuizIDExist(controller.db, userResponse.QuizID)
	if err != nil {
		return err
	}

	userQuiz, err := controller.getUserQuiz(userResponse.UserQuizAttemptID)
	if err != nil {
		return err
	}

	err = controller.isQuestionAnswered(userQuiz, userResponse.QuestionID)
	if err != nil {
		return err
	}

	err = controller.doesQuestionExistForQuiz(userResponse.QuizID, userResponse.QuestionID)
	if err != nil {
		return err
	}

	quiz, err := controller.getQuizByID(userResponse.QuizID)
	if err != nil {
		return err
	}

	question, err := controller.getQuestionByID(quiz, userResponse.QuestionID)
	if err != nil {
		return err
	}

	option, err := controller.getOptionByID(question, userResponse.SelectedOptionID)
	if err != nil {
		return err
	}

	if *option.IsCorrect {
		userResponse.IsCorrect = true
		controller.updateUserQuizScore(userResponse.UserQuizAttemptID)
	}

	userResponse.ID = uuid.New()

	for i, attempts := range controller.db.UserQuizAttempts {
		if attempts.ID == userResponse.UserQuizAttemptID {

			controller.db.UserQuizAttempts[i].UserResponses = append(controller.db.UserQuizAttempts[i].UserResponses, *userResponse)
			break
		}
	}

	fmt.Println("==============================================")
	fmt.Printf("%+v\n", controller.db.UserQuizAttempts)
	fmt.Println("==============================================")

	return nil
}

// updateUserQuizScore will updaate the total score of the UserQuizAttempts.
func (controller *userQuizController) updateUserQuizScore(userQuizAttemptID uuid.UUID) {
	for i, attempts := range controller.db.UserQuizAttempts {
		if attempts.ID == userQuizAttemptID {
			controller.db.UserQuizAttempts[i].TotalScore++
		}
	}
}

// getQuizByID will fetch quiz by given quizID.
func (controller *userQuizController) getQuizByID(quizID uuid.UUID) (*models.Quiz, error) {
	for _, quiz := range controller.db.Quiz {
		if quiz.ID == quizID {
			return &quiz, nil
		}
	}

	return nil, errors.New("quiz not found")
}

// getQuestionByID will fetch question by given questionID.
func (controller *userQuizController) getQuestionByID(quiz *models.Quiz, questionID uuid.UUID) (*models.Question, error) {
	for _, question := range quiz.Questions {
		if question.ID == questionID {
			return &question, nil
		}
	}

	return nil, errors.New("question not found")
}

// getOptionByID will fetch option by given optionID.
func (controller *userQuizController) getOptionByID(question *models.Question, optionID uuid.UUID) (*models.Option, error) {
	for _, option := range question.Options {
		if option.ID == optionID {
			return &option, nil
		}
	}

	return nil, errors.New("option not found")
}

// doesQuestionExistForQuiz will check if question exist in the given quiz.
func (controller *userQuizController) doesQuestionExistForQuiz(quizID, questionID uuid.UUID) error {

	quiz, err := controller.getQuizByID(quizID)
	if err != nil {
		return err
	}

	for _, question := range quiz.Questions {
		if question.ID == questionID {
			return nil
		}
	}
	return errors.New("question not found for specified quiz")
}

// getUserQuiz will check if quiz has started for a given user, if not then it will return an error
func (controller *userQuizController) getUserQuiz(userQuizAttemptID uuid.UUID) (*models.UserQuizAttempts, error) {
	for _, attempts := range controller.db.UserQuizAttempts {
		if attempts.ID == userQuizAttemptID {
			return &attempts, nil
		}
	}
	return nil, errors.New("please start quiz before submitting answers")
}

// isQuestionAnswered will check if all question has been answered for a given user, if yes then it will return an error
func (controller *userQuizController) isQuestionAnswered(userQuiz *models.UserQuizAttempts, questionID uuid.UUID) error {

	fmt.Println("================================================================")
	fmt.Printf("%+v\n", userQuiz.UserResponses)
	fmt.Printf("%+v\n", questionID)
	fmt.Println("================================================================")

	for _, repsonse := range userQuiz.UserResponses {
		if repsonse.QuestionID == questionID {
			return errors.New("question already answered")
		}
	}

	return nil
}
