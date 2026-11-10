package types

type GraveMemberResp struct {
	ID           uint   `json:"id"`
	GraveID      uint   `json:"graveID"`
	UserID       uint   `json:"userID"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	Role         int8   `json:"role"`
	RelationDesc string `json:"relationDesc"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type GraveMemberListResp struct {
	Items []GraveMemberResp `json:"items"`
}
