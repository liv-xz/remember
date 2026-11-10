package deceased

type CreateDeceasedReq struct {
	Name          string           `json:"name" binding:"required,max=50"`
	Gender        int8             `json:"gender" binding:"required,oneof=0 1 2"`
	BirthDate     string           `json:"birthDate" binding:"omitempty,max=10"`
	DeathDate     string           `json:"deathDate" binding:"omitempty,max=10"`
	Biography     string           `json:"biography" binding:"omitempty,max=65535"`
	PortraitURL   string           `json:"portraitUrl" binding:"omitempty,max=500"`
	Epitaph       string           `json:"epitaph" binding:"omitempty,max=500"`
	BurialOrder   int              `json:"burialOrder" binding:"omitempty"`
	SortOrder     int              `json:"sortOrder" binding:"omitempty"`
	OriginGraveID *uint            `json:"originGraveID"`
	UserRelation  *UserRelationReq `json:"userRelation"`
}

type UserRelationReq struct {
	RelationType int8   `json:"relationType" binding:"required,oneof=1 2 3 4 5 6 7 8 9 10 11 99"`
	Note         string `json:"note" binding:"max=255"`
}

type UpdateDeceasedReq struct {
	Name          *string `json:"name" binding:"omitempty,max=50"`
	Gender        *int8   `json:"gender" binding:"omitempty,oneof=0 1 2"`
	BirthDate     *string `json:"birthDate" binding:"omitempty,max=10"`
	DeathDate     *string `json:"deathDate" binding:"omitempty,max=10"`
	Biography     *string `json:"biography" binding:"omitempty,max=65535"`
	PortraitURL   *string `json:"portraitUrl" binding:"omitempty,max=500"`
	Epitaph       *string `json:"epitaph" binding:"omitempty,max=500"`
	BurialOrder   *int    `json:"burialOrder" binding:"omitempty"`
	SortOrder     *int    `json:"sortOrder" binding:"omitempty"`
	OriginGraveID *uint   `json:"originGraveID"`
}
