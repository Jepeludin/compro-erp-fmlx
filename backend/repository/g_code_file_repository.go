package repository

import (
    "ganttpro-backend/models"

    "gorm.io/gorm"
)

type GCodeFileRepository struct {
    db *gorm.DB
}

func NewGCodeFileRepository(db *gorm.DB) *GCodeFileRepository {
    return &GCodeFileRepository{db: db}
}

// Create creates a new G-code file record
func (r *GCodeFileRepository) Create(file *models.GCodeFile) error {
    return r.db.Create(file).Error
}

// FindByID finds a G-code file by ID
func (r *GCodeFileRepository) FindByID(id uint) (*models.GCodeFile, error) {
    var file models.GCodeFile
    err := r.db.Preload("Uploader").First(&file, id).Error
    if err != nil {
        return nil, err
    }
    return &file, nil
}

// FindByOperationPlanID finds all G-code files for an operation plan
func (r *GCodeFileRepository) FindByOperationPlanID(operationPlanID uint) ([]models.GCodeFile, error) {
    var files []models.GCodeFile
    err := r.db.Preload("Uploader").
        Where("operation_plan_id = ?", operationPlanID).
        Order("created_at DESC").
        Find(&files).Error

    return files, err
}

// Delete soft deletes a G-code file record
func (r *GCodeFileRepository) Delete(id uint) error {
    return r.db.Delete(&models.GCodeFile{}, id).Error
}

// DeleteByOperationPlanID deletes all G-code files for an operation plan
func (r *GCodeFileRepository) DeleteByOperationPlanID(operationPlanID uint) error {
    return r.db.Where("operation_plan_id = ?", operationPlanID).
        Delete(&models.GCodeFile{}).Error
}