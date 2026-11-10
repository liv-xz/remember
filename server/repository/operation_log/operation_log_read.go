package repository

import (
	oplogmodel "server/model/operation_log"
)

func (r *OperationLogRepository) ListByUserID(userID uint, page, size int) ([]oplogmodel.OperationLog, int64, error) {
	var logs []oplogmodel.OperationLog
	var total int64

	query := r.db.Model(&oplogmodel.OperationLog{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
