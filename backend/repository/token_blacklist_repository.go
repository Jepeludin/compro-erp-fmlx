package repository

import (
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

func (r *TokenBlacklistRepository) AddToBlacklist(token string, expiresAt time.Time) error {
	blacklistedToken := &models.TokenBlacklist{
		Token:     token,
		ExpiresAt: expiresAt,
	}
	return r.db.Create(blacklistedToken).Error
}

func (r *TokenBlacklistRepository) IsBlacklisted(token string) (bool, error) {
	var count int64
	err := r.db.Model(&models.TokenBlacklist{}).
		Where("token = ? AND expires_at > ?", token, time.Now()).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *TokenBlacklistRepository) CleanExpiredTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).
		Delete(&models.TokenBlacklist{}).Error
}
