package handlers

import (
	"net/http"
	"strconv"

	"ganttpro-backend/models"
	"ganttpro-backend/repository"
	"ganttpro-backend/services"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	emailService *services.EmailService
	opPlanRepo   *repository.OperationPlanRepository
	userRepo     *repository.UserRepository
}

func NewEmailHandler(
	emailService *services.EmailService,
	opPlanRepo *repository.OperationPlanRepository,
	userRepo *repository.UserRepository,
) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
		opPlanRepo:   opPlanRepo,
		userRepo:     userRepo,
	}
}

// SendApprovalReminder sends reminder emails to approvers who haven't approved yet
// @Summary Send approval reminder emails
// @Description Send reminder emails to approvers who haven't approved the operation plan yet (PEM role only)
// @Tags Email
// @Produce json
// @Security BearerAuth
// @Param id path int true "Operation Plan ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/operation-plans/{id}/send-reminder [post]
func (h *EmailHandler) SendApprovalReminder(c *gin.Context) {
	// Get operation plan ID
	planID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid operation plan ID"})
		return
	}

	// Get current user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "User not found in context"})
		return
	}
	currentUser := user.(*models.User)

	// Check if user has PEM role (only PEM can trigger reminders)
	if currentUser.Role != models.RolePEM {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only PEM users can send approval reminders"})
		return
	}

	// Check if email service is configured
	if !h.emailService.IsConfigured() {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Email service is not configured. Please set SMTP environment variables.",
		})
		return
	}

	// Get the operation plan
	plan, err := h.opPlanRepo.FindByID(uint(planID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Operation plan not found"})
		return
	}

	// Check if plan is in pending_approval status
	if plan.Status != models.StatusPendingApproval {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Can only send reminders for plans pending approval"})
		return
	}

	// Get pending approvers (those who haven't approved yet)
	pendingApprovers := h.getPendingApprovers(plan)
	if len(pendingApprovers) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "No pending approvers found. All approvals may be complete."})
		return
	}

	// Get users for pending approver roles
	pendingUsers, err := h.userRepo.FindByRoles(pendingApprovers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to fetch approvers"})
		return
	}

	if len(pendingUsers) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "No active users found for pending approver roles"})
		return
	}

	// Get creator info
	creatorName := "Unknown"
	if plan.Creator != nil {
		creatorName = plan.Creator.Username
	}

	// Get job order NJO and machine name from relations
	jobOrderNJO := ""
	if plan.JobOrder != nil {
		jobOrderNJO = plan.JobOrder.NJO
	}
	machineName := ""
	if plan.Machine != nil {
		machineName = plan.Machine.MachineName
	}

	// Build email data
	emailData := services.ApprovalReminderData{
		PlanID:          plan.ID,
		PlanDescription: plan.Description,
		JobOrderNJO:     jobOrderNJO,
		MachineName:     machineName,
		CreatorName:     creatorName,
		SubmittedAt:     plan.UpdatedAt.Format("2006-01-02 15:04"),
		FrontendURL:     h.emailService.GetFrontendURL(),
	}

	// Convert to user pointers
	var userPtrs []*models.User
	for i := range pendingUsers {
		userPtrs = append(userPtrs, &pendingUsers[i])
	}

	// Send emails
	errors := h.emailService.SendBulkApprovalReminders(userPtrs, emailData)

	// Build response
	successCount := len(pendingUsers) - len(errors)

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"message":       "Reminder emails processed",
		"emails_sent":   successCount,
		"emails_failed": len(errors),
		"pending_roles": pendingApprovers,
	})
}

// getPendingApprovers returns the roles that haven't approved yet
func (h *EmailHandler) getPendingApprovers(plan *models.OperationPlan) []string {
	// Get all approver roles
	allApproverRoles := models.ApproverRoles

	// Build a map of approved roles
	approvedRoles := make(map[string]bool)
	for _, approval := range plan.Approvals {
		if approval.Status == models.StatusApproved {
			approvedRoles[approval.ApproverRole] = true
		}
	}

	// Find roles that haven't approved
	var pendingRoles []string
	for _, role := range allApproverRoles {
		if !approvedRoles[role] {
			pendingRoles = append(pendingRoles, role)
		}
	}

	return pendingRoles
}

// CheckEmailConfig checks if email is properly configured
// @Summary Check email configuration
// @Description Check if SMTP email is properly configured
// @Tags Email
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/email/status [get]
func (h *EmailHandler) CheckEmailConfig(c *gin.Context) {
	isConfigured := h.emailService.IsConfigured()

	message := "Email service is not configured. Set SMTP_USERNAME, SMTP_PASSWORD, and SMTP_FROM_ADDR environment variables."
	if isConfigured {
		message = "Email service is configured and ready"
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"is_configured": isConfigured,
		"message":       message,
	})
}
