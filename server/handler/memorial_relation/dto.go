package memorial_relation

type CreateMemorialReq struct {
	DeceasedID   uint   `json:"deceasedID" binding:"required"`
	RelationType int8   `json:"relationType" binding:"required,oneof=1 2 3 4 5 6 7 8 9 10 11 99"`
	Note         string `json:"note" binding:"max=255"`
}
