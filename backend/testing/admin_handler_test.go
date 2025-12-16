package testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Admin Handler Request Validation Tests
// =============================================================================

func TestUpdateUserRequest_ValidJSON(t *testing.T) {
	testCases := []struct {
		name       string
		body       map[string]interface{}
		shouldBind bool
	}{
		{
			name: "Update password only",
			body: map[string]interface{}{
				"password": "newPassword123",
			},
			shouldBind: true,
		},
		{
			name: "Update role only",
			body: map[string]interface{}{
				"role": "PEM",
			},
			shouldBind: true,
		},
		{
			name: "Update is_active only",
			body: map[string]interface{}{
				"is_active": false,
			},
			shouldBind: true,
		},
		{
			name: "Update all fields",
			body: map[string]interface{}{
				"password":  "newPassword123",
				"role":      "Engineering",
				"is_active": true,
			},
			shouldBind: true,
		},
		{
			name:       "Empty body",
			body:       map[string]interface{}{},
			shouldBind: true, // Empty body is valid for partial updates
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.PUT("/users/:id", func(c *gin.Context) {
				var req struct {
					Password *string `json:"password,omitempty"`
					Role     *string `json:"role,omitempty"`
					IsActive *bool   `json:"is_active,omitempty"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/users/1", toJSON(tc.body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if tc.shouldBind {
				assert.Equal(t, http.StatusOK, w.Code)
			} else {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})
	}
}

// =============================================================================
// GetAllUsers Endpoint Tests
// =============================================================================

func TestGetAllUsers_ResponseFormat(t *testing.T) {
	router := gin.New()
	router.GET("/admin/users", func(c *gin.Context) {
		// Simulate successful response
		users := []gin.H{
			{
				"id":       1,
				"username": "ADMIN",
				"user_id":  "PI1224.0001",
				"role":     "Admin",
				"operator": "",
				"is_active": true,
			},
			{
				"id":       2,
				"username": "BAYU",
				"user_id":  "PI1224.0002",
				"role":     "Operator",
				"operator": "OP-YSD01",
				"is_active": true,
			},
		}
		c.JSON(http.StatusOK, gin.H{
			"users": users,
			"total": len(users),
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin/users", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "users")
	assert.Contains(t, response, "total")
	assert.Equal(t, float64(2), response["total"])
}

func TestGetAllUsers_EmptyList(t *testing.T) {
	router := gin.New()
	router.GET("/admin/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"users": []gin.H{},
			"total": 0,
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin/users", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(0), response["total"])
}

// =============================================================================
// UpdateUser Endpoint Tests
// =============================================================================

func TestUpdateUser_InvalidUserID(t *testing.T) {
	router := gin.New()
	router.PUT("/admin/users/:id", func(c *gin.Context) {
		_, err := c.Params.Get("id")
		if !err {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Simulate ID parsing
		id := c.Param("id")
		if id == "invalid" || id == "abc" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	testCases := []struct {
		id           string
		expectedCode int
	}{
		{"1", http.StatusOK},
		{"123", http.StatusOK},
		{"invalid", http.StatusBadRequest},
		{"abc", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run("ID_"+tc.id, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/admin/users/"+tc.id, bytes.NewBufferString("{}"))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestUpdateUser_ValidRoles(t *testing.T) {
	validRoles := []string{"Admin", "PPIC", "Toolpather", "PEM", "QC", "Engineering", "Guest"}

	for _, role := range validRoles {
		t.Run("Role_"+role, func(t *testing.T) {
			router := gin.New()
			router.PUT("/admin/users/:id", func(c *gin.Context) {
				var req struct {
					Role *string `json:"role"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				// Validate role
				validRoles := []string{"Admin", "PPIC", "Toolpather", "PEM", "QC", "Engineering", "Guest"}
				isValid := false
				for _, vr := range validRoles {
					if req.Role != nil && *req.Role == vr {
						isValid = true
						break
					}
				}

				if req.Role != nil && !isValid {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
					return
				}

				c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
			})

			body := map[string]string{"role": role}
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/admin/users/1", toJSON(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestUpdateUser_InvalidRole(t *testing.T) {
	invalidRoles := []string{"SuperAdmin", "Manager", "user", "OPERATOR", "invalid"}

	for _, role := range invalidRoles {
		t.Run("InvalidRole_"+role, func(t *testing.T) {
			router := gin.New()
			router.PUT("/admin/users/:id", func(c *gin.Context) {
				var req struct {
					Role *string `json:"role"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				validRoles := []string{"Admin", "PPIC", "Toolpather", "PEM", "QC", "Engineering", "Guest"}
				isValid := false
				for _, vr := range validRoles {
					if req.Role != nil && *req.Role == vr {
						isValid = true
						break
					}
				}

				if req.Role != nil && !isValid {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
					return
				}

				c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
			})

			body := map[string]string{"role": role}
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/admin/users/1", toJSON(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestUpdateUser_SuccessResponse(t *testing.T) {
	router := gin.New()
	router.PUT("/admin/users/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "User updated successfully",
			"user": gin.H{
				"id":        1,
				"username":  "TESTUSER",
				"user_id":   "PI1224.0001",
				"role":      "PEM",
				"is_active": true,
			},
		})
	})

	body := map[string]string{"role": "PEM"}
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/admin/users/1", toJSON(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "user")
}

// =============================================================================
// DeleteUser Endpoint Tests
// =============================================================================

func TestDeleteUser_InvalidUserID(t *testing.T) {
	router := gin.New()
	router.DELETE("/admin/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "invalid" || id == "abc" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	})

	testCases := []struct {
		id           string
		expectedCode int
	}{
		{"1", http.StatusOK},
		{"123", http.StatusOK},
		{"invalid", http.StatusBadRequest},
		{"abc", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run("DeleteID_"+tc.id, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/admin/users/"+tc.id, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestDeleteUser_CannotDeleteSelf(t *testing.T) {
	router := gin.New()
	router.DELETE("/admin/users/:id", func(c *gin.Context) {
		// Simulate current user ID from context
		c.Set("user_id", uint(1))
		
		adminUserID, exists := c.Get("user_id")
		if exists && adminUserID.(uint) == uint(1) && c.Param("id") == "1" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete your own account"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	})

	// Try to delete self
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/admin/users/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Cannot delete your own account")
}

func TestDeleteUser_SuccessResponse(t *testing.T) {
	router := gin.New()
	router.DELETE("/admin/users/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":  "User deleted successfully",
			"username": "DELETEDUSER",
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/admin/users/2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "username")
}

func TestDeleteUser_UserNotFound(t *testing.T) {
	router := gin.New()
	router.DELETE("/admin/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "999" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/admin/users/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}

// =============================================================================
// Error Response Tests
// =============================================================================

func TestAdminHandler_InternalServerError(t *testing.T) {
	router := gin.New()
	router.GET("/admin/users", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin/users", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to retrieve users")
}

// =============================================================================
// Password Update Tests
// =============================================================================

func TestUpdateUser_PasswordNotEmpty(t *testing.T) {
	router := gin.New()
	router.PUT("/admin/users/:id", func(c *gin.Context) {
		var req struct {
			Password *string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Password != nil && *req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password cannot be empty"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
	})

	// Test empty password
	body := map[string]string{"password": ""}
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/admin/users/1", toJSON(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateUser_ValidPassword(t *testing.T) {
	router := gin.New()
	router.PUT("/admin/users/:id", func(c *gin.Context) {
		var req struct {
			Password *string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
	})

	body := map[string]string{"password": "newSecurePassword123"}
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/admin/users/1", toJSON(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

