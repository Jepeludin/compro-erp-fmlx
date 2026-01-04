package repository

import (
	"errors"
	"fmt"
	"ganttpro-backend/models"
	"time"

	"gorm.io/gorm"
)

type PEMOperationPlanRepository struct {
	db *gorm.DB
}

func NewPEMOperationPlanRepository(db *gorm.DB) *PEMOperationPlanRepository {
	return &PEMOperationPlanRepository{db: db}
}

// GenerateFormNumber generates a unique form number: FRM-YYYYMMDD-XXX
func (r *PEMOperationPlanRepository) GenerateFormNumber() (string, error) {
	today := time.Now().Format("20060102")
	prefix := fmt.Sprintf("FRM-%s", today)

	var count int64
	err := r.db.Model(&models.PEMOperationPlan{}).
		Where("form_number LIKE ?", prefix+"%").
		Count(&count).Error

	if err != nil {
		return "", err
	}

	formNumber := fmt.Sprintf("%s-%03d", prefix, count+1)
	return formNumber, nil
}

// Create creates a new PEM operation plan with approval records for all 5 roles
func (r *PEMOperationPlanRepository) Create(plan *models.PEMOperationPlan) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Generate form number if not provided
		if plan.FormNumber == "" {
			formNumber, err := r.GenerateFormNumber()
			if err != nil {
				return err
			}
			plan.FormNumber = formNumber
		}

		// Create the operation plan
		if err := tx.Create(plan).Error; err != nil {
			return err
		}

		// Create approval records for all 5 required roles
		for _, role := range models.PEMApproverRoles {
			approval := &models.PEMApproval{
				OperationPlanID: plan.ID,
				ApproverRole:    role,
				Status:          models.ApprovalStatusPending,
			}
			if err := tx.Create(approval).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// FindByID finds a PEM operation plan by ID with all relations
func (r *PEMOperationPlanRepository) FindByID(id int64) (*models.PEMOperationPlan, error) {
	var plan models.PEMOperationPlan
	err := r.db.Preload("PPICSchedule").
		Preload("Creator").
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("step_number ASC")
		}).
		Preload("Approvals").
		Preload("Approvals.Approver").
		First(&plan, id).Error

	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// FindByFormNumber finds a PEM operation plan by form number
func (r *PEMOperationPlanRepository) FindByFormNumber(formNumber string) (*models.PEMOperationPlan, error) {
	var plan models.PEMOperationPlan
	err := r.db.Preload("PPICSchedule").
		Preload("Creator").
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("step_number ASC")
		}).
		Preload("Approvals").
		Preload("Approvals.Approver").
		Where("form_number = ?", formNumber).
		First(&plan).Error

	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// FindAll finds all PEM operation plans with optional filters
func (r *PEMOperationPlanRepository) FindAll(filters map[string]interface{}) ([]models.PEMOperationPlan, error) {
	var plans []models.PEMOperationPlan
	query := r.db.Preload("PPICSchedule").
		Preload("Creator").
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("step_number ASC")
		}).
		Preload("Approvals").
		Preload("Approvals.Approver")

	// Apply filters
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if createdBy, ok := filters["created_by"].(int64); ok && createdBy > 0 {
		query = query.Where("created_by = ?", createdBy)
	}
	if ppicScheduleID, ok := filters["ppic_schedule_id"].(int64); ok && ppicScheduleID > 0 {
		query = query.Where("ppic_schedule_id = ?", ppicScheduleID)
	}

	err := query.Order("created_at DESC").Find(&plans).Error
	return plans, err
}

// Update updates a PEM operation plan
func (r *PEMOperationPlanRepository) Update(plan *models.PEMOperationPlan) error {
	return r.db.Save(plan).Error
}

