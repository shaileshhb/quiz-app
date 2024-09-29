package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shaileshhb/quiz/src/db"
	"github.com/shaileshhb/quiz/src/db/models"
	"github.com/shaileshhb/quiz/src/db/validations"
)

// UserQuizService will consist of service methods that would be implemented by userQuizService
type UserQuizService interface {
	StartQuiz(*models.UserQuizAttempts) error
	SubmitAnswer(*models.UserResponse) (*models.Option, error)
	GetUserQuizResults(uuid.UUID, uuid.UUID) (*models.UserQuizAttempts, error)
}

// userQuizService will contain reference to db.
type userQuizService struct {
	db *db.Database
}

// NewUserQuizService will create new instance of userQuizService
func NewUserQuizService(db *db.Database) UserQuizService {
	return &userQuizService{
		db: db,
	}
}

// StartQuiz will start a quiz for a user.
func (service *userQuizService) StartQuiz(userQuiz *models.UserQuizAttempts) error {

	err := validations.DoesUserIDExist(service.db, userQuiz.UserID)
	if err != nil {
		return err
	}

	err = validations.DoesQuizIDExist(service.db, userQuiz.QuizID)
	if err != nil {
		return err
	}

	// check if user has attempted this quiz
	for _, attempt := range service.db.UserQuizAttempts {
		if attempt.UserID == userQuiz.UserID && attempt.QuizID == userQuiz.QuizID {
			return errors.New("user has already attempted this quiz")
		}
	}

	startTime := time.Now()
	userQuiz.StartedAt = &startTime
	userQuiz.TotalScore = 0
	userQuiz.ID = uuid.New()

	service.db.UserQuizAttempts = append(service.db.UserQuizAttempts, *userQuiz)

	return nil
}

// SubmitAnswer will submit user's answer for a given question and return correct answer and error if any.
func (service *userQuizService) SubmitAnswer(userResponse *models.UserResponse) (*models.Option, error) {
	err := validations.DoesUserIDExist(service.db, userResponse.UserID)
	if err != nil {
		return nil, err
	}

	err = validations.DoesQuizIDExist(service.db, userResponse.QuizID)
	if err != nil {
		return nil, err
	}

	err = service.isQuizCompleted(userResponse)
	if err != nil {
		return nil, err
	}

	userQuiz, err := service.getUserQuiz(userResponse.UserQuizAttemptID)
	if err != nil {
		return nil, err
	}

	err = service.isQuestionAnswered(userQuiz, userResponse.QuestionID)
	if err != nil {
		return nil, err
	}

	err = service.doesQuestionExistForQuiz(userResponse.QuizID, userResponse.QuestionID)
	if err != nil {
		return nil, err
	}

	quiz, err := service.getQuizByID(userResponse.QuizID)
	if err != nil {
		return nil, err
	}

	question, err := service.getQuestionByID(quiz, userResponse.QuestionID)
	if err != nil {
		return nil, err
	}

	correctOption := &models.Option{}
	option, err := service.getOptionByID(question, userResponse.SelectedOptionID)
	if err != nil {
		return nil, err
	}

	if *option.IsCorrect {
		userResponse.IsCorrect = true
		service.updateUserQuizScore(userResponse.UserQuizAttemptID)
		correctOption = option
	} else {
		// get correct option
		correctOption, err = service.getCorrectOption(question)
		if err != nil {
			return nil, err
		}
	}

	userResponse.ID = uuid.New()

	for i, attempts := range service.db.UserQuizAttempts {
		if attempts.ID == userResponse.UserQuizAttemptID {
			service.db.UserQuizAttempts[i].UserResponses = append(service.db.UserQuizAttempts[i].UserResponses, *userResponse)

			if len(quiz.Questions) == len(service.db.UserQuizAttempts[i].UserResponses) {
				endedAt := time.Now()
				service.db.UserQuizAttempts[i].EndedAt = &endedAt
				break
			}
			break
		}
	}

	return correctOption, nil
}

