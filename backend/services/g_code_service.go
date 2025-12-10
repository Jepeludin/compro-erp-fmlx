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

type GCodeService struct {
    gcodeRepo    *repository.GCodeFileRepository
    opPlanRepo   *repository.OperationPlanRepository
    uploadDir    string
}

func NewGCodeService(
    gcodeRepo *repository.GCodeFileRepository,
    opPlanRepo *repository.OperationPlanRepository,
    uploadDir string,
) *GCodeService {
    // Create upload directory if not exists
    if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
        panic(fmt.Sprintf("Failed to create upload directory: %v", err))
    }

    return &GCodeService{
        gcodeRepo:  gcodeRepo,
        opPlanRepo: opPlanRepo,
        uploadDir:  uploadDir,
    }
}

// UploadGCode uploads a G-code file for an operation plan
func (s *GCodeService) UploadGCode(planID uint, file *multipart.FileHeader, uploadedBy uint) (*models.GCodeFile, error) {
    // Validate operation plan exists and is in draft status
    plan, err := s.opPlanRepo.FindByID(planID)
    if err != nil {
        return nil, errors.New("operation plan not found")
    }

    if plan.Status != models.StatusPendingApproval {
        return nil, errors.New("can only upload files to plans in draft status")
    }

    // Validate file extension (.txt only)
    ext := strings.ToLower(filepath.Ext(file.Filename))
    if ext != ".txt" {
        return nil, errors.New("only .txt files are allowed")
    }

    // Validate file size (max 10MB)
    maxSize := int64(10 * 1024 * 1024) // 10MB
    if file.Size > maxSize {
        return nil, errors.New("file size exceeds 10MB limit")
    }

    // Generate unique filename: {plan_number}_{timestamp}.txt
    timestamp := time.Now().Unix()
    fileName := fmt.Sprintf("%s_%d.txt", plan.PlanNumber, timestamp)
    filePath := filepath.Join(s.uploadDir, fileName)

    // Save file to disk
    src, err := file.Open()
    if err != nil {
        return nil, errors.New("failed to open uploaded file")
    }
    defer src.Close()

    dst, err := os.Create(filePath)
    if err != nil {
        return nil, errors.New("failed to create file on disk")
    }
    defer dst.Close()

    if _, err := io.Copy(dst, src); err != nil {
        return nil, errors.New("failed to save file")
    }

    // Create database record
    gcodeFile := &models.GCodeFile{
        OperationPlanID: planID,
        FileName:        fileName,
        OriginalName:    file.Filename,
        FilePath:        filePath,
        FileSize:        file.Size,
        UploadedBy:      uploadedBy,
    }

    if err := s.gcodeRepo.Create(gcodeFile); err != nil {
        // Rollback: delete file from disk
        os.Remove(filePath)
        return nil, errors.New("failed to save file record")
    }

    // Load uploader info
    gcodeFile, _ = s.gcodeRepo.FindByID(gcodeFile.ID)
    return gcodeFile, nil
}

// GetGCodeFile retrieves a G-code file by ID
func (s *GCodeService) GetGCodeFile(id uint) (*models.GCodeFile, error) {
    return s.gcodeRepo.FindByID(id)
}

// GetGCodeFilesByPlan retrieves all G-code files for an operation plan
func (s *GCodeService) GetGCodeFilesByPlan(planID uint) ([]models.GCodeFile, error) {
    return s.gcodeRepo.FindByOperationPlanID(planID)
}

// DeleteGCodeFile deletes a G-code file (only if plan is in draft status)
func (s *GCodeService) DeleteGCodeFile(id uint, userID uint) error {
    file, err := s.gcodeRepo.FindByID(id)
    if err != nil {
        return errors.New("file not found")
    }

    // Check if uploader is the same user
    if file.UploadedBy != userID {
        return errors.New("only the uploader can delete this file")
    }

    // Check operation plan status
    plan, err := s.opPlanRepo.FindByID(file.OperationPlanID)
    if err != nil {
        return errors.New("operation plan not found")
    }

    if plan.Status != models.StatusDraft {
        return errors.New("can only delete files from plans in draft status")
    }

    // Delete from database
    if err := s.gcodeRepo.Delete(id); err != nil {
        return errors.New("failed to delete file record")
    }

    // Delete from disk
    if err := os.Remove(file.FilePath); err != nil {
        // Log error but don't fail the operation
        fmt.Printf("Warning: Failed to delete file from disk: %v\n", err)
    }

    return nil
}

// GetFilePath returns the file path for download
func (s *GCodeService) GetFilePath(id uint) (string, string, error) {
    file, err := s.gcodeRepo.FindByID(id)
    if err != nil {
        return "", "", errors.New("file not found")
    }

    // Check if file exists on disk
    if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
        return "", "", errors.New("file not found on disk")
    }

    return file.FilePath, file.OriginalName, nil
}