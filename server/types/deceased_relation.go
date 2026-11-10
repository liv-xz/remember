package types

type DeceasedRelationResp struct {
	ID           uint   `json:"id"`
	GraveID      uint   `json:"graveID"`
	DeceasedID1  uint   `json:"deceasedID1"`
	DeceasedID2  uint   `json:"deceasedID2"`
	RelationType int8   `json:"relationType"`
	Direction    int8   `json:"direction"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type DeceasedRelationListResp struct {
	Items []DeceasedRelationResp `json:"items"`
}
