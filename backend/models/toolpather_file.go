package models

import (
	"time"

	"gorm.io/gorm"
)

// ToolpatherFile represents a file uploaded by toolpather for a specific order
type ToolpatherFile struct {
	ID             int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	PPICScheduleID *int64         `gorm:"index" json:"ppic_schedule_id,omitempty"`
	PPICSchedule   *PPICSchedule  `gorm:"foreignKey:PPICScheduleID" json:"ppic_schedule,omitempty"`
	OrderNumber    string         `gorm:"size:50;index" json:"order_number"`      // NJO
	PartName       string         `gorm:"size:255" json:"part_name"`
	FileName       string         `gorm:"size:255;not null" json:"file_name"`     // Original filename
	FilePath       string         `gorm:"size:500;not null" json:"file_path"`     // Relative path from upload directory
	FileSize       int64          `json:"file_size"`                              // File size in bytes
	UploadedBy     int64          `gorm:"index;not null" json:"uploaded_by"`
	Uploader       *User          `gorm:"foreignKey:UploadedBy" json:"uploader,omitempty"`
	Notes          string         `gorm:"type:text" json:"notes"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ToolpatherFile) TableName() string {
	return "toolpather_files"
}

// Request DTOs

type UploadToolpatherFileRequest struct {
	PPICScheduleID *int64 `json:"ppic_schedule_id"`
	OrderNumber    string `json:"order_number" binding:"required"`
	PartName       string `json:"part_name"`
	Notes          string `json:"notes"`
}

// Response DTOs

type ToolpatherFileResponse struct {
	ID             int64          `json:"id"`
	PPICScheduleID *int64         `json:"ppic_schedule_id,omitempty"`
	OrderNumber    string         `json:"order_number"`
	PartName       string         `json:"part_name"`
	FileName       string         `json:"file_name"`
	FilePath       string         `json:"file_path"`
	FileSize       int64          `json:"file_size"`
	UploadedBy     int64          `json:"uploaded_by"`
	Uploader       *UserResponse  `json:"uploader,omitempty"`
	Notes          string         `json:"notes"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
