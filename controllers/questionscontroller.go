package controllers

import (
	"net/http"
	"strconv"

	"financial-framework/database"
)

func QuestionsPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}

	JobBand := r.Form["JobBand"]
	JobFamily := r.Form["JobFamily"]

	family, err := database.GetJobFamilyById(JobFamily[0])

	i, err := strconv.Atoi(JobBand[0])
	if err != nil {
		// ... handle error
		panic(err)
	}
	band, err := database.GetJobBandById(i)

	questions, err := database.GetQuestionsForFamilyAndBand(*family, *band)
	if err != nil {
		return
	}

	m := map[string]interface{}{
		"family":    family,
		"band":      band,
		"questions": questions,
	}

	err = templates.ExecuteTemplate(w, "questions.html", m)
	if err != nil {
		return
	}
}
