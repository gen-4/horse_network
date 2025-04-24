package handlers

import (
	"net/http"
	"time"

	middleware "api/api/middlewares"
	"api/api/models"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

func GenerateJWT(c *gin.Context) {
	var loginRequest models.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		models.ResponseJSON(c, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	if loginRequest.Mail != "admin@a.com" || loginRequest.Password != "password" {
		models.ResponseJSON(c, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"expiration": expirationTime.Unix(), // TODO: I will have to add id and role, at least, here
	})
	// Sign the token
	tokenString, err := token.SignedString(middleware.JwtSecret)
	if err != nil {
		models.ResponseJSON(c, http.StatusInternalServerError, "Could not generate token", nil)
		return
	}

	models.ResponseJSON(c, http.StatusOK, "Token generated successfully", gin.H{"token": tokenString})
}
