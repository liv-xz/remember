package model

import (
	"server/model"
	deceasedmodel "server/model/deceased"
	gravemodel "server/model/grave"
	usermodel "server/model/user"
)

type MemorialRecord struct {
	model.Model
	GraveID    uint                   `gorm:"index" json:"graveID"`
	Grave      gravemodel.Grave       `gorm:"foreignKey:GraveID" json:"-"`
	DeceasedID *uint                  `gorm:"index" json:"deceasedID"`
	Deceased   deceasedmodel.Deceased `gorm:"foreignKey:DeceasedID" json:"-"`
	UserID     uint                   `gorm:"index" json:"userID"`
	User       usermodel.User         `gorm:"foreignKey:UserID" json:"-"`
	Type       int8                   `gorm:"size:1" json:"type"`
	Message    string                 `gorm:"type:text" json:"message"`
}
