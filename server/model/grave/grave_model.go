package model

import (
	"server/model"
	usermodel "server/model/user"

	"gorm.io/gorm"
)

type Grave struct {
	model.Model
	Name        string         `gorm:"size:100" json:"name"`
	GraveType   int8           `gorm:"size:1" json:"graveType"`
	LayoutStyle int8           `gorm:"size:1" json:"layoutStyle"`
	Description string         `gorm:"type:text" json:"description"`
	CoverImage  string         `gorm:"size:500" json:"coverImage"`
	Status      int8           `gorm:"size:1" json:"status"`
	Latitude    *float64       `gorm:"type:decimal(10,7)" json:"latitude"`
	Longitude   *float64       `gorm:"type:decimal(10,7)" json:"longitude"`
	Address     string         `gorm:"size:255" json:"address"`
	GraveMode   int8           `gorm:"size:1;not null;default:1" json:"graveMode"`
	CreatedBy   uint           `json:"createdBy"`
	User        usermodel.User `gorm:"foreignKey:CreatedBy" json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
