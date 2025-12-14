package testing

import (
	"testing"
	"time"

	"ganttpro-backend/models"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Status Constants Tests
// =============================================================================

func TestOperationPlanStatusConstants(t *testing.T) {
	// Verify status constants are defined correctly
	assert.Equal(t, "draft", models.StatusDraft)
	assert.Equal(t, "pending_approval", models.StatusPendingApproval)
	assert.Equal(t, "approved", models.StatusApproved)
}

func TestStatusWorkflow(t *testing.T) {
	// Document the expected workflow: draft -> pending_approval -> approved
	statuses := []string{models.StatusDraft, models.StatusPendingApproval, models.StatusApproved}
	
	// Verify all statuses are different
	assert.NotEqual(t, statuses[0], statuses[1])
	assert.NotEqual(t, statuses[1], statuses[2])
	assert.NotEqual(t, statuses[0], statuses[2])
}

// =============================================================================
// OperationPlan Model Tests
// =============================================================================

func TestOperationPlan_TableName(t *testing.T) {
	plan := models.OperationPlan{}
	assert.Equal(t, "operation_plans", plan.TableName())
}

func TestOperationPlan_ToResponse_BasicFields(t *testing.T) {
	// Arrange
	now := time.Now()
	startTime := now.Add(-2 * time.Hour)
	finishTime := now
	
	plan := &models.OperationPlan{
		ID:           1,
		PlanNumber:   "OP-20241210-001",
		JobOrderID:   10,
		MachineID:    5,
		PartQuantity: 100,
		Description:  "Manufacturing plan for Base Plate",
		Status:       models.StatusApproved,
		CreatedBy:    3,
		StartTime:    &startTime,
		FinishTime:   &finishTime,
		CreatedAt:    now.Add(-24 * time.Hour),
		UpdatedAt:    now,
	}

	// Act
	response := plan.ToResponse()

	// Assert
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "OP-20241210-001", response.PlanNumber)
	assert.Equal(t, uint(10), response.JobOrderID)
	assert.Equal(t, uint(5), response.MachineID)
	assert.Equal(t, 100, response.PartQuantity)
	assert.Equal(t, "Manufacturing plan for Base Plate", response.Description)
	assert.Equal(t, models.StatusApproved, response.Status)
	assert.Equal(t, uint(3), response.CreatedBy)
	assert.Equal(t, &startTime, response.StartTime)
	assert.Equal(t, &finishTime, response.FinishTime)
}

func TestOperationPlan_ToResponse_WithJobOrder(t *testing.T) {
	// Arrange
	jobOrder := &models.JobOrder{
		ID:      10,
		NJO:     "NJO-2024-001",
		Project: "Mold Base Type A",
	}
	
	plan := &models.OperationPlan{
		ID:         1,
		PlanNumber: "OP-20241210-001",
		JobOrderID: 10,
		JobOrder:   jobOrder,
		Status:     models.StatusDraft,
	}

	// Act
	response := plan.ToResponse()

	// Assert
	assert.NotNil(t, response.JobOrder)
	assert.Equal(t, int64(10), response.JobOrder.ID)
	assert.Equal(t, "NJO-2024-001", response.JobOrder.NJO)
}

func TestOperationPlan_ToResponse_WithMachine(t *testing.T) {
	// Arrange
	machine := &models.Machine{
		ID:          5,
		MachineCode: "YSD01",
		MachineName: "Yasda YSD01",
	}
	
	plan := &models.OperationPlan{
		ID:        1,
		MachineID: 5,
		Machine:   machine,
		Status:    models.StatusDraft,
	}

	// Act
	response := plan.ToResponse()

	// Assert
	assert.NotNil(t, response.Machine)
	assert.Equal(t, int64(5), response.Machine.ID)
	assert.Equal(t, "YSD01", response.Machine.MachineCode)
}

