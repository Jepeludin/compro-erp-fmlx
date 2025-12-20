package handlers

import (
	"ganttpro-backend/models"
	"ganttpro-backend/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type GanttHandler struct {
	service *services.GanttService
}

func NewGanttHandler(service *services.GanttService) *GanttHandler {
	return &GanttHandler{service: service}
}

// GetGanttChart returns Gantt chart data with filters
// @Summary Get Gantt chart data
// @Description Get formatted Gantt chart data with optional filters and grouping
// @Tags Gantt
// @Produce json
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Param priority query string false "Filter by priority (Low, Medium, Urgent, Top Urgent)"
// @Param status query string false "Filter by status (pending, in_progress, completed)"
// @Param machine_id query int false "Filter by machine ID"
// @Param group_by query string false "Group by: priority, machine, or empty for all"
// @Success 200 {object} models.GanttChartResponse
// @Router /api/v1/gantt-chart [get]
func (h *GanttHandler) GetGanttChart(c *gin.Context) {
	var filter models.GanttFilterRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid filter parameters"})
		return
	}

	response, err := h.service.GetGanttChartData(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// GetAllPPICSchedules returns all PPIC schedules
// @Summary Get all PPIC schedules
// @Description Get all PPIC schedule entries
// @Tags PPIC Schedules
// @Produce json
// @Success 200 {array} models.PPICSchedule
// @Router /api/v1/ppic-schedules [get]
func (h *GanttHandler) GetAllPPICSchedules(c *gin.Context) {
	schedules, err := h.service.GetAllPPICSchedules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": schedules, "count": len(schedules)})
}

// GetPPICSchedule returns a single PPIC schedule
// @Summary Get PPIC schedule by ID
// @Description Get a single PPIC schedule entry
// @Tags PPIC Schedules
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} models.PPICSchedule
// @Router /api/v1/ppic-schedules/{id} [get]
func (h *GanttHandler) GetPPICSchedule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid ID"})
		return
	}

	schedule, err := h.service.GetPPICSchedule(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": schedule})
}

// CreatePPICSchedule creates a new PPIC schedule
// @Summary Create PPIC schedule
// @Description Create a new PPIC schedule entry
// @Tags PPIC Schedules
// @Accept json
// @Produce json
// @Param request body models.CreatePPICScheduleRequest true "Schedule details"
// @Success 201 {object} models.PPICSchedule
// @Router /api/v1/ppic-schedules [post]
func (h *GanttHandler) CreatePPICSchedule(c *gin.Context) {
	var req models.CreatePPICScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request", "details": err.Error()})
		return
	}

	// Get user ID from context
	userID := getUserIDFromContext(c)

	schedule, err := h.service.CreatePPICSchedule(&req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Schedule created successfully", "data": schedule})
}

// UpdatePPICSchedule updates an existing PPIC schedule
// @Summary Update PPIC schedule
// @Description Update an existing PPIC schedule entry
// @Tags PPIC Schedules
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param request body models.UpdatePPICScheduleRequest true "Update details"
// @Success 200 {object} models.PPICSchedule
// @Router /api/v1/ppic-schedules/{id} [put]
func (h *GanttHandler) UpdatePPICSchedule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid ID"})
		return
	}

	var req models.UpdatePPICScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request", "details": err.Error()})
		return
	}

	schedule, err := h.service.UpdatePPICSchedule(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Schedule updated successfully", "data": schedule})
}

// DeletePPICSchedule deletes a PPIC schedule
// @Summary Delete PPIC schedule
// @Description Delete a PPIC schedule entry
// @Tags PPIC Schedules
// @Param id path int true "Schedule ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/ppic-schedules/{id} [delete]
func (h *GanttHandler) DeletePPICSchedule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid ID"})
		return
	}

	if err := h.service.DeletePPICSchedule(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Schedule deleted successfully"})
}

// GetSchedulesByMachine returns schedules for a specific machine
// @Summary Get schedules by machine
// @Description Get all PPIC schedules assigned to a specific machine
// @Tags PPIC Schedules
// @Produce json
// @Param machine_id path int true "Machine ID"
// @Success 200 {array} models.PPICSchedule
// @Router /api/v1/ppic-schedules/machine/{machine_id} [get]
func (h *GanttHandler) GetSchedulesByMachine(c *gin.Context) {
	machineID, err := strconv.ParseInt(c.Param("machine_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid machine ID"})
		return
	}

	schedules, err := h.service.GetSchedulesByMachine(machineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": schedules, "count": len(schedules)})
}

// AddMachineAssignment adds a machine to a schedule
// @Summary Add machine assignment
// @Description Add a machine to an existing PPIC schedule
// @Tags PPIC Schedules
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param request body models.CreateMachineAssignmentRequest true "Machine assignment details"
// @Success 201 {object} models.MachineAssignment
// @Router /api/v1/ppic-schedules/{id}/machines [post]
func (h *GanttHandler) AddMachineAssignment(c *gin.Context) {
	scheduleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid schedule ID"})
		return
	}

	var req models.CreateMachineAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request", "details": err.Error()})
		return
	}

	assignment, err := h.service.AddMachineAssignment(scheduleID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Machine assignment added", "data": assignment})
}

// RemoveMachineAssignment removes a machine from a schedule
// @Summary Remove machine assignment
// @Description Remove a machine assignment from a PPIC schedule
// @Tags PPIC Schedules
// @Param id path int true "Schedule ID"
// @Param assignment_id path int true "Assignment ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/ppic-schedules/{id}/machines/{assignment_id} [delete]
func (h *GanttHandler) RemoveMachineAssignment(c *gin.Context) {
	assignmentID, err := strconv.ParseInt(c.Param("assignment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid assignment ID"})
		return
	}

	if err := h.service.RemoveMachineAssignment(assignmentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Machine assignment removed"})
}

// UpdateMachineAssignmentStatus updates status of a machine assignment
// @Summary Update machine assignment status
// @Description Update the status of a machine assignment (pending, in_progress, completed)
// @Tags PPIC Schedules
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param assignment_id path int true "Assignment ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/ppic-schedules/{id}/machines/{assignment_id}/status [put]
func (h *GanttHandler) UpdateMachineAssignmentStatus(c *gin.Context) {
	scheduleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid schedule ID"})
		return
	}

	assignmentID, err := strconv.ParseInt(c.Param("assignment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid assignment ID"})
		return
	}

	var req struct {
		Status      string     `json:"status" binding:"required"`
		ActualStart *time.Time `json:"actual_start"`
		ActualEnd   *time.Time `json:"actual_end"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request", "details": err.Error()})
		return
	}

	if err := h.service.UpdateMachineAssignmentStatus(scheduleID, assignmentID, req.Status, req.ActualStart, req.ActualEnd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Machine assignment status updated"})
}

// Helper function to get user ID from context
func getUserIDFromContext(c *gin.Context) int64 {
	if user, exists := c.Get("user"); exists {
		if u, ok := user.(*models.User); ok {
			return int64(u.ID)
		}
	}
	return 0
}
