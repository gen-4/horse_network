package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"api/api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Secret key for signing JWT
var JwtSecret = []byte(os.Getenv("SECRET_TOKEN"))

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			models.ResponseJSON(context, http.StatusUnauthorized, "Authorization token required", nil)
			context.Abort()
			return
		}
		tokenString, ok := strings.CutPrefix(tokenString, "Bearer ")
		if !ok {
			models.ResponseJSON(context, http.StatusUnauthorized, "Invalid token", nil)
			context.Abort()
			return
		}

		// parse and validate the token
		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return JwtSecret, nil
		})
		if err != nil {
			models.ResponseJSON(context, http.StatusUnauthorized, "Invalid token", nil)
			context.Abort()
			return
		}
		// TODO: I will have to check the expiration time here and return the id and roles of the user

		// Token is valid, proceed to the next handler
		context.Next()
	}
}
