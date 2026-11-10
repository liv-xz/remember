package model

import (
	"server/model"
	deceasedmodel "server/model/deceased"
	gravemodel "server/model/grave"
)

type DeceasedRelation struct {
	model.Model
	GraveID      uint                   `gorm:"index" json:"graveID"`
	Grave        gravemodel.Grave       `gorm:"foreignKey:GraveID" json:"-"`
	DeceasedID1  uint                   `gorm:"column:deceased_id1;index;uniqueIndex:idx_kdr_pair" json:"deceasedID1"`
	Deceased1    deceasedmodel.Deceased `gorm:"foreignKey:DeceasedID1;references:ID" json:"-"`
	DeceasedID2  uint                   `gorm:"column:deceased_id2;index;uniqueIndex:idx_kdr_pair" json:"deceasedID2"`
	Deceased2    deceasedmodel.Deceased `gorm:"foreignKey:DeceasedID2;references:ID" json:"-"`
	RelationType int8                   `gorm:"size:1" json:"relationType"`
	Direction    int8                   `gorm:"size:1;default:0" json:"direction"`
}

func (DeceasedRelation) TableName() string {
	return "kinship_deceased_relations"
}
