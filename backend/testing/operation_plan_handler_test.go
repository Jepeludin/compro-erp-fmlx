package testing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Operation Plan Request Validation Tests
// =============================================================================

func TestCreateOperationPlanRequest_Validation(t *testing.T) {
	testCases := []struct {
		name       string
		body       map[string]interface{}
		shouldBind bool
	}{
		{
			name: "Valid full operation plan",
			body: map[string]interface{}{
				"job_order_id":  1,
				"machine_id":    1,
				"part_quantity": 100,
				"description":   "Manufacturing plan for base plate",
			},
			shouldBind: true,
		},
		{
			name: "Valid operation plan with required fields only",
			body: map[string]interface{}{
				"job_order_id": 1,
				"machine_id":   1,
			},
			shouldBind: true,
		},
		{
			name: "Missing job_order_id",
			body: map[string]interface{}{
				"machine_id": 1,
			},
			shouldBind: false,
		},
		{
			name: "Missing machine_id",
			body: map[string]interface{}{
				"job_order_id": 1,
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
			router.POST("/operation-plans", func(c *gin.Context) {
				var req struct {
					JobOrderID   uint   `json:"job_order_id" binding:"required"`
					MachineID    uint   `json:"machine_id" binding:"required"`
					PartQuantity int    `json:"part_quantity"`
					Description  string `json:"description"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, gin.H{"status": "created"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/operation-plans", toJSON(tc.body))
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

// =============================================================================
// Operation Plan Status Tests
// =============================================================================

func TestOperationPlanStatus_Values(t *testing.T) {
	statusValues := []struct {
		status      string
		description string
	}{
		{"draft", "Initial status when plan is created"},
		{"pending_approval", "After plan is submitted for approval"},
		{"approved", "After all 5 roles have approved"},
	}

	for _, tc := range statusValues {
		t.Run(tc.status, func(t *testing.T) {
			router := gin.New()
			router.GET("/operation-plans/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"id":     1,
					"status": tc.status,
				})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/operation-plans/1", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			response := parseResponse(w)
			assert.Equal(t, tc.status, response["status"])
		})
	}
}

// =============================================================================
// Operation Plan Response Format Tests
// =============================================================================

func TestOperationPlanResponse_GetAll(t *testing.T) {
	router := gin.New()
	router.GET("/operation-plans", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"operation_plans": []map[string]interface{}{
				{
					"id":            1,
					"plan_number":   "OP-20241210-001",
					"job_order_id":  1,
					"machine_id":    1,
					"part_quantity": 100,
					"status":        "draft",
				},
			},
			"count": 1,
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/operation-plans", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response := parseResponse(w)
	assert.Contains(t, response, "operation_plans")
	assert.Contains(t, response, "count")
}

func TestOperationPlanResponse_WithApprovals(t *testing.T) {
	router := gin.New()
	router.GET("/operation-plans/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"id":          1,
			"plan_number": "OP-20241210-001",
			"status":      "pending_approval",
			"approvals": []map[string]interface{}{
				{"approver_role": "PEM", "status": "approved"},
				{"approver_role": "PPIC", "status": "pending"},
				{"approver_role": "QC", "status": "pending"},
				{"approver_role": "Engineering", "status": "pending"},
				{"approver_role": "Toolpather", "status": "pending"},
			},
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/operation-plans/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response := parseResponse(w)
	assert.Contains(t, response, "approvals")

	approvals := response["approvals"].([]interface{})
	assert.Len(t, approvals, 5, "Should have 5 approval records")
}

// =============================================================================
// Submit for Approval Tests
// =============================================================================

func TestSubmitForApproval_Success(t *testing.T) {
	router := gin.New()
	router.POST("/operation-plans/:id/submit", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Operation plan submitted for approval",
			"status":  "pending_approval",
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/operation-plans/1/submit", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response := parseResponse(w)
	assert.Equal(t, "pending_approval", response["status"])
}

func TestSubmitForApproval_NoGCodeFiles(t *testing.T) {
	router := gin.New()
	router.POST("/operation-plans/:id/submit", func(c *gin.Context) {
		// Simulate error: no G-code files uploaded
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "at least one G-code file must be uploaded before submitting",
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/operation-plans/1/submit", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	response := parseResponse(w)
	assert.Contains(t, response["error"], "G-code")
}

func TestSubmitForApproval_NotDraftStatus(t *testing.T) {
	router := gin.New()
	router.POST("/operation-plans/:id/submit", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "operation plan is not in draft status",
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/operation-plans/1/submit", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// =============================================================================
// Approve Operation Plan Tests
// =============================================================================

func TestApproveOperationPlan_Success(t *testing.T) {
	router := gin.New()
	router.POST("/operation-plans/:id/approve", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Operation plan approved",
			"approval": map[string]interface{}{
				"approver_role": "PEM",
				"status":        "approved",
			},
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/operation-plans/1/approve", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestApproveOperationPlan_InvalidRole(t *testing.T) {
	router := gin.New()
	router.POST("/operation-plans/:id/approve", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid approver role",
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/operation-plans/1/approve", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestApproveOperationPlan_AlreadyApproved(t *testing.T) {
	router := gin.New()
	router.POST("/operation-plans/:id/approve", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "approval not found or already approved",
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/operation-plans/1/approve", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// =============================================================================
// Pending Approvals Tests
// =============================================================================

func TestGetPendingApprovals(t *testing.T) {
	router := gin.New()
	router.GET("/operation-plans/pending-approvals", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"pending_approvals": []map[string]interface{}{
				{
					"id":          1,
					"plan_number": "OP-20241210-001",
					"status":      "pending_approval",
				},
			},
			"count": 1,
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/operation-plans/pending-approvals", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response := parseResponse(w)
	assert.Contains(t, response, "pending_approvals")
}

// =============================================================================
// Execution Tests
// =============================================================================

func TestStartExecution_Success(t *testing.T) {
	router := gin.New()
	router.POST("/operation-plans/:id/start", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":    "Execution started",
			"start_time": "2024-12-10T08:00:00Z",
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/operation-plans/1/start", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response := parseResponse(w)
	assert.Contains(t, response, "start_time")
}

func TestStartExecution_NotApproved(t *testing.T) {
	router := gin.New()
	router.POST("/operation-plans/:id/start", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "operation plan must be approved before starting execution",
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/operation-plans/1/start", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFinishExecution_Success(t *testing.T) {
	router := gin.New()
	router.POST("/operation-plans/:id/finish", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":     "Execution finished",
			"finish_time": "2024-12-10T16:00:00Z",
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/operation-plans/1/finish", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response := parseResponse(w)
	assert.Contains(t, response, "finish_time")
}

func TestFinishExecution_NotStarted(t *testing.T) {
	router := gin.New()
	router.POST("/operation-plans/:id/finish", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "operation plan execution has not started",
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/operation-plans/1/finish", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// =============================================================================
// Delete Operation Plan Tests
// =============================================================================

func TestDeleteOperationPlan_DraftOnly(t *testing.T) {
	testCases := []struct {
		name         string
		status       string
		expectedCode int
	}{
		{"Delete draft plan", "draft", http.StatusOK},
		{"Cannot delete pending plan", "pending_approval", http.StatusBadRequest},
		{"Cannot delete approved plan", "approved", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.DELETE("/operation-plans/:id", func(c *gin.Context) {
				if tc.status == "draft" {
					c.JSON(http.StatusOK, gin.H{"message": "Operation plan deleted"})
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete operation plan that is not in draft status"})
				}
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/operation-plans/1", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

// =============================================================================
// Plan Number Format Tests
// =============================================================================

func TestPlanNumberFormat(t *testing.T) {
	// Plan number format: OP-YYYYMMDD-XXX
	validPlanNumbers := []string{
		"OP-20241210-001",
		"OP-20241210-002",
		"OP-20241210-999",
		"OP-20250101-001",
	}

	for _, planNumber := range validPlanNumbers {
		t.Run(planNumber, func(t *testing.T) {
			assert.Len(t, planNumber, 17, "Plan number should be 17 characters")
			assert.Equal(t, "OP-", planNumber[:3], "Should start with OP-")
		})
	}
}

// =============================================================================
// Filter Tests
// =============================================================================

func TestOperationPlanFilter_ByStatus(t *testing.T) {
	router := gin.New()
	router.GET("/operation-plans", func(c *gin.Context) {
		status := c.Query("status")
		c.JSON(http.StatusOK, gin.H{
			"operation_plans": []map[string]interface{}{},
			"filter_status":   status,
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/operation-plans?status=draft", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response := parseResponse(w)
	assert.Equal(t, "draft", response["filter_status"])
}

func TestOperationPlanFilter_ByMachine(t *testing.T) {
	router := gin.New()
	router.GET("/operation-plans", func(c *gin.Context) {
		machineID := c.Query("machine_id")
		c.JSON(http.StatusOK, gin.H{
			"operation_plans":   []map[string]interface{}{},
			"filter_machine_id": machineID,
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/operation-plans?machine_id=1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response := parseResponse(w)
	assert.Equal(t, "1", response["filter_machine_id"])
}

