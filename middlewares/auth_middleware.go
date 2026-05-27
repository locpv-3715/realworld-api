package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"message": "Missing authorization header"}})
			return
		}

		var tokenString string
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && (parts[0] == "Token" || parts[0] == "Bearer") {
			tokenString = parts[1]
		} else if len(parts) == 1 {
			// raw JWT passed directly
			tokenString = parts[0]
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"message": "Invalid token format. Use 'Token <jwt>', 'Bearer <jwt>', or pass the JWT directly"}})
			return
		}
		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"message": "Token is invalid or expired"}})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userID", claims["user_id"])
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"message": "Could not read token claims"}})
			return
		}
	}
}

func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Set("userID", float64(0))
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		var rawToken string
		if len(parts) == 2 && (parts[0] == "Token" || parts[0] == "Bearer") {
			rawToken = parts[1]
		} else if len(parts) == 1 {
			rawToken = parts[0]
		}
		if rawToken != "" {
			secret := os.Getenv("JWT_SECRET")
			token, err := jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err == nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					c.Set("userID", claims["user_id"])
				}
			}
		}

		if _, exists := c.Get("userID"); !exists {
			c.Set("userID", float64(0))
		}
		c.Next()
	}
}
