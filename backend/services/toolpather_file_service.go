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

type ToolpatherFileService struct {
	repo      *repository.ToolpatherFileRepository
	userRepo  *repository.UserRepository
	uploadDir string
}

func NewToolpatherFileService(
	repo *repository.ToolpatherFileRepository,
	userRepo *repository.UserRepository,
	uploadDir string,
) *ToolpatherFileService {
	// Create upload directory if not exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		panic(fmt.Sprintf("Failed to create toolpather upload directory: %v", err))
	}

	return &ToolpatherFileService{
		repo:      repo,
		userRepo:  userRepo,
		uploadDir: uploadDir,
	}
}

// UploadFiles uploads multiple .txt files for an order number
func (s *ToolpatherFileService) UploadFiles(
	request models.UploadToolpatherFileRequest,
	files []*multipart.FileHeader,
	uploaderID int64,
) ([]models.ToolpatherFile, error) {
	if len(files) == 0 {
		return nil, errors.New("no files provided")
	}

	// Verify uploader exists
	_, err := s.userRepo.FindByID(uint(uploaderID))
	if err != nil {
		return nil, errors.New("uploader not found")
	}

	uploadedFiles := make([]models.ToolpatherFile, 0)

	// Create order-specific subdirectory
	orderDir := filepath.Join(s.uploadDir, sanitizeFileName(request.OrderNumber))
	if err := os.MkdirAll(orderDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create order directory: %w", err)
	}

	for _, file := range files {
		// Validate file extension (only .txt files)
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".txt" {
			return nil, fmt.Errorf("file %s: only .txt files are allowed", file.Filename)
		}

		// Validate file size (max 10MB)
		maxSize := int64(10 * 1024 * 1024) // 10MB
		if file.Size > maxSize {
			return nil, fmt.Errorf("file %s: file size exceeds 10MB limit", file.Filename)
		}

		// Generate unique filename with timestamp
		timestamp := time.Now().Unix()
		baseName := strings.TrimSuffix(filepath.Base(file.Filename), ext)
		uniqueFileName := fmt.Sprintf("%s_%d%s", sanitizeFileName(baseName), timestamp, ext)
		filePath := filepath.Join(orderDir, uniqueFileName)

		// Save file to disk
		src, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", file.Filename, err)
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create file %s: %w", file.Filename, err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return nil, fmt.Errorf("failed to save file %s: %w", file.Filename, err)
		}

		// Create database record
		// Store relative path from upload directory
		relPath := filepath.Join(sanitizeFileName(request.OrderNumber), uniqueFileName)

		toolpatherFile := models.ToolpatherFile{
			PPICScheduleID: request.PPICScheduleID,
			OrderNumber:    request.OrderNumber,
			PartName:       request.PartName,
			FileName:       file.Filename, // Original filename
			FilePath:       relPath,       // Relative path
			FileSize:       file.Size,
			UploadedBy:     uploaderID,
			Notes:          request.Notes,
		}

		if err := s.repo.Create(&toolpatherFile); err != nil {
			// Try to delete the file if database insert fails
			os.Remove(filePath)
			return nil, fmt.Errorf("failed to create database record for %s: %w", file.Filename, err)
		}

		// Load relations for response
		uploadedFile, _ := s.repo.FindByID(toolpatherFile.ID)
		if uploadedFile != nil {
			uploadedFiles = append(uploadedFiles, *uploadedFile)
		}
	}

	return uploadedFiles, nil
}

// GetFileByID retrieves a file by ID
func (s *ToolpatherFileService) GetFileByID(id int64) (*models.ToolpatherFile, error) {
	return s.repo.FindByID(id)
}

// GetAllFiles retrieves all files with optional filters
func (s *ToolpatherFileService) GetAllFiles(filters map[string]interface{}) ([]models.ToolpatherFile, error) {
	return s.repo.FindAll(filters)
}

// GetFilesByOrderNumber retrieves all files for a specific order number
func (s *ToolpatherFileService) GetFilesByOrderNumber(orderNumber string) ([]models.ToolpatherFile, error) {
	return s.repo.FindByOrderNumber(orderNumber)
}

// GetFilesByUploader retrieves all files uploaded by a specific user
func (s *ToolpatherFileService) GetFilesByUploader(uploaderID int64) ([]models.ToolpatherFile, error) {
	return s.repo.FindByUploader(uploaderID)
}

// DeleteFile deletes a file (only uploader or admin can delete)
func (s *ToolpatherFileService) DeleteFile(id int64, userID int64) error {
	file, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Get user to check if admin
	user, err := s.userRepo.FindByID(uint(userID))
	if err != nil {
		return errors.New("user not found")
	}

	// Only uploader or admin can delete
	if file.UploadedBy != userID && user.Role != "Admin" {
		return errors.New("only the uploader or admin can delete this file")
	}

	// Delete file from disk
	if file.FilePath != "" {
		filePath := filepath.Join(s.uploadDir, file.FilePath)
		if err := os.Remove(filePath); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Warning: failed to delete file %s: %v\n", filePath, err)
		}
	}

	return s.repo.Delete(id)
}

// GetFilePath returns the full path to a file
func (s *ToolpatherFileService) GetFilePath(file *models.ToolpatherFile) string {
	return filepath.Join(s.uploadDir, file.FilePath)
}

// sanitizeFileName removes unsafe characters from filename
func sanitizeFileName(filename string) string {
	// Replace unsafe characters with underscores
	unsafe := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", " "}
	result := filename
	for _, char := range unsafe {
		result = strings.ReplaceAll(result, char, "_")
	}
	return result
}
