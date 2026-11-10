package repository

import (
	"server/errors"
	eventmodel "server/model/memorial_event"

	"gorm.io/gorm"
)

type MemorialEventRepository struct {
	db *gorm.DB
}

func NewMemorialEventRepository(db *gorm.DB) *MemorialEventRepository {
	return &MemorialEventRepository{db: db}
}

func (r *MemorialEventRepository) Create(event *eventmodel.MemorialEvent) error {
	return r.db.Create(event).Error
}

func (r *MemorialEventRepository) Update(event *eventmodel.MemorialEvent) error {
	result := r.db.Save(event)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrEventNotFound
	}
	return nil
}

func (r *MemorialEventRepository) Delete(id uint) error {
	result := r.db.Delete(&eventmodel.MemorialEvent{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrEventNotFound
	}
	return nil
}
