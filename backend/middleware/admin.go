package middleware

import (
	"net/http"

	"ganttpro-backend/models"

	"github.com/gin-gonic/gin"
)

// RequireAdmin middleware ensures the user has Admin role
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := GetUserRole(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Unauthorized - user role not found",
			})
			c.Abort()
			return
		}

		if role != models.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Access denied. Admin role required.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
