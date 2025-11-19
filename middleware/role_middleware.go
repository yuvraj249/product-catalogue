package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func RoleAllowed(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "role missing (auth required)"})
			c.Abort()
			return
		}
		if slices.Contains(allowedRoles, role) {
			c.Next()
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		c.Abort()
	}
}

func RoleSystemAdmin() gin.HandlerFunc {
	return RoleAllowed("system_admin")
}

func RoleSupplierAdmin() gin.HandlerFunc {
	return RoleAllowed("supplier_admin")
}
