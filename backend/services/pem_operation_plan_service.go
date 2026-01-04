package services

import (
	"errors"
	"fmt"
	"ganttpro-backend/models"
	"ganttpro-backend/repository"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type PEMOperationPlanService struct {
	repo             *repository.PEMOperationPlanRepository
	userRepo         *repository.UserRepository
	ppicScheduleRepo *repository.PPICScheduleRepository
	emailService     *EmailService
	uploadDir        string
}

func NewPEMOperationPlanService(
	repo *repository.PEMOperationPlanRepository,
	userRepo *repository.UserRepository,
	ppicScheduleRepo *repository.PPICScheduleRepository,
	emailService *EmailService,
	uploadDir string,
) *PEMOperationPlanService {
	// Create upload directory if not exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		panic(fmt.Sprintf("Failed to create PEM upload directory: %v", err))
	}

	return &PEMOperationPlanService{
		repo:             repo,
		userRepo:         userRepo,
		ppicScheduleRepo: ppicScheduleRepo,
		emailService:     emailService,
		uploadDir:        uploadDir,
	}
}

// CreatePlan creates a new PEM operation plan
func (s *PEMOperationPlanService) CreatePlan(request models.CreatePEMPlanRequest, createdBy int64) (*models.PEMOperationPlan, error) {
	// Create plan
	plan := &models.PEMOperationPlan{
		PPICScheduleID: request.PPICScheduleID,
		PartName:       request.PartName,
		Material:       request.Material,
		DialSize:       request.DialSize,
		Quantity:       request.Quantity,
		Revision:       request.Revision,
		NoWP:           request.NoWP,
		Page:           request.Page,
		Status:         models.PEMStatusDraft,
		CreatedBy:      createdBy,
	}

	if err := s.repo.Create(plan); err != nil {
		return nil, fmt.Errorf("failed to create operation plan: %w", err)
	}

	// Create steps if provided
	for _, stepReq := range request.Steps {
		step := &models.OperationPlanStep{
			OperationPlanID: plan.ID,
			StepNumber:      stepReq.StepNumber,
			ClampingSystem:  stepReq.ClampingSystem,
			RawMaterial:     stepReq.RawMaterial,
			Setting:         stepReq.Setting,
			Process:         stepReq.Process,
			Note:            stepReq.Note,
			CheckingMethod:  stepReq.CheckingMethod,
		}

		if err := s.repo.CreateStep(step); err != nil {
			return nil, fmt.Errorf("failed to create step: %w", err)
		}
	}

	// Reload plan with all relations
	return s.repo.FindByID(plan.ID)
}

// GetPlanByID retrieves a plan by ID
func (s *PEMOperationPlanService) GetPlanByID(id int64) (*models.PEMOperationPlan, error) {
	plan, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}
	return plan, nil
}

// GetAllPlans retrieves all plans with optional filters
func (s *PEMOperationPlanService) GetAllPlans(filters map[string]interface{}) ([]models.PEMOperationPlan, error) {
	return s.repo.FindAll(filters)
}

// GetPlansByPPICSchedule retrieves plans for a specific PPIC schedule
func (s *PEMOperationPlanService) GetPlansByPPICSchedule(scheduleID int64) ([]models.PEMOperationPlan, error) {
	filters := map[string]interface{}{
		"ppic_schedule_id": scheduleID,
	}
	return s.repo.FindAll(filters)
}

// UpdatePlan updates an existing plan
func (s *PEMOperationPlanService) UpdatePlan(id int64, request models.UpdatePEMPlanRequest, userID int64) error {
	plan, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("plan not found: %w", err)
	}

	// Only allow updates to draft plans by the creator
	if plan.Status != models.PEMStatusDraft {
		return errors.New("can only update plans in draft status")
	}

	if plan.CreatedBy != userID {
		return errors.New("only the creator can update this plan")
	}

	// Update fields
	if request.PartName != "" {
		plan.PartName = request.PartName
	}
	if request.Material != "" {
		plan.Material = request.Material
	}
	if request.DialSize != "" {
		plan.DialSize = request.DialSize
	}
	if request.Quantity != nil {
		plan.Quantity = *request.Quantity
	}
	if request.Revision != "" {
		plan.Revision = request.Revision
	}
	if request.NoWP != "" {
		plan.NoWP = request.NoWP
	}
	if request.Page != "" {
		plan.Page = request.Page
	}

	return s.repo.Update(plan)
}

