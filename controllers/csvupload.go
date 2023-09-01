package controllers

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"financial-framework/database"
	"financial-framework/entities"
)

func ProcessCSVUploads(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)

	data := make(map[string]interface{})

	logLines := make([]string, 0)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
		return
	}

	logLines = append(logLines, processCSVUpload(r, "job_bands")...)
	logLines = append(logLines, processCSVUpload(r, "job_families")...)
	logLines = append(logLines, processCSVUpload(r, "job_function")...)
	logLines = append(logLines, processCSVUpload(r, "questions")...)
	logLines = append(logLines, processCSVUpload(r, "competencies")...)
	logLines = append(logLines, processCSVUpload(r, "training")...)

	data["logLines"] = logLines

	err := templates.ExecuteTemplate(w, "uploadResults.html", data)
	if err != nil {
		return
	}
}

func processCSVUpload(r *http.Request, fieldName string) []string {
	var logLines []string
	file, _, err := r.FormFile(fieldName)
	if err != nil && err == http.ErrMissingFile {
		logLines = append(logLines, logLine(fieldName, "file not provided"))

		return logLines
	} else if err != nil {
		logLines = append(logLines, logLine(fieldName, "Error Retrieving the File"))
		log.Println(err)
		return logLines
	}
	defer file.Close()

	if !validateFileType(file) {
		logLines = append(logLines, logLine(fieldName, "bad file type"))
		return logLines
	} else {
		logLines = append(logLines, logLine(fieldName, "valid file type"))

	}

	reader := csv.NewReader(file)

	// skip first line
	if _, err := reader.Read(); err != nil {
		return logLines
	}

	records, err := reader.ReadAll()
	if err != nil {
		log.Println(err)
		logLines = append(logLines, logLine(fieldName, "Error"))
		return logLines
	}
	switch fieldName {
	case "job_bands":
		logLines = append(logLines, loadCSVJobBands(records)...)
	case "job_families":
		logLines = append(logLines, loadCSVJobFamilies(records)...)
	case "job_function":
		logLines = append(logLines, loadCSVJobFunction(records)...)
	case "questions":
		logLines = append(logLines, loadCSVQuestions(records)...)
	case "competencies":
		logLines = append(logLines, loadCSVJobCompetencies(records)...)
	case "training":
		logLines = append(logLines, loadCSVJobTraining(records)...)

	default:
		return logLines
	}

	return logLines
}

func loadCSVJobTraining(rec [][]string) []string {
	var logLines []string
	var trainings []entities.Training
	var failed = false

	for num, row := range rec {
		var training entities.Training

		if len(row) < 5 {
			logLines = append(logLines, logLine("training", "row:"+strconv.Itoa(num)+" is missing a column"))
			failed = true
		}

		intVar2, err := strconv.Atoi(row[2])
		if err != nil {
			logLines = append(logLines, logLine("job_function", "could not convert the value on row:"+strconv.Itoa(num)+"("+row[2]+") to a string"))
			failed = true
			continue
		}
		if !database.CheckIfQuestionExists(row[1]) {
			logLines = append(logLines, logLine("job_function", "WARNING - row ("+row[0]+") could not find matching questionID("+row[1]+")"))
		}
		if !database.CheckIfJobBandExists(intVar2) {
			logLines = append(logLines, logLine("job_function", "WARNING - row ("+row[0]+") could not find matching Job Band ID("+row[2]+")"))
		}

		training.ID = row[0]
		training.QuestionID = row[1]
		training.JobBandID = intVar2
		training.TrainingName = row[3]
		training.TrainingLink = row[4]

		trainings = append(trainings, training)
	}
	if !failed {
		err := database.DeleteAllTrainings()
		if err != nil {
			logLines = append(logLines, logLine("training", "error clearing training table"))
			return logLines
		} else {
			logLines = append(logLines, logLine("training", "successfully cleared training table"))
		}
		err = database.AddAllTrainings(trainings)
		if err != nil {
			logLines = append(logLines, logLine("training", "error adding training to table"))
			return logLines
		} else {
			logLines = append(logLines, logLine("training", "successfully added all("+strconv.Itoa(len(trainings))+") training to table"))
		}
	} else {
		logLines = append(logLines, logLine("training", "didnt try to add rows to the db as data validation failed."))
	}
	return logLines
}

