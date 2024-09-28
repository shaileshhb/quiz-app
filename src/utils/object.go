package utils

import "github.com/shaileshhb/quiz/src/db/models"

func CopyQuiz(q models.Quiz) models.Quiz {
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
