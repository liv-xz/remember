package types

type MemorialEventResp struct {
	ID         uint   `json:"id"`
	GraveID    uint   `json:"graveID"`
	DeceasedID uint   `json:"deceasedID"`
	EventType  int8   `json:"eventType"`
	EventDate  string `json:"eventDate"`
	Title      string `json:"title"`
	Lunar      int8   `json:"lunar"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

type MemorialEventListResp struct {
	Items []MemorialEventResp `json:"items"`
}
