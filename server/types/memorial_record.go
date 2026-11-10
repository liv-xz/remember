package types

type MemorialRecordResp struct {
	ID         uint   `json:"id"`
	GraveID    uint   `json:"graveID"`
	DeceasedID *uint  `json:"deceasedID"`
	UserID     uint   `json:"userID"`
	Type       int8   `json:"type"`
	Message    string `json:"message"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

type MemorialRecordListResp struct {
	Items []MemorialRecordResp `json:"items"`
}
