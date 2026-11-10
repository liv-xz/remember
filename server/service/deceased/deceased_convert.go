package deceased

import (
	"time"

	deceasedmodel "server/model/deceased"
	"server/types"
)

func toDeceasedResp(m *deceasedmodel.Deceased) *types.DeceasedResp {
	return &types.DeceasedResp{
		ID:            m.ID,
		GraveID:       m.GraveID,
		Name:          m.Name,
		Gender:        m.Gender,
		BirthDate:     m.BirthDate,
		DeathDate:     m.DeathDate,
		Biography:     m.Biography,
		PortraitURL:   m.PortraitURL,
		Epitaph:       m.Epitaph,
		BurialOrder:   m.BurialOrder,
		SortOrder:     m.SortOrder,
		OriginGraveID: m.OriginGraveID,
		CreatedAt:     m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     m.UpdatedAt.Format(time.RFC3339),
	}
}

func toDeceasedListResp(list []deceasedmodel.Deceased) *types.DeceasedListResp {
	items := make([]types.DeceasedResp, 0, len(list))
	for _, d := range list {
		items = append(items, *toDeceasedResp(&d))
	}
	return &types.DeceasedListResp{Items: items}
}
