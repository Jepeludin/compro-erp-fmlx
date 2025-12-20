package handlers

import (
	"net/http"
	"strconv"

	"ganttpro-backend/models"
	"ganttpro-backend/services"

	"github.com/gin-gonic/gin"
)

type GCodeHandler struct {
	service *services.GCodeService
}

func NewGCodeHandler(service *services.GCodeService) *GCodeHandler {
	return &GCodeHandler{service: service}
}

// UploadGCode uploads a G-code file for an operation plan
// @Summary Upload G-code file
// @Description Upload a .txt G-code file for an operation plan
// @Tags g-codes
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param operation_plan_id formData int true "Operation Plan ID"
// @Param file formData file true "G-code file (.txt only)"
// @Success 201 {object} models.GCodeFile
// @Failure 400 {object} map[string]string
// @Router /api/v1/g-codes/upload [post]
func (h *GCodeHandler) UploadGCode(c *gin.Context) {
	// Get operation plan ID
	planIDStr := c.PostForm("operation_plan_id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation plan ID"})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Get user from context
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	// Upload file
	gcodeFile, err := h.service.UploadGCode(uint(planID), file, userObj.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "File uploaded successfully",
		"data":    gcodeFile,
	})
}

// GetGCodeFilesByPlan retrieves all G-code files for an operation plan
// @Summary Get G-code files by plan
// @Description Retrieve all G-code files for an operation plan
// @Tags g-codes
// @Produce json
// @Security BearerAuth
// @Param plan_id path int true "Operation Plan ID"
// @Success 200 {array} models.GCodeFile
// @Router /api/v1/g-codes/plan/{plan_id} [get]
func (h *GCodeHandler) GetGCodeFilesByPlan(c *gin.Context) {
	planID, err := strconv.ParseUint(c.Param("plan_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation plan ID"})
		return
	}

	files, err := h.service.GetGCodeFilesByPlan(uint(planID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": files})
}

// DownloadGCode downloads a G-code file
// @Summary Download G-code file
// @Description Download a G-code file by ID
// @Tags g-codes
// @Produce application/octet-stream
// @Security BearerAuth
// @Param id path int true "G-code File ID"
// @Success 200 {file} file
// @Failure 404 {object} map[string]string
// @Router /api/v1/g-codes/{id}/download [get]
func (h *GCodeHandler) DownloadGCode(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	filePath, originalName, err := h.service.GetFilePath(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.FileAttachment(filePath, originalName)
}

// DeleteGCode deletes a G-code file
// @Summary Delete G-code file
// @Description Delete a G-code file (draft status only)
// @Tags g-codes
// @Produce json
// @Security BearerAuth
// @Param id path int true "G-code File ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/g-codes/{id} [delete]
func (h *GCodeHandler) DeleteGCode(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	user, _ := c.Get("user")
	userObj := user.(*models.User)

	if err := h.service.DeleteGCodeFile(uint(id), userObj.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
