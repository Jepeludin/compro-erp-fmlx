package handlers

import (
	"ganttpro-backend/models"
	"ganttpro-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PPICLinkHandler struct {
	service *services.PPICLinkService
}

func NewPPICLinkHandler(service *services.PPICLinkService) *PPICLinkHandler {
	return &PPICLinkHandler{service: service}
}

// CreatePPICLink creates a new link between two schedules
// @Summary Create PPIC link
// @Description Create a dependency link between two PPIC schedules (Gantt arrow)
// @Tags PPIC Links
// @Accept json
// @Produce json
// @Param request body models.CreatePPICLinkRequest true "Link details"
// @Success 201 {object} models.PPICLink
// @Router /api/v1/ppic-links [post]
func (h *PPICLinkHandler) CreatePPICLink(c *gin.Context) {
	var req models.CreatePPICLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request", "details": err.Error()})
		return
	}

	link, err := h.service.CreateLink(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Link created successfully", "data": link})
}

// GetAllPPICLinks returns all PPIC links
// @Summary Get all PPIC links
// @Description Get all dependency links between PPIC schedules
// @Tags PPIC Links
// @Produce json
// @Success 200 {array} models.PPICLink
// @Router /api/v1/ppic-links [get]
func (h *PPICLinkHandler) GetAllPPICLinks(c *gin.Context) {
	links, err := h.service.GetAllLinks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": links, "count": len(links)})
}

// DeletePPICLink deletes a PPIC link
// @Summary Delete PPIC link
// @Description Delete a dependency link between PPIC schedules
// @Tags PPIC Links
// @Param id path int true "Link ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/ppic-links/{id} [delete]
func (h *PPICLinkHandler) DeletePPICLink(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteLink(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Link deleted successfully"})
}
