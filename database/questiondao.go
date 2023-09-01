package database

import (
	"errors"
	"strconv"
	"strings"

	"financial-framework/entities"
)

func GetQuestionsForFamilyAndBand(family entities.JobFamily, band entities.JobBand) ([]entities.Question, error) {
	if CheckIfJobFamilyExists(family.ID) == false {
		return nil, errors.New("JobFamily not found")
	}
	if CheckIfJobBandExists(band.ID) == false {
		return nil, errors.New("JobBand not found")
	}
	var questionLinks []entities.Competency
	var questions []entities.Question
	Instance.Debug().Where("job_band_id = ? AND job_family_id = ? and lower(level) != 'x'", strconv.Itoa(band.ID), family.ID).Find(&questionLinks)

	for _, question := range questionLinks {
		question, err := GetQuestionById(question.QuestionID)
		question.TempCriteria = splitNewlines(question.Criteria)
		if err != nil {
			return nil, err
		}
		questions = append(questions, *question)
	}

	return questions, nil
}

func splitNewlines(s string) []string {
	return strings.Split(s, "\\n")
}

func GetQuestionById(id string) (*entities.Question, error) {
	if CheckIfQuestionExists(id) == false {
		return nil, errors.New("question not found")
	}
	var Question entities.Question
	Instance.Debug().First(&Question, "id = ?", id)
	return &Question, nil
}

func GetQuestions() []entities.Question {
	var Questions []entities.Question
	Instance.Debug().Find(&Questions)
	return Questions
}

func AddAllQuestions(q []entities.Question) error {
	Instance.Debug().CreateInBatches(q, 100)
	return nil
}

func DeleteAllQuestions() error {
	Instance.Debug().Exec("DELETE FROM questions")
	return nil
}

func CheckIfQuestionExists(QuestionId string) bool {
	var Question entities.Question
	Instance.Debug().First(&Question, "id = ?", QuestionId)
	if Question.ID == "" {
		return false
	}
	return true
}