// UpdateStatus updates the status of a PEM operation plan
func (r *PEMOperationPlanRepository) UpdateStatus(id int64, status string) error {
	return r.db.Model(&models.PEMOperationPlan{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete soft deletes a PEM operation plan
func (r *PEMOperationPlanRepository) Delete(id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete all steps first
		if err := tx.Where("operation_plan_id = ?", id).Delete(&models.OperationPlanStep{}).Error; err != nil {
			return err
		}

		// Delete approval records
		if err := tx.Where("operation_plan_id = ?", id).Delete(&models.PEMApproval{}).Error; err != nil {
			return err
		}

		// Delete the plan itself
		if err := tx.Delete(&models.PEMOperationPlan{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

// Step Management Methods

// CreateStep creates a new operation plan step
func (r *PEMOperationPlanRepository) CreateStep(step *models.OperationPlanStep) error {
	return r.db.Create(step).Error
}

// FindStepByID finds a step by ID
func (r *PEMOperationPlanRepository) FindStepByID(stepID int64) (*models.OperationPlanStep, error) {
	var step models.OperationPlanStep
	err := r.db.First(&step, stepID).Error
	if err != nil {
		return nil, err
	}
	return &step, nil
}

// UpdateStep updates an operation plan step
func (r *PEMOperationPlanRepository) UpdateStep(step *models.OperationPlanStep) error {
	return r.db.Save(step).Error
}

// DeleteStep deletes an operation plan step
func (r *PEMOperationPlanRepository) DeleteStep(stepID int64) error {
	return r.db.Delete(&models.OperationPlanStep{}, stepID).Error
}

// Approval Management Methods

// AssignApprovers assigns approvers to all 5 roles
func (r *PEMOperationPlanRepository) AssignApprovers(planID int64, approvers map[string]int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for role, userID := range approvers {
			result := tx.Model(&models.PEMApproval{}).
				Where("operation_plan_id = ? AND approver_role = ?", planID, role).
				Updates(map[string]interface{}{
					"approver_id": userID,
					"status":      models.ApprovalStatusPending,
				})

			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return fmt.Errorf("approval record not found for role: %s", role)
			}
		}
		return nil
	})
}

// SubmitForApproval changes status from draft to pending_approval
func (r *PEMOperationPlanRepository) SubmitForApproval(id int64) error {
	return r.db.Model(&models.PEMOperationPlan{}).
		Where("id = ? AND status = ?", id, models.PEMStatusDraft).
		Update("status", models.PEMStatusPendingApproval).Error
}

// ApprovePlan approves the operation plan by a specific role
func (r *PEMOperationPlanRepository) ApprovePlan(planID int64, approverID int64, role string, comments string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Update the approval record for this role
		result := tx.Model(&models.PEMApproval{}).
			Where("operation_plan_id = ? AND approver_role = ? AND approver_id = ? AND status = ?",
				planID, role, approverID, models.ApprovalStatusPending).
			Updates(map[string]interface{}{
				"status":      models.ApprovalStatusApproved,
				"approved_at": now,
				"comments":    comments,
			})

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("approval not found or already processed")
		}

		// Check if all approvals are completed
		var pendingCount int64
		err := tx.Model(&models.PEMApproval{}).
			Where("operation_plan_id = ? AND status = ?", planID, models.ApprovalStatusPending).
			Count(&pendingCount).Error

		if err != nil {
			return err
		}

		// If all approved, update operation plan status
		if pendingCount == 0 {
			err = tx.Model(&models.PEMOperationPlan{}).
				Where("id = ?", planID).
				Update("status", models.PEMStatusApproved).Error

			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RejectPlan rejects the operation plan
func (r *PEMOperationPlanRepository) RejectPlan(planID int64, approverID int64, role string, comments string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Update the approval record
		result := tx.Model(&models.PEMApproval{}).
			Where("operation_plan_id = ? AND approver_role = ? AND approver_id = ? AND status = ?",
				planID, role, approverID, models.ApprovalStatusPending).
			Updates(map[string]interface{}{
				"status":      models.ApprovalStatusRejected,
				"approved_at": now,
				"comments":    comments,
			})

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("approval not found or already processed")
		}

		// Update operation plan status to rejected
		err := tx.Model(&models.PEMOperationPlan{}).
			Where("id = ?", planID).
			Update("status", models.PEMStatusRejected).Error

		if err != nil {
			return err
		}

		return nil
	})
}

// GetPendingApprovalsByApprover gets all plans pending approval for a specific approver
func (r *PEMOperationPlanRepository) GetPendingApprovalsByApprover(approverID int64) ([]models.PEMOperationPlan, error) {
	var plans []models.PEMOperationPlan

	err := r.db.Joins("JOIN pem_approvals ON pem_approvals.operation_plan_id = pem_operation_plans.id").
		Where("pem_approvals.approver_id = ? AND pem_approvals.status = ?", approverID, models.ApprovalStatusPending).
		Preload("PPICSchedule").
		Preload("Creator").
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("step_number ASC")
		}).
		Preload("Approvals").
		Preload("Approvals.Approver").
		Group("pem_operation_plans.id").
		Order("pem_operation_plans.created_at DESC").
		Find(&plans).Error

	return plans, err
}

// GetApprovalByPlanAndRole gets approval record for a specific plan and role
func (r *PEMOperationPlanRepository) GetApprovalByPlanAndRole(planID int64, role string) (*models.PEMApproval, error) {
	var approval models.PEMApproval
	err := r.db.Preload("Approver").
		Where("operation_plan_id = ? AND approver_role = ?", planID, role).
		First(&approval).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &approval, nil
}
