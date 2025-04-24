package handlers

import (
	"fmt"
	"log"
	"os"

	"api/api/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func checkError(err error) {
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}

func InitDB() {
	err := godotenv.Load()
	checkError(err)
	err = godotenv.Load("docker/.env")
	checkError(err)

	dbUrl := os.Getenv("HORSE_DATABASE_URL")
	dbName := os.Getenv("HORSE_DATABASE_NAME")
	dbUser := os.Getenv("HORSE_DATABASE_USER")
	dbPassword := os.Getenv("HORSE_DATABASE_PASS")
	dbPort := os.Getenv("HORSE_DATABASE_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbUrl,
		dbUser,
		dbPassword,
		dbName,
		dbPort,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	checkError(err)

	if err := DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Horse{},
	); err != nil {
		log.Fatal("Failed to migrate schema:", err)
	}
}
