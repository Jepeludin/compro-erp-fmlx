package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	RoleAdmin       = "Admin"
	RolePPIC        = "PPIC"
	RoleToolpather  = "Toolpather"
	RolePEM         = "PEM"
	RoleQC          = "QC"
	RoleEngineering = "Engineering"
	RoleGuest       = "Guest"
	RoleOperator    = "Operator"
)

var ValidRoles = []string{
	RoleAdmin,
	RolePPIC,
	RoleToolpather,
	RolePEM,
	RoleQC,
	RoleEngineering,
	RoleGuest,
	RoleOperator,
}

func IsValidRole(role string) bool {
	for _, r := range ValidRoles {
		if r == role {
			return true
		}
	}
	return false
}

func GetRoleDisplayName(role string) string {
	switch role {
	case RoleAdmin:
		return "Administrator"
	case RolePPIC:
		return "PPIC Staff"
	case RoleToolpather:
		return "Toolpather"
	case RolePEM:
		return "PEM Staff"
	case RoleQC:
		return "Quality Control"
	case RoleEngineering:
		return "Engineering"
	case RoleOperator:
		return "Machine Operator"
	case RoleGuest:
		return "Guest"
	default:
		return "Unknown"
	}
}

var ApproverRoles = []string{
	RolePEM,
	RolePPIC,
	RoleQC,
	RoleEngineering,
	RoleToolpather,
}

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
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"`
	Operator  string    `json:"operator"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		UserID:    u.UserID,
		Role:      u.Role,
		Operator:  u.Operator,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}
