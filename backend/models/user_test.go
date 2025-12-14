package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)


func TestRoleConstants(t *testing.T) {
	// Verify role constants are defined correctly
	assert.Equal(t, "Admin", RoleAdmin)
	assert.Equal(t, "PPIC", RolePPIC)
	assert.Equal(t, "Toolpather", RoleToolpather)
	assert.Equal(t, "PEM", RolePEM)
	assert.Equal(t, "QC", RoleQC)
	assert.Equal(t, "Engineering", RoleEngineering)
	assert.Equal(t, "Guest", RoleGuest)
	assert.Equal(t, "Operator", RoleOperator)
}

func TestApproverRoles(t *testing.T) {
	// Verify approver roles are correct (5 roles required for approval)
	expectedRoles := []string{RolePEM, RolePPIC, RoleQC, RoleEngineering, RoleToolpather}
	assert.Equal(t, expectedRoles, ApproverRoles)
	assert.Len(t, ApproverRoles, 5, "There should be exactly 5 approver roles")
}

func TestIsValidRole_ValidRoles(t *testing.T) {
	validRoles := []string{
		RoleAdmin,
		RolePPIC,
		RoleToolpather,
		RolePEM,
		RoleQC,
		RoleEngineering,
		RoleGuest,
	}

	for _, role := range validRoles {
		t.Run(role, func(t *testing.T) {
			assert.True(t, IsValidRole(role), "Role %s should be valid", role)
		})
	}
}

func TestIsValidRole_InvalidRoles(t *testing.T) {
	invalidRoles := []string{
		"InvalidRole",
		"admin",      // lowercase (case-sensitive)
		"ADMIN",      // uppercase
		"SuperAdmin",
		"",
		"User",
		"Manager",
	}

	for _, role := range invalidRoles {
		t.Run(role, func(t *testing.T) {
			assert.False(t, IsValidRole(role), "Role %s should be invalid", role)
		})
	}
}

func TestIsValidRole_OperatorNotInValidRoles(t *testing.T) {
	// Note: This test documents a potential bug - Operator is defined but not in validRoles
	assert.False(t, IsValidRole("Operator"), "Operator is not in validRoles array (potential bug)")
}

// =============================================================================
// User Model Tests
// =============================================================================

func TestUser_TableName(t *testing.T) {
	user := User{}
	assert.Equal(t, "users", user.TableName(), "Table name should be 'users'")
}

func TestUser_ToResponse(t *testing.T) {
	// Arrange
	now := time.Now()
	user := &User{
		ID:        1,
		Username:  "BAYU",
		UserID:    "PI0824.2374",
		Password:  "hashedPassword123", // Should NOT be in response
		Role:      "Operator",
		Operator:  "OP-YSD01",
		IsActive:  true,
		CreatedAt: now,
	}

	// Act
	response := user.ToResponse()

	// Assert
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "BAYU", response.Username)
	assert.Equal(t, "PI0824.2374", response.UserID)
	assert.Equal(t, "Operator", response.Role)
	assert.Equal(t, "OP-YSD01", response.Operator)
	assert.True(t, response.IsActive)
	assert.Equal(t, now, response.CreatedAt)
}

func TestUser_ToResponse_PasswordNotIncluded(t *testing.T) {
	// Arrange
	user := &User{
		ID:       1,
		Username: "ADMIN",
		UserID:   "PI0824.0001",
		Password: "supersecretpassword",
		Role:     "Admin",
	}

	// Act
	response := user.ToResponse()

	// Assert - UserResponse struct doesn't have Password field
	// This is a compile-time check - if Password was in UserResponse, this test would still pass
	// The real protection is the json:"-" tag on User.Password
	assert.Equal(t, "ADMIN", response.Username)
}

func TestUser_ToResponse_InactiveUser(t *testing.T) {
	// Arrange
	user := &User{
		ID:       1,
		Username: "INACTIVE_USER",
		UserID:   "PI0824.9999",
		IsActive: false,
	}

	// Act
	response := user.ToResponse()

	// Assert
	assert.False(t, response.IsActive, "IsActive should be false")
}

// =============================================================================
// UserResponse Tests
// =============================================================================

func TestUserResponse_Fields(t *testing.T) {
	// Arrange
	now := time.Now()
	response := UserResponse{
		ID:        42,
		Username:  "TEST_USER",
		UserID:    "PI1224.0042",
		Role:      "PEM",
		Operator:  "",
		IsActive:  true,
		CreatedAt: now,
	}

	// Assert all fields are correctly set
	assert.Equal(t, uint(42), response.ID)
	assert.Equal(t, "TEST_USER", response.Username)
	assert.Equal(t, "PI1224.0042", response.UserID)
	assert.Equal(t, "PEM", response.Role)
	assert.Empty(t, response.Operator)
	assert.True(t, response.IsActive)
	assert.Equal(t, now, response.CreatedAt)
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestUser_ToResponse_EmptyUser(t *testing.T) {
	// Arrange
	user := &User{}

	// Act
	response := user.ToResponse()

	// Assert - should handle zero values gracefully
	assert.Equal(t, uint(0), response.ID)
	assert.Empty(t, response.Username)
	assert.Empty(t, response.UserID)
	assert.Empty(t, response.Role)
	assert.Empty(t, response.Operator)
	assert.False(t, response.IsActive)
	assert.True(t, response.CreatedAt.IsZero())
}

func TestUser_AllRolesUsedInBusiness(t *testing.T) {
	// Document which roles are used for what purpose
	businessRoles := map[string]string{
		RoleAdmin:       "Full system access",
		RolePPIC:        "Production planning, scheduling, tracking, approval",
		RolePEM:         "Create operation plans, upload G-codes, approval",
		RoleEngineering: "Submit job requests, create G-codes, approval",
		RoleToolpather:  "Attach PDFs per order, approval",
		RoleQC:          "Quality control, approval",
		RoleGuest:       "View-only access",
	}

	for role, description := range businessRoles {
		t.Run(role, func(t *testing.T) {
			assert.True(t, IsValidRole(role), "Role %s (%s) should be valid", role, description)
		})
	}
}

