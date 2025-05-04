package tests

import (
	"api/api/handlers"
	"api/api/middlewares"
	"api/api/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const UNEXISTENT_ID = -1

func SetupTestDB() {
	var err error
	handlers.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect test database")
	}
	handlers.DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Horse{},
		&models.Group{},
	)

	roles := []models.Role{
		models.Role{Name: "ADMIN"},
		models.Role{Name: "USER"},
	}
	handlers.DB.Create(&roles)
}

func AddHorse(name string, age uint, breed string, owner models.User) models.Horse {
	horse := models.Horse{
		Name:  name,
		Age:   age,
		Breed: breed,
		Owner: &owner.ID,
	}
	handlers.DB.Create(&horse)
	return horse
}

func AddUser(
	username string,
	age uint,
	country string,
	mail string,
	gender string,
	password string,
	roles []models.Role,
) models.User {
	user := models.User{
		Age:      age,
		Username: username,
		Mail:     mail,
		Gender:   gender,
		Password: EncryptPass(password),
		Country:  country,
		Roles:    roles,
	}
	handlers.DB.Create(&user)
	user.Password = password
	return user
}

func EncryptPass(pass string) string {
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(pass), handlers.HASH_COST)
	return string(encryptedPass)
}

func GetRoleByName(name models.RoleEnum) models.Role {
	var role models.Role
	handlers.DB.Where("name = ?", name).First(&role)
	return role
}

func GetRoutes() *gin.Engine {
	router := gin.Default()

	router.POST("/login", handlers.Login)
	router.PUT("/signup", handlers.SignUp)

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

	return router
}

func LogUser(router *gin.Engine, user models.User) string {
	loginData := models.LoginRequest{
		Mail:     user.Mail,
		Password: user.Password,
	}
	jsonValue, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	var response models.JsonResponse
	json.NewDecoder(recorder.Body).Decode(&response)
	return response.Data.(map[string]any)["token"].(string)
}