// DeletePlan deletes a plan (only drafts can be deleted)
func (s *PEMOperationPlanService) DeletePlan(id int64, userID int64) error {
	plan, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("plan not found: %w", err)
	}

	// Only allow deletion of draft plans by the creator
	if plan.Status != models.PEMStatusDraft {
		return errors.New("can only delete plans in draft status")
	}

	if plan.CreatedBy != userID {
		return errors.New("only the creator can delete this plan")
	}

	// Delete step images from disk
	for _, step := range plan.Steps {
		if step.PictureURL != "" {
			s.deleteStepImageFile(step.PictureURL)
		}
	}

	return s.repo.Delete(id)
}

// Step Management

// AddStep adds a new step to a plan
func (s *PEMOperationPlanService) AddStep(planID int64, request models.CreateStepRequest, userID int64) (*models.OperationPlanStep, error) {
	plan, err := s.repo.FindByID(planID)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}

	// Only allow adding steps to draft plans by the creator
	if plan.Status != models.PEMStatusDraft {
		return nil, errors.New("can only add steps to plans in draft status")
	}

	if plan.CreatedBy != userID {
		return nil, errors.New("only the creator can modify this plan")
	}

	step := &models.OperationPlanStep{
		OperationPlanID: planID,
		StepNumber:      request.StepNumber,
		ClampingSystem:  request.ClampingSystem,
		RawMaterial:     request.RawMaterial,
		Setting:         request.Setting,
		Process:         request.Process,
		Note:            request.Note,
		CheckingMethod:  request.CheckingMethod,
	}

	if err := s.repo.CreateStep(step); err != nil {
		return nil, fmt.Errorf("failed to create step: %w", err)
	}

	return step, nil
}

// UpdateStep updates an existing step
func (s *PEMOperationPlanService) UpdateStep(stepID int64, request models.UpdateStepRequest, userID int64) error {
	step, err := s.repo.FindStepByID(stepID)
	if err != nil {
		return fmt.Errorf("step not found: %w", err)
	}

	// Get the plan to check permissions
	plan, err := s.repo.FindByID(step.OperationPlanID)
	if err != nil {
		return fmt.Errorf("plan not found: %w", err)
	}

	if plan.Status != models.PEMStatusDraft {
		return errors.New("can only update steps in plans with draft status")
	}

	if plan.CreatedBy != userID {
		return errors.New("only the creator can modify this plan")
	}

	// Update step fields
	if request.StepNumber > 0 {
		step.StepNumber = request.StepNumber
	}
	step.ClampingSystem = request.ClampingSystem
	step.RawMaterial = request.RawMaterial
	step.Setting = request.Setting
	step.Process = request.Process
	step.Note = request.Note
	step.CheckingMethod = request.CheckingMethod

	return s.repo.UpdateStep(step)
}

// DeleteStep deletes a step
func (s *PEMOperationPlanService) DeleteStep(stepID int64, userID int64) error {
	step, err := s.repo.FindStepByID(stepID)
	if err != nil {
		return fmt.Errorf("step not found: %w", err)
	}

	// Get the plan to check permissions
	plan, err := s.repo.FindByID(step.OperationPlanID)
	if err != nil {
		return fmt.Errorf("plan not found: %w", err)
	}

	if plan.Status != models.PEMStatusDraft {
		return errors.New("can only delete steps from plans in draft status")
	}

	if plan.CreatedBy != userID {
		return errors.New("only the creator can modify this plan")
	}

	// Delete image file if exists
	if step.PictureURL != "" {
		s.deleteStepImageFile(step.PictureURL)
	}

	return s.repo.DeleteStep(stepID)
}

// Image Upload Management

// UploadStepImage uploads an image for a step
func (s *PEMOperationPlanService) UploadStepImage(stepID int64, file *multipart.FileHeader, userID int64) error {
	step, err := s.repo.FindStepByID(stepID)
	if err != nil {
		return fmt.Errorf("step not found: %w", err)
	}

	// Get the plan to check permissions
	plan, err := s.repo.FindByID(step.OperationPlanID)
	if err != nil {
		return fmt.Errorf("plan not found: %w", err)
	}

	if plan.Status != models.PEMStatusDraft {
		return errors.New("can only upload images to plans in draft status")
	}

	if plan.CreatedBy != userID {
		return errors.New("only the creator can modify this plan")
	}

	// Validate file extension (PNG, JPG, PDF)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".pdf" {
		return errors.New("only PNG, JPG, and PDF files are allowed")
	}

	// Validate file size (max 5MB)
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if file.Size > maxSize {
		return errors.New("file size exceeds 5MB limit")
	}

	// Delete old image if exists
	if step.PictureURL != "" {
		s.deleteStepImageFile(step.PictureURL)
	}

	// Create plan-specific subdirectory
	planDir := filepath.Join(s.uploadDir, fmt.Sprintf("plan-%d", plan.ID))
	if err := os.MkdirAll(planDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create plan directory: %w", err)
	}

	// Generate unique filename: step-{stepNumber}-{timestamp}{ext}
	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("step-%d-%d%s", step.StepNumber, timestamp, ext)
	filePath := filepath.Join(planDir, fileName)

	// Save file to disk
	src, err := file.Open()
	if err != nil {
		return errors.New("failed to open uploaded file")
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return errors.New("failed to create file on disk")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return errors.New("failed to save file")
	}

	// Update step record with file info
	// Store relative path from upload directory
	relPath := filepath.Join(fmt.Sprintf("plan-%d", plan.ID), fileName)
	step.PictureURL = relPath
	step.PictureFilename = file.Filename

	return s.repo.UpdateStep(step)
}

