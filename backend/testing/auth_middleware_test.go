package testing

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"ganttpro-backend/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Test Setup
// =============================================================================

func init() {
	gin.SetMode(gin.TestMode)
}

// Helper to create a test context
func createTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

// Helper to parse JSON response
func parseJSONResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	return response
}

// =============================================================================
// RequireRole Tests
// =============================================================================
/*
func TestRequireRole_UserHasRequiredRole(t *testing.T) {
	// Arrange
	c, w := createTestContext()
	c.Set("role", "Admin")

	middleware := middleware.RequireRole("Admin")

	c.Request = httptest.NewRequest("GET", "/test", nil)

	// Create a test handler chain
	router := gin.New()
	router.GET("/test", middleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Reset the context for router
	w = httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)

	// Simulate setting the role in context (this would normally be done by AuthMiddleware)
	router.Use(func(c *gin.Context) {
		c.Set("role", "Admin")
		c.Next()
	})

	// Recreate router with proper middleware chain
	router = gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", "Admin")
		c.Next()
	})
	router.GET("/test", middleware.RequireRole("Admin"), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
	response := parseJSONResponse(w)
	assert.Equal(t, "Insufficient permissions", response["error"])
}
*/
func TestRequireRole_MultipleRolesAllowed(t *testing.T) {
	testCases := []struct {
		name         string
		userRole     string
		allowedRoles []string
		expectedCode int
	}{
		{
			name:         "Admin accessing Admin+PEM route",
			userRole:     "Admin",
			allowedRoles: []string{"Admin", "PEM"},
			expectedCode: http.StatusOK,
		},
		{
			name:         "PEM accessing Admin+PEM route",
			userRole:     "PEM",
			allowedRoles: []string{"Admin", "PEM"},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Operator accessing Admin+PEM route",
			userRole:     "Operator",
			allowedRoles: []string{"Admin", "PEM"},
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "QC accessing QC+Engineering+Toolpather route",
			userRole:     "QC",
			allowedRoles: []string{"QC", "Engineering", "Toolpather"},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("role", tc.userRole)
				c.Next()
			})
			router.GET("/test", middleware.RequireRole(tc.allowedRoles...), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestRequireRole_NoRoleInContext(t *testing.T) {
	// Arrange
	router := gin.New()
	// Note: NOT setting role in context
	router.GET("/test", middleware.RequireRole("Admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
	response := parseJSONResponse(w)
	assert.Equal(t, "User role not found", response["error"])
}

func TestRequireRole_EmptyRolesList(t *testing.T) {
	// Arrange
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", "Admin")
		c.Next()
	})
	router.GET("/test", middleware.RequireRole(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert - No roles specified means no one has permission
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// =============================================================================
// Authorization Header Parsing Tests
// =============================================================================

func TestParseAuthorizationHeader(t *testing.T) {
	testCases := []struct {
		name           string
		header         string
		expectedParts  int
		expectedPrefix string
		isValid        bool
	}{
		{
			name:           "Valid Bearer token",
			header:         "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectedParts:  2,
			expectedPrefix: "Bearer",
			isValid:        true,
		},
		{
			name:           "Missing Bearer prefix",
			header:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectedParts:  1,
			expectedPrefix: "",
			isValid:        false,
		},
		{
			name:           "Wrong prefix",
			header:         "Basic eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectedParts:  2,
			expectedPrefix: "Basic",
			isValid:        false,
		},
		{
			name:           "Empty header",
			header:         "",
			expectedParts:  1,
			expectedPrefix: "",
			isValid:        false,
		},
		{
			name:           "Bearer with no token",
			header:         "Bearer",
			expectedParts:  1,
			expectedPrefix: "Bearer",
			isValid:        false,
		},
		{
			name:           "Bearer with extra spaces",
			header:         "Bearer  token",
			expectedParts:  2, // SplitN with n=2 keeps remaining as second part
			expectedPrefix: "Bearer",
			isValid:        true, // Will have " token" with leading space
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate the parsing logic from AuthMiddleware
			parts := splitAuthHeader(tc.header)

			isValid := len(parts) == 2 && parts[0] == "Bearer"
			assert.Equal(t, tc.isValid, isValid, "Validity check should match")
		})
	}
}

// Helper function to split auth header (mirrors middleware logic)
func splitAuthHeader(header string) []string {
	if header == "" {
		return []string{""}
	}
	parts := make([]string, 0)
	for i, part := range splitN(header, " ", 2) {
		_ = i
		parts = append(parts, part)
	}
	return parts
}

// Simple splitN implementation
func splitN(s, sep string, n int) []string {
	if n == 0 {
		return nil
	}
	if n < 0 {
		n = len(s) + 1
	}

	result := make([]string, 0, n)
	for i := 0; i < n-1; i++ {
		idx := indexString(s, sep)
		if idx < 0 {
			break
		}
		result = append(result, s[:idx])
		s = s[idx+len(sep):]
	}
	result = append(result, s)
	return result
}

func indexString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// =============================================================================
// Role-Based Access Control Tests (RBAC)
// =============================================================================

func TestRBAC_AdminHasFullAccess(t *testing.T) {
	// Admin should be able to access all role-specific routes
	routes := []struct {
		path         string
		allowedRoles []string
	}{
		{"/admin/users", []string{"Admin"}},
		{"/operation-plans", []string{"Admin", "PEM", "PPIC"}},
		{"/machines", []string{"Admin", "PPIC", "Operator"}},
	}

	for _, route := range routes {
		t.Run(route.path, func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("role", "Admin")
				c.Next()
			})
			router.GET(route.path, middleware.RequireRole(route.allowedRoles...), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", route.path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Admin should have access to %s", route.path)
		})
	}
}

func TestRBAC_OperatorLimitedAccess(t *testing.T) {
	testCases := []struct {
		path         string
		allowedRoles []string
		expectedCode int
	}{
		{"/admin/users", []string{"Admin"}, http.StatusForbidden},
		{"/machines", []string{"Admin", "PPIC", "Operator"}, http.StatusOK},
		{"/operation-plans/approve", []string{"PEM", "PPIC", "QC", "Engineering", "Toolpather"}, http.StatusForbidden},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("role", "Operator")
				c.Next()
			})
			router.GET(tc.path, middleware.RequireRole(tc.allowedRoles...), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", tc.path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

// =============================================================================
// Approval Roles Tests
// =============================================================================

func TestApprovalRoles(t *testing.T) {
	// Test that all 5 approval roles can approve
	approverRoles := []string{"PEM", "PPIC", "QC", "Engineering", "Toolpather"}

	for _, role := range approverRoles {
		t.Run(role+"_CanApprove", func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("role", role)
				c.Next()
			})
			router.POST("/approve", middleware.RequireRole(approverRoles...), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "approved"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/approve", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "%s should be able to approve", role)
		})
	}
}

func TestNonApproverCannotApprove(t *testing.T) {
	nonApproverRoles := []string{"Operator", "Guest"}
	approverRoles := []string{"PEM", "PPIC", "QC", "Engineering", "Toolpather"}

	for _, role := range nonApproverRoles {
		t.Run(role+"_CannotApprove", func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("role", role)
				c.Next()
			})
			router.POST("/approve", middleware.RequireRole(approverRoles...), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "approved"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/approve", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusForbidden, w.Code, "%s should NOT be able to approve", role)
		})
	}
}
