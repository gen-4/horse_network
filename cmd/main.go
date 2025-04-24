package main

import (
	"api/api/handlers"
	middleware "api/api/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	handlers.InitDB()
	r := gin.Default()

	// Public routes
	r.POST("/login", handlers.Login)
	r.POST("/signup", handlers.SignUp)

	// protected routes
	protected := r.Group("/", middleware.JWTAuthMiddleware())
	{
		protected.POST("/horse", handlers.CreateHorse)
		protected.GET("/horses", handlers.GetHorses)
	}

	r.Run("localhost:8080")
}
