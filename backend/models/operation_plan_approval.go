package models

import (
    "time"
)

type OperationPlanApproval struct {
    ID              uint       `gorm:"primaryKey" json:"id"`
    OperationPlanID uint       `gorm:"not null;index" json:"operation_plan_id"`
    ApproverRole    string     `gorm:"size:20;not null" json:"approver_role"` // PEM, PPIC, QC, Engineering, Toolpather
    ApproverID      *uint      `json:"approver_id"`                           // Who approved (null if pending)
    Approver        *User      `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
    Status          string     `gorm:"size:20;default:'pending'" json:"status"` // pending, approved
    ApprovedAt      *time.Time `json:"approved_at"`
    CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (OperationPlanApproval) TableName() string {
    return "operation_plan_approvals"
}

// ToResponse converts to API response
func (a *OperationPlanApproval) ToResponse() OperationPlanApprovalResponse {
    response := OperationPlanApprovalResponse{
        ID:              a.ID,
        OperationPlanID: a.OperationPlanID,
        ApproverRole:    a.ApproverRole,
        ApproverID:      a.ApproverID,
        Status:          a.Status,
        ApprovedAt:      a.ApprovedAt,
        CreatedAt:       a.CreatedAt,
    }

    if a.Approver != nil {
        approverResponse := a.Approver.ToResponse()
        response.Approver = &approverResponse
    }

    return response
}

type OperationPlanApprovalResponse struct {
    ID              uint          `json:"id"`
    OperationPlanID uint          `json:"operation_plan_id"`
    ApproverRole    string        `json:"approver_role"`
    ApproverID      *uint         `json:"approver_id"`
    Approver        *UserResponse `json:"approver,omitempty"`
    Status          string        `json:"status"`
    ApprovedAt      *time.Time    `json:"approved_at"`
    CreatedAt       time.Time     `json:"created_at"`
}