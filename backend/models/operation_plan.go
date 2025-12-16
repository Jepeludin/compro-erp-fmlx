package models

import "time"
import "gorm.io/gorm"

const (
    StatusDraft           = "draft"
    StatusPendingApproval = "pending_approval"
    StatusApproved        = "approved"
)

type OperationPlan struct {
    ID           uint           `gorm:"primaryKey" json:"id"`
    PlanNumber   string         `gorm:"size:50;uniqueIndex;not null" json:"plan_number"` // Auto-generated: OP-YYYYMMDD-XXX
    JobOrderID   uint           `gorm:"not null;index" json:"job_order_id"`
    JobOrder     *JobOrder      `gorm:"foreignKey:JobOrderID" json:"job_order,omitempty"`
    MachineID    uint           `gorm:"not null;index" json:"machine_id"`
    Machine      *Machine       `gorm:"foreignKey:MachineID" json:"machine,omitempty"`
    PartQuantity int            `gorm:"default:1" json:"part_quantity"`
    Description  string         `gorm:"type:text" json:"description"`
    Status       string         `gorm:"size:20;default:'draft'" json:"status"` // draft, pending_approval, approved
    CreatedBy    uint           `gorm:"not null" json:"created_by"`
    Creator      *User          `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
    StartTime    *time.Time     `json:"start_time"`
    FinishTime   *time.Time     `json:"finish_time"`
    CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

    // Relations
    Approvals  []OperationPlanApproval `gorm:"foreignKey:OperationPlanID" json:"approvals,omitempty"`
    GCodeFiles []GCodeFile             `gorm:"foreignKey:OperationPlanID" json:"g_code_files,omitempty"`
}

func (OperationPlan) TableName() string {
    return "operation_plans"
}

// ToResponse converts to API response
func (op *OperationPlan) ToResponse() OperationPlanResponse {
    response := OperationPlanResponse{
        ID:           op.ID,
        PlanNumber:   op.PlanNumber,
        JobOrderID:   op.JobOrderID,
        MachineID:    op.MachineID,
        PartQuantity: op.PartQuantity,
        Description:  op.Description,
        Status:       op.Status,
        CreatedBy:    op.CreatedBy,
        StartTime:    op.StartTime,
        FinishTime:   op.FinishTime,
        CreatedAt:    op.CreatedAt,
        UpdatedAt:    op.UpdatedAt,
    }

    if op.JobOrder != nil {
        response.JobOrder = op.JobOrder
    }
    if op.Machine != nil {
        response.Machine = op.Machine
    }
    if op.Creator != nil {
        creatorResponse := op.Creator.ToResponse()
        response.Creator = &creatorResponse
    }

    return response
}

type OperationPlanResponse struct {
    ID           uint          `json:"id"`
    PlanNumber   string        `json:"plan_number"`
    JobOrderID   uint          `json:"job_order_id"`
    JobOrder     *JobOrder     `json:"job_order,omitempty"`
    MachineID    uint          `json:"machine_id"`
    Machine      *Machine      `json:"machine,omitempty"`
    PartQuantity int           `json:"part_quantity"`
    Description  string        `json:"description"`
    Status       string        `json:"status"`
    CreatedBy    uint          `json:"created_by"`
    Creator      *UserResponse `json:"creator,omitempty"`
    StartTime    *time.Time    `json:"start_time"`
    FinishTime   *time.Time    `json:"finish_time"`
    CreatedAt    time.Time     `json:"created_at"`
    UpdatedAt    time.Time     `json:"updated_at"`
}