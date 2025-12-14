package testing

import (
	"encoding/json"
	"testing"
	"time"

	"ganttpro-backend/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// JobOrder Structure Tests
// =============================================================================

func TestJobOrder_Structure(t *testing.T) {
	now := time.Now()
	operatorID := int64(5)
	jobOrder := models.JobOrder{
		ID:           1,
		MachineID:    2,
		MachineName:  "Yasda YSD01",
		NJO:          "NJO-2024-001",
		Project:      "Mold Base Type A",
		Item:         "Base Plate",
		Note:         "Priority order",
		Deadline:     "2024-12-31",
		OperatorID:   &operatorID,
		OperatorName: "BAYU",
		Status:       "pending",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	assert.Equal(t, int64(1), jobOrder.ID)
	assert.Equal(t, int64(2), jobOrder.MachineID)
	assert.Equal(t, "Yasda YSD01", jobOrder.MachineName)
	assert.Equal(t, "NJO-2024-001", jobOrder.NJO)
	assert.Equal(t, "Mold Base Type A", jobOrder.Project)
	assert.Equal(t, "Base Plate", jobOrder.Item)
	assert.Equal(t, "Priority order", jobOrder.Note)
	assert.Equal(t, "2024-12-31", jobOrder.Deadline)
	assert.NotNil(t, jobOrder.OperatorID)
	assert.Equal(t, int64(5), *jobOrder.OperatorID)
	assert.Equal(t, "BAYU", jobOrder.OperatorName)
	assert.Equal(t, "pending", jobOrder.Status)
}

func TestJobOrder_JSONSerialization(t *testing.T) {
	jobOrder := models.JobOrder{
		ID:          1,
		MachineID:   2,
		MachineName: "Yasda YSD01",
		NJO:         "NJO-2024-001",
		Project:     "Test Project",
		Item:        "Test Item",
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	jsonData, err := json.Marshal(jobOrder)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	assert.Contains(t, result, "id")
	assert.Contains(t, result, "machine_id")
	assert.Contains(t, result, "njo")
	assert.Contains(t, result, "project")
	assert.Contains(t, result, "status")
}

func TestJobOrder_JSONDeserialization(t *testing.T) {
	jsonData := `{
		"id": 1,
		"machine_id": 2,
		"njo": "NJO-2024-002",
		"project": "Test Project",
		"item": "Test Item",
		"note": "Test note",
		"deadline": "2024-12-15",
		"status": "in_progress"
	}`

	var jobOrder models.JobOrder
	err := json.Unmarshal([]byte(jsonData), &jobOrder)
	require.NoError(t, err)

	assert.Equal(t, int64(1), jobOrder.ID)
	assert.Equal(t, int64(2), jobOrder.MachineID)
	assert.Equal(t, "NJO-2024-002", jobOrder.NJO)
	assert.Equal(t, "Test Project", jobOrder.Project)
	assert.Equal(t, "Test Item", jobOrder.Item)
	assert.Equal(t, "in_progress", jobOrder.Status)
}

func TestJobOrder_WithProcessStages(t *testing.T) {
	jobOrder := models.JobOrder{
		ID:     1,
		NJO:    "NJO-2024-001",
		Status: "in_progress",
		ProcessStages: []models.ProcessStage{
			{ID: 1, JobOrderID: 1, StageName: "setting"},
			{ID: 2, JobOrderID: 1, StageName: "proses"},
			{ID: 3, JobOrderID: 1, StageName: "cmm"},
			{ID: 4, JobOrderID: 1, StageName: "kalibrasi"},
		},
	}

	assert.Len(t, jobOrder.ProcessStages, 4)
	assert.Equal(t, "setting", jobOrder.ProcessStages[0].StageName)
	assert.Equal(t, "proses", jobOrder.ProcessStages[1].StageName)
	assert.Equal(t, "cmm", jobOrder.ProcessStages[2].StageName)
	assert.Equal(t, "kalibrasi", jobOrder.ProcessStages[3].StageName)
}

func TestJobOrder_CompletedAt(t *testing.T) {
	completedAt := time.Now()
	jobOrder := models.JobOrder{
		ID:          1,
		Status:      "completed",
		CompletedAt: &completedAt,
	}

	assert.NotNil(t, jobOrder.CompletedAt)
	assert.Equal(t, "completed", jobOrder.Status)
}

func TestJobOrder_NilOptionalFields(t *testing.T) {
	jobOrder := models.JobOrder{
		ID:     1,
		NJO:    "NJO-2024-001",
		Status: "pending",
	}

	assert.Nil(t, jobOrder.OperatorID)
	assert.Nil(t, jobOrder.CompletedAt)
	assert.Nil(t, jobOrder.DeletedAt)
}

// =============================================================================
// ProcessStage Structure Tests
// =============================================================================

func TestProcessStage_Structure(t *testing.T) {
	now := time.Now()
	startTime := now.Add(-2 * time.Hour)
	finishTime := now.Add(-1 * time.Hour)
	duration := float64(60)
	operatorID := int64(5)

	stage := models.ProcessStage{
		ID:              1,
		JobOrderID:      1,
		StageName:       "proses",
		StartTime:       &startTime,
		FinishTime:      &finishTime,
		DurationMinutes: &duration,
		OperatorID:      &operatorID,
		OperatorName:    "BAYU",
		Notes:           "Completed successfully",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	assert.Equal(t, int64(1), stage.ID)
	assert.Equal(t, int64(1), stage.JobOrderID)
	assert.Equal(t, "proses", stage.StageName)
	assert.NotNil(t, stage.StartTime)
	assert.NotNil(t, stage.FinishTime)
	assert.NotNil(t, stage.DurationMinutes)
	assert.Equal(t, float64(60), *stage.DurationMinutes)
	assert.Equal(t, "BAYU", stage.OperatorName)
}

func TestProcessStage_ValidStageNames(t *testing.T) {
	validStageNames := []string{"setting", "proses", "cmm", "kalibrasi"}

	for _, stageName := range validStageNames {
		t.Run("Stage_"+stageName, func(t *testing.T) {
			stage := models.ProcessStage{
				StageName: stageName,
			}
			assert.Equal(t, stageName, stage.StageName)
		})
	}
}

func TestProcessStage_JSONSerialization(t *testing.T) {
	startTime := time.Now()
	stage := models.ProcessStage{
		ID:         1,
		JobOrderID: 1,
		StageName:  "setting",
		StartTime:  &startTime,
		Notes:      "Started",
	}

	jsonData, err := json.Marshal(stage)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	assert.Contains(t, result, "id")
	assert.Contains(t, result, "job_order_id")
	assert.Contains(t, result, "stage_name")
	assert.Contains(t, result, "start_time")
}

// =============================================================================
// CreateJobOrderRequest Tests
// =============================================================================

func TestCreateJobOrderRequest_Structure(t *testing.T) {
	operatorID := int64(5)
	req := models.CreateJobOrderRequest{
		MachineID:  1,
		NJO:        "NJO-2024-001",
		Project:    "Test Project",
		Item:       "Test Item",
		Note:       "Test Note",
		Deadline:   "2024-12-31",
		OperatorID: &operatorID,
	}

	assert.Equal(t, int64(1), req.MachineID)
	assert.Equal(t, "NJO-2024-001", req.NJO)
	assert.Equal(t, "Test Project", req.Project)
	assert.NotNil(t, req.OperatorID)
	assert.Equal(t, int64(5), *req.OperatorID)
}

func TestCreateJobOrderRequest_JSONDeserialization(t *testing.T) {
	testCases := []struct {
		name     string
		jsonData string
		expected models.CreateJobOrderRequest
	}{
		{
			name: "Full request",
			jsonData: `{
				"machine_id": 1,
				"njo": "NJO-2024-001",
				"project": "Test Project",
				"item": "Test Item",
				"note": "Test Note",
				"deadline": "2024-12-31",
				"operator_id": 5
			}`,
			expected: models.CreateJobOrderRequest{
				MachineID: 1,
				NJO:       "NJO-2024-001",
				Project:   "Test Project",
				Item:      "Test Item",
				Note:      "Test Note",
				Deadline:  "2024-12-31",
			},
		},
		{
			name: "Minimal request",
			jsonData: `{
				"machine_id": 1,
				"njo": "NJO-2024-002"
			}`,
			expected: models.CreateJobOrderRequest{
				MachineID: 1,
				NJO:       "NJO-2024-002",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var req models.CreateJobOrderRequest
			err := json.Unmarshal([]byte(tc.jsonData), &req)
			require.NoError(t, err)

			assert.Equal(t, tc.expected.MachineID, req.MachineID)
			assert.Equal(t, tc.expected.NJO, req.NJO)
		})
	}
}

// =============================================================================
// UpdateJobOrderRequest Tests
// =============================================================================

func TestUpdateJobOrderRequest_Structure(t *testing.T) {
	operatorID := int64(5)
	req := models.UpdateJobOrderRequest{
		Project:    "Updated Project",
		Item:       "Updated Item",
		Note:       "Updated Note",
		Deadline:   "2024-12-25",
		OperatorID: &operatorID,
		Status:     "in_progress",
	}

	assert.Equal(t, "Updated Project", req.Project)
	assert.Equal(t, "Updated Item", req.Item)
	assert.Equal(t, "Updated Note", req.Note)
	assert.Equal(t, "2024-12-25", req.Deadline)
	assert.NotNil(t, req.OperatorID)
	assert.Equal(t, "in_progress", req.Status)
}

func TestUpdateJobOrderRequest_PartialUpdate(t *testing.T) {
	jsonData := `{"status": "completed"}`

	var req models.UpdateJobOrderRequest
	err := json.Unmarshal([]byte(jsonData), &req)
	require.NoError(t, err)

	assert.Equal(t, "completed", req.Status)
	assert.Empty(t, req.Project)
	assert.Empty(t, req.Item)
	assert.Nil(t, req.OperatorID)
}

// =============================================================================
// UpdateProcessStageRequest Tests
// =============================================================================

func TestUpdateProcessStageRequest_Structure(t *testing.T) {
	startTime := time.Now()
	finishTime := startTime.Add(2 * time.Hour)
	operatorID := int64(5)

	req := models.UpdateProcessStageRequest{
		StartTime:  &startTime,
		FinishTime: &finishTime,
		OperatorID: &operatorID,
		Notes:      "Stage completed",
	}

	assert.NotNil(t, req.StartTime)
	assert.NotNil(t, req.FinishTime)
	assert.NotNil(t, req.OperatorID)
	assert.Equal(t, "Stage completed", req.Notes)
}

func TestUpdateProcessStageRequest_PartialUpdate(t *testing.T) {
	startTime := time.Now()
	req := models.UpdateProcessStageRequest{
		StartTime: &startTime,
	}

	assert.NotNil(t, req.StartTime)
	assert.Nil(t, req.FinishTime)
	assert.Nil(t, req.OperatorID)
	assert.Empty(t, req.Notes)
}

// =============================================================================
// JobOrderWithStages Tests
// =============================================================================

func TestJobOrderWithStages_Structure(t *testing.T) {
	jobOrderWithStages := models.JobOrderWithStages{
		JobOrder: models.JobOrder{
			ID:     1,
			NJO:    "NJO-2024-001",
			Status: "in_progress",
		},
		Stages: []models.ProcessStage{
			{ID: 1, StageName: "setting"},
			{ID: 2, StageName: "proses"},
		},
	}

	assert.Equal(t, int64(1), jobOrderWithStages.ID)
	assert.Equal(t, "NJO-2024-001", jobOrderWithStages.NJO)
	assert.Len(t, jobOrderWithStages.Stages, 2)
}

// =============================================================================
// Status Tests
// =============================================================================

func TestJobOrder_ValidStatuses(t *testing.T) {
	validStatuses := []string{"pending", "in_progress", "completed", "on_hold", "cancelled"}

	for _, status := range validStatuses {
		t.Run("Status_"+status, func(t *testing.T) {
			jobOrder := models.JobOrder{Status: status}
			assert.Equal(t, status, jobOrder.Status)
		})
	}
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestJobOrder_EmptyFields(t *testing.T) {
	jobOrder := models.JobOrder{}

	assert.Equal(t, int64(0), jobOrder.ID)
	assert.Equal(t, int64(0), jobOrder.MachineID)
	assert.Empty(t, jobOrder.NJO)
	assert.Empty(t, jobOrder.Project)
	assert.Empty(t, jobOrder.Status)
	assert.Nil(t, jobOrder.ProcessStages)
}

func TestProcessStage_EmptyFields(t *testing.T) {
	stage := models.ProcessStage{}

	assert.Equal(t, int64(0), stage.ID)
	assert.Empty(t, stage.StageName)
	assert.Nil(t, stage.StartTime)
	assert.Nil(t, stage.FinishTime)
	assert.Nil(t, stage.DurationMinutes)
}

func TestJobOrder_SpecialCharactersInNJO(t *testing.T) {
	jobOrder := models.JobOrder{
		NJO: "NJO-2024/001-A",
	}

	assert.Equal(t, "NJO-2024/001-A", jobOrder.NJO)
}

func TestJobOrder_LongNote(t *testing.T) {
	longNote := "This is a very long note that might contain detailed instructions about the job order including specifications, materials needed, and any special requirements that the operator should be aware of."
	jobOrder := models.JobOrder{
		Note: longNote,
	}

	assert.Equal(t, longNote, jobOrder.Note)
}

