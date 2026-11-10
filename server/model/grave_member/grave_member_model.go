package model

import (
	"server/model"
	gravemodel "server/model/grave"
	usermodel "server/model/user"
)

type GraveMember struct {
	model.Model
	GraveID      uint             `gorm:"index" json:"graveID"`
	Grave        gravemodel.Grave `gorm:"foreignKey:GraveID" json:"-"`
	UserID       uint             `gorm:"index" json:"userID"`
	User         usermodel.User   `gorm:"foreignKey:UserID" json:"-"`
	Role         int8             `gorm:"size:1" json:"role"`
	RelationDesc string           `gorm:"size:50" json:"relationDesc"`
}
