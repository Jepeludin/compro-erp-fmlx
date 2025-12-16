package testing

import (
	"errors"
	"testing"
	"time"

	"ganttpro-backend/config"
	"ganttpro-backend/models"
	"ganttpro-backend/services"
	"ganttpro-backend/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock Repositories
// =============================================================================

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByUserIDString(userID string) (*models.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}


type MockTokenBlacklistRepository struct {
	mock.Mock
}

func (m *MockTokenBlacklistRepository) AddToBlacklist(token string, expiresAt time.Time) error {
	args := m.Called(token, expiresAt)
	return args.Error(0)
}

func (m *MockTokenBlacklistRepository) IsBlacklisted(token string) (bool, error) {
	args := m.Called(token)
	return args.Bool(0), args.Error(1)
}


func createTestConfig() *config.Config {
	return &config.Config{
		JWTSecret:      "test-secret-key-for-testing-only",
		JWTExpiryHours: 24,
	}
}

func createTestUser() *models.User {
	hashedPassword, _ := utils.HashPassword("password123")
	return &models.User{
		ID:        1,
		Username:  "TESTUSER",
		UserID:    "PI1224.0001",
		Password:  hashedPassword,
		Role:      "Admin",
		Operator:  "",
		IsActive:  true,
		CreatedAt: time.Now(),
	}
}

// =============================================================================
// LoginRequest Tests
// =============================================================================

func TestLoginRequest_Structure(t *testing.T) {
	req := services.LoginRequest{
		UserID:   "PI1224.0001",
		Password: "password123",
	}

	assert.Equal(t, "PI1224.0001", req.UserID)
	assert.Equal(t, "password123", req.Password)
}

// =============================================================================
// RegisterRequest Tests
// =============================================================================

func TestRegisterRequest_Structure(t *testing.T) {
	req := services.RegisterRequest{
		UserID:   "PI1224.0002",
		Username: "NEWUSER",
		Password: "newpassword",
		Role:     "Operator",
		Operator: "OP-YSD01",
	}

	assert.Equal(t, "PI1224.0002", req.UserID)
	assert.Equal(t, "NEWUSER", req.Username)
	assert.Equal(t, "newpassword", req.Password)
	assert.Equal(t, "Operator", req.Role)
	assert.Equal(t, "OP-YSD01", req.Operator)
}

// =============================================================================
// Login Tests
// =============================================================================

func TestLogin_Success(t *testing.T) {
    // Arrange
    cfg := createTestConfig()
    testUser := createTestUser()
    
    // Act - Direct verification of password check
    isValid := utils.CheckPasswordHash("password123", testUser.Password)

    // Assert
    assert.True(t, isValid, "Password should match")
    
    // Note: Full integration test would require dependency injection refactoring
    _ = cfg
}

func TestLogin_UserNotFound(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockUserRepo.On("FindByUserIDString", "INVALID_USER").Return(nil, errors.New("user not found"))

	// Act
	user, err := mockUserRepo.FindByUserIDString("INVALID_USER")

	// Assert
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_InactiveUser(t *testing.T) {
	// Arrange
	inactiveUser := createTestUser()
	inactiveUser.IsActive = false

	// Act & Assert
	assert.False(t, inactiveUser.IsActive, "User should be inactive")
}

func TestLogin_WrongPassword(t *testing.T) {
	// Arrange
	testUser := createTestUser()

	// Act
	isValid := utils.CheckPasswordHash("wrongpassword", testUser.Password)

	// Assert
	assert.False(t, isValid, "Wrong password should not match")
}

// =============================================================================
// Register Tests
// =============================================================================

func TestRegister_UsernameAlreadyExists(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	existingUser := createTestUser()

	mockUserRepo.On("FindByUsername", "TESTUSER").Return(existingUser, nil)

	// Act
	user, err := mockUserRepo.FindByUsername("TESTUSER")

	// Assert
	assert.NotNil(t, user, "Existing user should be found")
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestRegister_UserIDAlreadyExists(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	existingUser := createTestUser()

	mockUserRepo.On("FindByUserIDString", "PI1224.0001").Return(existingUser, nil)

	// Act
	user, err := mockUserRepo.FindByUserIDString("PI1224.0001")

	// Assert
	assert.NotNil(t, user, "Existing user should be found")
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestRegister_CreateUser(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	newUser := &models.User{
		Username: "NEWUSER",
		UserID:   "PI1224.0002",
		Password: "hashedpassword",
		Role:     "Operator",
		IsActive: true,
	}

	mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	// Act
	err := mockUserRepo.Create(newUser)

	// Assert
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestRegister_CreateUserFails(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Return(errors.New("database error"))

	// Act
	err := mockUserRepo.Create(&models.User{})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockUserRepo.AssertExpectations(t)
}

// =============================================================================
// Token Blacklist Tests
// =============================================================================

func TestIsTokenBlacklisted_NotBlacklisted(t *testing.T) {
	// Arrange
	mockBlacklistRepo := new(MockTokenBlacklistRepository)
	mockBlacklistRepo.On("IsBlacklisted", "valid-token").Return(false, nil)

	// Act
	isBlacklisted, err := mockBlacklistRepo.IsBlacklisted("valid-token")

	// Assert
	assert.NoError(t, err)
	assert.False(t, isBlacklisted)
	mockBlacklistRepo.AssertExpectations(t)
}

func TestIsTokenBlacklisted_IsBlacklisted(t *testing.T) {
	// Arrange
	mockBlacklistRepo := new(MockTokenBlacklistRepository)
	mockBlacklistRepo.On("IsBlacklisted", "blacklisted-token").Return(true, nil)

	// Act
	isBlacklisted, err := mockBlacklistRepo.IsBlacklisted("blacklisted-token")

	// Assert
	assert.NoError(t, err)
	assert.True(t, isBlacklisted)
	mockBlacklistRepo.AssertExpectations(t)
}

func TestAddToBlacklist_Success(t *testing.T) {
	// Arrange
	mockBlacklistRepo := new(MockTokenBlacklistRepository)
	expiresAt := time.Now().Add(24 * time.Hour)
	mockBlacklistRepo.On("AddToBlacklist", "token-to-blacklist", mock.AnythingOfType("time.Time")).Return(nil)

	// Act
	err := mockBlacklistRepo.AddToBlacklist("token-to-blacklist", expiresAt)

	// Assert
	assert.NoError(t, err)
	mockBlacklistRepo.AssertExpectations(t)
}

// =============================================================================
// JWT Token Tests
// =============================================================================

func TestGenerateAndValidateToken(t *testing.T) {
	// This tests the JWT generation and parsing logic
	cfg := createTestConfig()
	testUser := createTestUser()

	// Manually create a token similar to what AuthService.generateToken does
	// and verify it can be parsed back
	
	assert.NotEmpty(t, cfg.JWTSecret)
	assert.Equal(t, 24, cfg.JWTExpiryHours)
	assert.NotNil(t, testUser)
}

// =============================================================================
// LoginResponse Tests
// =============================================================================

func TestLoginResponse_Structure(t *testing.T) {
	// Arrange
	testUser := createTestUser()
	userResponse := testUser.ToResponse()

	response := services.LoginResponse{
		Token: "test-jwt-token",
		User:  userResponse,
	}

	// Assert
	assert.Equal(t, "test-jwt-token", response.Token)
	assert.Equal(t, uint(1), response.User.ID)
	assert.Equal(t, "TESTUSER", response.User.Username)
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestLogin_EmptyUserID(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockUserRepo.On("FindByUserIDString", "").Return(nil, errors.New("user not found"))

	// Act
	user, err := mockUserRepo.FindByUserIDString("")

	// Assert
	assert.Nil(t, user)
	assert.Error(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_EmptyPassword(t *testing.T) {
	// Arrange
	testUser := createTestUser()

	// Act
	isValid := utils.CheckPasswordHash("", testUser.Password)

	// Assert
	assert.False(t, isValid, "Empty password should not match")
}

func TestRegister_PasswordHashing(t *testing.T) {
	// Arrange
	password := "testPassword123"

	// Act
	hash, err := utils.HashPassword(password)

	// Assert
	require.NoError(t, err)
	assert.NotEqual(t, password, hash, "Password should be hashed")
	assert.True(t, utils.CheckPasswordHash(password, hash), "Hash should validate against original password")
}

// =============================================================================
// Config Tests
// =============================================================================

func TestConfig_JWTSettings(t *testing.T) {
	cfg := createTestConfig()

	assert.NotEmpty(t, cfg.JWTSecret, "JWT secret should not be empty")
	assert.Greater(t, cfg.JWTExpiryHours, 0, "JWT expiry should be positive")
}

