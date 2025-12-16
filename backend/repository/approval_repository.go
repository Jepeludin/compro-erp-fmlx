package repository

import (
    "ganttpro-backend/models"
    "gorm.io/gorm"
)

type ApprovalRepository struct {
    db *gorm.DB
}

func NewApprovalRepository(db *gorm.DB) *ApprovalRepository {
    return &ApprovalRepository{db: db}
}

// CreateBatch creates multiple approval records
func (r *ApprovalRepository) CreateBatch(approvals []models.OperationPlanApproval) error {
    return r.db.Create(&approvals).Error
}

// FindByOperationPlanID finds all approvals for an operation plan
func (r *ApprovalRepository) FindByOperationPlanID(planID uint) ([]models.OperationPlanApproval, error) {
    var approvals []models.OperationPlanApproval
    err := r.db.Preload("Approver").
        Where("operation_plan_id = ?", planID).
        Find(&approvals).Error
    return approvals, err
}

// FindPendingByRoleAndPlanID finds pending approval for specific role
func (r *ApprovalRepository) FindPendingByRoleAndPlanID(planID uint, role string) (*models.OperationPlanApproval, error) {
    var approval models.OperationPlanApproval
    err := r.db.Where("operation_plan_id = ? AND approver_role = ? AND status = ?", 
        planID, role, models.StatusPendingApproval).
        First(&approval).Error
    return &approval, err
}

// Update updates an approval record
func (r *ApprovalRepository) Update(approval *models.OperationPlanApproval) error {
    return r.db.Save(approval).Error
}

// CountApprovedByPlanID counts approved approvals for a plan
func (r *ApprovalRepository) CountApprovedByPlanID(planID uint) (int64, error) {
    var count int64
    err := r.db.Model(&models.OperationPlanApproval{}).
        Where("operation_plan_id = ? AND status = ?", planID, models.StatusApproved).
        Count(&count).Error
    return count, err
}