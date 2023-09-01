package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Port             string `mapstructure:"port"`
	ConnectionString string `mapstructure:"connection_string"`
	DevelopmentEnv   string `mapstructure:"development_env"`
	AWSSecretname    string `mapstructure:"aws_secretname"`
}

var AppConfig *Config
var ConnectionString string

func LoadAppConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Getting Application Deployment Env
	developmentEnv := os.Getenv("DEVELOPMENT_ENV")
	if developmentEnv == "" {
		log.Fatal("no DEVELOPMENT_ENV specified")
	}

	//Getting AWS Secrete Name
	awsSecretname := os.Getenv("AWS_SECRETE_NAME")
	if awsSecretname == "" {
		log.Fatal("no AWS_SECRETE_NAME specified")
	}

	//Server Information
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	// Getting DB Info
	postgres_User := os.Getenv("POSTGRES_USER")
	if postgres_User == "" {
		log.Fatal("no POSTGRES_USER specified")
	}

	postgres_Password := os.Getenv("POSTGRES_PASSWORD")
	if postgres_Password == "" {
		log.Fatal("no POSTGRES_PASSWORD specified")
	}

	postgres_DB := os.Getenv("POSTGRES_DB")
	if postgres_DB == "" {
		log.Fatal("no POSTGRES_DB specified")
	}

	postgres_Port := os.Getenv("POSTGRES_PORT")
	if postgres_Port == "" {
		log.Fatal("no POSTGRES_PORT specified")
	}

	postgres_Host := os.Getenv("POSTGRES_HOST")
	if postgres_Host == "" {
		log.Fatal("no POSTGRES_HOST specified")
	}

	//Forming Connection String
	ConnectionString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", postgres_Host, postgres_Port, postgres_User, postgres_Password, postgres_DB)

	if ConnectionString == "" {
		log.Println("Loading Server Configurations...")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("json")
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal(err)
		}
		err = viper.Unmarshal(&AppConfig)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		AppConfig = &Config{
			ConnectionString: ConnectionString,
			Port:             port,
			DevelopmentEnv:   developmentEnv,
			AWSSecretname:    awsSecretname,
		}
	}
}
