package testing

import (
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

func setupAdminTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// Helper to set user role in context
func setRoleInContext(c *gin.Context, role string) {
	c.Set("role", role)
}

// =============================================================================
// RequireAdmin Middleware Tests
// =============================================================================

func TestRequireAdmin_AdminRole(t *testing.T) {
	router := setupAdminTestRouter()
	
	// Setup route with middleware
	router.GET("/admin", func(c *gin.Context) {
		c.Set("role", "Admin") // Simulate auth middleware setting role
		c.Next()
	}, middleware.RequireAdmin(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireAdmin_NonAdminRole(t *testing.T) {
	nonAdminRoles := []string{"PPIC", "PEM", "Engineering", "Toolpather", "QC", "Operator", "Guest"}

	for _, role := range nonAdminRoles {
		t.Run("Role_"+role, func(t *testing.T) {
			router := setupAdminTestRouter()
			
			router.GET("/admin", func(c *gin.Context) {
				c.Set("role", role)
				c.Next()
			}, middleware.RequireAdmin(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "success"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/admin", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusForbidden, w.Code)
		})
	}
}

func TestRequireAdmin_NoRoleSet(t *testing.T) {
	router := setupAdminTestRouter()
	
	// Don't set role in context
	router.GET("/admin", middleware.RequireAdmin(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRequireAdmin_EmptyRole(t *testing.T) {
	router := setupAdminTestRouter()
	
	router.GET("/admin", func(c *gin.Context) {
		c.Set("role", "") // Empty role
		c.Next()
	}, middleware.RequireAdmin(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequireAdmin_InvalidRoleType(t *testing.T) {
	router := setupAdminTestRouter()
	
	router.GET("/admin", func(c *gin.Context) {
		c.Set("role", 123) // Wrong type (int instead of string)
		c.Next()
	}, middleware.RequireAdmin(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// =============================================================================
// Error Response Tests
// =============================================================================

func TestRequireAdmin_UnauthorizedErrorResponse(t *testing.T) {
	router := setupAdminTestRouter()
	
	router.GET("/admin", middleware.RequireAdmin(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
}

func TestRequireAdmin_ForbiddenErrorResponse(t *testing.T) {
	router := setupAdminTestRouter()
	
	router.GET("/admin", func(c *gin.Context) {
		c.Set("role", "PPIC")
		c.Next()
	}, middleware.RequireAdmin(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Access denied")
	assert.Contains(t, w.Body.String(), "Admin role required")
}

// =============================================================================
// Middleware Chain Tests
// =============================================================================

func TestRequireAdmin_MiddlewareChain(t *testing.T) {
	router := setupAdminTestRouter()
	
	callOrder := []string{}

	router.GET("/admin",
		func(c *gin.Context) {
			callOrder = append(callOrder, "first")
			c.Set("role", "Admin")
			c.Next()
		},
		middleware.RequireAdmin(),
		func(c *gin.Context) {
			callOrder = append(callOrder, "handler")
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		},
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, []string{"first", "handler"}, callOrder)
}

func TestRequireAdmin_AbortsProperly(t *testing.T) {
	router := setupAdminTestRouter()
	
	handlerCalled := false

	router.GET("/admin",
		func(c *gin.Context) {
			c.Set("role", "PPIC") // Non-admin
			c.Next()
		},
		middleware.RequireAdmin(),
		func(c *gin.Context) {
			handlerCalled = true
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		},
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.False(t, handlerCalled, "Handler should not be called when middleware aborts")
}

// =============================================================================
// Case Sensitivity Tests
// =============================================================================

func TestRequireAdmin_CaseSensitiveRole(t *testing.T) {
	testCases := []struct {
		role           string
		shouldSucceed  bool
	}{
		{"Admin", true},
		{"admin", false},    // Lowercase
		{"ADMIN", false},    // Uppercase
		{"AdMiN", false},    // Mixed case
		{" Admin", false},   // Leading space
		{"Admin ", false},   // Trailing space
	}

	for _, tc := range testCases {
		t.Run("Role_"+tc.role, func(t *testing.T) {
			router := setupAdminTestRouter()
			
			router.GET("/admin", func(c *gin.Context) {
				c.Set("role", tc.role)
				c.Next()
			}, middleware.RequireAdmin(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "success"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/admin", nil)
			router.ServeHTTP(w, req)

			if tc.shouldSucceed {
				assert.Equal(t, http.StatusOK, w.Code)
			} else {
				assert.Equal(t, http.StatusForbidden, w.Code)
			}
		})
	}
}

// =============================================================================
// Multiple Routes Tests
// =============================================================================

func TestRequireAdmin_MultipleRoutes(t *testing.T) {
	router := setupAdminTestRouter()
	
	// Protected admin routes
	adminGroup := router.Group("/admin")
	adminGroup.Use(func(c *gin.Context) {
		c.Set("role", "Admin")
		c.Next()
	})
	adminGroup.Use(middleware.RequireAdmin())
	{
		adminGroup.GET("/users", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"route": "users"})
		})
		adminGroup.GET("/settings", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"route": "settings"})
		})
	}

	// Test /admin/users
	w1 := httptest.NewRecorder()
	req1 := httptest.NewRequest("GET", "/admin/users", nil)
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Test /admin/settings
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/admin/settings", nil)
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

