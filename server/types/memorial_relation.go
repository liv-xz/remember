package types

type MemorialRelationResp struct {
	ID                 uint   `json:"id"`
	GraveID            uint   `json:"graveID"`
	UserID             uint   `json:"userID"`
	UserName           string `json:"userName"`
	UserAvatar         string `json:"userAvatar"`
	DeceasedID         uint   `json:"deceasedID"`
	DeceasedName       string `json:"deceasedName"`
	DeceasedAvatar     string `json:"deceasedAvatar"`
	RelationType       int8   `json:"relationType"`
	ConfirmationStatus int8   `json:"confirmationStatus"`
	ConfirmedBy        *uint  `json:"confirmedBy"`
	Note               string `json:"note"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

type MemorialRelationListResp struct {
	Items []MemorialRelationResp `json:"items"`
}
