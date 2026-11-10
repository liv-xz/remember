package repository

import (
	loginlogmodel "server/model/login_log"
)

func (r *LoginLogRepository) ListByUserID(userID uint, page, size int) ([]loginlogmodel.LoginLog, int64, error) {
	var logs []loginlogmodel.LoginLog
	var total int64

	query := r.db.Model(&loginlogmodel.LoginLog{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
