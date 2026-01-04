package handlers

import (
	"net/http"
	"strconv"

	"ganttpro-backend/models"
	"ganttpro-backend/services"

	"github.com/gin-gonic/gin"
)

type PEMOperationPlanHandler struct {
	service *services.PEMOperationPlanService
}

func NewPEMOperationPlanHandler(service *services.PEMOperationPlanService) *PEMOperationPlanHandler {
	return &PEMOperationPlanHandler{service: service}
}

// CreatePEMPlan creates a new PEM operation plan
func (h *PEMOperationPlanHandler) CreatePEMPlan(c *gin.Context) {
	var request models.CreatePEMPlanRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Get user from context
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	plan, err := h.service.CreatePlan(request, int64(userObj.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create operation plan", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Operation plan created successfully",
		"data":    plan,
	})
}

// GetAllPEMPlans retrieves all PEM operation plans with optional filters
func (h *PEMOperationPlanHandler) GetAllPEMPlans(c *gin.Context) {
	filters := make(map[string]interface{})

	// Apply filters from query params
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if createdByStr := c.Query("created_by"); createdByStr != "" {
		if createdBy, err := strconv.ParseInt(createdByStr, 10, 64); err == nil {
			filters["created_by"] = createdBy
		}
	}
	if scheduleIDStr := c.Query("ppic_schedule_id"); scheduleIDStr != "" {
		if scheduleID, err := strconv.ParseInt(scheduleIDStr, 10, 64); err == nil {
			filters["ppic_schedule_id"] = scheduleID
		}
	}

	plans, err := h.service.GetAllPlans(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve operation plans", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    plans,
	})
}

// GetPEMPlan retrieves a single PEM operation plan by ID
func (h *PEMOperationPlanHandler) GetPEMPlan(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	plan, err := h.service.GetPlanByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    plan,
	})
}

// UpdatePEMPlan updates an existing PEM operation plan
func (h *PEMOperationPlanHandler) UpdatePEMPlan(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	var request models.UpdatePEMPlanRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Get user from context
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	if err := h.service.UpdatePlan(id, request, int64(userObj.ID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update plan", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Plan updated successfully",
	})
}

// DeletePEMPlan deletes a PEM operation plan
func (h *PEMOperationPlanHandler) DeletePEMPlan(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	// Get user from context
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	if err := h.service.DeletePlan(id, int64(userObj.ID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete plan", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Plan deleted successfully",
	})
}

// Step Management

// AddPlanStep adds a new step to an operation plan
func (h *PEMOperationPlanHandler) AddPlanStep(c *gin.Context) {
	planID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	var request models.CreateStepRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Get user from context
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	step, err := h.service.AddStep(planID, request, int64(userObj.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to add step", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Step added successfully",
		"data":    step,
	})
}

// UpdatePlanStep updates an existing step
func (h *PEMOperationPlanHandler) UpdatePlanStep(c *gin.Context) {
	stepID, err := strconv.ParseInt(c.Param("step_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid step ID"})
		return
	}

	var request models.UpdateStepRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Get user from context
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	if err := h.service.UpdateStep(stepID, request, int64(userObj.ID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update step", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Step updated successfully",
	})
}

// DeletePlanStep deletes a step from an operation plan
func (h *PEMOperationPlanHandler) DeletePlanStep(c *gin.Context) {
	stepID, err := strconv.ParseInt(c.Param("step_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid step ID"})
		return
	}

	// Get user from context
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	if err := h.service.DeleteStep(stepID, int64(userObj.ID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete step", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Step deleted successfully",
	})
}

// Image Upload Management

// UploadStepImage uploads an image for a step
func (h *PEMOperationPlanHandler) UploadStepImage(c *gin.Context) {
	stepID, err := strconv.ParseInt(c.Param("step_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid step ID"})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	// Get user from context
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	if err := h.service.UploadStepImage(stepID, file, int64(userObj.ID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload image", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Image uploaded successfully",
	})
}

// DeleteStepImage deletes the image for a step
func (h *PEMOperationPlanHandler) DeleteStepImage(c *gin.Context) {
	stepID, err := strconv.ParseInt(c.Param("step_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid step ID"})
		return
	}

	// Get user from context
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	if err := h.service.DeleteStepImage(stepID, int64(userObj.ID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete image", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Image deleted successfully",
	})
}

// Approval Workflow

// AssignApprovers assigns approvers to all 5 roles
func (h *PEMOperationPlanHandler) AssignApprovers(c *gin.Context) {
	planID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	var request models.AssignApproversRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Get user from context
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	if err := h.service.AssignApprovers(planID, request.Approvers, int64(userObj.ID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to assign approvers", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Approvers assigned successfully",
	})
}

// SubmitPlanForApproval submits a plan for approval
func (h *PEMOperationPlanHandler) SubmitPlanForApproval(c *gin.Context) {
	planID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	// Get user from context
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	if err := h.service.SubmitForApproval(planID, int64(userObj.ID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to submit for approval", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Plan submitted for approval successfully",
	})
}

// ApprovePlan approves a plan
func (h *PEMOperationPlanHandler) ApprovePlan(c *gin.Context) {
	planID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	// Get role from query param or body
	role := c.Query("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role parameter is required"})
		return
	}

	var request models.ApprovalActionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		// Comments are optional, so we can proceed without error
		request.Comments = ""
	}

	// Get user from context
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	if err := h.service.ApprovePlan(planID, int64(userObj.ID), role, request.Comments); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to approve plan", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Plan approved successfully",
	})
}

// RejectPlan rejects a plan
func (h *PEMOperationPlanHandler) RejectPlan(c *gin.Context) {
	planID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	// Get role from query param
	role := c.Query("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role parameter is required"})
		return
	}

	var request models.ApprovalActionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		// Comments are optional, so we can proceed without error
		request.Comments = ""
	}

	// Get user from context
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	if err := h.service.RejectPlan(planID, int64(userObj.ID), role, request.Comments); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to reject plan", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Plan rejected",
	})
}

// GetPlansByPPICSchedule retrieves plans for a specific PPIC schedule
func (h *PEMOperationPlanHandler) GetPlansByPPICSchedule(c *gin.Context) {
	scheduleID, err := strconv.ParseInt(c.Param("schedule_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	plans, err := h.service.GetPlansByPPICSchedule(scheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve plans", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    plans,
	})
}

// GetPendingApprovals retrieves plans pending approval for the current user
func (h *PEMOperationPlanHandler) GetPendingApprovals(c *gin.Context) {
	// Get user from context
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	plans, err := h.service.GetPendingApprovals(int64(userObj.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve pending approvals", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    plans,
	})
}