func loadCSVJobCompetencies(rec [][]string) []string {
	var logLines []string
	var competencies []entities.Competency
	var failed = false

	for num, row := range rec {
		var competency entities.Competency

		intVar2, err := strconv.Atoi(row[2])
		if err != nil {
			logLines = append(logLines, logLine("job_function", "could not convert the value on row:"+strconv.Itoa(num)+"("+row[2]+") to a string"))
			failed = true
			continue
		}

		if len(row) < 6 {
			logLines = append(logLines, logLine("competencies", "row:"+strconv.Itoa(num)+" is missing a column"))
			failed = true
		}

		if !database.CheckIfJobBandExists(intVar2) {
			logLines = append(logLines, logLine("job_function", "WARNING - row ("+row[0]+") could not find matching Job Band ID("+row[2]+")"))
		}

		if !database.CheckIfJobFamilyExists(row[3]) {
			logLines = append(logLines, logLine("job_function", "WARNING - row ("+row[0]+") could not find matching Job Family ID("+row[3]+")"))
		}

		if !database.CheckIfQuestionExists(row[4]) {
			logLines = append(logLines, logLine("job_function", "WARNING - row ("+row[0]+") could not find matching Question ID("+row[4]+")"))
		}

		competency.ID = row[0]
		competency.Level = row[1]
		competency.JobBandID = row[2]
		competency.JobFamilyID = row[3]
		competency.QuestionID = row[4]
		competency.DesiredProficiency = row[5]

		competencies = append(competencies, competency)
	}
	if !failed {
		err := database.DeleteAllCompetencies()
		if err != nil {
			logLines = append(logLines, logLine("competencies", "error clearing competencies table"))
			return logLines
		} else {
			logLines = append(logLines, logLine("competencies", "successfully cleared competencies table"))
		}
		err = database.AddAllCompetencies(competencies)
		if err != nil {
			logLines = append(logLines, logLine("competencies", "error adding competencies to table"))
			return logLines
		} else {
			logLines = append(logLines, logLine("competencies", "successfully added all("+strconv.Itoa(len(competencies))+") competencies to table"))
		}
	} else {
		logLines = append(logLines, logLine("competencies", "didnt try to add rows to the db as data validation failed."))
	}
	return logLines
}

func loadCSVQuestions(rec [][]string) []string {
	var logLines []string
	var questions []entities.Question
	var failed = false

	for num, row := range rec {
		var question entities.Question

		if len(row) < 5 {
			logLines = append(logLines, logLine("questions", "row:"+strconv.Itoa(num)+" is missing a column"))
			failed = true
		}

		question.ID = row[0]
		question.Description = row[1]
		question.Criteria = row[2]
		question.Category = row[3]
		question.Capability = row[4]

		questions = append(questions, question)
	}

	if !failed {
		err := database.DeleteAllQuestions()
		if err != nil {
			logLines = append(logLines, logLine("questions", "error clearing questions table"))
			return logLines
		} else {
			logLines = append(logLines, logLine("questions", "successfully cleared questions table"))
		}

		err = database.AddAllQuestions(questions)
		if err != nil {
			logLines = append(logLines, logLine("questions", "error adding questions to table"))
			return logLines
		} else {
			logLines = append(logLines, logLine("questions", "successfully added all("+strconv.Itoa(len(questions))+") questions to table"))
		}
	} else {
		logLines = append(logLines, logLine("questions", "didnt try to add rows to the db as data validation failed."))
	}
	return logLines
}

func loadCSVJobFunction(rec [][]string) []string {
	var logLines []string
	var jobFunctions []entities.JobFunction
	var failed = false

	for num, row := range rec {
		var jobFunction entities.JobFunction

		if len(row) < 3 {
			logLines = append(logLines, logLine("job_function", "row:"+strconv.Itoa(num)+" is missing a column"))
			failed = true
		}

		intVar2, err := strconv.Atoi(row[2])
		if err != nil {
			logLines = append(logLines, logLine("job_function", "could not convert the value on row:"+strconv.Itoa(num)+"("+row[2]+") to a string"))
			failed = true
			continue
		}

		jobFunction.ID = row[0]
		jobFunction.Name = row[1]
		jobFunction.DisplayOrder = intVar2

		jobFunctions = append(jobFunctions, jobFunction)
	}
	if !failed {
		err := database.DeleteAllJobFunctions()
		if err != nil {
			logLines = append(logLines, logLine("job_function", "error clearing job_function table"))
			return logLines
		} else {
			logLines = append(logLines, logLine("job_function", "successfully cleared job_function table"))
		}
		err = database.AddAllJobFunctions(jobFunctions)
		if err != nil {
			logLines = append(logLines, logLine("job_function", "error adding job_function to table"))
			return logLines
		} else {
			logLines = append(logLines, logLine("job_function", "successfully added all("+strconv.Itoa(len(jobFunctions))+") job_function to table"))
		}
	} else {
		logLines = append(logLines, logLine("questions", "didnt try to add rows to the db as data validation failed."))
	}
	return logLines
}

