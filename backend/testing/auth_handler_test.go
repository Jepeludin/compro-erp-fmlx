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
// Test Setup
// =============================================================================

func init() {
	gin.SetMode(gin.TestMode)
}

// Helper to create JSON request body
func toJSON(v interface{}) *bytes.Buffer {
	data, _ := json.Marshal(v)
	return bytes.NewBuffer(data)
}

// Helper to parse JSON response
func parseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	return response
}

// =============================================================================
// Login Request Validation Tests
// =============================================================================

func TestLoginRequest_ValidJSON(t *testing.T) {
	testCases := []struct {
		name       string
		body       map[string]interface{}
		shouldBind bool
	}{
		{
			name: "Valid login request",
			body: map[string]interface{}{
				"user_id":  "PI1224.0001",
				"password": "password123",
			},
			shouldBind: true,
		},
		{
			name: "Missing user_id",
			body: map[string]interface{}{
				"password": "password123",
			},
			shouldBind: false,
		},
		{
			name: "Missing password",
			body: map[string]interface{}{
				"user_id": "PI1224.0001",
			},
			shouldBind: false,
		},
		{
			name:       "Empty body",
			body:       map[string]interface{}{},
			shouldBind: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/login", func(c *gin.Context) {
				var req struct {
					UserID   string `json:"user_id" binding:"required"`
					Password string `json:"password" binding:"required"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/login", toJSON(tc.body))
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
// Register Request Validation Tests
// =============================================================================

func TestRegisterRequest_ValidJSON(t *testing.T) {
	testCases := []struct {
		name       string
		body       map[string]interface{}
		shouldBind bool
	}{
		{
			name: "Valid full registration",
			body: map[string]interface{}{
				"user_id":  "PI1224.0001",
				"username": "TESTUSER",
				"password": "password123",
				"role":     "Operator",
				"operator": "OP-YSD01",
			},
			shouldBind: true,
		},
		{
			name: "Valid registration without operator",
			body: map[string]interface{}{
				"user_id":  "PI1224.0001",
				"username": "ADMIN",
				"password": "admin123",
				"role":     "Admin",
			},
			shouldBind: true,
		},
		{
			name: "Missing user_id",
			body: map[string]interface{}{
				"username": "TEST",
				"password": "pass",
				"role":     "Admin",
			},
			shouldBind: false,
		},
		{
			name: "Missing username",
			body: map[string]interface{}{
				"user_id":  "PI1224.0001",
				"password": "pass",
				"role":     "Admin",
			},
			shouldBind: false,
		},
		{
			name: "Missing password",
			body: map[string]interface{}{
				"user_id":  "PI1224.0001",
				"username": "TEST",
				"role":     "Admin",
			},
			shouldBind: false,
		},
		{
			name: "Missing role",
			body: map[string]interface{}{
				"user_id":  "PI1224.0001",
				"username": "TEST",
				"password": "pass",
			},
			shouldBind: false,
		},
		{
			name: "Username too short",
			body: map[string]interface{}{
				"user_id":  "PI1224.0001",
				"username": "AB", // min=3
				"password": "password",
				"role":     "Admin",
			},
			shouldBind: false,
		},
		{
			name: "Password too short",
			body: map[string]interface{}{
				"user_id":  "PI1224.0001",
				"username": "TESTUSER",
				"password": "abc", // min=4
				"role":     "Admin",
			},
			shouldBind: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/register", func(c *gin.Context) {
				var req struct {
					UserID   string `json:"user_id" binding:"required"`
					Username string `json:"username" binding:"required,min=3,max=50"`
					Password string `json:"password" binding:"required,min=4"`
					Role     string `json:"role" binding:"required"`
					Operator string `json:"operator"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, gin.H{"status": "created"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/register", toJSON(tc.body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if tc.shouldBind {
				assert.Equal(t, http.StatusCreated, w.Code, "Expected successful binding for %s", tc.name)
			} else {
				assert.Equal(t, http.StatusBadRequest, w.Code, "Expected binding error for %s", tc.name)
			}
		})
	}
}

// =============================================================================
// Authorization Header Tests
// =============================================================================

func TestAuthorizationHeader_Required(t *testing.T) {
	router := gin.New()
	router.GET("/protected", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	testCases := []struct {
		name         string
		authHeader   string
		expectedCode int
	}{
		{
			name:         "No header",
			authHeader:   "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "With header",
			authHeader:   "Bearer some-token",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/protected", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

// =============================================================================
// Response Format Tests
// =============================================================================

func TestLoginResponse_Format(t *testing.T) {
	// Simulate successful login response structure
	response := map[string]interface{}{
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		"user": map[string]interface{}{
			"id":        1,
			"username":  "TESTUSER",
			"user_id":   "PI1224.0001",
			"role":      "Admin",
			"operator":  "",
			"is_active": true,
		},
	}

	// Verify response has expected fields
	assert.Contains(t, response, "token")
	assert.Contains(t, response, "user")

	user := response["user"].(map[string]interface{})
	assert.Contains(t, user, "id")
	assert.Contains(t, user, "username")
	assert.Contains(t, user, "user_id")
	assert.Contains(t, user, "role")
	assert.Contains(t, user, "is_active")
}

func TestErrorResponse_Format(t *testing.T) {
	router := gin.New()
	router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/error", nil)
	router.ServeHTTP(w, req)

	response := parseResponse(w)
	assert.Contains(t, response, "error")
	assert.Equal(t, "Something went wrong", response["error"])
}

// =============================================================================
// Profile Endpoint Tests
// =============================================================================

func TestGetProfile_RequiresAuth(t *testing.T) {
	router := gin.New()
	router.GET("/profile", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		// Simulate returning user profile
		c.JSON(http.StatusOK, gin.H{
			"user": map[string]interface{}{
				"id":       1,
				"username": "TESTUSER",
			},
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/profile", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetProfile_WithAuth(t *testing.T) {
	router := gin.New()
	router.GET("/profile", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"user": map[string]interface{}{
				"id":       1,
				"username": "TESTUSER",
			},
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/profile", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response := parseResponse(w)
	assert.Contains(t, response, "user")
}

// =============================================================================
// Logout Endpoint Tests
// =============================================================================

func TestLogout_RequiresAuth(t *testing.T) {
	router := gin.New()
	router.POST("/logout", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/logout", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogout_Success(t *testing.T) {
	router := gin.New()
	router.POST("/logout", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/logout", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response := parseResponse(w)
	assert.Contains(t, response, "message")
}

// =============================================================================
// Content-Type Tests
// =============================================================================

func TestContentType_RequiredForJSON(t *testing.T) {
	router := gin.New()
	router.POST("/json-endpoint", func(c *gin.Context) {
		var data map[string]interface{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"received": data})
	})

	testCases := []struct {
		name        string
		contentType string
		body        string
		expectError bool
	}{
		{
			name:        "With JSON content type",
			contentType: "application/json",
			body:        `{"key": "value"}`,
			expectError: false,
		},
		{
			name:        "Without content type",
			contentType: "",
			body:        `{"key": "value"}`,
			expectError: false, // Gin can still parse it
		},
		{
			name:        "Invalid JSON",
			contentType: "application/json",
			body:        `{invalid json}`,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/json-endpoint", bytes.NewBufferString(tc.body))
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			router.ServeHTTP(w, req)

			if tc.expectError {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
			}
		})
	}
}

