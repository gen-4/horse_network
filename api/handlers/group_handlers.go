package handlers

import (
	"api/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateGroup(context *gin.Context) {
	userId := context.Keys["user_id"].(uint)
	var group models.Group
	var user models.User

	if err := context.ShouldBindJSON(&group); err != nil {
		models.ResponseJSON(context, http.StatusBadRequest, "Invalid input: "+err.Error(), nil)
		return
	}

	if err := DB.First(&user, userId).Error; err != nil {
		models.ResponseJSON(context, http.StatusNotFound, "User not found", nil)
		return
	}
	group.Users = append(group.Users, user)

	if err := DB.Create(&group).Error; err != nil {
		models.ResponseJSON(context, http.StatusInternalServerError, "Error creating group", nil)
		return
	}

	models.ResponseJSON(context, http.StatusCreated, "Group creation successful", nil)
}
