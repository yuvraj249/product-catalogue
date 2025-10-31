package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v5"
)

func RoleAllowed(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims1, ok := c.Get("claims")
		if !ok {
			token1, exists := c.Get("token")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication Required"})
				c.Abort()
				return
			}

			token, ok := token1.(*jwt.Token)
			if !ok || token == nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token!!"})
				c.Abort()
				return
			}
			claims2, ok2 := token.Claims.(jwt.MapClaims)
			if !ok2 {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims!!"})
				c.Abort()
				return
			}
			claims1 = claims2

		}

		claims, ok := claims1.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims!!"})
			c.Abort()
			return
		}

		roleValue, ok := claims["role"]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "role absent in token!!"})
			c.Abort()
			return
		}
		role, _ := roleValue.(string)
		if slices.Contains(allowedRoles, role) {
			c.Set("role", role)
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access denied"})
		c.Abort()
	}
}

func RoleSystemAdmin() gin.HandlerFunc {
	return RoleAllowed("system_admin")
}

func RoleSupplierAdmin() gin.HandlerFunc {
	return RoleAllowed("supplier_admin")
}
