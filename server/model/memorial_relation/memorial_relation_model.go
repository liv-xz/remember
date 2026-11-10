package model

import (
	"server/model"
	deceasedmodel "server/model/deceased"
	gravemodel "server/model/grave"
	usermodel "server/model/user"
)

type MemorialRelation struct {
	model.Model
	GraveID            uint                   `gorm:"index" json:"graveID"`
	Grave              gravemodel.Grave       `gorm:"foreignKey:GraveID" json:"-"`
	UserID             uint                   `gorm:"index" json:"userID"`
	User               usermodel.User         `gorm:"foreignKey:UserID" json:"-"`
	DeceasedID         uint                   `gorm:"index" json:"deceasedID"`
	Deceased           deceasedmodel.Deceased `gorm:"foreignKey:DeceasedID" json:"-"`
	RelationType       int8                   `gorm:"size:1;not null" json:"relationType"`
	ConfirmationStatus int8                   `gorm:"size:1;default:0" json:"confirmationStatus"`
	ConfirmedBy        *uint                  `gorm:"index" json:"confirmedBy"`
	ConfirmedUser      *usermodel.User        `gorm:"foreignKey:ConfirmedBy" json:"-"`
	Note               string                 `gorm:"size:255" json:"note"`
}

func (MemorialRelation) TableName() string {
	return "memorial_relations"
}
