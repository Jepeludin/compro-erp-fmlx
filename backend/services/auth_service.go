package services

import (
	"errors"
	"time"

	"ganttpro-backend/config"
	"ganttpro-backend/models"
	"ganttpro-backend/repository"
	"ganttpro-backend/utils"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo      *repository.UserRepository
	blacklistRepo *repository.TokenBlacklistRepository
	config        *config.Config
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo *repository.UserRepository, blacklistRepo *repository.TokenBlacklistRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		blacklistRepo: blacklistRepo,
		config:        cfg,
	}
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	UserID   string `json:"user_id" binding:"required"`  // User ID (e.g., PI0824.0001, PI0824.2374)
	Password string `json:"password" binding:"required"` // Plain password to be verified
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string              `json:"token"`
	User  models.UserResponse `json:"user"`
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	// Find user by UserID (e.g., PI0824.0001)
	user, err := s.userRepo.FindByUserIDString(req.UserID)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is disabled")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

// generateToken generates a JWT token for the user
func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":        user.ID,
		"username":       user.Username,
		"user_id_string": user.UserID, // PI0824.2374 format
		"role":           user.Role,
		"operator":       user.Operator,
		"exp":            time.Now().Add(time.Hour * time.Duration(s.config.JWTExpiryHours)).Unix(),
		"iat":            time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	UserID   string `json:"user_id" binding:"required"`               // e.g., PI0824.2374
	Username string `json:"username" binding:"required,min=3,max=50"` // e.g., BAYU
	Password string `json:"password" binding:"required,min=4"`        // Password (will be hashed)
	Role     string `json:"role" binding:"required"`                  // Admin, PPIC, etc.
	Operator string `json:"operator"`                                 // Operator field
}

// Register creates a new user account
func (s *AuthService) Register(req RegisterRequest) (*models.UserResponse, error) {
	// Check if username already exists
	existingUser, _ := s.userRepo.FindByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if user_id already exists
	existingUserID, _ := s.userRepo.FindByUserIDString(req.UserID)
	if existingUserID != nil {
		return nil, errors.New("user_id already exists")
	} // Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	// Create user with updated structure
	user := &models.User{
		Username: req.Username,
		UserID:   req.UserID,
		Password: hashedPassword,
		Role:     req.Role,
		Operator: req.Operator,
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	response := user.ToResponse()
	return &response, nil
}

// ValidateToken validates a JWT token and returns the user
func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		return s.userRepo.FindByID(userID)
	}

	return nil, errors.New("invalid token")
}

func (s *AuthService) Logout(tokenString string) error {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("invalid expiration time")
	}

	expiresAt := time.Unix(int64(exp), 0)

	if err := s.blacklistRepo.AddToBlacklist(tokenString, expiresAt); err != nil {
		return errors.New("failed to logout")
	}

	return nil
}

func (s *AuthService) IsTokenBlacklisted(tokenString string) (bool, error) {
	return s.blacklistRepo.IsBlacklisted(tokenString)
}
