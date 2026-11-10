package repository

import (
	"server/errors"
	gslmodel "server/model/grave_share_link"
	"server/utils/paginator"
)

func (r *GraveShareLinkRepository) GetByShareCode(code string) (*gslmodel.GraveShareLink, error) {
	var link gslmodel.GraveShareLink
	err := r.db.Where("share_code = ?", code).First(&link).Error
	if err != nil {
		if errors.IsRecordNotFound(err) {
			return nil, errors.ErrShareLinkNotFound
		}
		return nil, err
	}
	return &link, nil
}

func (r *GraveShareLinkRepository) GetByID(id uint) (*gslmodel.GraveShareLink, error) {
	var link gslmodel.GraveShareLink
	err := r.db.First(&link, id).Error
	if err != nil {
		if errors.IsRecordNotFound(err) {
			return nil, errors.ErrShareLinkNotFound
		}
		return nil, err
	}
	return &link, nil
}

func (r *GraveShareLinkRepository) ListByGrave(graveID uint, page, size int) ([]gslmodel.GraveShareLink, int64, error) {
	query := r.db.Where("grave_id = ?", graveID)
	return paginator.Paginate[gslmodel.GraveShareLink](query, page, size)
}
