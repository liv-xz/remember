package repository

import (
	"server/errors"
	membermodel "server/model/grave_member"
	"server/utils/paginator"
)

func (r *GraveMemberRepository) ListByGrave(graveID uint, page, size int) ([]membermodel.GraveMember, int64, error) {
	query := r.db.Where("grave_id = ?", graveID)
	return paginator.Paginate[membermodel.GraveMember](query, page, size)
}

func (r *GraveMemberRepository) ListByUser(userID uint) ([]membermodel.GraveMember, error) {
	var list []membermodel.GraveMember
	err := r.db.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}

func (r *GraveMemberRepository) GetByUserAndGrave(userID, graveID uint) (*membermodel.GraveMember, error) {
	var member membermodel.GraveMember
	err := r.db.Where("user_id = ? AND grave_id = ?", userID, graveID).First(&member).Error
	if err != nil {
		if errors.IsRecordNotFound(err) {
			return nil, errors.ErrMemberNotFound
		}
		return nil, err
	}
	return &member, nil
}

func (r *GraveMemberRepository) GetByID(id uint) (*membermodel.GraveMember, error) {
	var member membermodel.GraveMember
	err := r.db.First(&member, id).Error
	if err != nil {
		if errors.IsRecordNotFound(err) {
			return nil, errors.ErrMemberNotFound
		}
		return nil, err
	}
	return &member, nil
}

func (r *GraveMemberRepository) IsMember(graveID, userID uint) bool {
	var count int64
	r.db.Model(&membermodel.GraveMember{}).
		Where("grave_id = ? AND user_id = ?", graveID, userID).Count(&count)
	return count > 0
}

func (r *GraveMemberRepository) HasRole(graveID, userID uint, minRole int8) bool {
	var member membermodel.GraveMember
	err := r.db.Where("grave_id = ? AND user_id = ?", graveID, userID).First(&member).Error
	if err != nil {
		return false
	}
	return member.Role >= minRole
}
