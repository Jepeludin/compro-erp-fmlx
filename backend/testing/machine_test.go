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
// Machine Structure Tests
// =============================================================================

func TestMachine_Structure(t *testing.T) {
	now := time.Now()
	machine := models.Machine{
		ID:          1,
		MachineCode: "YSD01",
		MachineName: "Yasda YSD01",
		MachineType: "CNC Milling",
		Location:    "Building A",
		Status:      "active",
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   nil,
	}

	assert.Equal(t, int64(1), machine.ID)
	assert.Equal(t, "YSD01", machine.MachineCode)
	assert.Equal(t, "Yasda YSD01", machine.MachineName)
	assert.Equal(t, "CNC Milling", machine.MachineType)
	assert.Equal(t, "Building A", machine.Location)
	assert.Equal(t, "active", machine.Status)
	assert.Nil(t, machine.DeletedAt)
}

func TestMachine_JSONSerialization(t *testing.T) {
	now := time.Now()
	machine := models.Machine{
		ID:          1,
		MachineCode: "YSD01",
		MachineName: "Yasda YSD01",
		MachineType: "CNC Milling",
		Location:    "Building A",
		Status:      "active",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(machine)
	require.NoError(t, err)

	// Verify JSON contains expected fields
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	assert.Contains(t, result, "id")
	assert.Contains(t, result, "machine_code")
	assert.Contains(t, result, "machine_name")
	assert.Contains(t, result, "machine_type")
	assert.Contains(t, result, "location")
	assert.Contains(t, result, "status")
	assert.Contains(t, result, "created_at")
	assert.Contains(t, result, "updated_at")
}

func TestMachine_JSONDeserialization(t *testing.T) {
	jsonData := `{
		"id": 1,
		"machine_code": "V33I",
		"machine_name": "Makino V33i",
		"machine_type": "5-Axis CNC",
		"location": "Building B",
		"status": "maintenance"
	}`

	var machine models.Machine
	err := json.Unmarshal([]byte(jsonData), &machine)
	require.NoError(t, err)

	assert.Equal(t, int64(1), machine.ID)
	assert.Equal(t, "V33I", machine.MachineCode)
	assert.Equal(t, "Makino V33i", machine.MachineName)
	assert.Equal(t, "5-Axis CNC", machine.MachineType)
	assert.Equal(t, "Building B", machine.Location)
	assert.Equal(t, "maintenance", machine.Status)
}

func TestMachine_WithDeletedAt(t *testing.T) {
	now := time.Now()
	deletedAt := now.Add(-1 * time.Hour)
	machine := models.Machine{
		ID:          1,
		MachineCode: "DEL01",
		MachineName: "Deleted Machine",
		DeletedAt:   &deletedAt,
	}

	assert.NotNil(t, machine.DeletedAt)
	assert.True(t, machine.DeletedAt.Before(now))
}

// =============================================================================
// CreateMachineRequest Tests
// =============================================================================

func TestCreateMachineRequest_Structure(t *testing.T) {
	req := models.CreateMachineRequest{
		MachineCode: "NEW01",
		MachineName: "New Machine",
		MachineType: "CNC Lathe",
		Location:    "Building C",
		Status:      "active",
	}

	assert.Equal(t, "NEW01", req.MachineCode)
	assert.Equal(t, "New Machine", req.MachineName)
	assert.Equal(t, "CNC Lathe", req.MachineType)
	assert.Equal(t, "Building C", req.Location)
	assert.Equal(t, "active", req.Status)
}

func TestCreateMachineRequest_JSONSerialization(t *testing.T) {
	req := models.CreateMachineRequest{
		MachineCode: "NEW01",
		MachineName: "New Machine",
		MachineType: "CNC",
		Location:    "Building A",
		Status:      "active",
	}

	jsonData, err := json.Marshal(req)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	assert.Equal(t, "NEW01", result["machine_code"])
	assert.Equal(t, "New Machine", result["machine_name"])
}

func TestCreateMachineRequest_JSONDeserialization(t *testing.T) {
	testCases := []struct {
		name     string
		jsonData string
		expected models.CreateMachineRequest
	}{
		{
			name: "Full request",
			jsonData: `{
				"machine_code": "MC01",
				"machine_name": "Machine 1",
				"machine_type": "CNC",
				"location": "Floor 1",
				"status": "active"
			}`,
			expected: models.CreateMachineRequest{
				MachineCode: "MC01",
				MachineName: "Machine 1",
				MachineType: "CNC",
				Location:    "Floor 1",
				Status:      "active",
			},
		},
		{
			name: "Minimal request",
			jsonData: `{
				"machine_code": "MC02",
				"machine_name": "Machine 2"
			}`,
			expected: models.CreateMachineRequest{
				MachineCode: "MC02",
				MachineName: "Machine 2",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var req models.CreateMachineRequest
			err := json.Unmarshal([]byte(tc.jsonData), &req)
			require.NoError(t, err)

			assert.Equal(t, tc.expected.MachineCode, req.MachineCode)
			assert.Equal(t, tc.expected.MachineName, req.MachineName)
			assert.Equal(t, tc.expected.MachineType, req.MachineType)
			assert.Equal(t, tc.expected.Location, req.Location)
			assert.Equal(t, tc.expected.Status, req.Status)
		})
	}
}

// =============================================================================
// UpdateMachineRequest Tests
// =============================================================================

func TestUpdateMachineRequest_Structure(t *testing.T) {
	req := models.UpdateMachineRequest{
		MachineName: "Updated Name",
		MachineType: "Updated Type",
		Location:    "Updated Location",
		Status:      "maintenance",
	}

	assert.Equal(t, "Updated Name", req.MachineName)
	assert.Equal(t, "Updated Type", req.MachineType)
	assert.Equal(t, "Updated Location", req.Location)
	assert.Equal(t, "maintenance", req.Status)
}

func TestUpdateMachineRequest_PartialUpdate(t *testing.T) {
	// Test that partial updates work (only some fields set)
	jsonData := `{"status": "offline"}`

	var req models.UpdateMachineRequest
	err := json.Unmarshal([]byte(jsonData), &req)
	require.NoError(t, err)

	assert.Equal(t, "offline", req.Status)
	assert.Empty(t, req.MachineName)
	assert.Empty(t, req.MachineType)
	assert.Empty(t, req.Location)
}

// =============================================================================
// DeleteMachineRequest Tests
// =============================================================================

func TestDeleteMachineRequest_Structure(t *testing.T) {
	req := models.DeleteMachineRequest{
		MachineCode: "DEL01",
		Reason:      "Machine is obsolete",
	}

	assert.Equal(t, "DEL01", req.MachineCode)
	assert.Equal(t, "Machine is obsolete", req.Reason)
}

func TestDeleteMachineRequest_JSONSerialization(t *testing.T) {
	req := models.DeleteMachineRequest{
		MachineCode: "DEL01",
		Reason:      "End of life",
	}

	jsonData, err := json.Marshal(req)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	assert.Equal(t, "DEL01", result["machine_code"])
	assert.Equal(t, "End of life", result["reason"])
}

// =============================================================================
// Machine Status Tests
// =============================================================================

func TestMachine_StatusValues(t *testing.T) {
	validStatuses := []string{"active", "maintenance", "offline", "inactive"}

	for _, status := range validStatuses {
		t.Run("Status_"+status, func(t *testing.T) {
			machine := models.Machine{
				Status: status,
			}
			assert.Equal(t, status, machine.Status)
		})
	}
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestMachine_EmptyFields(t *testing.T) {
	machine := models.Machine{}

	assert.Equal(t, int64(0), machine.ID)
	assert.Empty(t, machine.MachineCode)
	assert.Empty(t, machine.MachineName)
	assert.Empty(t, machine.MachineType)
	assert.Empty(t, machine.Location)
	assert.Empty(t, machine.Status)
	assert.Nil(t, machine.DeletedAt)
}

func TestCreateMachineRequest_EmptyFields(t *testing.T) {
	req := models.CreateMachineRequest{}

	assert.Empty(t, req.MachineCode)
	assert.Empty(t, req.MachineName)
	assert.Empty(t, req.MachineType)
	assert.Empty(t, req.Location)
	assert.Empty(t, req.Status)
}

func TestMachine_SpecialCharactersInFields(t *testing.T) {
	machine := models.Machine{
		MachineCode: "MC-01/A",
		MachineName: "Machine (Special) #1",
		MachineType: "CNC/VMC",
		Location:    "Building A - Floor 1 & 2",
	}

	assert.Equal(t, "MC-01/A", machine.MachineCode)
	assert.Equal(t, "Machine (Special) #1", machine.MachineName)
	assert.Equal(t, "CNC/VMC", machine.MachineType)
	assert.Equal(t, "Building A - Floor 1 & 2", machine.Location)
}

func TestMachine_LongFields(t *testing.T) {
	longString := "This is a very long machine name that might exceed typical display limits but should still be valid"
	machine := models.Machine{
		MachineName: longString,
	}

	assert.Equal(t, longString, machine.MachineName)
}

