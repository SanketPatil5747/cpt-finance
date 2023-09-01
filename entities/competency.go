package entities

type Competency struct {
	ID                 string    `json:"id"`
	Level              string    `json:"level"`
	JobBandID          string    `json:"JobBandID"`
	JobBand            JobBand   `json:"JobBand" gorm:"-"`
	JobFamilyID        string    `json:"JobFamilyID" `
	JobFamily          JobFamily `json:"JobFamily" gorm:"-"`
	QuestionID         string    `json:"QuestionID"`
	Question           Question  `json:"Question" gorm:"-"`
	DesiredProficiency string    `json:"DesiredProficiency"`
}
