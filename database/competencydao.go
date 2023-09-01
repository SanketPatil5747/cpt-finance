package database

import (
	"errors"
	"strconv"

	"financial-framework/entities"
)

func GetCompetencyForFamilyAndBandAndQuestion(family entities.JobFamily, band entities.JobBand, question entities.Question) (*entities.Competency, error) {
	if CheckIfJobFamilyExists(family.ID) == false {
		return nil, errors.New("JobFamily not found")
	}
	if CheckIfJobBandExists(band.ID) == false {
		return nil, errors.New("JobBand not found")
	}
	if CheckIfQuestionExists(question.ID) == false {
		return nil, errors.New("Question not found")
	}
	var Competency entities.Competency

	Instance.Debug().Where("job_band_id = ? AND job_family_id = ? AND question_id = ?", strconv.Itoa(band.ID), family.ID, question.ID).First(&Competency)

	return &Competency, nil
}

func GetCompetencys() []entities.Competency {
	var Competencys []entities.Competency
	Instance.Debug().Find(&Competencys)
	return Competencys
}

func AddAllCompetencies(c []entities.Competency) error {
	Instance.CreateInBatches(c, 100)
	return nil
}

func DeleteAllCompetencies() error {
	Instance.Debug().Exec("DELETE FROM competencies")
	return nil
}

func CheckIfCompetencyExists(CompetencyId string) bool {
	var Competency entities.Competency
	Instance.Debug().First(&Competency, "id = ?", CompetencyId)
	if Competency.ID == "" {
		return false
	}
	return true
}