func loadCSVJobBands(rec [][]string) []string {
	var logLines []string
	var jobBands []entities.JobBand
	var failed = false

	for num, row := range rec {
		var jobBand entities.JobBand

		if len(row) < 3 {
			logLines = append(logLines, logLine("job_bands", "row:"+strconv.Itoa(num)+" is missing a column"))
			failed = true
		}

		intVar1, err := strconv.Atoi(row[0])
		if err != nil {
			logLines = append(logLines, logLine("job_bands", "could not convert the value on row:"+strconv.Itoa(num)+"("+row[0]+") to a string"))
			failed = true
			continue
		}

		intVar2, err := strconv.Atoi(row[2])
		if err != nil {
			logLines = append(logLines, logLine("job_bands", "could not convert the value on row:"+strconv.Itoa(num)+"("+row[2]+") to a string"))
			failed = true
			continue
		}

		jobBand.ID = intVar1
		jobBand.Name = row[1]
		jobBand.DisplayOrder = intVar2

		jobBands = append(jobBands, jobBand)
	}
	if !failed {
		err := database.DeleteAllJobBands()
		if err != nil {
			logLines = append(logLines, logLine("job_bands", "error clearing job_bands table"))
			return logLines
		} else {
			logLines = append(logLines, logLine("job_bands", "successfully cleared job_bands table"))
		}
		err = database.AddAllJobBands(jobBands)
		if err != nil {
			logLines = append(logLines, logLine("job_bands", "error adding job_bands to table"))
			return logLines
		} else {
			logLines = append(logLines, logLine("job_bands", "successfully added all("+strconv.Itoa(len(jobBands))+") job_bands to table"))
		}
	} else {
		logLines = append(logLines, logLine("job_bands", "didnt try to add rows to the db as data validation failed."))
	}
	return logLines
}

func loadCSVJobFamilies(rec [][]string) []string {
	var logLines []string
	var jobFamilies []entities.JobFamily
	var failed = false

	for num, row := range rec {
		var jobFamily entities.JobFamily

		if len(row) < 4 {
			logLines = append(logLines, logLine("job_families", "row:"+strconv.Itoa(num)+" is missing a column"))
			failed = true
		}

		intVar, err := strconv.Atoi(row[2])
		if err != nil {
			logLines = append(logLines, logLine("job_families", "could not convert the value on row:"+strconv.Itoa(num)+"("+row[2]+") to a string"))
			failed = true
			continue
		}

		if !database.CheckIfJobFunctionExists(row[3]) {
			logLines = append(logLines, logLine("job_function", "WARNING - row ("+row[0]+") could not find matching Job Function("+row[3]+")"))
		}

		jobFamily.ID = row[0]
		jobFamily.Name = row[1]
		jobFamily.DisplayOrder = intVar
		jobFamily.JobFunctionID = row[3]

		jobFamilies = append(jobFamilies, jobFamily)
	}
	if !failed {
		err := database.DeleteAllJobFamilies()
		if err != nil {
			logLines = append(logLines, logLine("job_families", "error clearing job_families table"))
			return logLines
		} else {
			logLines = append(logLines, logLine("job_families", "successfully cleared job_families table"))
		}
		err = database.AddAllJobFamilies(jobFamilies)
		if err != nil {
			logLines = append(logLines, logLine("job_families", "error adding job_families to table"))
			return logLines
		} else {
			logLines = append(logLines, logLine("job_families", "successfully added all("+strconv.Itoa(len(jobFamilies))+") job_families to table"))
		}
	} else {
		logLines = append(logLines, logLine("job_bands", "didnt try to add rows to the db as data validation failed."))
	}
	return logLines
}

func validateFileType(file multipart.File) bool {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return false
	}
	detectedFileType := http.DetectContentType(buf.Bytes())
	_, err := file.Seek(0, 0)
	if err != nil {
		return false
	}
	return strings.Contains(detectedFileType, "text/plain")
}
