package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireAdmin middleware memastikan user memiliki role Admin
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by auth middleware)
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		// Check if role is Admin
		roleStr, ok := role.(string)
		if !ok || roleStr != "Admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied. Admin role required.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
