package database

import (
	"financial-framework/entities"
)

func GetJobFunctions() []entities.JobFunction {
	var JobFunctions []entities.JobFunction
	Instance.Debug().Find(&JobFunctions).Order("display_order ASC")
	return JobFunctions
}

func AddAllJobFunctions(jf []entities.JobFunction) error {
	Instance.Debug().CreateInBatches(jf, 100)
	return nil
}

func DeleteAllJobFunctions() error {
	Instance.Debug().Exec("DELETE FROM job_functions")
	return nil
}

func CheckIfJobFunctionExists(JobFunctionId string) bool {
	var JobFunction entities.JobFunction
	Instance.Debug().First(&JobFunction, "id = ?", JobFunctionId)
	if JobFunction.ID == "" {
		return false
	}
	return true
}
