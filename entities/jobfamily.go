package entities

type JobFamily struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	DisplayOrder  int         `json:"DisplayOrder"`
	JobFunctionID string      `json:"JobFunctionID"`
	JobFunction   JobFunction `json:"JobFunction" gorm:"-"`
}
