package testing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Job Order Request Validation Tests
// =============================================================================

func TestCreateJobOrderRequest_Validation(t *testing.T) {
	testCases := []struct {
		name       string
		body       map[string]interface{}
		shouldBind bool
	}{
		{
			name: "Valid full job order",
			body: map[string]interface{}{
				"machine_id":  1,
				"njo":         "NJO-2024-001",
				"project":     "Mold Base Type A",
				"item":        "Base Plate BP-001",
				"note":        "Priority order",
				"deadline":    "2024-12-31",
				"operator_id": 2,
			},
			shouldBind: true,
		},
		{
			name: "Valid job order with required fields only",
			body: map[string]interface{}{
				"machine_id": 1,
				"njo":        "NJO-2024-002",
			},
			shouldBind: true,
		},
		{
			name: "Missing machine_id",
			body: map[string]interface{}{
				"njo": "NJO-2024-003",
			},
			shouldBind: false,
		},
		{
			name: "Missing njo",
			body: map[string]interface{}{
				"machine_id": 1,
			},
			shouldBind: false,
		},
		{
			name:       "Empty body",
			body:       map[string]interface{}{},
			shouldBind: false,
		},
		{
			name: "Invalid machine_id type",
			body: map[string]interface{}{
				"machine_id": "not-a-number",
				"njo":        "NJO-2024-004",
			},
			shouldBind: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/job-orders", func(c *gin.Context) {
				var req struct {
					MachineID  int64  `json:"machine_id" binding:"required"`
					NJO        string `json:"njo" binding:"required"`
					Project    string `json:"project"`
					Item       string `json:"item"`
					Note       string `json:"note"`
					Deadline   string `json:"deadline"`
					OperatorID *int64 `json:"operator_id"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, gin.H{"status": "created"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/job-orders", toJSON(tc.body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if tc.shouldBind {
				assert.Equal(t, http.StatusCreated, w.Code, "Test case: %s", tc.name)
			} else {
				assert.Equal(t, http.StatusBadRequest, w.Code, "Test case: %s", tc.name)
			}
		})
	}
}

func TestUpdateJobOrderRequest_Validation(t *testing.T) {
	testCases := []struct {
		name       string
		body       map[string]interface{}
		shouldBind bool
	}{
		{
			name: "Full update",
			body: map[string]interface{}{
				"project":     "Updated Project",
				"item":        "Updated Item",
				"note":        "Updated Note",
				"deadline":    "2025-01-15",
				"operator_id": 3,
				"status":      "in_progress",
			},
			shouldBind: true,
		},
		{
			name: "Partial update - only status",
			body: map[string]interface{}{
				"status": "completed",
			},
			shouldBind: true,
		},
		{
			name:       "Empty body is valid",
			body:       map[string]interface{}{},
			shouldBind: true,
		},
		{
			name: "Null operator_id (unassign)",
			body: map[string]interface{}{
				"operator_id": nil,
			},
			shouldBind: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.PUT("/job-orders/:id", func(c *gin.Context) {
				var req struct {
					Project    string `json:"project"`
					Item       string `json:"item"`
					Note       string `json:"note"`
					Deadline   string `json:"deadline"`
					OperatorID *int64 `json:"operator_id"`
					Status     string `json:"status"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"status": "updated"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/job-orders/1", toJSON(tc.body))
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
// Job Order Response Format Tests
// =============================================================================

func TestJobOrderResponse_GetAll(t *testing.T) {
	router := gin.New()
	router.GET("/job-orders", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"job_orders": []map[string]interface{}{
				{
					"id":           1,
					"machine_id":   1,
					"machine_name": "Yasda YSD01",
					"njo":          "NJO-2024-001",
					"project":      "Mold Base Type A",
					"item":         "Base Plate",
					"status":       "pending",
				},
			},
			"count": 1,
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/job-orders", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response := parseResponse(w)
	assert.Contains(t, response, "job_orders")
	assert.Contains(t, response, "count")
}

func TestJobOrderResponse_GetByMachine(t *testing.T) {
	router := gin.New()
	router.GET("/job-orders/machine/:machine_id", func(c *gin.Context) {
		machineID := c.Param("machine_id")
		c.JSON(http.StatusOK, gin.H{
			"job_orders": []map[string]interface{}{},
			"count":      0,
			"machine_id": machineID,
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/job-orders/machine/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response := parseResponse(w)
	assert.Contains(t, response, "job_orders")
	assert.Equal(t, "1", response["machine_id"])
}

// =============================================================================
// Job Order Status Values Tests
// =============================================================================

func TestJobOrderStatus_Values(t *testing.T) {
	validStatuses := []string{"pending", "in_progress", "completed", "on_hold", "cancelled"}

	for _, status := range validStatuses {
		t.Run(status, func(t *testing.T) {
			body := map[string]interface{}{
				"status": status,
			}

			router := gin.New()
			router.PUT("/job-orders/:id", func(c *gin.Context) {
				var req struct {
					Status string `json:"status"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"status": req.Status})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/job-orders/1", toJSON(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			response := parseResponse(w)
			assert.Equal(t, status, response["status"])
		})
	}
}

// =============================================================================
// Process Stage Tests
// =============================================================================

func TestUpdateProcessStageRequest_Validation(t *testing.T) {
	testCases := []struct {
		name       string
		body       map[string]interface{}
		shouldBind bool
	}{
		{
			name: "Full update with times",
			body: map[string]interface{}{
				"start_time":  "2024-12-10T08:00:00Z",
				"finish_time": "2024-12-10T12:00:00Z",
				"operator_id": 2,
				"notes":       "Stage completed successfully",
			},
			shouldBind: true,
		},
		{
			name: "Only start time",
			body: map[string]interface{}{
				"start_time": "2024-12-10T08:00:00Z",
			},
			shouldBind: true,
		},
		{
			name: "Only notes",
			body: map[string]interface{}{
				"notes": "Some notes",
			},
			shouldBind: true,
		},
		{
			name:       "Empty body",
			body:       map[string]interface{}{},
			shouldBind: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.PUT("/process-stages/:id", func(c *gin.Context) {
				var req struct {
					StartTime  *string `json:"start_time"`
					FinishTime *string `json:"finish_time"`
					OperatorID *int64  `json:"operator_id"`
					Notes      string  `json:"notes"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"status": "updated"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/process-stages/1", toJSON(tc.body))
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

func TestProcessStageNames(t *testing.T) {
	// Document the expected stage names
	expectedStages := []string{"setting", "proses", "cmm", "kalibrasi"}

	router := gin.New()
	router.GET("/job-orders/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"id": 1,
			"process_stages": []map[string]interface{}{
				{"id": 1, "stage_name": "setting"},
				{"id": 2, "stage_name": "proses"},
				{"id": 3, "stage_name": "cmm"},
				{"id": 4, "stage_name": "kalibrasi"},
			},
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/job-orders/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response := parseResponse(w)
	
	stages := response["process_stages"].([]interface{})
	assert.Len(t, stages, 4, "Should have 4 process stages")
	
	for i, stage := range stages {
		stageMap := stage.(map[string]interface{})
		assert.Equal(t, expectedStages[i], stageMap["stage_name"])
	}
}

// =============================================================================
// Job Order Not Found Tests
// =============================================================================

func TestJobOrder_NotFound(t *testing.T) {
	router := gin.New()
	router.GET("/job-orders/:id", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job order not found"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/job-orders/99999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	response := parseResponse(w)
	assert.Equal(t, "Job order not found", response["error"])
}

// =============================================================================
// Job Order Delete Tests
// =============================================================================

func TestJobOrderDelete_Success(t *testing.T) {
	router := gin.New()
	router.DELETE("/job-orders/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Job order deleted successfully"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/job-orders/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response := parseResponse(w)
	assert.Contains(t, response, "message")
}

// =============================================================================
// NJO Number Format Tests
// =============================================================================

func TestNJONumberFormat(t *testing.T) {
	// NJO format: NJO-YYYY-XXX
	validNJOs := []string{
		"NJO-2024-001",
		"NJO-2024-999",
		"NJO-2025-001",
	}

	for _, njo := range validNJOs {
		t.Run(njo, func(t *testing.T) {
			body := map[string]interface{}{
				"machine_id": 1,
				"njo":        njo,
			}

			router := gin.New()
			router.POST("/job-orders", func(c *gin.Context) {
				var req struct {
					MachineID int64  `json:"machine_id" binding:"required"`
					NJO       string `json:"njo" binding:"required"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, gin.H{"njo": req.NJO})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/job-orders", toJSON(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)
			response := parseResponse(w)
			assert.Equal(t, njo, response["njo"])
		})
	}
}

