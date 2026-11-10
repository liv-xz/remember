package repository

import (
	"server/errors"
	deceasedmodel "server/model/deceased"
	"server/utils/paginator"
)

func (r *DeceasedRepository) GetByID(id uint) (*deceasedmodel.Deceased, error) {
	var deceased deceasedmodel.Deceased
	err := r.db.First(&deceased, id).Error
	if err != nil {
		if errors.IsRecordNotFound(err) {
			return nil, errors.ErrDeceasedNotFound
		}
		return nil, err
	}
	return &deceased, nil
}

func (r *DeceasedRepository) ListByGrave(graveID uint, page, size int) ([]deceasedmodel.Deceased, int64, error) {
	query := r.db.Where("grave_id = ?", graveID).Order("burial_order")
	return paginator.Paginate[deceasedmodel.Deceased](query, page, size)
}

func (r *DeceasedRepository) CountByGrave(graveID uint) (int64, error) {
	var count int64
	err := r.db.Model(&deceasedmodel.Deceased{}).Where("grave_id = ?", graveID).Count(&count).Error
	return count, err
}
