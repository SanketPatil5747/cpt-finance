package database

import (
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"financial-framework/entities"
	"financial-framework/middleware"
)

type AWSSecretData struct {
	Host     string `json:"db_host"`
	Port     string `json:"db_port"`
	DBName   string `json:"db_name"`
	Username string `json:"db_username"`
	Password string `json:"db_password"`
}

var (
	Instance *gorm.DB
	err      error
	secret   AWSSecretData
)

func Connect(connectionString string) {
	Instance, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Println(err)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to Database...")
}

func Migrate() {
	err := Instance.AutoMigrate(&entities.JobFamily{}, &entities.JobBand{}, &entities.Question{}, &entities.Competency{}, &entities.Training{}, &entities.JobFunction{})
	if err != nil {
		log.Println(err)
		panic("Database Migration Failed")
	}
	log.Println("Database Migration Completed...")
}

func ConnectToAWSRDS(aws_secretname string) {
	//Gtting secrete from AWS Secret Vault
	secretString := middleware.LoadSecretFromAWS(aws_secretname)

	//Unmarshaling secret to Secret structure
	jsonerr := json.Unmarshal([]byte(secretString), &secret)
	if jsonerr != nil {
		log.Println(jsonerr)
		panic("error unmarshaling a Secret String")
	}

	// Forming connecrion string
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", secret.Host, secret.Port, secret.Username, secret.Password, secret.DBName)

	//DB connection
	Instance, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Println(err)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to Database...")
}
