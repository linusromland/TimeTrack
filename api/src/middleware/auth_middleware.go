package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		jwtSecret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["email"] == nil || claims["userId"] == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userID := claims["userId"].(string)
		c.Set("user_id", userID)
		c.Next()
	}
}

func IntegrationAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		integrationTokenString := c.GetHeader("Authorization")
		if integrationTokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Integration token is required"})
			c.Abort()
			return
		}

		integrationTokenString = strings.Replace(integrationTokenString, "Bearer ", "", 1)

		integrationSecret := os.Getenv("INTEGRATION_SECRET")
		token, err := jwt.Parse(integrationTokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(integrationSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired integration token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
