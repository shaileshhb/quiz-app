package models

// Quiz will contain details related to quiz
type Quiz struct {
	ID        uint64     `json:"id"`
	Title     string     `json:"title"`
	Questions []Question `json:"questions"`
}

// Validate will validate if all fields of quiz are valid.
func (q *Quiz) Validate() error {
	return nil
}
