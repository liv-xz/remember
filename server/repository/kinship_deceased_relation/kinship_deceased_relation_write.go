package repository

import (
	"server/errors"
	kdr "server/model/kinship_deceased_relation"

	"gorm.io/gorm"
)

type DeceasedRelationRepository struct {
	db *gorm.DB
}

func NewDeceasedRelationRepository(db *gorm.DB) *DeceasedRelationRepository {
	return &DeceasedRelationRepository{db: db}
}

func (r *DeceasedRelationRepository) Tx(tx *gorm.DB) *DeceasedRelationRepository {
	return &DeceasedRelationRepository{db: tx}
}

func (r *DeceasedRelationRepository) Create(relation *kdr.DeceasedRelation) error {
	return r.db.Create(relation).Error
}

func (r *DeceasedRelationRepository) Delete(id uint) error {
	result := r.db.Delete(&kdr.DeceasedRelation{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrRelationNotFound
	}
	return nil
}

func (r *DeceasedRelationRepository) DeleteByDeceased(deceasedID uint) error {
	return r.db.Where("deceased_id1 = ? OR deceased_id2 = ?", deceasedID, deceasedID).Delete(&kdr.DeceasedRelation{}).Error
}
