package repository

import (
	"errors"

	"ganttpro-backend/models"

	"gorm.io/gorm"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByUsername finds a user by their Username
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ? AND is_active = ?", username, true).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// FindByUserIDString finds a user by their UserID string (e.g., PI0824.2374)
func (r *UserRepository) FindByUserIDString(userid string) (*models.User, error) {
	var user models.User
	err := r.db.Where("user_id = ? AND is_active = ?", userid, true).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// FindByID finds a user by their ID
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// FindByRole finds users by their role
func (r *UserRepository) FindByRole(role string) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("role = ? AND is_active = ?", role, true).Find(&users).Error
	return users, err
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// FindAll returns all active users
func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Where("is_active = ?", true).Find(&users).Error
	return users, err
}

// GetAllUsers returns all users (including inactive) for admin
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Order("created_at DESC").Find(&users).Error
	return users, err
}

// GetByID finds a user by ID (alias for FindByID for consistency)
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	return r.FindByID(id)
}
