package repository

import (
	"server/errors"
	usermodel "server/model/user"
)

func (r *UserRepository) GetByID(id uint) (*usermodel.User, error) {
	var user usermodel.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.IsRecordNotFound(err) {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*usermodel.User, error) {
	var user usermodel.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.IsRecordNotFound(err) {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*usermodel.User, error) {
	var user usermodel.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.IsRecordNotFound(err) {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByOpenID(openID string) (*usermodel.User, error) {
	var user usermodel.User
	err := r.db.Where("open_id = ?", openID).First(&user).Error
	if err != nil {
		if errors.IsRecordNotFound(err) {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetConfigByUserID(userID uint) (*usermodel.UserConfig, error) {
	var config usermodel.UserConfig
	err := r.db.Where("user_id = ?", userID).First(&config).Error
	if err != nil {
		if errors.IsRecordNotFound(err) {
			return nil, errors.ErrConfigNotFound
		}
		return nil, err
	}
	return &config, nil
}
