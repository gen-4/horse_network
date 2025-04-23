package api

import "github.com/gin-gonic/gin"

type Horse struct {
	ID uint `json:"id" gorm:"primaryKey"`
}

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
