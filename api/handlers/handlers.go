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

	userId := context.Keys["user_id"].(uint)
	horse.Owner = &userId

	{
		execution := DB.Create(&horse)
		if execution.Error != nil {
			models.ResponseJSON(context, http.StatusInternalServerError, "Unable to create the horse", nil)
			return
		}
	}
	models.ResponseJSON(context, http.StatusCreated, "Horse created successfully", horse)
}

func GetHorses(context *gin.Context) {
	userId := context.Keys["user_id"].(uint)
	var horses []models.Horse
	DB.Where("owner = ?", userId).Find(&horses)
	models.ResponseJSON(context, http.StatusOK, "Horses retrieved successfully", horses)
}
