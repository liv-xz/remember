package types

type DeceasedResp struct {
	ID            uint   `json:"id"`
	GraveID       uint   `json:"graveId"`
	Name          string `json:"name"`
	Gender        int8   `json:"gender"`
	BirthDate     string `json:"birthDate"`
	DeathDate     string `json:"deathDate"`
	Biography     string `json:"biography"`
	PortraitURL   string `json:"portraitUrl"`
	Epitaph       string `json:"epitaph"`
	BurialOrder   int    `json:"burialOrder"`
	SortOrder     int    `json:"sortOrder"`
	OriginGraveID *uint  `json:"originGraveID"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type DeceasedListResp struct {
	Items []DeceasedResp `json:"items"`
}

type GraveFullResp struct {
	Grave    GraveResp      `json:"grave"`
	Deceased []DeceasedResp `json:"deceased"`
}
