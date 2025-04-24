package handlers

import (
	"net/http"
	"time"

	middleware "api/api/middlewares"
	"api/api/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

const HASH_COST = 14
const EXPIRATION_MINUTES = 30

func generateJWT(context *gin.Context, user models.User) *string {
	var roleNames []string
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Name)
	}
	expirationTime := time.Now().Add(EXPIRATION_MINUTES * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"expiration": expirationTime.Unix(),
		"user_id":    user.ID,
		"roles":      roleNames,
	})
	// Sign the token
	{
		tokenString, err := token.SignedString(middleware.JwtSecret)
		if err != nil {
			models.ResponseJSON(context, http.StatusInternalServerError, "Could not generate token", nil)
			return nil
		}

		return &tokenString
	}
}

func Login(context *gin.Context) {
	var loginRequest models.LoginRequest
	var user models.User
	if err := context.ShouldBindJSON(&loginRequest); err != nil {
		models.ResponseJSON(context, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}
	if err := DB.Where("mail = ?", loginRequest.Mail).Preload("Roles").First(&user).Error; err != nil {
		models.ResponseJSON(context, http.StatusNotFound, "User not found", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(loginRequest.Password),
	); loginRequest.Mail != user.Mail || err != nil {
		models.ResponseJSON(context, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	if token := generateJWT(context, user); token != nil {
		models.ResponseJSON(context, http.StatusOK, "Login successful", gin.H{
			"user":  user,
			"token": token,
		})
	}
}

func SignUp(context *gin.Context) {
	var user models.User
	var userRole models.Role
	if err := context.ShouldBindJSON(&user); err != nil {
		models.ResponseJSON(context, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	DB.Where("name = ?", "USER").First(&userRole)
	user.Roles = []models.Role{userRole}
	if encryptedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), HASH_COST); err != nil {
		models.ResponseJSON(context, http.StatusInternalServerError, "Unable to hash password", nil)
		return
	} else {
		user.Password = string(encryptedPass)
	}

	{
		execution := DB.Create(&user)
		if execution.Error != nil {
			models.ResponseJSON(context, http.StatusInternalServerError, "Unable to create the user", nil)
			return
		}
	}
	if token := generateJWT(context, user); token != nil {
		models.ResponseJSON(context, http.StatusCreated, "Signup successful", gin.H{
			"user":  user,
			"token": token,
		})
	}
}
