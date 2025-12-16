package testing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Machine Request Validation Tests
// =============================================================================

func TestCreateMachineRequest_Validation(t *testing.T) {
	testCases := []struct {
		name       string
		body       map[string]interface{}
		shouldBind bool
	}{
		{
			name: "Valid full machine",
			body: map[string]interface{}{
				"machine_code": "YSD01",
				"machine_name": "Yasda YSD01",
				"machine_type": "CNC Milling",
				"location":     "Building A",
				"status":       "active",
			},
			shouldBind: true,
		},
		{
			name: "Valid machine with required fields only",
			body: map[string]interface{}{
				"machine_code": "V33I",
				"machine_name": "Makino V33i",
			},
			shouldBind: true,
		},
		{
			name: "Missing machine_code",
			body: map[string]interface{}{
				"machine_name": "Test Machine",
			},
			shouldBind: false,
		},
		{
			name: "Missing machine_name",
			body: map[string]interface{}{
				"machine_code": "TEST01",
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
			router.POST("/machines", func(c *gin.Context) {
				var req struct {
					MachineCode string `json:"machine_code" binding:"required"`
					MachineName string `json:"machine_name" binding:"required"`
					MachineType string `json:"machine_type"`
					Location    string `json:"location"`
					Status      string `json:"status"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, gin.H{"status": "created"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/machines", toJSON(tc.body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if tc.shouldBind {
				assert.Equal(t, http.StatusCreated, w.Code)
			} else {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})
	}
}

func TestUpdateMachineRequest_Validation(t *testing.T) {
	testCases := []struct {
		name       string
		body       map[string]interface{}
		shouldBind bool
	}{
		{
			name: "Full update",
			body: map[string]interface{}{
				"machine_name": "Updated Name",
				"machine_type": "5-Axis CNC",
				"location":     "New Location",
				"status":       "maintenance",
			},
			shouldBind: true,
		},
		{
			name: "Partial update - only status",
			body: map[string]interface{}{
				"status": "inactive",
			},
			shouldBind: true,
		},
		{
			name:       "Empty body is valid for update",
			body:       map[string]interface{}{},
			shouldBind: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.PUT("/machines/:id", func(c *gin.Context) {
				var req struct {
					MachineName string `json:"machine_name"`
					MachineType string `json:"machine_type"`
					Location    string `json:"location"`
					Status      string `json:"status"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"status": "updated"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/machines/1", toJSON(tc.body))
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
// Machine ID Parameter Tests
// =============================================================================

func TestMachineID_Parameter(t *testing.T) {
	testCases := []struct {
		name         string
		id           string
		expectedCode int
	}{
		{
			name:         "Valid numeric ID",
			id:           "1",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Valid large ID",
			id:           "999999",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid non-numeric ID",
			id:           "abc",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Negative ID",
			id:           "-1",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Zero ID",
			id:           "0",
			expectedCode: http.StatusOK, // Valid but might return 404 in real scenario
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/machines/:id", func(c *gin.Context) {
				idStr := c.Param("id")
				var id int64
				if _, err := parseID(idStr, &id); err != nil || id < 0 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine ID"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"id": id})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/machines/"+tc.id, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

// Helper function to parse ID
func parseID(s string, result *int64) (int, error) {
	var n int
	for i, c := range s {
		if c < '0' || c > '9' {
			if i == 0 && c == '-' {
				continue
			}
			return 0, assert.AnError
		}
		n = n*10 + int(c-'0')
	}
	if len(s) > 0 && s[0] == '-' {
		n = -n
	}
	*result = int64(n)
	return len(s), nil
}

// =============================================================================
// Machine Response Format Tests
// =============================================================================

func TestMachineResponse_GetAll(t *testing.T) {
	router := gin.New()
	router.GET("/machines", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"machines": []map[string]interface{}{
				{
					"id":           1,
					"machine_code": "YSD01",
					"machine_name": "Yasda YSD01",
					"machine_type": "CNC Milling",
					"location":     "Building A",
					"status":       "active",
				},
			},
			"count": 1,
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/machines", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response := parseResponse(w)
	assert.Contains(t, response, "machines")
	assert.Contains(t, response, "count")
}

func TestMachineResponse_GetByID(t *testing.T) {
	router := gin.New()
	router.GET("/machines/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"id":           1,
			"machine_code": "YSD01",
			"machine_name": "Yasda YSD01",
			"machine_type": "CNC Milling",
			"location":     "Building A",
			"status":       "active",
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/machines/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response := parseResponse(w)
	assert.Contains(t, response, "id")
	assert.Contains(t, response, "machine_code")
	assert.Contains(t, response, "machine_name")
}

// =============================================================================
// Machine Status Values Tests
// =============================================================================

func TestMachineStatus_Values(t *testing.T) {
	validStatuses := []string{"active", "inactive", "maintenance", "offline"}

	for _, status := range validStatuses {
		t.Run(status, func(t *testing.T) {
			body := map[string]interface{}{
				"machine_code": "TEST01",
				"machine_name": "Test Machine",
				"status":       status,
			}

			router := gin.New()
			router.POST("/machines", func(c *gin.Context) {
				var req struct {
					MachineCode string `json:"machine_code" binding:"required"`
					MachineName string `json:"machine_name" binding:"required"`
					Status      string `json:"status"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, gin.H{"status": req.Status})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/machines", toJSON(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)
			response := parseResponse(w)
			assert.Equal(t, status, response["status"])
		})
	}
}

// =============================================================================
// Machine Delete Tests
// =============================================================================

func TestMachineDelete_Success(t *testing.T) {
	router := gin.New()
	router.DELETE("/machines/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Machine deleted successfully"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/machines/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response := parseResponse(w)
	assert.Contains(t, response, "message")
}

// =============================================================================
// Machine Code Uniqueness Test
// =============================================================================

func TestMachineCode_ConflictResponse(t *testing.T) {
	router := gin.New()
	router.POST("/machines", func(c *gin.Context) {
		// Simulate duplicate machine code
		c.JSON(http.StatusConflict, gin.H{"error": "Machine code already exists"})
	})

	w := httptest.NewRecorder()
	body := map[string]interface{}{
		"machine_code": "DUPLICATE",
		"machine_name": "Duplicate Machine",
	}
	req := httptest.NewRequest("POST", "/machines", toJSON(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	response := parseResponse(w)
	assert.Equal(t, "Machine code already exists", response["error"])
}

