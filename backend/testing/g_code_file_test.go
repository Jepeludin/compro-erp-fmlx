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
// GCodeFile Structure Tests
// =============================================================================

func TestGCodeFile_Structure(t *testing.T) {
	now := time.Now()
	gcodeFile := models.GCodeFile{
		ID:              1,
		OperationPlanID: 10,
		FileName:        "OP-20241210-001_1702185600.txt",
		OriginalName:    "program_v1.txt",
		FilePath:        "./uploads/gcodes/OP-20241210-001_1702185600.txt",
		FileSize:        1024,
		UploadedBy:      5,
		CreatedAt:       now,
	}

	assert.Equal(t, uint(1), gcodeFile.ID)
	assert.Equal(t, uint(10), gcodeFile.OperationPlanID)
	assert.Equal(t, "OP-20241210-001_1702185600.txt", gcodeFile.FileName)
	assert.Equal(t, "program_v1.txt", gcodeFile.OriginalName)
	assert.Equal(t, "./uploads/gcodes/OP-20241210-001_1702185600.txt", gcodeFile.FilePath)
	assert.Equal(t, int64(1024), gcodeFile.FileSize)
	assert.Equal(t, uint(5), gcodeFile.UploadedBy)
}

func TestGCodeFile_TableName(t *testing.T) {
	gcodeFile := models.GCodeFile{}
	tableName := gcodeFile.TableName()

	assert.Equal(t, "g_code_files", tableName)
}

func TestGCodeFile_JSONSerialization(t *testing.T) {
	gcodeFile := models.GCodeFile{
		ID:              1,
		OperationPlanID: 10,
		FileName:        "test.txt",
		OriginalName:    "original.txt",
		FilePath:        "./uploads/test.txt",
		FileSize:        2048,
		UploadedBy:      1,
		CreatedAt:       time.Now(),
	}

	jsonData, err := json.Marshal(gcodeFile)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	assert.Contains(t, result, "id")
	assert.Contains(t, result, "operation_plan_id")
	assert.Contains(t, result, "file_name")
	assert.Contains(t, result, "original_name")
	assert.Contains(t, result, "file_path")
	assert.Contains(t, result, "file_size")
	assert.Contains(t, result, "uploaded_by")
	assert.Contains(t, result, "created_at")
}

func TestGCodeFile_JSONDeserialization(t *testing.T) {
	jsonData := `{
		"id": 1,
		"operation_plan_id": 10,
		"file_name": "OP-20241210-001_123.txt",
		"original_name": "gcode.txt",
		"file_path": "./uploads/gcodes/OP-20241210-001_123.txt",
		"file_size": 4096,
		"uploaded_by": 3
	}`

	var gcodeFile models.GCodeFile
	err := json.Unmarshal([]byte(jsonData), &gcodeFile)
	require.NoError(t, err)

	assert.Equal(t, uint(1), gcodeFile.ID)
	assert.Equal(t, uint(10), gcodeFile.OperationPlanID)
	assert.Equal(t, "OP-20241210-001_123.txt", gcodeFile.FileName)
	assert.Equal(t, "gcode.txt", gcodeFile.OriginalName)
	assert.Equal(t, int64(4096), gcodeFile.FileSize)
}

// =============================================================================
// ToResponse Tests
// =============================================================================

func TestGCodeFile_ToResponse_WithoutUploader(t *testing.T) {
	gcodeFile := models.GCodeFile{
		ID:              1,
		OperationPlanID: 10,
		FileName:        "test.txt",
		OriginalName:    "original.txt",
		FilePath:        "./uploads/test.txt",
		FileSize:        1024,
		UploadedBy:      5,
		CreatedAt:       time.Now(),
	}

	response := gcodeFile.ToResponse()

	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, uint(10), response.OperationPlanID)
	assert.Equal(t, "test.txt", response.FileName)
	assert.Equal(t, "original.txt", response.OriginalName)
	assert.Equal(t, "./uploads/test.txt", response.FilePath)
	assert.Equal(t, int64(1024), response.FileSize)
	assert.Equal(t, uint(5), response.UploadedBy)
	assert.Nil(t, response.Uploader) // No uploader set
}

func TestGCodeFile_ToResponse_WithUploader(t *testing.T) {
	uploader := &models.User{
		ID:       5,
		Username: "BAYU",
		UserID:   "PI1224.0001",
		Role:     "PEM",
		IsActive: true,
	}

	gcodeFile := models.GCodeFile{
		ID:              1,
		OperationPlanID: 10,
		FileName:        "test.txt",
		OriginalName:    "original.txt",
		FilePath:        "./uploads/test.txt",
		FileSize:        1024,
		UploadedBy:      5,
		Uploader:        uploader,
		CreatedAt:       time.Now(),
	}

	response := gcodeFile.ToResponse()

	assert.NotNil(t, response.Uploader)
	assert.Equal(t, uint(5), response.Uploader.ID)
	assert.Equal(t, "BAYU", response.Uploader.Username)
	assert.Equal(t, "PEM", response.Uploader.Role)
}

// =============================================================================
// GCodeFileResponse Tests
// =============================================================================

func TestGCodeFileResponse_Structure(t *testing.T) {
	response := models.GCodeFileResponse{
		ID:              1,
		OperationPlanID: 10,
		FileName:        "test.txt",
		OriginalName:    "original.txt",
		FilePath:        "./uploads/test.txt",
		FileSize:        2048,
		UploadedBy:      3,
		CreatedAt:       time.Now(),
	}

	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, uint(10), response.OperationPlanID)
	assert.Equal(t, "test.txt", response.FileName)
	assert.Equal(t, "original.txt", response.OriginalName)
	assert.Equal(t, int64(2048), response.FileSize)
}

