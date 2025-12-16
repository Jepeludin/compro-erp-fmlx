package middleware

import (
	"net/http"

	"ganttpro-backend/models"

	"github.com/gin-gonic/gin"
)

// Permission constants define granular access controls
const (
	// User management permissions
	PermissionUserRead   = "user:read"
	PermissionUserCreate = "user:create"
	PermissionUserUpdate = "user:update"
	PermissionUserDelete = "user:delete"

	// Machine permissions
	PermissionMachineRead   = "machine:read"
	PermissionMachineCreate = "machine:create"
	PermissionMachineUpdate = "machine:update"
	PermissionMachineDelete = "machine:delete"

	// Job Order permissions
	PermissionJobOrderRead   = "job_order:read"
	PermissionJobOrderCreate = "job_order:create"
	PermissionJobOrderUpdate = "job_order:update"
	PermissionJobOrderDelete = "job_order:delete"

	// Operation Plan permissions
	PermissionOpPlanRead    = "op_plan:read"
	PermissionOpPlanCreate  = "op_plan:create"
	PermissionOpPlanUpdate  = "op_plan:update"
	PermissionOpPlanDelete  = "op_plan:delete"
	PermissionOpPlanApprove = "op_plan:approve"
	PermissionOpPlanExecute = "op_plan:execute"

	// G-Code permissions
	PermissionGCodeRead   = "gcode:read"
	PermissionGCodeUpload = "gcode:upload"
	PermissionGCodeDelete = "gcode:delete"

	// PPIC Schedule permissions
	PermissionPPICScheduleRead   = "ppic_schedule:read"
	PermissionPPICScheduleCreate = "ppic_schedule:create"
	PermissionPPICScheduleUpdate = "ppic_schedule:update"
	PermissionPPICScheduleDelete = "ppic_schedule:delete"

	// Gantt Chart permissions
	PermissionGanttRead = "gantt:read"
)

// RolePermissions maps roles to their allowed permissions
var RolePermissions = map[string][]string{
	models.RoleAdmin: {
		// Admin has all permissions
		PermissionUserRead, PermissionUserCreate, PermissionUserUpdate, PermissionUserDelete,
		PermissionMachineRead, PermissionMachineCreate, PermissionMachineUpdate, PermissionMachineDelete,
		PermissionJobOrderRead, PermissionJobOrderCreate, PermissionJobOrderUpdate, PermissionJobOrderDelete,
		PermissionOpPlanRead, PermissionOpPlanCreate, PermissionOpPlanUpdate, PermissionOpPlanDelete, PermissionOpPlanApprove, PermissionOpPlanExecute,
		PermissionGCodeRead, PermissionGCodeUpload, PermissionGCodeDelete,
		PermissionPPICScheduleRead, PermissionPPICScheduleCreate, PermissionPPICScheduleUpdate, PermissionPPICScheduleDelete,
		PermissionGanttRead,
	},
	models.RolePPIC: {
		PermissionMachineRead,
		PermissionJobOrderRead, PermissionJobOrderCreate, PermissionJobOrderUpdate,
		PermissionOpPlanRead, PermissionOpPlanApprove,
		PermissionGCodeRead,
		PermissionPPICScheduleRead, PermissionPPICScheduleCreate, PermissionPPICScheduleUpdate, PermissionPPICScheduleDelete,
		PermissionGanttRead,
	},
	models.RolePEM: {
		PermissionMachineRead,
		PermissionJobOrderRead,
		PermissionOpPlanRead, PermissionOpPlanCreate, PermissionOpPlanUpdate, PermissionOpPlanApprove,
		PermissionGCodeRead, PermissionGCodeUpload,
		PermissionPPICScheduleRead,
		PermissionGanttRead,
	},
	models.RoleToolpather: {
		PermissionMachineRead,
		PermissionJobOrderRead,
		PermissionOpPlanRead, PermissionOpPlanApprove,
		PermissionGCodeRead, PermissionGCodeUpload,
		PermissionPPICScheduleRead,
		PermissionGanttRead,
	},
	models.RoleQC: {
		PermissionMachineRead,
		PermissionJobOrderRead,
		PermissionOpPlanRead, PermissionOpPlanApprove,
		PermissionGCodeRead,
		PermissionPPICScheduleRead,
		PermissionGanttRead,
	},
	models.RoleEngineering: {
		PermissionMachineRead,
		PermissionJobOrderRead,
		PermissionOpPlanRead, PermissionOpPlanApprove,
		PermissionGCodeRead,
		PermissionPPICScheduleRead,
		PermissionGanttRead,
	},
	models.RoleOperator: {
		PermissionMachineRead,
		PermissionJobOrderRead,
		PermissionOpPlanRead, PermissionOpPlanExecute,
		PermissionGCodeRead,
		PermissionPPICScheduleRead,
		PermissionGanttRead,
	},
	models.RoleGuest: {
		// Guest has read-only access
		PermissionMachineRead,
		PermissionJobOrderRead,
		PermissionOpPlanRead,
		PermissionGCodeRead,
		PermissionPPICScheduleRead,
		PermissionGanttRead,
	},
}

// HasPermission checks if a role has a specific permission
func HasPermission(role, permission string) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// GetUserRole safely extracts the user role from the Gin context
func GetUserRole(c *gin.Context) (string, bool) {
	roleValue, exists := c.Get("role")
	if !exists {
		return "", false
	}

	role, ok := roleValue.(string)
	if !ok {
		return "", false
	}

	return role, true
}

// GetUserFromContext safely extracts the user from the Gin context
func GetUserFromContext(c *gin.Context) (*models.User, bool) {
	userValue, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	user, ok := userValue.(*models.User)
	if !ok {
		return nil, false
	}

	return user, true
}

// GetUserIDFromContext safely extracts the user ID from the Gin context
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		return 0, false
	}

	return userID, true
}

// RequirePermission checks if the user has the required permission
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := GetUserRole(c)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "User role not found in context",
			})
			c.Abort()
			return
		}

		if !HasPermission(role, permission) {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "You do not have permission to perform this action",
				"required_permission": permission,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission checks if the user has at least one of the required permissions
func RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := GetUserRole(c)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "User role not found in context",
			})
			c.Abort()
			return
		}

		for _, permission := range permissions {
			if HasPermission(role, permission) {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "You do not have permission to perform this action",
			"required_permissions": permissions,
		})
		c.Abort()
	}
}

// RequireAllPermissions checks if the user has all the required permissions
func RequireAllPermissions(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := GetUserRole(c)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "User role not found in context",
			})
			c.Abort()
			return
		}

		for _, permission := range permissions {
			if !HasPermission(role, permission) {
				c.JSON(http.StatusForbidden, gin.H{
					"success": false,
					"error":   "You do not have permission to perform this action",
					"missing_permission": permission,
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}