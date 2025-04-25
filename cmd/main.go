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
	r.PUT("/signup", handlers.SignUp)

	// protected routes
	protected := r.Group("/", middleware.JWTAuthMiddleware())
	{
		protected.POST("/horse", handlers.CreateHorse)
		protected.GET("/horses", handlers.GetHorses)
		protected.PUT("/user", handlers.UpdateUser)
	}

	r.Run("localhost:8080")
}
