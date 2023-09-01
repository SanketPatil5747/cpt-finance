package database

import (
	"errors"

	"financial-framework/entities"
)

func GetJobBandById(id int) (*entities.JobBand, error) {
	if CheckIfJobBandExists(id) == false {
		return nil, errors.New("JobBand not found")
	}
	var JobBand entities.JobBand
	Instance.Debug().First(&JobBand, "id = ?", id)
	return &JobBand, nil
}

func GetJobBands() []entities.JobBand {
	var JobBands []entities.JobBand
	Instance.Debug().Find(&JobBands).Order("display_order ASC")
	return JobBands
}

func AddAllJobBands(jb []entities.JobBand) error {
	Instance.Debug().CreateInBatches(jb, 100)
	return nil
}

func DeleteAllJobBands() error {
	Instance.Debug().Exec("DELETE FROM job_bands")
	return nil
}

func CheckIfJobBandExists(JobBandId int) bool {
	var JobBand entities.JobBand
	Instance.Debug().First(&JobBand, "id = ?", JobBandId)
	if JobBand.ID == 0 {
		return false
	}
	return true
}
