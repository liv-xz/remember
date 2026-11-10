package model

import "server/model"

type OperationLog struct {
	model.AuditModel
	UserID         uint   `gorm:"index;not null" json:"userID"`
	Action         string `gorm:"size:20;not null" json:"action"`
	TargetType     string `gorm:"size:50;not null" json:"targetType"`
	TargetID       uint   `json:"targetID"`
	ResponseStatus int    `gorm:"size:3" json:"responseStatus"`
	Detail         string `gorm:"type:json" json:"detail"`
	IP             string `gorm:"size:45" json:"ip"`
}
