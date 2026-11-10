package repository

import (
	"server/errors"
	membermodel "server/model/grave_member"

	"gorm.io/gorm"
)

type GraveMemberRepository struct {
	db *gorm.DB
}

func NewGraveMemberRepository(db *gorm.DB) *GraveMemberRepository {
	return &GraveMemberRepository{db: db}
}

func (r *GraveMemberRepository) Tx(tx *gorm.DB) *GraveMemberRepository {
	return &GraveMemberRepository{db: tx}
}

func (r *GraveMemberRepository) Create(member *membermodel.GraveMember) error {
	return r.db.Create(member).Error
}

func (r *GraveMemberRepository) UpdateRole(id uint, role int8) error {
	result := r.db.Model(&membermodel.GraveMember{}).Where("id = ?", id).Update("role", role)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrMemberNotFound
	}
	return nil
}

func (r *GraveMemberRepository) Update(member *membermodel.GraveMember) error {
	result := r.db.Save(member)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrMemberNotFound
	}
	return nil
}

func (r *GraveMemberRepository) Delete(id uint) error {
	result := r.db.Delete(&membermodel.GraveMember{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrMemberNotFound
	}
	return nil
}
