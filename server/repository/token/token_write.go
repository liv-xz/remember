package repository

import (
	"server/errors"
	rtmodel "server/model/token"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(rt *rtmodel.RefreshToken) error {
	return r.db.Create(rt).Error
}

func (r *RefreshTokenRepository) DeleteInvalidByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).
		Where("expires_at < ? OR revoked = ?", time.Now(), true).
		Delete(&rtmodel.RefreshToken{}).Error
}

func (r *RefreshTokenRepository) Revoke(id uint) error {
	result := r.db.Model(&rtmodel.RefreshToken{}).Where("id = ?", id).Update("revoked", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrRefreshTokenNotFound
	}
	return nil
}

func (r *RefreshTokenRepository) RevokeAllByUserID(userID uint) error {
	return r.db.Model(&rtmodel.RefreshToken{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Update("revoked", true).Error
}

func (r *RefreshTokenRepository) RevokeByHash(hash string) error {
	result := r.db.Model(&rtmodel.RefreshToken{}).
		Where("token_hash = ? AND revoked = ?", hash, false).
		Update("revoked", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrRefreshTokenNotFound
	}
	return nil
}
