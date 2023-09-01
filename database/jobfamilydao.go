package database

import (
	"errors"

	"financial-framework/entities"
)

func GetJobFamilyById(id string) (*entities.JobFamily, error) {
	if CheckIfJobFamilyExists(id) == false {
		return nil, errors.New("JobFamily not found")
	}
	var JobFamily entities.JobFamily
	Instance.Debug().First(&JobFamily, "id = ?", id)
	return &JobFamily, nil
}

func GetJobFamilyies() []entities.JobFamily {
	var JobFamilyies []entities.JobFamily
	Instance.Debug().Find(&JobFamilyies).Order("display_order ASC")
	return JobFamilyies
}

func AddAllJobFamilies(jf []entities.JobFamily) error {
	Instance.Debug().CreateInBatches(jf, 100)
	return nil
}

func DeleteAllJobFamilies() error {
	Instance.Debug().Exec("DELETE FROM job_families")
	return nil
}

func CheckIfJobFamilyExists(JobFamilyId string) bool {
	var JobFamily entities.JobFamily
	Instance.Debug().First(&JobFamily, "id = ?", JobFamilyId)
	if JobFamily.ID == "" {
		return false
	}
	return true
}
