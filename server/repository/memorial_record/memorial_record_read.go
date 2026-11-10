package repository

import (
	recordmodel "server/model/memorial_record"
)

func (r *MemorialRecordRepository) ListByGrave(graveID uint, page, size int) ([]recordmodel.MemorialRecord, int64, error) {
	var list []recordmodel.MemorialRecord
	var total int64
	r.db.Model(&recordmodel.MemorialRecord{}).Where("grave_id = ?", graveID).Count(&total)
	err := r.db.Where("grave_id = ?", graveID).
		Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&list).Error
	return list, total, err
}

func (r *MemorialRecordRepository) ListByDeceased(deceasedID uint, page, size int) ([]recordmodel.MemorialRecord, int64, error) {
	var list []recordmodel.MemorialRecord
	var total int64
	r.db.Model(&recordmodel.MemorialRecord{}).Where("deceased_id = ?", deceasedID).Count(&total)
	err := r.db.Where("deceased_id = ?", deceasedID).
		Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&list).Error
	return list, total, err
}
