package memorial_record

type CreateRecordReq struct {
	Type       int8   `json:"type" binding:"required,oneof=1 2 3 4 5"`
	DeceasedID *uint  `json:"deceasedID"`
	Message    string `json:"message" binding:"max=65535"`
}
