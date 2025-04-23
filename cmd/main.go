package main

import (
	"api/api"

	"github.com/gin-gonic/gin"
)

func main() {
	api.InitDB()
	r := gin.Default()

	//routes
	r.POST("/horse", api.CreateHorse)
	r.GET("/horses", api.GetHorses)

	r.Run("localhost:8080")
}
