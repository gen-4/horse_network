package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const LOG_FILE = "ginlogs.log"

var fileDescriptor *os.File

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
		slog.Warn("Unable to read .env file")
	} else {
		environment = os.Getenv("ENVIRONMENT")
	}

	return environment
}

func Config() {
	var environment string = getEnv()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	switch environment {
	case "dev":
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		slog.SetDefault(logger)
		gin.SetMode(gin.DebugMode)

	case "test":
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		slog.SetDefault(logger)
		gin.SetMode(gin.TestMode)

	case "pro":
		f, err := os.OpenFile(LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			slog.Error("Error opening file: %v", err)
		}
		defer f.Close()
		logger := slog.New(slog.NewJSONHandler(f, nil))
		slog.SetDefault(logger)
		gin.SetMode(gin.ReleaseMode)
	}

	slog.Info(fmt.Sprintf("Running in %s environment", environment))
}

func CloseConfig() {
	if fileDescriptor != nil {
		fileDescriptor.Close()
	}
}
