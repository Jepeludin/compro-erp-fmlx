package models

import "time"

type TokenBlacklist struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Token     string    `gorm:"type:text;not null;index" json:"token"`
    ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (TokenBlacklist) TableName() string {
	return "token_blacklists"
}