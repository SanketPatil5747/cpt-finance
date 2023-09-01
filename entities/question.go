package entities

type Question struct {
	ID           string   `json:"id"`
	Description  string   `json:"Description"`
	Criteria     string   `json:"Criteria"`
	Category     string   `json:"Category"`
	Capability   string   `json:"Capability"`
	TempCriteria []string `gorm:"-"`
}
