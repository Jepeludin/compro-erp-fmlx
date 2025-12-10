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

type OperationPlanService struct {
    planRepo      *repository.OperationPlanRepository
    gCodeFileRepo *repository.GCodeFileRepository
    jobOrderRepo  *repository.JobOrderRepository
    uploadPath    string
}

func NewOperationPlanService(
    planRepo *repository.OperationPlanRepository,
    gCodeFileRepo *repository.GCodeFileRepository,
    jobOrderRepo *repository.JobOrderRepository,
) *OperationPlanService {
    // Create upload directory if not exists
    uploadPath := "./uploads/gcodes"
    if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
        panic("Failed to create upload directory: " + err.Error())
    }

    return &OperationPlanService{
        planRepo:      planRepo,
        gCodeFileRepo: gCodeFileRepo,
        jobOrderRepo:  jobOrderRepo,
        uploadPath:    uploadPath,
    }
}

// CreateOperationPlanRequest represents the request to create an operation plan
type CreateOperationPlanRequest struct {
    JobOrderID   uint   `json:"job_order_id" binding:"required"`
    MachineID    uint   `json:"machine_id" binding:"required"`
    PartQuantity int    `json:"part_quantity"`
    Description  string `json:"description"`
}

// Create creates a new operation plan
func (s *OperationPlanService) Create(req CreateOperationPlanRequest, createdBy uint) (*models.OperationPlan, error) {
    // Verify job order exists
    _, err := s.jobOrderRepo.GetByID(int64(req.JobOrderID))
    if err != nil {
        return nil, errors.New("job order not found")
    }

    // Check if operation plan already exists for this job order
    existingPlan, _ := s.planRepo.FindByJobOrderID(req.JobOrderID)
    if existingPlan != nil {
        return nil, errors.New("operation plan already exists for this job order")
    }

    plan := &models.OperationPlan{
        JobOrderID:   req.JobOrderID,
        MachineID:    req.MachineID,
        PartQuantity: req.PartQuantity,
        Description:  req.Description,
        Status:       models.StatusDraft,
        CreatedBy:    createdBy,
    }

    if err := s.planRepo.Create(plan); err != nil {
        return nil, errors.New("failed to create operation plan")
    }

    // Reload with relations
    return s.planRepo.FindByID(plan.ID)
}

// GetByID gets an operation plan by ID
func (s *OperationPlanService) GetByID(id uint) (*models.OperationPlan, error) {
    return s.planRepo.FindByID(id)
}

// GetByJobOrderID gets an operation plan by job order ID
func (s *OperationPlanService) GetByJobOrderID(jobOrderID uint) (*models.OperationPlan, error) {
    return s.planRepo.FindByJobOrderID(jobOrderID)
}

// GetAll gets all operation plans with optional filters
func (s *OperationPlanService) GetAll(status string, machineID uint) ([]models.OperationPlan, error) {
    return s.planRepo.FindAll(status, machineID)
}

// GetByCreator gets operation plans created by a user
func (s *OperationPlanService) GetByCreator(userID uint) ([]models.OperationPlan, error) {
    return s.planRepo.FindByCreator(userID)
}

// UpdateOperationPlanRequest represents the request to update an operation plan
type UpdateOperationPlanRequest struct {
    PartQuantity int    `json:"part_quantity"`
    Description  string `json:"description"`
}

// Update updates an operation plan (only if status is draft)
func (s *OperationPlanService) Update(id uint, req UpdateOperationPlanRequest) (*models.OperationPlan, error) {
    plan, err := s.planRepo.FindByID(id)
    if err != nil {
        return nil, errors.New("operation plan not found")
    }

    if plan.Status != models.StatusDraft {
        return nil, errors.New("cannot update operation plan that is not in draft status")
    }

    plan.PartQuantity = req.PartQuantity
    plan.Description = req.Description

    if err := s.planRepo.Update(plan); err != nil {
        return nil, errors.New("failed to update operation plan")
    }

    return s.planRepo.FindByID(id)
}

// Delete deletes an operation plan (only if status is draft)
func (s *OperationPlanService) Delete(id uint) error {
    plan, err := s.planRepo.FindByID(id)
    if err != nil {
        return errors.New("operation plan not found")
    }

    if plan.Status != models.StatusDraft {
        return errors.New("cannot delete operation plan that is not in draft status")
    }

    // Delete associated G-code files from filesystem
    for _, file := range plan.GCodeFiles {
        os.Remove(file.FilePath)
    }

    // Delete G-code file records
    if err := s.gCodeFileRepo.DeleteByOperationPlanID(id); err != nil {
        return errors.New("failed to delete G-code files")
    }

    return s.planRepo.Delete(id)
}

// SubmitForApproval submits an operation plan for approval
func (s *OperationPlanService) SubmitForApproval(id uint) (*models.OperationPlan, error) {
    plan, err := s.planRepo.FindByID(id)
    if err != nil {
        return nil, errors.New("operation plan not found")
    }

    if plan.Status != models.StatusDraft {
        return nil, errors.New("operation plan is not in draft status")
    }

    // Check if at least one G-code file is uploaded
    if len(plan.GCodeFiles) == 0 {
        return nil, errors.New("at least one G-code file must be uploaded before submitting")
    }

    if err := s.planRepo.SubmitForApproval(id); err != nil {
        return nil, errors.New("failed to submit for approval")
    }

    return s.planRepo.FindByID(id)
}

