package controllers

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shaileshhb/quiz/src/db"
	"github.com/shaileshhb/quiz/src/db/models"
	"github.com/shaileshhb/quiz/src/db/validations"
)

// UserQuizController will consist of controllers methods that would be implemented by userQuizController
type UserQuizController interface {
	StartQuiz(*models.UserQuizAttempts) error
	SubmitAnswer(*models.UserResponse) (*models.Option, error)
	GetUserQuizResults(uuid.UUID, uuid.UUID) (*models.UserQuizAttempts, error)
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

	startTime := time.Now()
	userQuiz.StartedAt = &startTime
	userQuiz.TotalScore = 0
	userQuiz.ID = uuid.New()

	controller.db.UserQuizAttempts = append(controller.db.UserQuizAttempts, *userQuiz)

	return nil
}

// SubmitAnswer will submit user's answer for a given question and return correct answer and error if any.
func (controller *userQuizController) SubmitAnswer(userResponse *models.UserResponse) (*models.Option, error) {
	err := validations.DoesUserIDExist(controller.db, userResponse.UserID)
	if err != nil {
		return nil, err
	}

	err = validations.DoesQuizIDExist(controller.db, userResponse.QuizID)
	if err != nil {
		return nil, err
	}

	err = controller.isQuizCompleted(userResponse)
	if err != nil {
		return nil, err
	}

	userQuiz, err := controller.getUserQuiz(userResponse.UserQuizAttemptID)
	if err != nil {
		return nil, err
	}

	err = controller.isQuestionAnswered(userQuiz, userResponse.QuestionID)
	if err != nil {
		return nil, err
	}

	err = controller.doesQuestionExistForQuiz(userResponse.QuizID, userResponse.QuestionID)
	if err != nil {
		return nil, err
	}

	quiz, err := controller.getQuizByID(userResponse.QuizID)
	if err != nil {
		return nil, err
	}

	question, err := controller.getQuestionByID(quiz, userResponse.QuestionID)
	if err != nil {
		return nil, err
	}

	correctOption := &models.Option{}
	option, err := controller.getOptionByID(question, userResponse.SelectedOptionID)
	if err != nil {
		return nil, err
	}

	if *option.IsCorrect {
		userResponse.IsCorrect = true
		controller.updateUserQuizScore(userResponse.UserQuizAttemptID)
		correctOption = option
	} else {
		// get correct option
		correctOption, err = controller.getCorrectOption(question)
		if err != nil {
			return nil, err
		}
	}

	userResponse.ID = uuid.New()

	for i, attempts := range controller.db.UserQuizAttempts {
		if attempts.ID == userResponse.UserQuizAttemptID {
			controller.db.UserQuizAttempts[i].UserResponses = append(controller.db.UserQuizAttempts[i].UserResponses, *userResponse)

			if len(quiz.Questions) == len(controller.db.UserQuizAttempts[i].UserResponses) {
				endedAt := time.Now()
				controller.db.UserQuizAttempts[i].EndedAt = &endedAt
				break
			}
			break
		}
	}

	return correctOption, nil
}

// GetUserQuizResults will return results for specific quiz for specified user.
func (controller *userQuizController) GetUserQuizResults(userID, quizID uuid.UUID) (*models.UserQuizAttempts, error) {

	err := validations.DoesUserIDExist(controller.db, userID)
	if err != nil {
		return nil, err
	}

	err = validations.DoesQuizIDExist(controller.db, quizID)
	if err != nil {
		return nil, err
	}

	for _, attempt := range controller.db.UserQuizAttempts {
		if attempt.UserID == userID && attempt.QuizID == quizID {
			return &attempt, nil
		}
	}

	return nil, errors.New("user not attempted specified quiz")
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

// getCorrectOption will fetch correct option.
func (controller *userQuizController) getCorrectOption(question *models.Question) (*models.Option, error) {
	for _, option := range question.Options {
		if *option.IsCorrect {
			return &option, nil
		}
	}

	return nil, errors.New("correct option not found")
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
	for _, repsonse := range userQuiz.UserResponses {
		if repsonse.QuestionID == questionID {
			return errors.New("question already answered")
		}
	}

	return nil
}

// isQuizCompleted will check if quiz has ended or max time is exceeded.
func (controller *userQuizController) isQuizCompleted(userResponse *models.UserResponse) error {
	userQuiz, err := controller.getUserQuiz(userResponse.UserQuizAttemptID)
	if err != nil {
		return err
	}

	if userQuiz.EndedAt != nil {
		return errors.New("cannot answer questions after quiz has ended")
	}

	quiz, err := controller.getQuizByID(userQuiz.QuizID)
	if err != nil {
		return err
	}

	if time.Since(*userQuiz.StartedAt) > time.Duration(quiz.MaxTime*uint64(time.Minute)) {
		return errors.New("maximum time exceeded for this quiz")
	}

	return nil
}