func TestOperationPlan_ToResponse_WithCreator(t *testing.T) {
	// Arrange
	creator := &models.User{
		ID:       3,
		Username: "PEM_USER",
		UserID:   "PI0824.0003",
		Role:     models.RolePEM,
		IsActive: true,
	}
	
	plan := &models.OperationPlan{
		ID:        1,
		CreatedBy: 3,
		Creator:   creator,
		Status:    models.StatusDraft,
	}

	// Act
	response := plan.ToResponse()

	// Assert
	assert.NotNil(t, response.Creator)
	assert.Equal(t, uint(3), response.Creator.ID)
	assert.Equal(t, "PEM_USER", response.Creator.Username)
	assert.Equal(t, models.RolePEM, response.Creator.Role)
}

func TestOperationPlan_ToResponse_NilRelations(t *testing.T) {
	// Arrange
	plan := &models.OperationPlan{
		ID:         1,
		PlanNumber: "OP-20241210-001",
		JobOrder:   nil,
		Machine:    nil,
		Creator:    nil,
		Status:     models.StatusDraft,
	}

	// Act
	response := plan.ToResponse()

	// Assert - nil relations should remain nil in response
	assert.Nil(t, response.JobOrder)
	assert.Nil(t, response.Machine)
	assert.Nil(t, response.Creator)
}

func TestOperationPlan_ToResponse_NilTimes(t *testing.T) {
	// Arrange
	plan := &models.OperationPlan{
		ID:         1,
		PlanNumber: "OP-20241210-001",
		StartTime:  nil,
		FinishTime: nil,
		Status:     models.StatusDraft,
	}

	// Act
	response := plan.ToResponse()

	// Assert
	assert.Nil(t, response.StartTime)
	assert.Nil(t, response.FinishTime)
}

// =============================================================================
// OperationPlanResponse Tests
// =============================================================================

func TestOperationPlanResponse_AllStatuses(t *testing.T) {
	testCases := []struct {
		name   string
		status string
	}{
		{"Draft", models.StatusDraft},
		{"PendingApproval", models.StatusPendingApproval},
		{"Approved", models.StatusApproved},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			plan := &models.OperationPlan{
				ID:     1,
				Status: tc.status,
			}
			response := plan.ToResponse()
			assert.Equal(t, tc.status, response.Status)
		})
	}
}

// =============================================================================
// Plan Number Format Tests
// =============================================================================

func TestPlanNumberForma(t *testing.T) {
	// Plan number format: OP-YYYYMMDD-XXX
	validPlanNumbers := []string{
		"OP-20241210-001",
		"OP-20241210-999",
		"OP-20240101-001",
		"OP-20251231-100",
	}

	for _, planNumber := range validPlanNumbers {
		t.Run(planNumber, func(t *testing.T) {
			plan := &models.OperationPlan{
				PlanNumber: planNumber,
			}
			response := plan.ToResponse()
			assert.Equal(t, planNumber, response.PlanNumber)
			assert.Len(t, planNumber, 17, "Plan number should be 17 characters: OP-YYYYMMDD-XXX")
		})
	}
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestOperationPlan_ToResponse_EmptyPlan(t *testing.T) {
	// Arrange
	plan := &models.OperationPlan{}

	// Act
	response := plan.ToResponse()

	// Assert - should handle zero values gracefully
	assert.Equal(t, uint(0), response.ID)
	assert.Empty(t, response.PlanNumber)
	assert.Equal(t, uint(0), response.JobOrderID)
	assert.Equal(t, uint(0), response.MachineID)
	assert.Equal(t, 0, response.PartQuantity)
	assert.Empty(t, response.Description)
	assert.Empty(t, response.Status)
}

func TestOperationPlan_ToResponse_ZeroPartQuantity(t *testing.T) {
	// Arrange
	plan := &models.OperationPlan{
		ID:           1,
		PartQuantity: 0, // Edge case: zero quantity
		Status:       models.StatusDraft,
	}

	// Act
	response := plan.ToResponse()

	// Assert
	assert.Equal(t, 0, response.PartQuantity)
}

func TestOperationPlan_ToResponse_LargePartQuantity(t *testing.T) {
	// Arrange
	plan := &models.OperationPlan{
		ID:           1,
		PartQuantity: 1000000, // Large quantity
		Status:       models.StatusDraft,
	}

	// Act
	response := plan.ToResponse()

	// Assert
	assert.Equal(t, 1000000, response.PartQuantity)
}

