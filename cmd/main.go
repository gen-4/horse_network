package main

import (
	"api/api/handlers"
	middleware "api/api/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	handlers.InitDB()
	router := gin.Default()

	// Public routes
	router.POST("/login", handlers.Login)
	router.PUT("/signup", handlers.SignUp)

	// protected routes
	protected := router.Group("/", middleware.JWTAuthMiddleware())
	{
		protected.POST("/horse", handlers.CreateHorse)
		protected.GET("/horses", handlers.GetHorses)
		protected.PUT("/user", handlers.UpdateUser)
		protected.DELETE("/horse/:id", handlers.DeleteHorse)
		protected.POST("/group", handlers.CreateGroup)
		protected.POST("/group/:id/join", handlers.JoinGroup)

	}

	router.Run("localhost:8080")
}
