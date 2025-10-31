package handlers

import (
	"ganttpro-backend/repository"
	"ganttpro-backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	userRepo *repository.UserRepository
}

func NewAdminHandler(userRepo *repository.UserRepository) *AdminHandler {
	return &AdminHandler{
		userRepo: userRepo,
	}
}

// GetAllUsers - Get all users (Admin only)
func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userRepo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users",
		})
		return
	}

	// Remove password from response
	var usersResponse []gin.H
	for _, user := range users {
		usersResponse = append(usersResponse, gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"user_id":    user.UserID,
			"role":       user.Role,
			"operator":   user.Operator,
			"is_active":  user.IsActive,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users": usersResponse,
		"total": len(usersResponse),
	})
}

// UpdateUser - Update user details (Admin only)
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	// Get user ID from URL
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Parse request body
	var req struct {
		Password *string `json:"password,omitempty"`
		Role     *string `json:"role,omitempty"`
		IsActive *bool   `json:"is_active,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Get existing user
	user, err := h.userRepo.GetByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Update password if provided
	if req.Password != nil && *req.Password != "" {
		hashedPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to hash password",
			})
			return
		}
		user.Password = hashedPassword
	}

	// Update role if provided
	if req.Role != nil {
		// Validate role
		validRoles := []string{"Admin", "PPIC", "Toolpather", "PEM", "QC", "Engineering", "Guest"}
		isValidRole := false
		for _, validRole := range validRoles {
			if *req.Role == validRole {
				isValidRole = true
				break
			}
		}
		if !isValidRole {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid role. Must be one of: Admin, PPIC, Toolpather, PEM, QC, Engineering, Guest",
			})
			return
		}
		user.Role = *req.Role
	}

	// Update is_active if provided
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Save to database
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"user_id":   user.UserID,
			"role":      user.Role,
			"is_active": user.IsActive,
		},
	})
}

// DeleteUser - Delete user (Admin only)
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	// Get user ID from URL
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Get current admin user ID to prevent self-deletion
	adminUserID, exists := c.Get("user_id")
	if exists && adminUserID.(uint) == uint(userID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot delete your own account",
		})
		return
	}

	// Check if user exists
	user, err := h.userRepo.GetByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Soft delete (or hard delete if you prefer)
	if err := h.userRepo.Delete(uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "User deleted successfully",
		"username": user.Username,
	})
}
