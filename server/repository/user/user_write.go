package repository

import (
	"server/errors"
	usermodel "server/model/user"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *usermodel.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *usermodel.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) Delete(id uint) error {
	result := r.db.Delete(&usermodel.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) CreateConfig(config *usermodel.UserConfig) error {
	return r.db.Create(config).Error
}

func (r *UserRepository) UpdateConfig(config *usermodel.UserConfig) error {
	result := r.db.Save(config)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrConfigNotFound
	}
	return nil
}

func (r *UserRepository) UpsertConfig(config *usermodel.UserConfig) error {
	var count int64
	r.db.Model(&usermodel.UserConfig{}).Where("user_id = ?", config.UserID).Count(&count)
	if count == 0 {
		return r.db.Create(config).Error
	}
	result := r.db.Where("user_id = ?", config.UserID).Updates(config)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrConfigNotFound
	}
	return nil
}
