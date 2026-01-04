package models

import (
	"time"

	"gorm.io/gorm"
)

// PEM Operation Plan status constants
const (
	PEMStatusDraft           = "draft"
	PEMStatusPendingApproval = "pending_approval"
	PEMStatusApproved        = "approved"
	PEMStatusRejected        = "rejected"
)

// PEM Approval status constants
const (
	ApprovalStatusPending  = "pending"
	ApprovalStatusApproved = "approved"
	ApprovalStatusRejected = "rejected"
)

// PEM Approval roles - 5 required approvers
const (
	ApproverRolePEM        = "PEM"
	ApproverRoleToolpather = "Toolpather"
	ApproverRoleQC         = "QC"
	ApproverRoleCustom1    = "Custom1" // Selected by PEM team
	ApproverRoleCustom2    = "Custom2" // Selected by PEM team
)

var PEMApproverRoles = []string{
	ApproverRolePEM,
	ApproverRoleToolpather,
	ApproverRoleQC,
	ApproverRoleCustom1,
	ApproverRoleCustom2,
}

// PEMOperationPlan represents a simplified operation plan for PEM workflow
type PEMOperationPlan struct {
	ID             int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	FormNumber     string         `gorm:"uniqueIndex;size:50;not null" json:"form_number"`                // Auto: FRM-YYYYMMDD-XXX
	PPICScheduleID *int64         `gorm:"index" json:"ppic_schedule_id,omitempty"`
	PPICSchedule   *PPICSchedule  `gorm:"foreignKey:PPICScheduleID" json:"ppic_schedule,omitempty"`
	PartName       string         `gorm:"size:255" json:"part_name"`
	Material       string         `gorm:"size:255" json:"material"`
	DialSize       string         `gorm:"size:100" json:"dial_size"`       // e.g., "16 x 135 mm"
	Quantity       int            `json:"quantity"`
	Revision       string         `gorm:"size:50" json:"revision"`
	NoWP           string         `gorm:"size:100" json:"no_wp"`
	Page           string         `gorm:"size:50" json:"page"`
	Status         string         `gorm:"size:50;default:'draft'" json:"status"` // draft, pending_approval, approved, rejected
	CreatedBy      int64          `gorm:"index;not null" json:"created_by"`
	Creator        *User          `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Steps          []OperationPlanStep `gorm:"foreignKey:OperationPlanID" json:"steps,omitempty"`
	Approvals      []PEMApproval  `gorm:"foreignKey:OperationPlanID" json:"approvals,omitempty"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PEMOperationPlan) TableName() string {
	return "pem_operation_plans"
}

// OperationPlanStep represents a single process step in the operation plan
type OperationPlanStep struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OperationPlanID int64     `gorm:"index;not null" json:"operation_plan_id"`
	StepNumber      int       `gorm:"not null" json:"step_number"`                // 1, 2, 3...
	PictureURL      string    `gorm:"size:500" json:"picture_url,omitempty"`      // Path to uploaded image
	PictureFilename string    `gorm:"size:255" json:"picture_filename,omitempty"` // Original filename
	ClampingSystem  string    `gorm:"type:text" json:"clamping_system"`
	RawMaterial     string    `gorm:"type:text" json:"raw_material"`
	Setting         string    `gorm:"type:text" json:"setting"`
	Process         string    `gorm:"type:text" json:"process"`
	Note            string    `gorm:"type:text" json:"note"`
	CheckingMethod  string    `gorm:"type:text" json:"checking_method"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (OperationPlanStep) TableName() string {
	return "operation_plan_steps"
}

// PEMApproval represents the approval record for a PEM operation plan
type PEMApproval struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	OperationPlanID int64      `gorm:"index;not null" json:"operation_plan_id"`
	ApproverRole    string     `gorm:"size:50;not null" json:"approver_role"` // PEM, Toolpather, QC, Custom1, Custom2
	ApproverID      *int64     `gorm:"index" json:"approver_id,omitempty"`
	Approver        *User      `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
	Status          string     `gorm:"size:50;default:'pending'" json:"status"` // pending, approved, rejected
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	Comments        string     `gorm:"type:text" json:"comments"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (PEMApproval) TableName() string {
	return "pem_approvals"
}

// Request DTOs

type CreatePEMPlanRequest struct {
	PPICScheduleID *int64                   `json:"ppic_schedule_id"`
	PartName       string                   `json:"part_name" binding:"required"`
	Material       string                   `json:"material"`
	DialSize       string                   `json:"dial_size"`
	Quantity       int                      `json:"quantity"`
	Revision       string                   `json:"revision"`
	NoWP           string                   `json:"no_wp"`
	Page           string                   `json:"page"`
	Steps          []CreateStepRequest      `json:"steps"`
}

type UpdatePEMPlanRequest struct {
	PartName       string                   `json:"part_name"`
	Material       string                   `json:"material"`
	DialSize       string                   `json:"dial_size"`
	Quantity       *int                     `json:"quantity"`
	Revision       string                   `json:"revision"`
	NoWP           string                   `json:"no_wp"`
	Page           string                   `json:"page"`
	Steps          []UpdateStepRequest      `json:"steps"`
}

type CreateStepRequest struct {
	StepNumber     int    `json:"step_number" binding:"required"`
	ClampingSystem string `json:"clamping_system"`
	RawMaterial    string `json:"raw_material"`
	Setting        string `json:"setting"`
	Process        string `json:"process"`
	Note           string `json:"note"`
	CheckingMethod string `json:"checking_method"`
}

type UpdateStepRequest struct {
	ID             int64  `json:"id"`
	StepNumber     int    `json:"step_number"`
	ClampingSystem string `json:"clamping_system"`
	RawMaterial    string `json:"raw_material"`
	Setting        string `json:"setting"`
	Process        string `json:"process"`
	Note           string `json:"note"`
	CheckingMethod string `json:"checking_method"`
}

type AssignApproversRequest struct {
	Approvers map[string]int64 `json:"approvers" binding:"required"` // Map of role to user ID
}

type ApprovalActionRequest struct {
	Comments string `json:"comments"`
}

// Response DTOs

type PEMOperationPlanResponse struct {
	ID             int64                      `json:"id"`
	FormNumber     string                     `json:"form_number"`
	PPICScheduleID *int64                     `json:"ppic_schedule_id,omitempty"`
	PPICSchedule   *PPICSchedule              `json:"ppic_schedule,omitempty"`
	PartName       string                     `json:"part_name"`
	Material       string                     `json:"material"`
	DialSize       string                     `json:"dial_size"`
	Quantity       int                        `json:"quantity"`
	Revision       string                     `json:"revision"`
	NoWP           string                     `json:"no_wp"`
	Page           string                     `json:"page"`
	Status         string                     `json:"status"`
	CreatedBy      int64                      `json:"created_by"`
	Creator        *UserResponse              `json:"creator,omitempty"`
	Steps          []OperationPlanStep        `json:"steps,omitempty"`
	Approvals      []PEMApproval              `json:"approvals,omitempty"`
	CreatedAt      time.Time                  `json:"created_at"`
	UpdatedAt      time.Time                  `json:"updated_at"`
}

// Validation functions

func ValidatePEMStatus(status string) bool {
	validStatuses := []string{PEMStatusDraft, PEMStatusPendingApproval, PEMStatusApproved, PEMStatusRejected}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

func ValidateApprovalStatus(status string) bool {
	validStatuses := []string{ApprovalStatusPending, ApprovalStatusApproved, ApprovalStatusRejected}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}
