package tests

import (
	"api/api"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() {
	var err error
	api.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect test database")
	}
	api.DB.AutoMigrate(&api.Horse{})
}

func AddHorse() api.Horse {
	horse := api.Horse{}
	api.DB.Create(&horse)
	return horse
}
