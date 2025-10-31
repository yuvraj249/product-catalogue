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
			token_string = strings.TrimSpace(cookie)
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
			token_string = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		}

		secret_k := os.Getenv("JWT_SECRET_KEY")
		if secret_k == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "jwt not configured"})
			return
		}
		token, err := jwt.Parse(token_string, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenUnverifiable
			}
			return []byte(secret_k), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims!!"})
			return
		}
		c.Set("token", token)
		c.Set("claims", claims)
		if _, err := c.Cookie("userCookie"); err == nil {
			c.Set("auth_source", "cookie")
			c.Set("auth_source", "header")

		}

		c.Next()

	}
}
