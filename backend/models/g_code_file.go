package models

import (
    "time"

    "gorm.io/gorm"
)

type GCodeFile struct {
    ID              uint           `gorm:"primaryKey" json:"id"`
    OperationPlanID uint           `gorm:"not null;index" json:"operation_plan_id"`
    FileName        string         `gorm:"size:255;not null" json:"file_name"`    // Auto-generated name
    OriginalName    string         `gorm:"size:255;not null" json:"original_name"` // Original upload name
    FilePath        string         `gorm:"type:text;not null" json:"file_path"`   
    FileSize        int64          `gorm:"not null" json:"file_size"`             // Size in bytes
    UploadedBy      uint           `gorm:"not null" json:"uploaded_by"`
    Uploader        *User          `gorm:"foreignKey:UploadedBy" json:"uploader,omitempty"`
    CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
    DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (GCodeFile) TableName() string {
    return "g_code_files"
}

// ToResponse converts to API response
func (g *GCodeFile) ToResponse() GCodeFileResponse {
    response := GCodeFileResponse{
        ID:              g.ID,
        OperationPlanID: g.OperationPlanID,
        FileName:        g.FileName,
        OriginalName:    g.OriginalName,
        FilePath:        g.FilePath,
        FileSize:        g.FileSize,
        UploadedBy:      g.UploadedBy,
        CreatedAt:       g.CreatedAt,
    }

    if g.Uploader != nil {
        uploaderResponse := g.Uploader.ToResponse()
        response.Uploader = &uploaderResponse
    }

    return response
}

type GCodeFileResponse struct {
    ID              uint          `json:"id"`
    OperationPlanID uint          `json:"operation_plan_id"`
    FileName        string        `json:"file_name"`
    OriginalName    string        `json:"original_name"`
    FilePath        string        `json:"file_path"`
    FileSize        int64         `json:"file_size"`
    UploadedBy      uint          `json:"uploaded_by"`
    Uploader        *UserResponse `json:"uploader,omitempty"`
    CreatedAt       time.Time     `json:"created_at"`
}