package repository

import (
	"server/errors"
	eventmodel "server/model/memorial_event"
	"server/utils/paginator"
	"time"
)

func (r *MemorialEventRepository) ListByGrave(graveID uint, page, size int) ([]eventmodel.MemorialEvent, int64, error) {
	query := r.db.Where("grave_id = ?", graveID)
	return paginator.Paginate[eventmodel.MemorialEvent](query, page, size)
}

func (r *MemorialEventRepository) ListByDate(graveID uint, month time.Time) ([]eventmodel.MemorialEvent, error) {
	var list []eventmodel.MemorialEvent
	err := r.db.Where("grave_id = ? AND MONTH(event_date) = ?", graveID, month.Month()).Find(&list).Error
	return list, err
}

func (r *MemorialEventRepository) GetByID(id uint) (*eventmodel.MemorialEvent, error) {
	var event eventmodel.MemorialEvent
	err := r.db.First(&event, id).Error
	if err != nil {
		if errors.IsRecordNotFound(err) {
			return nil, errors.ErrEventNotFound
		}
		return nil, err
	}
	return &event, nil
}

func (r *MemorialEventRepository) ListByDeceased(deceasedID uint, page, size int) ([]eventmodel.MemorialEvent, int64, error) {
	query := r.db.Where("deceased_id = ?", deceasedID)
	return paginator.Paginate[eventmodel.MemorialEvent](query, page, size)
}
