package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"ganttpro-backend/models"
	"time"

	"gorm.io/gorm"
)

type TokenBlacklistRepository struct {
	db *gorm.DB
}

func NewTokenBlacklistRepository(db *gorm.DB) *TokenBlacklistRepository {
	return &TokenBlacklistRepository{db: db}
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (r *TokenBlacklistRepository) AddToBlacklist(token string, expiresAt time.Time) error {
	blacklistedToken := &models.TokenBlacklist{
		Token:     hashToken(token),
		ExpiresAt: expiresAt,
	}
	return r.db.Create(blacklistedToken).Error
}

// IsBlacklisted checks if a token is blacklisted and not expired
func (r *TokenBlacklistRepository) IsBlacklisted(token string) (bool, error) {
	var count int64
	hashedToken := hashToken(token)

	err := r.db.Model(&models.TokenBlacklist{}).
		Where("token = ? AND expires_at > ?", hashedToken, time.Now()).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Returns the number of tokens removed
func (r *TokenBlacklistRepository) CleanExpiredTokens() error {
	result := r.db.Where("expires_at < ?", time.Now()).
		Delete(&models.TokenBlacklist{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *TokenBlacklistRepository) CleanExpiredTokensWithCount() (int64, error) {
	result := r.db.Where("expires_at < ?", time.Now()).
		Delete(&models.TokenBlacklist{})

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

// GetBlacklistCount returns the total number of blacklisted tokens
func (r *TokenBlacklistRepository) GetBlacklistCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.TokenBlacklist{}).Count(&count).Error
	return count, err
}

// GetExpiredCount returns the number of expired tokens waiting for cleanup
func (r *TokenBlacklistRepository) GetExpiredCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.TokenBlacklist{}).
		Where("expires_at < ?", time.Now()).
		Count(&count).Error
	return count, err
}
