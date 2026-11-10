package model

import "time"

type AuditModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
}
