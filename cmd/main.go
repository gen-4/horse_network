package main

import (
	"api/api/handlers"
	middleware "api/api/middlewares"
	"api/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Config()
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
		protected.POST("/group/:id/leave", handlers.LeaveGroup)
		protected.GET("/groups", handlers.DiscoverGroups)
	}

	router.Run("localhost:8080")
}
