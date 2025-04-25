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

func JoinGroup(context *gin.Context) {
	userId := context.Keys["user_id"].(uint)
	groupId := context.Param("id")
	var user models.User
	var group models.Group

	if err := DB.First(&user, userId).Error; err != nil {
		models.ResponseJSON(context, http.StatusNotFound, "User not found", nil)
		return
	}
	if err := DB.First(&group, groupId).Preload("Users").Error; err != nil {
		models.ResponseJSON(context, http.StatusNotFound, "Group not found", nil)
		return
	}

	group.Users = append(group.Users, user)
	if err := DB.Save(group).Error; err != nil {
		models.ResponseJSON(context, http.StatusInternalServerError, "Unable to join the group: "+err.Error(), nil)
		return
	}

	models.ResponseJSON(context, http.StatusOK, "Group joined", nil)
}

func LeaveGroup(context *gin.Context) {
	userId := context.Keys["user_id"].(uint)
	groupId := context.Param("id")
	var user models.User
	var group models.Group

	if err := DB.First(&user, userId).Error; err != nil {
		models.ResponseJSON(context, http.StatusNotFound, "User not found", nil)
		return
	}
	if err := DB.First(&group, groupId).Error; err != nil {
		models.ResponseJSON(context, http.StatusNotFound, "Group not found", nil)
		return
	}

	if err := DB.Model(&group).Association("Users").Delete(&user); err != nil {
		models.ResponseJSON(context, http.StatusInternalServerError, "Unable to leave group: "+err.Error(), nil)
		return
	}

	models.ResponseJSON(context, http.StatusOK, "Group left", nil)
}
