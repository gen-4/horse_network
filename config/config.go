package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func getEnv() string {
	var err error
	environment := "dev"

	if flag.Lookup("test.v") == nil {
		err = godotenv.Load()
	} else {
		envFileError := godotenv.Load(".test.env")
		if envFileError != nil {
			err = godotenv.Load("../.test.env")
		}
	}

	if err != nil {
		log.Println("WARNING: Unable to read .env file")
	} else {
		environment = os.Getenv("ENVIRONMENT")
	}

	return environment
}

func Config() {
	var environment string = getEnv()
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	switch environment {
	case "dev":
		gin.SetMode(gin.DebugMode)

	case "test":
		gin.SetMode(gin.TestMode)

	case "pro":
		f, err := os.OpenFile("ginlogs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
		gin.SetMode(gin.ReleaseMode)
	}

	log.Println(fmt.Sprintf("INFO: Running in %s environment", environment))
}
