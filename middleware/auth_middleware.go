package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token_string string
		cookie, err := c.Cookie("userCookie")
		if err == nil && cookie != "" {
			token_string = cookie
		}
		if token_string == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing in header"})
				c.Abort()
				return
			}
			if !strings.HasPrefix(authHeader, "Bearer ") {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
				c.Abort()
				return
			}
			token_string = strings.TrimPrefix(authHeader, "Bearer")
			token_string = strings.TrimSpace(token_string)
		}

		secret_k := os.Getenv("JWT_SECRET_KEY")
		token, err := jwt.Parse(token_string, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret_k), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("token", token)

		c.Next()

	}
}
