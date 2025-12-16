package middleware

import (
	"net/http"
	"strings"

	"ganttpro-backend/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]
		isBlacklisted, err := authService.IsTokenBlacklisted(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify token"})
			c.Abort()
			return
		}

		if isBlacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			c.Abort()
			return
		}

		user, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("userid", user.UserID)
		c.Set("role", user.Role)

		c.Next()
	}
}

// RequireRole checks if the user has one of the required roles
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, ok := GetUserRole(c)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "User role not found in context",
			})
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Insufficient permissions. Required role: " + strings.Join(roles, " or "),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireApprover checks if the user can approve operation plans
func RequireApprover() gin.HandlerFunc {
	return RequirePermission(PermissionOpPlanApprove)
}
