package types

type GraveResp struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	GraveType   int8     `json:"graveType"`
	GraveMode   int8     `json:"graveMode"`
	LayoutStyle int8     `json:"layoutStyle"`
	Description string   `json:"description"`
	CoverImage  string   `json:"coverImage"`
	Status      int8     `json:"status"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	Address     string   `json:"address"`
	CreatedBy   uint     `json:"createdBy"`
	CreatorName string   `json:"creatorName"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

type GraveListResp struct {
	Items []GraveResp `json:"items"`
}

type GraveWithDeceasedResp struct {
	Grave    GraveResp      `json:"grave"`
	Deceased []DeceasedResp `json:"deceased"`
}
