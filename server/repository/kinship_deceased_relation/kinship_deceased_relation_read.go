package repository

import (
	"server/errors"
	kdr "server/model/kinship_deceased_relation"
	"server/utils/paginator"
)

func (r *DeceasedRelationRepository) ListByGrave(graveID uint, page, size int) ([]kdr.DeceasedRelation, int64, error) {
	query := r.db.Where("grave_id = ?", graveID)
	return paginator.Paginate[kdr.DeceasedRelation](query, page, size)
}

func (r *DeceasedRelationRepository) ListByDeceased(deceasedID uint, page, size int) ([]kdr.DeceasedRelation, int64, error) {
	query := r.db.Where("deceased_id1 = ? OR deceased_id2 = ?", deceasedID, deceasedID)
	return paginator.Paginate[kdr.DeceasedRelation](query, page, size)
}

func (r *DeceasedRelationRepository) GetByID(id uint) (*kdr.DeceasedRelation, error) {
	var relation kdr.DeceasedRelation
	err := r.db.First(&relation, id).Error
	if err != nil {
		if errors.IsRecordNotFound(err) {
			return nil, errors.ErrRelationNotFound
		}
		return nil, err
	}
	return &relation, nil
}
