package handlers

import (
	"net/http"
	"strconv"

	"ganttpro-backend/models"
	"ganttpro-backend/services"

	"github.com/gin-gonic/gin"
)

type GanttHandler struct {
	ganttService *services.GanttService
}

func NewGanttHandler(ganttService *services.GanttService) *GanttHandler {
	return &GanttHandler{ganttService: ganttService}
}

// GetGanttChart godoc
// @Summary Get Gantt chart data
// @Description Returns formatted data for Gantt chart display with filtering options
// @Tags Gantt Chart
// @Accept json
// @Produce json
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Param machine_id query int false "Filter by machine ID"
// @Param priority query string false "Filter by priority (Low, Medium, Urgent, Top Urgent)"
// @Param status query string false "Filter by status (pending, in_progress, completed)"
// @Param group_by query string false "Group tasks by (priority, machine)"
// @Success 200 {object} models.GanttChartResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/gantt-chart [get]
func (h *GanttHandler) GetGanttChart(c *gin.Context) {
	var filter models.GanttFilterRequest

	// Parse query parameters
	filter.StartDate = c.Query("start_date")
	filter.EndDate = c.Query("end_date")
	filter.Priority = c.Query("priority")
	filter.Status = c.Query("status")
	filter.GroupBy = c.Query("group_by")

	if machineIDStr := c.Query("machine_id"); machineIDStr != "" {
		machineID, err := strconv.ParseInt(machineIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid machine_id"})
			return
		}
		filter.MachineID = machineID
	}

	response, err := h.ganttService.GetGanttChartData(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetAllPPICSchedules godoc
// @Summary Get all PPIC schedules
// @Description Returns all PPIC schedule entries
// @Tags PPIC Schedule
// @Produce json
// @Success 200 {array} models.PPICSchedule
// @Failure 500 {object} map[string]string
// @Router /api/v1/ppic-schedules [get]
func (h *GanttHandler) GetAllPPICSchedules(c *gin.Context) {
	schedules, err := h.ganttService.GetAllPPICSchedules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    schedules,
	})
}

// GetPPICSchedule godoc
// @Summary Get a PPIC schedule by ID
// @Description Returns a single PPIC schedule entry with machine assignments
// @Tags PPIC Schedule
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} models.PPICSchedule
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/ppic-schedules/{id} [get]
func (h *GanttHandler) GetPPICSchedule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	schedule, err := h.ganttService.GetPPICSchedule(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    schedule,
	})
}

// CreatePPICSchedule godoc
// @Summary Create a new PPIC schedule
// @Description Creates a new PPIC schedule entry with machine assignments
// @Tags PPIC Schedule
// @Accept json
// @Produce json
// @Param schedule body models.CreatePPICScheduleRequest true "Schedule data"
// @Success 201 {object} models.PPICSchedule
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/ppic-schedules [post]
func (h *GanttHandler) CreatePPICSchedule(c *gin.Context) {
	var req models.CreatePPICScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Convert uint to int64
	userIDInt64 := int64(userID.(uint))

	schedule, err := h.ganttService.CreatePPICSchedule(&req, userIDInt64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "PPIC schedule created successfully",
		"data":    schedule,
	})
}

// UpdatePPICSchedule godoc
// @Summary Update a PPIC schedule
// @Description Updates an existing PPIC schedule entry
// @Tags PPIC Schedule
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param schedule body models.UpdatePPICScheduleRequest true "Schedule data"
// @Success 200 {object} models.PPICSchedule
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/ppic-schedules/{id} [put]
func (h *GanttHandler) UpdatePPICSchedule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req models.UpdatePPICScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule, err := h.ganttService.UpdatePPICSchedule(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "PPIC schedule updated successfully",
		"data":    schedule,
	})
}

// DeletePPICSchedule godoc
// @Summary Delete a PPIC schedule
// @Description Soft deletes a PPIC schedule entry
// @Tags PPIC Schedule
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/ppic-schedules/{id} [delete]
func (h *GanttHandler) DeletePPICSchedule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	err = h.ganttService.DeletePPICSchedule(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "PPIC schedule deleted successfully",
	})
}

// GetSchedulesByMachine godoc
// @Summary Get schedules by machine
// @Description Returns all PPIC schedules assigned to a specific machine
// @Tags PPIC Schedule
// @Produce json
// @Param machine_id path int true "Machine ID"
// @Success 200 {array} models.PPICSchedule
// @Failure 500 {object} map[string]string
// @Router /api/v1/ppic-schedules/machine/{machine_id} [get]
func (h *GanttHandler) GetSchedulesByMachine(c *gin.Context) {
	machineID, err := strconv.ParseInt(c.Param("machine_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid machine_id"})
		return
	}

	schedules, err := h.ganttService.GetSchedulesByMachine(machineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    schedules,
	})
}

// AddMachineAssignment godoc
// @Summary Add machine to schedule
// @Description Adds a new machine assignment to an existing schedule (max 5 machines)
// @Tags PPIC Schedule
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param assignment body models.CreateMachineAssignmentRequest true "Machine assignment data"
// @Success 201 {object} models.MachineAssignment
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/ppic-schedules/{id}/machines [post]
func (h *GanttHandler) AddMachineAssignment(c *gin.Context) {
	scheduleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule ID"})
		return
	}

	var req models.CreateMachineAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assignment, err := h.ganttService.AddMachineAssignment(scheduleID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Machine assignment added successfully",
		"data":    assignment,
	})
}

// RemoveMachineAssignment godoc
// @Summary Remove machine from schedule
// @Description Removes a machine assignment from a schedule
// @Tags PPIC Schedule
// @Produce json
// @Param id path int true "Schedule ID"
// @Param assignment_id path int true "Assignment ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/ppic-schedules/{id}/machines/{assignment_id} [delete]
func (h *GanttHandler) RemoveMachineAssignment(c *gin.Context) {
	assignmentID, err := strconv.ParseInt(c.Param("assignment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignment ID"})
		return
	}

	err = h.ganttService.RemoveMachineAssignment(assignmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Machine assignment removed successfully",
	})
}

// UpdateMachineAssignmentStatus godoc
// @Summary Update machine assignment status
// @Description Updates the status and actual times of a machine assignment
// @Tags PPIC Schedule
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param assignment_id path int true "Assignment ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/ppic-schedules/{id}/machines/{assignment_id}/status [put]
func (h *GanttHandler) UpdateMachineAssignmentStatus(c *gin.Context) {
	scheduleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule ID"})
		return
	}

	assignmentID, err := strconv.ParseInt(c.Param("assignment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignment ID"})
		return
	}

	var req models.UpdateMachineAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.ganttService.UpdateMachineAssignmentStatus(scheduleID, assignmentID, req.Status, req.ActualStart, req.ActualEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Machine assignment status updated successfully",
	})
}
