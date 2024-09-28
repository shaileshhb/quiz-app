package models

// Question will contain question details for a quiz
type Question struct {
	ID      uint64   `json:"id"`
	Text    string   `json:"text"`
	Options []Option `json:"options"`
}

// Validate will validate if all fields for a question are correctly specified.
func (q *Question) Validate() error {
	return nil
}
