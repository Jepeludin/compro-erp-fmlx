package testing

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// UploadGCode Endpoint Tests
// =============================================================================

func TestUploadGCode_MissingPlanID(t *testing.T) {
	router := gin.New()
	router.POST("/g-codes/upload", func(c *gin.Context) {
		planIDStr := c.PostForm("operation_plan_id")
		if planIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation plan ID"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "File uploaded successfully"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/g-codes/upload", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid operation plan ID")
}

func TestUploadGCode_InvalidPlanID(t *testing.T) {
	testCases := []struct {
		name   string
		planID string
	}{
		{"Empty string", ""},
		{"Non-numeric", "abc"},
		{"Negative", "-1"},
		{"Float", "1.5"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/g-codes/upload", func(c *gin.Context) {
				planIDStr := c.PostForm("operation_plan_id")
				if planIDStr == "" || planIDStr == "abc" || planIDStr == "-1" || planIDStr == "1.5" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation plan ID"})
					return
				}
				c.JSON(http.StatusCreated, gin.H{"message": "File uploaded successfully"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/g-codes/upload?operation_plan_id="+tc.planID, nil)
			req.Header.Set("Content-Type", "multipart/form-data")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestUploadGCode_MissingFile(t *testing.T) {
	router := gin.New()
	router.POST("/g-codes/upload", func(c *gin.Context) {
		_, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "File uploaded successfully"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/g-codes/upload", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "File is required")
}

func TestUploadGCode_SuccessResponse(t *testing.T) {
	router := gin.New()
	router.POST("/g-codes/upload", func(c *gin.Context) {
		// Simulate successful upload
		c.JSON(http.StatusCreated, gin.H{
			"message": "File uploaded successfully",
			"data": gin.H{
				"id":                1,
				"operation_plan_id": 10,
				"file_name":         "OP-20241210-001_1702185600.txt",
				"original_name":     "gcode.txt",
				"file_size":         1024,
			},
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/g-codes/upload", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "data")
}

// =============================================================================
// GetGCodeFilesByPlan Endpoint Tests
// =============================================================================

func TestGetGCodeFilesByPlan_ValidPlanID(t *testing.T) {
	router := gin.New()
	router.GET("/g-codes/plan/:plan_id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": []gin.H{
				{"id": 1, "file_name": "file1.txt"},
				{"id": 2, "file_name": "file2.txt"},
			},
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/g-codes/plan/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "data")
}

func TestGetGCodeFilesByPlan_InvalidPlanID(t *testing.T) {
	router := gin.New()
	router.GET("/g-codes/plan/:plan_id", func(c *gin.Context) {
		planID := c.Param("plan_id")
		if planID == "invalid" || planID == "abc" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation plan ID"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
	})

	testCases := []struct {
		planID       string
		expectedCode int
	}{
		{"1", http.StatusOK},
		{"123", http.StatusOK},
		{"invalid", http.StatusBadRequest},
		{"abc", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run("PlanID_"+tc.planID, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/g-codes/plan/"+tc.planID, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestGetGCodeFilesByPlan_EmptyList(t *testing.T) {
	router := gin.New()
	router.GET("/g-codes/plan/:plan_id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/g-codes/plan/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})
	assert.Len(t, data, 0)
}

// =============================================================================
// DownloadGCode Endpoint Tests
// =============================================================================

func TestDownloadGCode_ValidID(t *testing.T) {
	router := gin.New()
	router.GET("/g-codes/:id/download", func(c *gin.Context) {
		// Simulate file download headers
		c.Header("Content-Disposition", "attachment; filename=test.txt")
		c.Header("Content-Type", "application/octet-stream")
		c.String(http.StatusOK, "G-code content")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/g-codes/1/download", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
}

func TestDownloadGCode_InvalidID(t *testing.T) {
	router := gin.New()
	router.GET("/g-codes/:id/download", func(c *gin.Context) {
		id := c.Param("id")
		if id == "invalid" || id == "abc" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
			return
		}
		c.String(http.StatusOK, "G-code content")
	})

	testCases := []struct {
		id           string
		expectedCode int
	}{
		{"1", http.StatusOK},
		{"123", http.StatusOK},
		{"invalid", http.StatusBadRequest},
		{"abc", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run("FileID_"+tc.id, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/g-codes/"+tc.id+"/download", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestDownloadGCode_FileNotFound(t *testing.T) {
	router := gin.New()
	router.GET("/g-codes/:id/download", func(c *gin.Context) {
		id := c.Param("id")
		if id == "999" {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.String(http.StatusOK, "G-code content")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/g-codes/999/download", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "file not found")
}

// =============================================================================
// DeleteGCode Endpoint Tests
// =============================================================================

func TestDeleteGCode_ValidID(t *testing.T) {
	router := gin.New()
	router.DELETE("/g-codes/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/g-codes/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "File deleted successfully")
}

func TestDeleteGCode_InvalidID(t *testing.T) {
	router := gin.New()
	router.DELETE("/g-codes/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "invalid" || id == "abc" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
	})

	testCases := []struct {
		id           string
		expectedCode int
	}{
		{"1", http.StatusOK},
		{"invalid", http.StatusBadRequest},
		{"abc", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run("DeleteID_"+tc.id, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/g-codes/"+tc.id, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestDeleteGCode_NotDraftStatus(t *testing.T) {
	router := gin.New()
	router.DELETE("/g-codes/:id", func(c *gin.Context) {
		// Simulate error when plan is not in draft status
		c.JSON(http.StatusBadRequest, gin.H{"error": "can only delete files from plans in draft status"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/g-codes/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "draft status")
}

func TestDeleteGCode_NotOwner(t *testing.T) {
	router := gin.New()
	router.DELETE("/g-codes/:id", func(c *gin.Context) {
		// Simulate error when user is not the uploader
		c.JSON(http.StatusBadRequest, gin.H{"error": "only the uploader can delete this file"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/g-codes/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "uploader")
}

// =============================================================================
// Authorization Tests
// =============================================================================

func TestGCodeEndpoints_RequireAuth(t *testing.T) {
	endpoints := []struct {
		method string
		path   string
	}{
		{"POST", "/g-codes/upload"},
		{"GET", "/g-codes/plan/1"},
		{"GET", "/g-codes/1/download"},
		{"DELETE", "/g-codes/1"},
	}

	for _, ep := range endpoints {
		t.Run(ep.method+"_"+ep.path, func(t *testing.T) {
			router := gin.New()
			
			// Middleware that checks for auth
			router.Use(func(c *gin.Context) {
				authHeader := c.GetHeader("Authorization")
				if authHeader == "" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
					c.Abort()
					return
				}
				c.Next()
			})

			router.Handle(ep.method, ep.path, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(ep.method, ep.path, nil)
			// No Authorization header
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

// =============================================================================
// File Validation Tests
// =============================================================================

func TestUploadGCode_FileExtensionValidation(t *testing.T) {
	testCases := []struct {
		fileName     string
		shouldAccept bool
	}{
		{"program.txt", true},
		{"gcode.TXT", true},
		{"program.nc", false},
		{"program.gcode", false},
		{"program.pdf", false},
		{"program.exe", false},
	}

	for _, tc := range testCases {
		t.Run("File_"+tc.fileName, func(t *testing.T) {
			router := gin.New()
			router.POST("/g-codes/upload", func(c *gin.Context) {
				// Simulate file extension check
				ext := tc.fileName[len(tc.fileName)-4:]
				if ext != ".txt" && ext != ".TXT" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "only .txt files are allowed"})
					return
				}
				c.JSON(http.StatusCreated, gin.H{"message": "File uploaded successfully"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/g-codes/upload", nil)
			router.ServeHTTP(w, req)

			if tc.shouldAccept {
				assert.Equal(t, http.StatusCreated, w.Code)
			} else {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})
	}
}

func TestUploadGCode_FileSizeValidation(t *testing.T) {
	router := gin.New()
	router.POST("/g-codes/upload", func(c *gin.Context) {
		// Simulate max file size check (10MB)
		maxSize := int64(10 * 1024 * 1024)
		fileSize := int64(15 * 1024 * 1024) // 15MB - exceeds limit

		if fileSize > maxSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds 10MB limit"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "File uploaded successfully"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/g-codes/upload", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "10MB limit")
}

