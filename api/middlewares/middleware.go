package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"api/api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Secret key for signing JWT
var JwtSecret = []byte(os.Getenv("SECRET_TOKEN"))

func extractClaims(token *jwt.Token) jwt.MapClaims {
	claims, _ := token.Claims.(jwt.MapClaims)
	return claims
}

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
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
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

		claims := extractClaims(token)

		expiration := time.Unix(int64(claims["expiration"].(float64)), 0)
		if time.Now().UTC().After(expiration) {
			models.ResponseJSON(context, http.StatusUnauthorized, "Token expired", nil)
			context.Abort()
			return
		}

		// Token is valid, proceed to the next handler
		var roles []string
		for _, role := range claims["roles"].([]any) {
			roles = append(roles, role.(string))
		}
		context.Set("user_id", uint(claims["user_id"].(float64)))
		context.Set("roles", roles)
		context.Next()
	}
}
