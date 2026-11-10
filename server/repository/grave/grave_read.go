package repository

import (
	"server/errors"
	gravemodel "server/model/grave"
	"server/utils/paginator"
)

func (r *GraveRepository) GetByID(id uint) (*gravemodel.Grave, error) {
	var grave gravemodel.Grave
	err := r.db.First(&grave, id).Error
	if err != nil {
		if errors.IsRecordNotFound(err) {
			return nil, errors.ErrGraveNotFound
		}
		return nil, err
	}
	return &grave, nil
}

func (r *GraveRepository) ListByUser(userID uint) ([]gravemodel.Grave, error) {
	var graves []gravemodel.Grave
	err := r.db.Where("created_by = ?", userID).Find(&graves).Error
	return graves, err
}

func (r *GraveRepository) ListByUserOrMember(userID uint, page, size int) ([]gravemodel.Grave, int64, error) {
	query := r.db.Where("created_by = ? OR id IN (SELECT grave_id FROM grave_members WHERE user_id = ?)", userID, userID)
	return paginator.Paginate[gravemodel.Grave](query, page, size)
}