// DeleteStepImage deletes the image for a step
func (s *PEMOperationPlanService) DeleteStepImage(stepID int64, userID int64) error {
	step, err := s.repo.FindStepByID(stepID)
	if err != nil {
		return fmt.Errorf("step not found: %w", err)
	}

	// Get the plan to check permissions
	plan, err := s.repo.FindByID(step.OperationPlanID)
	if err != nil {
		return fmt.Errorf("plan not found: %w", err)
	}

	if plan.Status != models.PEMStatusDraft {
		return errors.New("can only delete images from plans in draft status")
	}

	if plan.CreatedBy != userID {
		return errors.New("only the creator can modify this plan")
	}

	// Delete file from disk
	if step.PictureURL != "" {
		s.deleteStepImageFile(step.PictureURL)
	}

	// Clear step image fields
	step.PictureURL = ""
	step.PictureFilename = ""

	return s.repo.UpdateStep(step)
}

// deleteStepImageFile deletes an image file from disk
func (s *PEMOperationPlanService) deleteStepImageFile(relPath string) {
	if relPath == "" {
		return
	}

	filePath := filepath.Join(s.uploadDir, relPath)
	if err := os.Remove(filePath); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to delete image file %s: %v\n", filePath, err)
	}
}

// Approval Workflow

// AssignApprovers assigns approvers to all 5 roles
func (s *PEMOperationPlanService) AssignApprovers(planID int64, approvers map[string]int64, userID int64) error {
	plan, err := s.repo.FindByID(planID)
	if err != nil {
		return fmt.Errorf("plan not found: %w", err)
	}

	// Only creator can assign approvers
	if plan.CreatedBy != userID {
		return errors.New("only the creator can assign approvers")
	}

	// Validate all 5 approvers are provided
	if len(approvers) != 5 {
		return errors.New("all 5 approvers must be assigned")
	}

	// Verify all required roles are present
	for _, role := range models.PEMApproverRoles {
		if _, exists := approvers[role]; !exists {
			return fmt.Errorf("missing approver for role: %s", role)
		}
	}

	// Verify all approvers exist and are active
	for role, userID := range approvers {
		approver, err := s.userRepo.FindByID(uint(userID))
		if err != nil {
			return fmt.Errorf("approver not found for role %s: %w", role, err)
		}
		if !approver.IsActive {
			return fmt.Errorf("approver account is not active for role: %s", role)
		}
	}

	return s.repo.AssignApprovers(planID, approvers)
}

// SubmitForApproval submits a plan for approval and sends email notifications
func (s *PEMOperationPlanService) SubmitForApproval(planID int64, userID int64) error {
	plan, err := s.repo.FindByID(planID)
	if err != nil {
		return fmt.Errorf("plan not found: %w", err)
	}

	// Only creator can submit
	if plan.CreatedBy != userID {
		return errors.New("only the creator can submit for approval")
	}

	// Validate plan is in draft status
	if plan.Status != models.PEMStatusDraft {
		return errors.New("only draft plans can be submitted for approval")
	}

	// Check if all approvers are assigned
	for _, role := range models.PEMApproverRoles {
		approval, err := s.repo.GetApprovalByPlanAndRole(planID, role)
		if err != nil {
			return fmt.Errorf("failed to check approval for role %s: %w", role, err)
		}
		if approval == nil || approval.ApproverID == nil {
			return fmt.Errorf("please assign approver for role: %s", role)
		}
	}

	// Submit for approval
	if err := s.repo.SubmitForApproval(planID); err != nil {
		return fmt.Errorf("failed to submit for approval: %w", err)
	}

	// Send email notifications to all approvers
	if s.emailService.IsConfigured() {
		for _, role := range models.PEMApproverRoles {
			approval, _ := s.repo.GetApprovalByPlanAndRole(planID, role)
			if approval != nil && approval.ApproverID != nil {
				approver, err := s.userRepo.FindByID(uint(*approval.ApproverID))
				if err == nil {
					s.sendApprovalNotification(plan, approver, role)
				}
			}
		}
	}

	return nil
}

