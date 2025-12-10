package repository

import (
    "fmt"
    "ganttpro-backend/models"
    "time"

    "gorm.io/gorm"
)

type OperationPlanRepository struct {
    db *gorm.DB
}

func NewOperationPlanRepository(db *gorm.DB) *OperationPlanRepository {
    return &OperationPlanRepository{db: db}
}

// GeneratePlanNumber generates a unique plan number: OP-YYYYMMDD-XXX
func (r *OperationPlanRepository) GeneratePlanNumber() (string, error) {
    today := time.Now().Format("20060102")
    prefix := fmt.Sprintf("OP-%s", today)

    var count int64
    err := r.db.Model(&models.OperationPlan{}).
        Where("plan_number LIKE ?", prefix+"%").
        Count(&count).Error

    if err != nil {
        return "", err
    }

    planNumber := fmt.Sprintf("%s-%03d", prefix, count+1)
    return planNumber, nil
}

// Create creates a new operation plan with approval records
func (r *OperationPlanRepository) Create(plan *models.OperationPlan) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        // Generate plan number
        planNumber, err := r.GeneratePlanNumber()
        if err != nil {
            return err
        }
        plan.PlanNumber = planNumber

        // Create the operation plan
        if err := tx.Create(plan).Error; err != nil {
            return err
        }

        // Create approval records for all 5 required roles
        for _, role := range models.ApproverRoles {
            approval := &models.OperationPlanApproval{
                OperationPlanID: plan.ID,
                ApproverRole:    role,
                Status:          "pending",
            }
            if err := tx.Create(approval).Error; err != nil {
                return err
            }
        }

        return nil
    })
}

// FindByID finds an operation plan by ID with all relations
func (r *OperationPlanRepository) FindByID(id uint) (*models.OperationPlan, error) {
    var plan models.OperationPlan
    err := r.db.Preload("JobOrder").
        Preload("Machine").
        Preload("Creator").
        Preload("Approvals").
        Preload("Approvals.Approver").
        Preload("GCodeFiles").
        Preload("GCodeFiles.Uploader").
        First(&plan, id).Error

    if err != nil {
        return nil, err
    }
    return &plan, nil
}

// FindByPlanNumber finds an operation plan by plan number
func (r *OperationPlanRepository) FindByPlanNumber(planNumber string) (*models.OperationPlan, error) {
    var plan models.OperationPlan
    err := r.db.Preload("JobOrder").
        Preload("Machine").
        Preload("Creator").
        Preload("Approvals").
        Preload("Approvals.Approver").
        Preload("GCodeFiles").
        Where("plan_number = ?", planNumber).
        First(&plan).Error

    if err != nil {
        return nil, err
    }
    return &plan, nil
}

// FindByJobOrderID finds operation plan by job order ID
func (r *OperationPlanRepository) FindByJobOrderID(jobOrderID uint) (*models.OperationPlan, error) {
    var plan models.OperationPlan
    err := r.db.Preload("JobOrder").
        Preload("Machine").
        Preload("Creator").
        Preload("Approvals").
        Preload("Approvals.Approver").
        Preload("GCodeFiles").
        Where("job_order_id = ?", jobOrderID).
        First(&plan).Error

    if err != nil {
        return nil, err
    }
    return &plan, nil
}

// FindAll finds all operation plans with optional filters
func (r *OperationPlanRepository) FindAll(status string, machineID uint) ([]models.OperationPlan, error) {
    var plans []models.OperationPlan
    query := r.db.Preload("JobOrder").
        Preload("Machine").
        Preload("Creator").
        Preload("Approvals").
        Preload("GCodeFiles")

    if status != "" {
        query = query.Where("status = ?", status)
    }
    if machineID > 0 {
        query = query.Where("machine_id = ?", machineID)
    }

    err := query.Order("created_at DESC").Find(&plans).Error
    return plans, err
}

// FindByCreator finds operation plans created by a specific user
func (r *OperationPlanRepository) FindByCreator(userID uint) ([]models.OperationPlan, error) {
    var plans []models.OperationPlan
    err := r.db.Preload("JobOrder").
        Preload("Machine").
        Preload("Approvals").
        Preload("GCodeFiles").
        Where("created_by = ?", userID).
        Order("created_at DESC").
        Find(&plans).Error

    return plans, err
}

// Update updates an operation plan
func (r *OperationPlanRepository) Update(plan *models.OperationPlan) error {
    return r.db.Save(plan).Error
}

// UpdateStatus updates the status of an operation plan
func (r *OperationPlanRepository) UpdateStatus(id uint, status string) error {
    return r.db.Model(&models.OperationPlan{}).
        Where("id = ?", id).
        Update("status", status).Error
}

// Delete soft deletes an operation plan
func (r *OperationPlanRepository) Delete(id uint) error {
    return r.db.Delete(&models.OperationPlan{}, id).Error
}

// SubmitForApproval changes status from draft to pending_approval
func (r *OperationPlanRepository) SubmitForApproval(id uint) error {
    return r.db.Model(&models.OperationPlan{}).
        Where("id = ? AND status = ?", id, models.StatusDraft).
        Update("status", models.StatusPendingApproval).Error
}

// Approve approves an operation plan by a specific role
func (r *OperationPlanRepository) Approve(planID uint, approverID uint, approverRole string) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        now := time.Now()

        // Update the approval record
        result := tx.Model(&models.OperationPlanApproval{}).
            Where("operation_plan_id = ? AND approver_role = ? AND status = ?", planID, approverRole, "pending").
            Updates(map[string]interface{}{
                "approver_id": approverID,
                "status":      "approved",
                "approved_at": now,
            })

        if result.Error != nil {
            return result.Error
        }

        if result.RowsAffected == 0 {
            return fmt.Errorf("approval not found or already approved")
        }

        // Check if all approvals are completed
        var pendingCount int64
        err := tx.Model(&models.OperationPlanApproval{}).
            Where("operation_plan_id = ? AND status = ?", planID, "pending").
            Count(&pendingCount).Error

        if err != nil {
            return err
        }

        // If all approved, update operation plan status
        if pendingCount == 0 {
            err = tx.Model(&models.OperationPlan{}).
                Where("id = ?", planID).
                Update("status", models.StatusApproved).Error

            if err != nil {
                return err
            }
        }

        return nil
    })
}

// GetApprovalStatus gets the approval status for an operation plan
func (r *OperationPlanRepository) GetApprovalStatus(planID uint) ([]models.OperationPlanApproval, error) {
    var approvals []models.OperationPlanApproval
    err := r.db.Preload("Approver").
        Where("operation_plan_id = ?", planID).
        Find(&approvals).Error

    return approvals, err
}

// IsFullyApproved checks if all approvals are completed
func (r *OperationPlanRepository) IsFullyApproved(planID uint) (bool, error) {
    var pendingCount int64
    err := r.db.Model(&models.OperationPlanApproval{}).
        Where("operation_plan_id = ? AND status = ?", planID, "pending").
        Count(&pendingCount).Error

    if err != nil {
        return false, err
    }

    return pendingCount == 0, nil
}