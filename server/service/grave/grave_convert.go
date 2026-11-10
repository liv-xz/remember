package grave

import (
	"time"

	gravemodel "server/model/grave"
	"server/types"
)

func toGraveResp(grave *gravemodel.Grave, creatorName string) *types.GraveResp {
	return &types.GraveResp{
		ID:          grave.ID,
		Name:        grave.Name,
		GraveType:   grave.GraveType,
		GraveMode:   grave.GraveMode,
		LayoutStyle: grave.LayoutStyle,
		Description: grave.Description,
		CoverImage:  grave.CoverImage,
		Status:      grave.Status,
		Latitude:    grave.Latitude,
		Longitude:   grave.Longitude,
		Address:     grave.Address,
		CreatedBy:   grave.CreatedBy,
		CreatorName: creatorName,
		CreatedAt:   grave.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   grave.UpdatedAt.Format(time.RFC3339),
	}
}
