package controllers

import (
	"bytes"
	"encoding/csv"
	"log"
	"net/http"
	"strconv"

	"financial-framework/database"
)

func DownloadTableAsCSV(w http.ResponseWriter, r *http.Request) {
	f := r.URL.Query().Get("f")
	if f == "" {
		log.Println("missing param")
	}

	//check "f" is one of the actual tables
	if IsValidFile(f) == "" {
		log.Println("not a valid file")
	}

	//get the data

	//write the data
	b := &bytes.Buffer{}
	wr := csv.NewWriter(b)

	switch IsValidFile(f) {
	case "job_bands":
		writeJobBands(wr)
	case "job_families":
		writeJobFamilies(wr)
	case "job_function":
		writeJobFunction(wr)
	case "questions":
		writeQuestions(wr)
	case "competencies":
		writeCompetencies(wr)
	case "training":
		writeTraining(wr)
	}

	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Disposition", "attachment; filename="+IsValidFile(f)+".csv")
	w.Header().Set("Content-Type", "text/csv")
	//w.Header().Set("Content-Length", string(rune(len(b.Bytes()))))

	_, err := w.Write(b.Bytes())
	if err != nil {
		log.Fatalln("error writing bytes to response:", err)
	}
}

func writeTraining(wr *csv.Writer) {
	if err := wr.Write([]string{"id", "question_id", "job_band_id", "training_name", "training_link"}); err != nil {
		log.Fatalln("error writing header record to csv:", err)
	}

	for _, training := range database.GetTrainings() {
		var record []string
		record = append(record, training.ID)
		record = append(record, training.QuestionID)
		record = append(record, strconv.FormatInt(int64(training.JobBandID), 10))
		record = append(record, training.TrainingName)
		record = append(record, training.TrainingLink)
		if err := wr.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	wr.Flush()
	if err := wr.Error(); err != nil {
		log.Fatal(err)
	}
}

func writeCompetencies(wr *csv.Writer) {
	if err := wr.Write([]string{"id", "level", "job_band_id", "job_family_id", "question_id"}); err != nil {
		log.Fatalln("error writing header record to csv:", err)
	}

	for _, competency := range database.GetCompetencys() {
		var record []string
		record = append(record, competency.ID)
		record = append(record, competency.Level)
		record = append(record, competency.JobBandID)
		record = append(record, competency.JobFamilyID)
		record = append(record, competency.QuestionID)
		if err := wr.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	wr.Flush()
	if err := wr.Error(); err != nil {
		log.Fatal(err)
	}
}

func writeQuestions(wr *csv.Writer) {
	if err := wr.Write([]string{"id", "description", "criteria", "category", "capability"}); err != nil {
		log.Fatalln("error writing header record to csv:", err)
	}

	for _, question := range database.GetQuestions() {
		var record []string
		record = append(record, question.ID)
		record = append(record, question.Description)
		record = append(record, question.Criteria)
		record = append(record, question.Category)
		record = append(record, question.Capability)
		if err := wr.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	wr.Flush()
	if err := wr.Error(); err != nil {
		log.Fatal(err)
	}
}

func writeJobFunction(wr *csv.Writer) {
	if err := wr.Write([]string{"id", "name", "display_order"}); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	for _, function := range database.GetJobFunctions() {
		var record []string
		record = append(record, function.ID)
		record = append(record, function.Name)
		record = append(record, strconv.FormatInt(int64(function.DisplayOrder), 10))
		if err := wr.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	wr.Flush()
	if err := wr.Error(); err != nil {
		log.Fatal(err)
	}
}

func writeJobFamilies(wr *csv.Writer) {
	if err := wr.Write([]string{"id", "name", "display_order", "job_function_id"}); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	for _, families := range database.GetJobFamilyies() {
		var record []string
		record = append(record, families.ID)
		record = append(record, families.Name)
		record = append(record, strconv.FormatInt(int64(families.DisplayOrder), 10))
		record = append(record, families.JobFunctionID)
		if err := wr.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	wr.Flush()
	if err := wr.Error(); err != nil {
		log.Fatal(err)
	}
}

func writeJobBands(wr *csv.Writer) {
	if err := wr.Write([]string{"id", "name", "display_order"}); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	for _, band := range database.GetJobBands() {
		var record []string
		record = append(record, strconv.FormatInt(int64(band.ID), 10))
		record = append(record, band.Name)
		record = append(record, strconv.FormatInt(int64(band.DisplayOrder), 10))
		if err := wr.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	wr.Flush()
	if err := wr.Error(); err != nil {
		log.Fatal(err)
	}
}

func IsValidFile(file string) string {
	switch file {
	case "job_bands":
		return "job_bands"
	case "job_families":
		return "job_families"
	case "job_function":
		return "job_function"
	case "questions":
		return "questions"
	case "competencies":
		return "competencies"
	case "training":
		return "training"
	}
	return ""
}
