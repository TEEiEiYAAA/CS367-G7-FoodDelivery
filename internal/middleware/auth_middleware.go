package middleware

import (
	"CS367-G7-FoodDelivery/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Expect: Bearer TOKEN
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		// Save user info in context
		c.Set("username", claims["username"])
		c.Set("role", claims["role"])

		c.Next()
	}
}
