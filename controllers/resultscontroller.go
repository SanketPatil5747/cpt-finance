package controllers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"financial-framework/database"
	"financial-framework/entities"
)

type Result struct {
	Question       entities.Question
	Competency     entities.Competency
	Answer         uint
	Gap            bool
	AnswerText     string
	CompetencyText string
	CompetencyNum  uint
	GapSize        int
	Trainings      []entities.Training
}

func ResultsPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}

	JobBand := r.Form["band_id"]
	JobFamily := r.Form["family_id"]
	family, err := database.GetJobFamilyById(JobFamily[0])

	i, err := strconv.Atoi(JobBand[0])
	if err != nil {
		// ... handle error
		panic(err)
	}
	band, err := database.GetJobBandById(i)
	var results []Result
	for key, value := range r.Form {
		if database.CheckIfQuestionExists(key) {
			var result Result
			question, err := database.GetQuestionById(key)
			if err != nil {
				return
			}

			result.Question = *question
			competency, err := database.GetCompetencyForFamilyAndBandAndQuestion(*family, *band, *question)
			if err != nil {
				return
			}

			result.Competency = *competency
			result.Answer = answerToInt(value[0])
			result.AnswerText = answerToString(answerToInt(value[0]))
			result.CompetencyNum = answerToInt(competency.Level)
			result.CompetencyText = competency.DesiredProficiency

			result.Gap = isThereAGap(competency.Level, value[0])
			if result.Gap {
				result.GapSize = int(result.CompetencyNum) - int(result.Answer)
				training, err := database.GetTrainingForTrainingGap(key, int(result.Answer), int(answerToInt(result.Competency.Level)))
				if err != nil {
					return
				}

				result.Trainings = *training
			}
			results = append(results, result)
		}
	}
	m := map[string]interface{}{
		"family":  family,
		"band":    band,
		"results": results,
	}
	err = templates.ExecuteTemplate(w, "results.html", m)
	if err != nil {
		log.Println(err)
	}
}

func isThereAGap(competencyLevel string, answer string) bool {
	if answerToInt(answer) < answerToInt(competencyLevel) {
		return true
	}
	return false
}

func answerToString(answer uint) string {
	switch answer {
	case 1:
		return "Very little or no knowledge"
	case 2:
		return "Beginner"
	case 3:
		return "Able to work unsupervised"
	case 4:
		return "Advanced knowledge"
	case 5:
		return "Expert and capable of training others"
	default:
		return ""
	}
}

func answerToInt(answer string) uint {
	switch strings.ToLower(answer) {
	case "bf":
		return 1
	case "f":
		return 2
	case "i":
		return 3
	case "a":
		return 4
	case "e":
		return 5
	default:
		return 0
	}
}
