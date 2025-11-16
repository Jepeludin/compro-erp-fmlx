package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string         `gorm:"uniqueIndex;not null;size:50" json:"username"` // e.g., BAYU, Amelia
	UserID    string         `gorm:"uniqueIndex;not null;size:50" json:"user_id"`  // e.g., PI0824.2374
	Password  string         `gorm:"not null" json:"-"`                            // Hashed password
	Role      string         `gorm:"size:20;default:'PPIC'" json:"role"`           // Admin, PPIC, etc.
	Operator  string         `gorm:"size:50" json:"operator"`                      // Operator field
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	Role     string `json:"role"`
	Operator string `json:"operator"`
	IsActive bool   `json:"is_active"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{

		Username: u.Username,
		UserID:   u.UserID,
		Role:     u.Role,
		Operator: u.Operator,
		IsActive: u.IsActive,
	}
}
