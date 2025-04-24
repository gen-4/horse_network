package tests

import (
	"api/api/handlers"
	"api/api/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() {
	var err error
	handlers.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect test database")
	}
	handlers.DB.AutoMigrate(&models.Horse{})
}

func AddHorse() models.Horse {
	horse := models.Horse{}
	handlers.DB.Create(&horse)
	return horse
}