func TestGCodeFileResponse_JSONSerialization(t *testing.T) {
	response := models.GCodeFileResponse{
		ID:              1,
		OperationPlanID: 10,
		FileName:        "output.txt",
		OriginalName:    "input.txt",
		FilePath:        "./uploads/output.txt",
		FileSize:        1024,
		UploadedBy:      1,
		CreatedAt:       time.Now(),
	}

	jsonData, err := json.Marshal(response)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	assert.Contains(t, result, "id")
	assert.Contains(t, result, "operation_plan_id")
	assert.Contains(t, result, "file_name")
	assert.Contains(t, result, "original_name")
	assert.Contains(t, result, "file_path")
	assert.Contains(t, result, "file_size")
	assert.Contains(t, result, "uploaded_by")
}

func TestGCodeFileResponse_WithUploader(t *testing.T) {
	uploader := &models.UserResponse{
		ID:       3,
		Username: "AMELIA",
		UserID:   "PI1224.0003",
		Role:     "Engineering",
		IsActive: true,
	}

	response := models.GCodeFileResponse{
		ID:       1,
		Uploader: uploader,
	}

	assert.NotNil(t, response.Uploader)
	assert.Equal(t, "AMELIA", response.Uploader.Username)
}

// =============================================================================
// File Name Generation Tests
// =============================================================================

func TestGCodeFile_FileNamePattern(t *testing.T) {
	testCases := []struct {
		name         string
		fileName     string
		isValidName  bool
	}{
		{
			name:        "Standard pattern",
			fileName:    "OP-20241210-001_1702185600.txt",
			isValidName: true,
		},
		{
			name:        "With timestamp",
			fileName:    "OP-20241210-002_20241210_153045.txt",
			isValidName: true,
		},
		{
			name:        "Simple name",
			fileName:    "test_file.txt",
			isValidName: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gcodeFile := models.GCodeFile{
				FileName: tc.fileName,
			}
			assert.NotEmpty(t, gcodeFile.FileName)
		})
	}
}

// =============================================================================
// File Size Tests
// =============================================================================

func TestGCodeFile_FileSizes(t *testing.T) {
	testCases := []struct {
		name     string
		fileSize int64
	}{
		{"Small file (1KB)", 1024},
		{"Medium file (100KB)", 102400},
		{"Large file (1MB)", 1048576},
		{"Max allowed (10MB)", 10485760},
		{"Zero size", 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gcodeFile := models.GCodeFile{
				FileSize: tc.fileSize,
			}
			assert.Equal(t, tc.fileSize, gcodeFile.FileSize)
		})
	}
}

// =============================================================================
// File Extension Tests
// =============================================================================

func TestGCodeFile_ValidExtension(t *testing.T) {
	// Only .txt files are allowed
	validFiles := []string{
		"program.txt",
		"gcode.txt",
		"OP-20241210-001.txt",
		"TEST_FILE.TXT",
	}

	for _, fileName := range validFiles {
		t.Run(fileName, func(t *testing.T) {
			gcodeFile := models.GCodeFile{
				OriginalName: fileName,
			}
			assert.Contains(t, gcodeFile.OriginalName, ".txt") // Case insensitive handled by service
		})
	}
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestGCodeFile_EmptyFields(t *testing.T) {
	gcodeFile := models.GCodeFile{}

	assert.Equal(t, uint(0), gcodeFile.ID)
	assert.Equal(t, uint(0), gcodeFile.OperationPlanID)
	assert.Empty(t, gcodeFile.FileName)
	assert.Empty(t, gcodeFile.OriginalName)
	assert.Empty(t, gcodeFile.FilePath)
	assert.Equal(t, int64(0), gcodeFile.FileSize)
	assert.Equal(t, uint(0), gcodeFile.UploadedBy)
	assert.Nil(t, gcodeFile.Uploader)
}

func TestGCodeFile_SpecialCharactersInOriginalName(t *testing.T) {
	gcodeFile := models.GCodeFile{
		OriginalName: "my file (v1.0) - final.txt",
	}

	assert.Equal(t, "my file (v1.0) - final.txt", gcodeFile.OriginalName)
}

func TestGCodeFile_LongFilePath(t *testing.T) {
	longPath := "./uploads/gcodes/very/long/nested/directory/structure/OP-20241210-001_1702185600.txt"
	gcodeFile := models.GCodeFile{
		FilePath: longPath,
	}

	assert.Equal(t, longPath, gcodeFile.FilePath)
}

func TestGCodeFile_UnicodeOriginalName(t *testing.T) {
	gcodeFile := models.GCodeFile{
		OriginalName: "プログラム_v1.txt", // Japanese characters
	}

	assert.NotEmpty(t, gcodeFile.OriginalName)
}

// =============================================================================
// Relationship Tests
// =============================================================================

func TestGCodeFile_BelongsToOperationPlan(t *testing.T) {
	gcodeFile := models.GCodeFile{
		ID:              1,
		OperationPlanID: 100,
	}

	assert.Equal(t, uint(100), gcodeFile.OperationPlanID)
}

func TestGCodeFile_BelongsToUploader(t *testing.T) {
	gcodeFile := models.GCodeFile{
		ID:         1,
		UploadedBy: 50,
	}

	assert.Equal(t, uint(50), gcodeFile.UploadedBy)
}

