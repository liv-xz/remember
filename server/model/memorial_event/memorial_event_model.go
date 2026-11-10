package model

import (
	"server/model"
	deceasedmodel "server/model/deceased"
	gravemodel "server/model/grave"
)

type MemorialEvent struct {
	model.Model
	GraveID    uint                   `gorm:"index" json:"graveID"`
	Grave      gravemodel.Grave       `gorm:"foreignKey:GraveID" json:"-"`
	DeceasedID uint                   `gorm:"index" json:"deceasedID"`
	Deceased   deceasedmodel.Deceased `gorm:"foreignKey:DeceasedID" json:"-"`
	EventType  int8                   `gorm:"size:1" json:"eventType"`
	EventDate  string                 `gorm:"size:10" json:"eventDate"`
	Title      string                 `gorm:"size:100" json:"title"`
	Lunar      int8                   `gorm:"size:1" json:"lunar"`
}
