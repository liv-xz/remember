package repository

import (
	"gorm.io/gorm"

	"server/errors"
	mrmodel "server/model/memorial_relation"
)

type MemorialRelationRepository struct {
	db *gorm.DB
}

func NewMemorialRelationRepository(db *gorm.DB) *MemorialRelationRepository {
	return &MemorialRelationRepository{db: db}
}

func (r *MemorialRelationRepository) Create(relation *mrmodel.MemorialRelation) error {
	return r.db.Create(relation).Error
}

func (r *MemorialRelationRepository) Update(id uint, relation *mrmodel.MemorialRelation) error {
	result := r.db.Model(&mrmodel.MemorialRelation{}).Where("id = ?", id).Updates(relation)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrRelationNotFound
	}
	return nil
}

func (r *MemorialRelationRepository) Delete(id uint) error {
	result := r.db.Delete(&mrmodel.MemorialRelation{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrRelationNotFound
	}
	return nil
}

func (r *MemorialRelationRepository) DeleteByDeceased(deceasedID uint) error {
	return r.db.Where("deceased_id = ?", deceasedID).Delete(&mrmodel.MemorialRelation{}).Error
}
