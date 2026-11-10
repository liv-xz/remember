package repository

import (
	loginlogmodel "server/model/login_log"

	"gorm.io/gorm"
)

type LoginLogRepository struct {
	db *gorm.DB
}

func NewLoginLogRepository(db *gorm.DB) *LoginLogRepository {
	return &LoginLogRepository{db: db}
}

func (r *LoginLogRepository) Create(log *loginlogmodel.LoginLog) error {
	return r.db.Create(log).Error
}
