package repository

import (
	"server/errors"
	gravemodel "server/model/grave"

	"gorm.io/gorm"
)

type GraveRepository struct {
	db *gorm.DB
}

func NewGraveRepository(db *gorm.DB) *GraveRepository {
	return &GraveRepository{db: db}
}

func (r *GraveRepository) Tx(tx *gorm.DB) *GraveRepository {
	return &GraveRepository{db: tx}
}

func (r *GraveRepository) Create(grave *gravemodel.Grave) error {
	return r.db.Create(grave).Error
}

func (r *GraveRepository) Update(grave *gravemodel.Grave) error {
	result := r.db.Save(grave)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrGraveNotFound
	}
	return nil
}

func (r *GraveRepository) Delete(id uint) error {
	result := r.db.Delete(&gravemodel.Grave{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrGraveNotFound
	}
	return nil
}
