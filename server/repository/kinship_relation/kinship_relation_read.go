package repository

import (
	"server/errors"
	krmodel "server/model/kinship_relation"
	"server/utils/paginator"

	"gorm.io/gorm"
)

func (r *KinshipRelationRepository) GetByID(id uint) (*krmodel.KinshipRelation, error) {
	var relation krmodel.KinshipRelation
	err := r.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username", "nickname", "avatar")
	}).Preload("Deceased", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "portrait_url")
	}).First(&relation, id).Error
	if err != nil {
		if errors.IsRecordNotFound(err) {
			return nil, errors.ErrRelationNotFound
		}
		return nil, err
	}
	return &relation, nil
}

func (r *KinshipRelationRepository) ListByGrave(graveID uint, page, size int) ([]krmodel.KinshipRelation, int64, error) {
	query := r.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username", "nickname", "avatar")
	}).Preload("Deceased", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "portrait_url")
	}).Where("grave_id = ?", graveID)
	return paginator.Paginate[krmodel.KinshipRelation](query, page, size)
}

func (r *KinshipRelationRepository) ListByDeceased(deceasedID uint, page, size int) ([]krmodel.KinshipRelation, int64, error) {
	query := r.db.Preload("User").Where("deceased_id = ?", deceasedID)
	return paginator.Paginate[krmodel.KinshipRelation](query, page, size)
}

func (r *KinshipRelationRepository) ListByUser(userID uint) ([]krmodel.KinshipRelation, error) {
	var list []krmodel.KinshipRelation
	err := r.db.Preload("Deceased").
		Where("user_id = ?", userID).
		Find(&list).Error
	return list, err
}

func (r *KinshipRelationRepository) CheckExists(userID, deceasedID uint) (bool, error) {
	var count int64
	err := r.db.Model(&krmodel.KinshipRelation{}).
		Where("user_id = ? AND deceased_id = ?", userID, deceasedID).
		Count(&count).Error
	return count > 0, err
}
