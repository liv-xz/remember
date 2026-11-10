package repository

import (
	"server/errors"
	mrmodel "server/model/memorial_relation"
	"server/utils/paginator"

	"gorm.io/gorm"
)

func (r *MemorialRelationRepository) GetByID(id uint) (*mrmodel.MemorialRelation, error) {
	var relation mrmodel.MemorialRelation
	err := r.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username", "nickname", "avatar")
	}).Preload("Deceased", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "portrait_url")
	}).First(&relation, id).Error
	if err != nil {
		if errors.IsRecordNotFound(err) {
			return nil, errors.ErrRelationNotFound
		}
		return nil, err
	}
	return &relation, nil
}

func (r *MemorialRelationRepository) ListByGrave(graveID uint, page, size int) ([]mrmodel.MemorialRelation, int64, error) {
	query := r.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username", "nickname", "avatar")
	}).Preload("Deceased", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "portrait_url")
	}).Where("grave_id = ?", graveID)
	return paginator.Paginate[mrmodel.MemorialRelation](query, page, size)
}

func (r *MemorialRelationRepository) ListByDeceased(deceasedID uint) ([]mrmodel.MemorialRelation, error) {
	var list []mrmodel.MemorialRelation
	err := r.db.Preload("User").
		Where("deceased_id = ?", deceasedID).
		Find(&list).Error
	return list, err
}

func (r *MemorialRelationRepository) ListByUser(userID uint) ([]mrmodel.MemorialRelation, error) {
	var list []mrmodel.MemorialRelation
	err := r.db.Preload("Deceased").
		Where("user_id = ?", userID).
		Find(&list).Error
	return list, err
}

func (r *MemorialRelationRepository) CheckExists(userID, deceasedID uint) (bool, error) {
	var count int64
	err := r.db.Model(&mrmodel.MemorialRelation{}).
		Where("user_id = ? AND deceased_id = ?", userID, deceasedID).
		Count(&count).Error
	return count > 0, err
}
