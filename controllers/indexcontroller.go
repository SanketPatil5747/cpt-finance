package controllers

import (
	"html/template"
	"net/http"

	"financial-framework/database"
)

var templates *template.Template

func Init(t *template.Template) {
	templates = t
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}

	data := make(map[string]interface{})

	JobFunction := r.Form["JobFunction"]

	if len(JobFunction) != 0 && database.CheckIfJobFunctionExists(JobFunction[0]) {

		familys := database.GetJobFamilyies()
		bands := database.GetJobBands()

		data["families"] = familys
		data["bands"] = bands

		err = templates.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			return
		}
	} else {
		functions := database.GetJobFunctions()

		data["functions"] = functions

		err = templates.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			return
		}
	}
}