// GetUserQuizResults will return results for specific quiz for specified user.
func (service *userQuizService) GetUserQuizResults(userID, quizID uuid.UUID) (*models.UserQuizAttempts, error) {

	err := validations.DoesUserIDExist(service.db, userID)
	if err != nil {
		return nil, err
	}

	err = validations.DoesQuizIDExist(service.db, quizID)
	if err != nil {
		return nil, err
	}

	for _, attempt := range service.db.UserQuizAttempts {
		if attempt.UserID == userID && attempt.QuizID == quizID {
			return &attempt, nil
		}
	}

	return nil, errors.New("user not attempted specified quiz")
}

// updateUserQuizScore will updaate the total score of the UserQuizAttempts.
func (service *userQuizService) updateUserQuizScore(userQuizAttemptID uuid.UUID) {
	for i, attempts := range service.db.UserQuizAttempts {
		if attempts.ID == userQuizAttemptID {
			service.db.UserQuizAttempts[i].TotalScore++
		}
	}
}

// getQuizByID will fetch quiz by given quizID.
func (service *userQuizService) getQuizByID(quizID uuid.UUID) (*models.Quiz, error) {
	for _, quiz := range service.db.Quiz {
		if quiz.ID == quizID {
			return &quiz, nil
		}
	}

	return nil, errors.New("quiz not found")
}

// getQuestionByID will fetch question by given questionID.
func (service *userQuizService) getQuestionByID(quiz *models.Quiz, questionID uuid.UUID) (*models.Question, error) {
	for _, question := range quiz.Questions {
		if question.ID == questionID {
			return &question, nil
		}
	}

	return nil, errors.New("question not found")
}

// getOptionByID will fetch option by given optionID.
func (service *userQuizService) getOptionByID(question *models.Question, optionID uuid.UUID) (*models.Option, error) {
	for _, option := range question.Options {
		if option.ID == optionID {
			return &option, nil
		}
	}

	return nil, errors.New("option not found")
}

// getCorrectOption will fetch correct option.
func (service *userQuizService) getCorrectOption(question *models.Question) (*models.Option, error) {
	for _, option := range question.Options {
		if *option.IsCorrect {
			return &option, nil
		}
	}

	return nil, errors.New("correct option not found")
}

// doesQuestionExistForQuiz will check if question exist in the given quiz.
func (service *userQuizService) doesQuestionExistForQuiz(quizID, questionID uuid.UUID) error {

	quiz, err := service.getQuizByID(quizID)
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
func (service *userQuizService) getUserQuiz(userQuizAttemptID uuid.UUID) (*models.UserQuizAttempts, error) {
	for _, attempts := range service.db.UserQuizAttempts {
		if attempts.ID == userQuizAttemptID {
			return &attempts, nil
		}
	}
	return nil, errors.New("please start quiz before submitting answers")
}

// isQuestionAnswered will check if all question has been answered for a given user, if yes then it will return an error
func (service *userQuizService) isQuestionAnswered(userQuiz *models.UserQuizAttempts, questionID uuid.UUID) error {
	for _, repsonse := range userQuiz.UserResponses {
		if repsonse.QuestionID == questionID {
			return errors.New("question already answered")
		}
	}

	return nil
}

// isQuizCompleted will check if quiz has ended or max time is exceeded.
func (service *userQuizService) isQuizCompleted(userResponse *models.UserResponse) error {
	userQuiz, err := service.getUserQuiz(userResponse.UserQuizAttemptID)
	if err != nil {
		return err
	}

	if userQuiz.EndedAt != nil {
		return errors.New("cannot answer questions after quiz has ended")
	}

	quiz, err := service.getQuizByID(userQuiz.QuizID)
	if err != nil {
		return err
	}

	if time.Since(*userQuiz.StartedAt) > time.Duration(quiz.MaxTime*uint64(time.Minute)) {
		return errors.New("maximum time exceeded for this quiz")
	}

	return nil
}
