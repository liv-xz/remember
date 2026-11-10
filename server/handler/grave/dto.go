package grave

type CreateGraveReq struct {
	Name        string   `json:"name" binding:"required,max=100"`
	GraveType   int8     `json:"graveType" binding:"required,oneof=1 2 3"`
	GraveMode   int8     `json:"graveMode" binding:"required,oneof=1 2"`
	LayoutStyle int8     `json:"layoutStyle" binding:"required,oneof=1 2 3"`
	Description string   `json:"description" binding:"max=65535"`
	CoverImage  string   `json:"coverImage" binding:"max=500"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	Address     string   `json:"address" binding:"max=255"`
}

type DeceasedItem struct {
	Name          string `json:"name" binding:"required,max=50"`
	Gender        int8   `json:"gender" binding:"required,oneof=0 1 2"`
	BirthDate     string `json:"birthDate" binding:"max=10"`
	DeathDate     string `json:"deathDate" binding:"max=10"`
	Biography     string `json:"biography" binding:"max=65535"`
	PortraitURL   string `json:"portraitUrl" binding:"max=500"`
	Epitaph       string `json:"epitaph" binding:"max=500"`
	BurialOrder   int    `json:"burialOrder"`
	OriginGraveID *uint  `json:"originGraveID"`
}

type CreateGraveWithDeceasedReq struct {
	Name        string         `json:"name" binding:"required,max=100"`
	GraveType   int8           `json:"graveType" binding:"required,oneof=1 2 3"`
	GraveMode   int8           `json:"graveMode" binding:"required,oneof=1 2"`
	LayoutStyle int8           `json:"layoutStyle" binding:"required,oneof=1 2 3"`
	Description string         `json:"description" binding:"max=65535"`
	CoverImage  string         `json:"coverImage" binding:"max=500"`
	Latitude    *float64       `json:"latitude"`
	Longitude   *float64       `json:"longitude"`
	Address     string         `json:"address" binding:"max=255"`
	Deceased    []DeceasedItem `json:"deceased" binding:"required,min=1"`
}

type UpdateGraveReq struct {
	Name          string   `json:"name" binding:"omitempty,max=100"`
	GraveType     int8     `json:"graveType" binding:"omitempty,oneof=1 2 3"`
	GraveMode     int8     `json:"graveMode" binding:"omitempty,oneof=1 2"`
	LayoutStyle   int8     `json:"layoutStyle" binding:"omitempty,oneof=1 2 3"`
	Description   string   `json:"description" binding:"omitempty,max=65535"`
	CoverImage    string   `json:"coverImage" binding:"omitempty,max=500"`
	Status        *int8    `json:"status" binding:"omitempty,oneof=0 1"`
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
	Address       string   `json:"address" binding:"omitempty,max=255"`
	ClearLocation *bool    `json:"clearLocation"`
}
