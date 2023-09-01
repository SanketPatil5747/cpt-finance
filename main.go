package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"financial-framework/controllers"
	"financial-framework/database"
)

var DB *gorm.DB

var templates *template.Template

func main() {
	// Load Configurations from config.json using Viper
	LoadAppConfig()

	// Initialize Database
	if AppConfig.DevelopmentEnv == "local" {
		//Passing DB ConnectionString to connect local DB
		database.Connect(AppConfig.ConnectionString)
		database.Migrate()
	} else {
		//Getting DB details from cloud
		database.ConnectToAWSRDS(AppConfig.AWSSecretname)
		database.Migrate()
	}

	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)
	templates = template.Must(templates.ParseGlob("templates/*.html"))
	controllers.Init(templates)

	// Register Routes
	RegisterRoutes(router)

	// Start the server
	log.Println(fmt.Sprintf("Starting Server on port %s...", AppConfig.Port))

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("Running on Hostname: %s\n", hostname)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", AppConfig.Port), router))
}

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/", controllers.IndexPage).Methods("GET")
	router.HandleFunc("/", controllers.IndexPage).Methods("POST")

	router.HandleFunc("/admin", controllers.AdminPage).Methods("GET")
	router.HandleFunc("/admin", controllers.AdminPage).Methods("POST")
	router.HandleFunc("/admin/download", controllers.DownloadTableAsCSV).Methods("GET")
	router.HandleFunc("/admin/upload", controllers.ProcessCSVUploads).Methods("POST")

	router.HandleFunc("/questions", controllers.QuestionsPage).Methods("POST")
	router.HandleFunc("/questions", controllers.QuestionsPage).Methods("GET")

	router.HandleFunc("/results", controllers.ResultsPage).Methods("POST")
	router.HandleFunc("/results", controllers.ResultsPage).Methods("GET")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	http.Handle("/", router)
}
