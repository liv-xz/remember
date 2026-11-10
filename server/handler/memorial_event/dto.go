package memorial_event

type CreateEventReq struct {
	DeceasedID uint   `json:"deceasedID" binding:"required"`
	EventType  int8   `json:"eventType" binding:"required,oneof=1 2 3 4"`
	EventDate  string `json:"eventDate" binding:"required,max=10"`
	Title      string `json:"title" binding:"required,max=100"`
	Lunar      int8   `json:"lunar" binding:"oneof=0 1"`
}

type UpdateEventReq struct {
	DeceasedID *uint   `json:"deceasedID"`
	EventType  *int8   `json:"eventType" binding:"omitempty,oneof=1 2 3 4"`
	EventDate  *string `json:"eventDate" binding:"omitempty,max=10"`
	Title      *string `json:"title" binding:"omitempty,max=100"`
	Lunar      *int8   `json:"lunar" binding:"omitempty,oneof=0 1"`
}
