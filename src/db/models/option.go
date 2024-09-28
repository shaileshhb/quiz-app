package models

// Option will contain options for specific question.
type Option struct {
	ID         uint64 `json:"id"`
	QuestionID uint64 `json:"questionID"`
	Answer     string `json:"answer"`
	IsCorrect  bool   `json:"isCorrect"`
}

// Validate will validate if all fields for a option are correctly specified.
func (o *Option) Validate() error {
	return nil
}
