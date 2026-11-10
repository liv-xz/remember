package model

import "server/model"

type LoginLog struct {
	model.AuditModel
	UserID    uint   `gorm:"index;not null" json:"userID"`
	IP        string `gorm:"size:45" json:"ip"`
	UserAgent string `gorm:"size:512" json:"userAgent"`
	Status    string `gorm:"size:10;not null" json:"status"`
	Reason    string `gorm:"size:255" json:"reason"`
}
