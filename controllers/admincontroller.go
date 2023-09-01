package controllers

import (
	"fmt"
	"log"
	"net/http"
)

const MaxUploadSize = 1024 * 1024 // 1MB

func AdminPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}

	data := make(map[string]interface{})
	err = templates.ExecuteTemplate(w, "admin.html", data)
	if err != nil {
		return
	}
}

func logLine(fieldName string, message string) string {
	log.Println(fmt.Sprintf("%s : %s", fieldName, message))
	return fmt.Sprintf("%s : %s \n", fieldName, message)
}
