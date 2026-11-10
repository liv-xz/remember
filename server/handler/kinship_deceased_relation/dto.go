package kinship_deceased_relation

type CreateRelationReq struct {
	DeceasedID1  uint `json:"deceasedID1" binding:"required"`
	DeceasedID2  uint `json:"deceasedID2" binding:"required"`
	RelationType int8 `json:"relationType" binding:"required,oneof=1 2"`
	Direction    int8 `json:"direction" binding:"oneof=0 1 2"`
}
