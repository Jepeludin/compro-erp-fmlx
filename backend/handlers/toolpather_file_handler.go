package handlers

import (
	"ganttpro-backend/models"
	"ganttpro-backend/services"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ToolpatherFileHandler struct {
	service *services.ToolpatherFileService
}

func NewToolpatherFileHandler(service *services.ToolpatherFileService) *ToolpatherFileHandler {
	return &ToolpatherFileHandler{service: service}
}

// UploadFiles handles multiple file uploads
func (h *ToolpatherFileHandler) UploadFiles(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid form data"})
		return
	}

	files := form.File["files"] // Array of files
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "No files provided"})
		return
	}

	// Get form fields
	orderNumber := c.PostForm("order_number")
	if orderNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Order number is required"})
		return
	}

	partName := c.PostForm("part_name")
	notes := c.PostForm("notes")

	var ppicScheduleID *int64
	if scheduleIDStr := c.PostForm("ppic_schedule_id"); scheduleIDStr != "" {
		id, err := strconv.ParseInt(scheduleIDStr, 10, 64)
		if err == nil {
			ppicScheduleID = &id
		}
	}

	request := models.UploadToolpatherFileRequest{
		PPICScheduleID: ppicScheduleID,
		OrderNumber:    orderNumber,
		PartName:       partName,
		Notes:          notes,
	}

	// Upload files
	uploadedFiles, err := h.service.UploadFiles(request, files, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Files uploaded successfully",
		"data":    uploadedFiles,
	})
}

// GetAllFiles retrieves all files with optional filters
func (h *ToolpatherFileHandler) GetAllFiles(c *gin.Context) {
	filters := make(map[string]interface{})

	if orderNumber := c.Query("order_number"); orderNumber != "" {
		filters["order_number"] = orderNumber
	}

	if uploaderIDStr := c.Query("uploaded_by"); uploaderIDStr != "" {
		uploaderID, err := strconv.ParseInt(uploaderIDStr, 10, 64)
		if err == nil {
			filters["uploaded_by"] = uploaderID
		}
	}

	if scheduleIDStr := c.Query("ppic_schedule_id"); scheduleIDStr != "" {
		scheduleID, err := strconv.ParseInt(scheduleIDStr, 10, 64)
		if err == nil {
			filters["ppic_schedule_id"] = scheduleID
		}
	}

	files, err := h.service.GetAllFiles(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": files})
}

// GetFileByID retrieves a single file by ID
func (h *ToolpatherFileHandler) GetFileByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid file ID"})
		return
	}

	file, err := h.service.GetFileByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "File not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": file})
}

// GetFilesByOrderNumber retrieves all files for a specific order
func (h *ToolpatherFileHandler) GetFilesByOrderNumber(c *gin.Context) {
	orderNumber := c.Param("orderNumber")
	if orderNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Order number is required"})
		return
	}

	files, err := h.service.GetFilesByOrderNumber(orderNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": files})
}

// GetMyFiles retrieves all files uploaded by current user
func (h *ToolpatherFileHandler) GetMyFiles(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	files, err := h.service.GetFilesByUploader(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": files})
}

// DownloadFile serves a file for download
func (h *ToolpatherFileHandler) DownloadFile(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid file ID"})
		return
	}

	file, err := h.service.GetFileByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "File not found"})
		return
	}

	// Get full file path
	filePath := h.service.GetFilePath(file)

	// Set headers for download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(file.FileName))
	c.Header("Content-Type", "application/octet-stream")

	// Serve the file
	c.File(filePath)
}

// DeleteFile deletes a file
func (h *ToolpatherFileHandler) DeleteFile(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid file ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	err = h.service.DeleteFile(id, userID.(int64))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "File deleted successfully"})
}