// Approve approves an operation plan
func (s *OperationPlanService) Approve(planID uint, approverID uint, approverRole string) (*models.OperationPlan, error) {
    plan, err := s.planRepo.FindByID(planID)
    if err != nil {
        return nil, errors.New("operation plan not found")
    }

    if plan.Status != models.StatusPendingApproval {
        return nil, errors.New("operation plan is not pending approval")
    }

    // Verify approver role is valid
    validRole := false
    for _, role := range models.ApproverRoles {
        if role == approverRole {
            validRole = true
            break
        }
    }
    if !validRole {
        return nil, errors.New("invalid approver role")
    }

    if err := s.planRepo.Approve(planID, approverID, approverRole); err != nil {
        return nil, err
    }

    // Reload plan to get updated status
    plan, _ = s.planRepo.FindByID(planID)

    // TODO: Send email notification if fully approved
    if plan.Status == models.StatusApproved {
        // s.sendApprovalNotification(plan)
    }

    return plan, nil
}

// GetApprovalStatus gets the approval status
func (s *OperationPlanService) GetApprovalStatus(planID uint) ([]models.OperationPlanApproval, error) {
    return s.planRepo.GetApprovalStatus(planID)
}

// UploadGCodeFile uploads a G-code file
func (s *OperationPlanService) UploadGCodeFile(planID uint, file *multipart.FileHeader, uploadedBy uint) (*models.GCodeFile, error) {
    // Verify operation plan exists
    plan, err := s.planRepo.FindByID(planID)
    if err != nil {
        return nil, errors.New("operation plan not found")
    }

    // Only allow uploads for draft plans
    if plan.Status != models.StatusDraft {
        return nil, errors.New("cannot upload files to operation plan that is not in draft status")
    }

    // Validate file extension
    ext := strings.ToLower(filepath.Ext(file.Filename))
    if ext != ".txt" {
        return nil, errors.New("only .txt files are allowed")
    }

    // Generate unique filename
    timestamp := time.Now().Format("20060102_150405")
    generatedName := fmt.Sprintf("%s_%s%s", plan.PlanNumber, timestamp, ext)
    filePath := filepath.Join(s.uploadPath, generatedName)

    // Open uploaded file
    src, err := file.Open()
    if err != nil {
        return nil, errors.New("failed to open uploaded file")
    }
    defer src.Close()

    // Create destination file
    dst, err := os.Create(filePath)
    if err != nil {
        return nil, errors.New("failed to create file on server")
    }
    defer dst.Close()

    // Copy file content
    if _, err := io.Copy(dst, src); err != nil {
        return nil, errors.New("failed to save file")
    }

    // Create database record
    gCodeFile := &models.GCodeFile{
        OperationPlanID: planID,
        FileName:        generatedName,
        OriginalName:    file.Filename,
        FilePath:        filePath,
        FileSize:        file.Size,
        UploadedBy:      uploadedBy,
    }

    if err := s.gCodeFileRepo.Create(gCodeFile); err != nil {
        os.Remove(filePath) // Clean up file on error
        return nil, errors.New("failed to save file record")
    }

    return s.gCodeFileRepo.FindByID(gCodeFile.ID)
}

// GetGCodeFile gets a G-code file by ID
func (s *OperationPlanService) GetGCodeFile(fileID uint) (*models.GCodeFile, error) {
    return s.gCodeFileRepo.FindByID(fileID)
}

// GetGCodeFiles gets all G-code files for an operation plan
func (s *OperationPlanService) GetGCodeFiles(planID uint) ([]models.GCodeFile, error) {
    return s.gCodeFileRepo.FindByOperationPlanID(planID)
}

// DeleteGCodeFile deletes a G-code file
func (s *OperationPlanService) DeleteGCodeFile(fileID uint, userID uint) error {
    file, err := s.gCodeFileRepo.FindByID(fileID)
    if err != nil {
        return errors.New("file not found")
    }

    // Get the operation plan to check status
    plan, err := s.planRepo.FindByID(file.OperationPlanID)
    if err != nil {
        return errors.New("operation plan not found")
    }

    if plan.Status != models.StatusDraft {
        return errors.New("cannot delete files from operation plan that is not in draft status")
    }

    // Delete file from filesystem
    if err := os.Remove(file.FilePath); err != nil {
        // Log error but continue with database deletion
        fmt.Printf("Warning: failed to delete file from filesystem: %v\n", err)
    }

    return s.gCodeFileRepo.Delete(fileID)
}

// StartExecution marks the start of operation plan execution
func (s *OperationPlanService) StartExecution(planID uint) (*models.OperationPlan, error) {
    plan, err := s.planRepo.FindByID(planID)
    if err != nil {
        return nil, errors.New("operation plan not found")
    }

    if plan.Status != models.StatusApproved {
        return nil, errors.New("operation plan must be approved before starting execution")
    }

    now := time.Now()
    plan.StartTime = &now

    if err := s.planRepo.Update(plan); err != nil {
        return nil, errors.New("failed to update start time")
    }

    return s.planRepo.FindByID(planID)
}

// FinishExecution marks the finish of operation plan execution
func (s *OperationPlanService) FinishExecution(planID uint) (*models.OperationPlan, error) {
    plan, err := s.planRepo.FindByID(planID)
    if err != nil {
        return nil, errors.New("operation plan not found")
    }

    if plan.StartTime == nil {
        return nil, errors.New("operation plan execution has not started")
    }

    now := time.Now()
    plan.FinishTime = &now

    if err := s.planRepo.Update(plan); err != nil {
        return nil, errors.New("failed to update finish time")
    }

    return s.planRepo.FindByID(planID)
}