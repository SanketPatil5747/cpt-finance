package database

import (
	"financial-framework/entities"
)

func GetTrainingForTrainingGap(QuestionID string, Answer int, Level int) (*[]entities.Training, error) {
	var Trainings []entities.Training
	Instance.Debug().Where("question_id = ? AND job_band_id > ? AND job_band_id <= ? AND training_name != '' ", QuestionID, Answer, Level).Find(&Trainings)
	return &Trainings, nil
}

func GetTrainings() []entities.Training {
	var Trainings []entities.Training
	Instance.Debug().Find(&Trainings)
	return Trainings
}

func AddAllTrainings(t []entities.Training) error {
	Instance.Debug().CreateInBatches(t, 100)
	return nil
}

func DeleteAllTrainings() error {
	Instance.Debug().Exec("DELETE FROM trainings")
	return nil
}

func CheckIfTrainingExists(TrainingId string) bool {
	var Training entities.Training
	Instance.Debug().First(&Training, "id = ?", TrainingId)
	if Training.ID == "" {
		return false
	}
	return true
}
