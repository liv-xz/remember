package repository

import (
	"gorm.io/gorm"

	"server/errors"
	krmodel "server/model/kinship_relation"
)

type KinshipRelationRepository struct {
	db *gorm.DB
}

func NewKinshipRelationRepository(db *gorm.DB) *KinshipRelationRepository {
	return &KinshipRelationRepository{db: db}
}

func (r *KinshipRelationRepository) Tx(tx *gorm.DB) *KinshipRelationRepository {
	return &KinshipRelationRepository{db: tx}
}

func (r *KinshipRelationRepository) Create(relation *krmodel.KinshipRelation) error {
	return r.db.Create(relation).Error
}

func (r *KinshipRelationRepository) Update(id uint, relation *krmodel.KinshipRelation) error {
	result := r.db.Model(&krmodel.KinshipRelation{}).Where("id = ?", id).Updates(relation)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrRelationNotFound
	}
	return nil
}

func (r *KinshipRelationRepository) Delete(id uint) error {
	result := r.db.Delete(&krmodel.KinshipRelation{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrRelationNotFound
	}
	return nil
}

func (r *KinshipRelationRepository) DeleteByDeceased(deceasedID uint) error {
	return r.db.Where("deceased_id = ?", deceasedID).Delete(&krmodel.KinshipRelation{}).Error
}
