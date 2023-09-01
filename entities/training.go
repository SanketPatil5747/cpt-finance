package entities

type Training struct {
	ID string `json:"id"`

	QuestionID string   `json:"QuestionID"`
	Question   Question `json:"Question" gorm:"-"`

	JobBandID int     `json:"JobBandID"`
	JobBand   JobBand `json:"JobBand" gorm:"-"`

	TrainingName string `json:"TrainingName"`
	TrainingLink string `json:"TrainingLink"`
}
