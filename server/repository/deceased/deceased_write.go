package repository

import (
	"server/errors"
	deceasedmodel "server/model/deceased"

	"gorm.io/gorm"
)

type DeceasedRepository struct {
	db *gorm.DB
}

func NewDeceasedRepository(db *gorm.DB) *DeceasedRepository {
	return &DeceasedRepository{db: db}
}

func (r *DeceasedRepository) Create(deceased *deceasedmodel.Deceased) error {
	return r.db.Create(deceased).Error
}

func (r *DeceasedRepository) Update(deceased *deceasedmodel.Deceased) error {
	result := r.db.Save(deceased)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrDeceasedNotFound
	}
	return nil
}

func (r *DeceasedRepository) Delete(id uint) error {
	result := r.db.Delete(&deceasedmodel.Deceased{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrDeceasedNotFound
	}
	return nil
}
