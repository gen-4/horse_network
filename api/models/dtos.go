package models

import "github.com/gin-gonic/gin"

type JsonResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func ResponseJSON(context *gin.Context, status int, message string, data any) {
	response := JsonResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}

	context.JSON(status, response)
}

type LoginRequest struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

type UserDto struct {
	ID       uint       `json:"id"`
	Username string     `json:"username"`
	Mail     string     `json:"mail"`
	Age      uint       `json:"age"`
	Gender   string     `json:"gender"`
	Country  string     `json:"country"`
	Roles    []Role     `json:"roles"`
	Horses   []HorseDto `json:"horses"`
}

type HorseDto struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Breed  string `json:"breed"`
	Age    uint   `json:"age"`
	Gender string `json:"gender"`
}

type UpdateUser struct {
	Username string `json:"username" binding:"omitempty,min=3"`
	Mail     string `json:"mail" binding:"omitempty,email"`
	Password string `json:"omitempty,password"`
	Age      uint   `json:"age"`
	Gender   string `json:"gender" binding:"omitempty,oneof=m f"`
	Country  string `json:"country" binding:"omitempty,iso3166_1_alpha2"`
}

type GroupDto struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Users uint   `json:"user_number"`
}
