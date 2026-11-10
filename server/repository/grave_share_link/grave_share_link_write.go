package repository

import (
	"gorm.io/gorm"

	gslmodel "server/model/grave_share_link"
)

type GraveShareLinkRepository struct {
	db *gorm.DB
}

func NewGraveShareLinkRepository(db *gorm.DB) *GraveShareLinkRepository {
	return &GraveShareLinkRepository{db: db}
}

func (r *GraveShareLinkRepository) Create(link *gslmodel.GraveShareLink) error {
	return r.db.Create(link).Error
}

func (r *GraveShareLinkRepository) UpdateUsedCount(id uint) error {
	return r.db.Model(&gslmodel.GraveShareLink{}).
		Where("id = ?", id).
		UpdateColumn("used_count", gorm.Expr("used_count + ?", 1)).Error
}

func (r *GraveShareLinkRepository) Deactivate(id uint) error {
	return r.db.Model(&gslmodel.GraveShareLink{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}
