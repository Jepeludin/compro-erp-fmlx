package testing

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"ganttpro-backend/middleware"
	"ganttpro-backend/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Test Setup
// =============================================================================

func init() {
	gin.SetMode(gin.TestMode)
}

func createTestConfigWithOrigins(origins []string) *config.Config {
	return &config.Config{
		AllowedOrigins: origins,
	}
}

// =============================================================================
// CORS Middleware Tests
// =============================================================================

func TestCORS_AllowedOrigin(t *testing.T) {
	cfg := createTestConfigWithOrigins([]string{"http://localhost:5173", "http://example.com"})

	router := gin.New()
	router.Use(middleware.CORS(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_DisallowedOrigin(t *testing.T) {
	cfg := createTestConfigWithOrigins([]string{"http://localhost:5173"})

	router := gin.New()
	router.Use(middleware.CORS(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://malicious-site.com")
	router.ServeHTTP(w, req)

	// Request should still succeed, but no CORS header set
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_WildcardOrigin(t *testing.T) {
	cfg := createTestConfigWithOrigins([]string{"*"})

	router := gin.New()
	router.Use(middleware.CORS(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://any-origin.com")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://any-origin.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_OptionsPreflightRequest(t *testing.T) {
	cfg := createTestConfigWithOrigins([]string{"http://localhost:5173"})

	router := gin.New()
	router.Use(middleware.CORS(cfg))
	router.POST("/api/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("OPTIONS", "/api/data", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")
	router.ServeHTTP(w, req)

	// OPTIONS should return 204 No Content
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestCORS_AllowedHeaders(t *testing.T) {
	cfg := createTestConfigWithOrigins([]string{"http://localhost:5173"})

	router := gin.New()
	router.Use(middleware.CORS(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	router.ServeHTTP(w, req)

	allowedHeaders := w.Header().Get("Access-Control-Allow-Headers")
	
	// Check that important headers are allowed
	assert.Contains(t, allowedHeaders, "Content-Type")
	assert.Contains(t, allowedHeaders, "Authorization")
	assert.Contains(t, allowedHeaders, "X-CSRF-Token")
}

func TestCORS_AllowedMethods(t *testing.T) {
	cfg := createTestConfigWithOrigins([]string{"http://localhost:5173"})

	router := gin.New()
	router.Use(middleware.CORS(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	router.ServeHTTP(w, req)

	allowedMethods := w.Header().Get("Access-Control-Allow-Methods")

	// Check that common HTTP methods are allowed
	assert.Contains(t, allowedMethods, "GET")
	assert.Contains(t, allowedMethods, "POST")
	assert.Contains(t, allowedMethods, "PUT")
	assert.Contains(t, allowedMethods, "DELETE")
	assert.Contains(t, allowedMethods, "OPTIONS")
	assert.Contains(t, allowedMethods, "PATCH")
}

func TestCORS_AllowCredentials(t *testing.T) {
	cfg := createTestConfigWithOrigins([]string{"http://localhost:5173"})

	router := gin.New()
	router.Use(middleware.CORS(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	router.ServeHTTP(w, req)

	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCORS_NoOriginHeader(t *testing.T) {
	cfg := createTestConfigWithOrigins([]string{"http://localhost:5173"})

	router := gin.New()
	router.Use(middleware.CORS(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	// No Origin header set
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Should not set CORS origin header when no origin in request
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_MultipleAllowedOrigins(t *testing.T) {
	cfg := createTestConfigWithOrigins([]string{
		"http://localhost:3000",
		"http://localhost:5173",
		"https://example.com",
	})

	testCases := []struct {
		name           string
		origin         string
		expectedOrigin string
	}{
		{"First origin", "http://localhost:3000", "http://localhost:3000"},
		{"Second origin", "http://localhost:5173", "http://localhost:5173"},
		{"Third origin", "https://example.com", "https://example.com"},
		{"Disallowed origin", "http://hacker.com", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.Use(middleware.CORS(cfg))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Origin", tc.origin)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedOrigin, w.Header().Get("Access-Control-Allow-Origin"))
		})
	}
}

// =============================================================================
// HTTP Methods Tests
// =============================================================================

func TestCORS_AllHTTPMethods(t *testing.T) {
	cfg := createTestConfigWithOrigins([]string{"http://localhost:5173"})

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run("Method_"+method, func(t *testing.T) {
			router := gin.New()
			router.Use(middleware.CORS(cfg))
			router.Handle(method, "/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(method, "/test", nil)
			req.Header.Set("Origin", "http://localhost:5173")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
		})
	}
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestCORS_EmptyAllowedOrigins(t *testing.T) {
	cfg := createTestConfigWithOrigins([]string{})

	router := gin.New()
	router.Use(middleware.CORS(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_OriginWithTrailingSlash(t *testing.T) {
	cfg := createTestConfigWithOrigins([]string{"http://localhost:5173"})

	router := gin.New()
	router.Use(middleware.CORS(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173/") // With trailing slash
	router.ServeHTTP(w, req)

	// Won't match exact origin
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_CaseSensitiveOrigin(t *testing.T) {
	cfg := createTestConfigWithOrigins([]string{"http://localhost:5173"})

	router := gin.New()
	router.Use(middleware.CORS(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "HTTP://LOCALHOST:5173") // Different case
	router.ServeHTTP(w, req)

	// Origin comparison is case-sensitive
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

// =============================================================================
// Next Handler Execution Tests
// =============================================================================

func TestCORS_NextHandlerExecuted(t *testing.T) {
	cfg := createTestConfigWithOrigins([]string{"http://localhost:5173"})
	handlerCalled := false

	router := gin.New()
	router.Use(middleware.CORS(cfg))
	router.GET("/test", func(c *gin.Context) {
		handlerCalled = true
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	router.ServeHTTP(w, req)

	assert.True(t, handlerCalled, "Next handler should be called for non-OPTIONS requests")
}

func TestCORS_OptionsDoesNotCallNext(t *testing.T) {
	cfg := createTestConfigWithOrigins([]string{"http://localhost:5173"})
	handlerCalled := false

	router := gin.New()
	router.Use(middleware.CORS(cfg))
	router.POST("/test", func(c *gin.Context) {
		handlerCalled = true
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	router.ServeHTTP(w, req)

	assert.False(t, handlerCalled, "Next handler should NOT be called for OPTIONS requests")
	assert.Equal(t, http.StatusNoContent, w.Code)
}

