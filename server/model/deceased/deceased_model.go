package model

import (
	"server/model"
	gravemodel "server/model/grave"
)

type Deceased struct {
	model.Model
	GraveID       uint             `gorm:"index" json:"graveID"`
	Grave         gravemodel.Grave `gorm:"foreignKey:GraveID" json:"-"`
	Name          string           `gorm:"size:50" json:"name"`
	Gender        int8             `gorm:"size:1" json:"gender"`
	BirthDate     string           `gorm:"size:10" json:"birthDate"`
	DeathDate     string           `gorm:"size:10" json:"deathDate"`
	Biography     string           `gorm:"type:text" json:"biography"`
	PortraitURL   string           `gorm:"size:500" json:"portraitUrl"`
	Epitaph       string           `gorm:"size:500" json:"epitaph"`
	BurialOrder   int              `json:"burialOrder"`
	SortOrder     int              `json:"sortOrder"`
	OriginGraveID *uint            `gorm:"index" json:"originGraveID"`
}
