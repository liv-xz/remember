package repository

import (
	"server/errors"
	recordmodel "server/model/memorial_record"

	"gorm.io/gorm"
)

type MemorialRecordRepository struct {
	db *gorm.DB
}

func NewMemorialRecordRepository(db *gorm.DB) *MemorialRecordRepository {
	return &MemorialRecordRepository{db: db}
}

func (r *MemorialRecordRepository) Create(record *recordmodel.MemorialRecord) error {
	return r.db.Create(record).Error
}

func (r *MemorialRecordRepository) Delete(id, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&recordmodel.MemorialRecord{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrRecordNotFound
	}
	return nil
}
