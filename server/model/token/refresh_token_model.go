package model

import (
	"server/model"
	"time"
)

type RefreshToken struct {
	model.Model
	UserID    uint      `gorm:"index;not null" json:"userID"`
	TokenHash string    `gorm:"size:64;not null;uniqueIndex" json:"-"`
	ExpiresAt time.Time `gorm:"not null" json:"expiresAt"`
	Revoked   bool      `gorm:"default:false" json:"-"`
}
