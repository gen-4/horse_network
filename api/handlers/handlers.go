package handlers

import (
	"net/http"

	"api/api/models"

	"github.com/gin-gonic/gin"
)

func CreateHorse(context *gin.Context) {
	var horse models.Horse

	//bind the request body
	if err := context.ShouldBindJSON(&horse); err != nil {
		models.ResponseJSON(context, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	DB.Create(&horse)
	models.ResponseJSON(context, http.StatusCreated, "Horse created successfully", horse)
}

func GetHorses(context *gin.Context) {
	var horses []models.Horse
	DB.Find(&horses)
	models.ResponseJSON(context, http.StatusOK, "Horses retrieved successfully", horses)
}