// ApprovePlan approves a plan by a specific approver
func (s *PEMOperationPlanService) ApprovePlan(planID int64, approverID int64, role string, comments string) error {
	plan, err := s.repo.FindByID(planID)
	if err != nil {
		return fmt.Errorf("plan not found: %w", err)
	}

	// Validate plan is pending approval
	if plan.Status != models.PEMStatusPendingApproval {
		return errors.New("only plans pending approval can be approved")
	}

	// Approve the plan
	if err := s.repo.ApprovePlan(planID, approverID, role, comments); err != nil {
		return fmt.Errorf("failed to approve plan: %w", err)
	}

	// Reload plan to check if all approved
	plan, _ = s.repo.FindByID(planID)
	if plan.Status == models.PEMStatusApproved {
		// All approvals complete - send email notification to creator
		creator, err := s.userRepo.FindByID(uint(plan.CreatedBy))
		if err == nil && s.emailService.IsConfigured() {
			s.sendApprovedNotification(plan, creator, true)
		}

		// Auto-update PPIC schedule status to "in_progress"
		if plan.PPICScheduleID != nil {
			updateReq := &models.UpdatePPICScheduleRequest{
				Status: "in_progress",
			}
			_, err := s.ppicScheduleRepo.Update(*plan.PPICScheduleID, updateReq, nil, nil)
			if err != nil {
				// Log error but don't fail the approval
				fmt.Printf("Warning: failed to update PPIC schedule status: %v\n", err)
			} else {
				fmt.Printf("PPIC schedule %d status updated to 'in_progress'\n", *plan.PPICScheduleID)
			}
		}
	}

	return nil
}

// RejectPlan rejects a plan
func (s *PEMOperationPlanService) RejectPlan(planID int64, approverID int64, role string, comments string) error {
	plan, err := s.repo.FindByID(planID)
	if err != nil {
		return fmt.Errorf("plan not found: %w", err)
	}

	// Validate plan is pending approval
	if plan.Status != models.PEMStatusPendingApproval {
		return errors.New("only plans pending approval can be rejected")
	}

	// Reject the plan
	if err := s.repo.RejectPlan(planID, approverID, role, comments); err != nil {
		return fmt.Errorf("failed to reject plan: %w", err)
	}

	// Send email notification to creator
	creator, err := s.userRepo.FindByID(uint(plan.CreatedBy))
	if err == nil && s.emailService.IsConfigured() {
		s.sendApprovedNotification(plan, creator, false)
	}

	return nil
}

// GetPendingApprovals gets all plans pending approval for a user
func (s *PEMOperationPlanService) GetPendingApprovals(approverID int64) ([]models.PEMOperationPlan, error) {
	return s.repo.GetPendingApprovalsByApprover(approverID)
}

// Email Notification Helpers

func (s *PEMOperationPlanService) sendApprovalNotification(plan *models.PEMOperationPlan, approver *models.User, role string) {
	// Create email data (simplified version, expand as needed)
	subject := fmt.Sprintf("[COMPRO ERP] Operation Plan Approval Required - %s", plan.FormNumber)

	body := fmt.Sprintf(`
Hi %s,

You have been designated as the %s approver for Operation Plan %s.

Plan Details:
- Form Number: %s
- Part Name: %s
- Material: %s
- Created Date: %s

Please review and approve the plan at your earliest convenience.

Thank you!
`, approver.Username, role, plan.FormNumber, plan.FormNumber, plan.PartName, plan.Material, plan.CreatedAt.Format("2006-01-02"))

	// Send email (basic implementation)
	if err := s.emailService.sendEmail(approver.Email, subject, body); err != nil {
		fmt.Printf("Warning: failed to send approval email: %v\n", err)
	}
}

func (s *PEMOperationPlanService) sendApprovedNotification(plan *models.PEMOperationPlan, creator *models.User, approved bool) {
	status := "Approved"
	if !approved {
		status = "Rejected"
	}

	subject := fmt.Sprintf("[COMPRO ERP] Operation Plan %s - %s", status, plan.FormNumber)

	body := fmt.Sprintf(`
Hi %s,

Your Operation Plan %s has been %s.

The plan is now ready for production execution.

Thank you!
`, creator.Username, plan.FormNumber, strings.ToLower(status))

	// Send email
	if err := s.emailService.sendEmail(creator.Email, subject, body); err != nil {
		fmt.Printf("Warning: failed to send %s email: %v\n", status, err)
	}
}
