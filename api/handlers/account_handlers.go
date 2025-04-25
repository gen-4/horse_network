package handlers

import (
	"net/http"

	"api/api/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateHorse(context *gin.Context) {
	var horse models.Horse

	//bind the request body
	if err := context.ShouldBindJSON(&horse); err != nil {
		models.ResponseJSON(context, http.StatusBadRequest, "Invalid input: "+err.Error(), nil)
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
	models.ResponseJSON(context, http.StatusCreated, "Horse created successfully", models.ToHorseDto(horse))
}

func GetHorses(context *gin.Context) {
	userId := context.Keys["user_id"].(uint)
	var horses []models.Horse
	DB.Where("owner = ?", userId).Find(&horses)
	models.ResponseJSON(context, http.StatusOK, "Horses retrieved successfully", models.ToHorseDtos(horses))
}

func UpdateUser(context *gin.Context) {
	userId := context.Keys["user_id"].(uint)
	var user models.User
	var updateUser models.UpdateUser
	if err := context.ShouldBindJSON(&updateUser); err != nil {
		models.ResponseJSON(context, http.StatusBadRequest, "Invalid input: "+err.Error(), nil)
		return
	}
	user = models.UpdateUserToUser(updateUser)
	if user.Password != "" {
		if encryptedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), HASH_COST); err != nil {
			models.ResponseJSON(context, http.StatusInternalServerError, "Unable to hash password", nil)
			return
		} else {
			user.Password = string(encryptedPass)
		}
	}
	user.ID = userId

	if err := DB.Model(&user).Updates(user).Error; err != nil {
		models.ResponseJSON(context, http.StatusInternalServerError, "Unable to update user", nil)
		return
	}
	DB.First(&user, user.ID)
	models.ResponseJSON(context, http.StatusOK, "User updated successfully", models.ToUserDto(user))
}

func DeleteHorse(context *gin.Context) {
	userId := context.Keys["user_id"].(uint)
	horseId := context.Param("id")
	var horse models.Horse

	if err := DB.Where("id = ? AND owner = ?", horseId, userId).First(&horse).Error; err != nil {
		models.ResponseJSON(context, http.StatusNotFound, "Horse not found", nil)
		return
	}
	DB.Delete(&horse)
	models.ResponseJSON(context, http.StatusOK, "Horse deleted successfully", nil)
}
