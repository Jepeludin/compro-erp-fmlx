package repository

import (
	"ganttpro-backend/models"

	"gorm.io/gorm"
)

type ToolpatherFileRepository struct {
	db *gorm.DB
}

func NewToolpatherFileRepository(db *gorm.DB) *ToolpatherFileRepository {
	return &ToolpatherFileRepository{db: db}
}

// Create creates a new toolpather file record
func (r *ToolpatherFileRepository) Create(file *models.ToolpatherFile) error {
	return r.db.Create(file).Error
}

// FindByID finds a file by ID with relations
func (r *ToolpatherFileRepository) FindByID(id int64) (*models.ToolpatherFile, error) {
	var file models.ToolpatherFile
	err := r.db.
		Preload("Uploader").
		Preload("PPICSchedule").
		First(&file, id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// FindAll retrieves all files with optional filters
func (r *ToolpatherFileRepository) FindAll(filters map[string]interface{}) ([]models.ToolpatherFile, error) {
	var files []models.ToolpatherFile
	query := r.db.
		Preload("Uploader").
		Preload("PPICSchedule")

	// Apply filters
	if orderNumber, ok := filters["order_number"].(string); ok && orderNumber != "" {
		query = query.Where("order_number = ?", orderNumber)
	}

	if uploadedBy, ok := filters["uploaded_by"].(int64); ok && uploadedBy > 0 {
		query = query.Where("uploaded_by = ?", uploadedBy)
	}

	if ppicScheduleID, ok := filters["ppic_schedule_id"].(int64); ok && ppicScheduleID > 0 {
		query = query.Where("ppic_schedule_id = ?", ppicScheduleID)
	}

	err := query.Order("created_at DESC").Find(&files).Error
	return files, err
}

// FindByOrderNumber retrieves all files for a specific order number
func (r *ToolpatherFileRepository) FindByOrderNumber(orderNumber string) ([]models.ToolpatherFile, error) {
	var files []models.ToolpatherFile
	err := r.db.
		Preload("Uploader").
		Preload("PPICSchedule").
		Where("order_number = ?", orderNumber).
		Order("created_at DESC").
		Find(&files).Error
	return files, err
}

// FindByUploader retrieves all files uploaded by a specific user
func (r *ToolpatherFileRepository) FindByUploader(uploaderID int64) ([]models.ToolpatherFile, error) {
	var files []models.ToolpatherFile
	err := r.db.
		Preload("Uploader").
		Preload("PPICSchedule").
		Where("uploaded_by = ?", uploaderID).
		Order("created_at DESC").
		Find(&files).Error
	return files, err
}

// Update updates a file record
func (r *ToolpatherFileRepository) Update(file *models.ToolpatherFile) error {
	return r.db.Save(file).Error
}

// Delete deletes a file record
func (r *ToolpatherFileRepository) Delete(id int64) error {
	return r.db.Delete(&models.ToolpatherFile{}, id).Error
}
