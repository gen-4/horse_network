package main

import (
	"api/api/handlers"
	middleware "api/api/middlewares"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := godotenv.Load(); err != nil {
		log.Println("WARNING: Unable to read .env file")
	} else {
		switch os.Getenv("ENVIRONMENT") {
		case "dev":
			gin.SetMode(gin.DebugMode)

		case "test":
			gin.SetMode(gin.TestMode)

		case "pro":
			f, err := os.OpenFile("ginlogs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("Error opening file: %v", err)
			}
			defer f.Close()
			log.SetOutput(f)
			gin.SetMode(gin.ReleaseMode)
		}
	}

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
	}

	router.Run("localhost:8080")
}
