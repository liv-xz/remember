package repository

import (
	oplogmodel "server/model/operation_log"

	"gorm.io/gorm"
)

type OperationLogRepository struct {
	db *gorm.DB
}

func NewOperationLogRepository(db *gorm.DB) *OperationLogRepository {
	return &OperationLogRepository{db: db}
}

func (r *OperationLogRepository) Create(log *oplogmodel.OperationLog) error {
	return r.db.Create(log).Error
}
