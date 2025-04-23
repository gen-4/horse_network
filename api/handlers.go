package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

	db_url := os.Getenv("HORSE_DATABASE_URL")
	db_name := os.Getenv("HORSE_DATABASE_NAME")
	db_user := os.Getenv("HORSE_DATABASE_USER")
	db_password := os.Getenv("HORSE_DATABASE_PASS")
	db_port := os.Getenv("HORSE_DATABASE_PORT")

	dsn := "host=" + db_url + " user=" +
		db_user + " password=" + db_password + " dbname=" +
		db_name + " port=" + db_port + " sslmode=disable"

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	checkError(err)

	// migrate the schema
	if err := DB.AutoMigrate(&Horse{}); err != nil {
		log.Fatal("Failed to migrate schema:", err)
	}
}

func CreateHorse(context *gin.Context) {
	var horse Horse

	//bind the request body
	if err := context.ShouldBindJSON(&horse); err != nil {
		ResponseJSON(context, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	DB.Create(&horse)
	ResponseJSON(context, http.StatusCreated, "Horse created successfully", horse)
}

func GetHorses(context *gin.Context) {
	var horses []Horse
	DB.Find(&horses)
	ResponseJSON(context, http.StatusOK, "Horses retrieved successfully", horses)
}
