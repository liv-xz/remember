package model

import (
	"server/model"
	gravemodel "server/model/grave"
	usermodel "server/model/user"
	"time"
)

type GraveShareLink struct {
	model.Model
	GraveID     uint             `gorm:"index" json:"graveID"`
	Grave       gravemodel.Grave `gorm:"foreignKey:GraveID" json:"-"`
	CreatedBy   uint             `gorm:"index" json:"createdBy"`
	Creator     usermodel.User   `gorm:"foreignKey:CreatedBy" json:"-"`
	ShareCode   string           `gorm:"size:32;uniqueIndex;not null" json:"shareCode"`
	ExpireAt    *time.Time       `gorm:"index" json:"expireAt"`
	MaxUseCount int              `gorm:"default:0" json:"maxUseCount"`
	UsedCount   int              `gorm:"default:0" json:"usedCount"`
	IsActive    bool             `gorm:"default:true;index" json:"isActive"`
}

func (GraveShareLink) TableName() string {
	return "grave_share_links"
}
