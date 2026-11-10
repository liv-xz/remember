package memorial_event

import (
	"time"

	eventmodel "server/model/memorial_event"
	"server/types"
)

func toEventResp(m *eventmodel.MemorialEvent) *types.MemorialEventResp {
	return &types.MemorialEventResp{
		ID:         m.ID,
		GraveID:    m.GraveID,
		DeceasedID: m.DeceasedID,
		EventType:  m.EventType,
		EventDate:  m.EventDate,
		Title:      m.Title,
		Lunar:      m.Lunar,
		CreatedAt:  m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  m.UpdatedAt.Format(time.RFC3339),
	}
}

func toEventListResp(list []eventmodel.MemorialEvent) *types.MemorialEventListResp {
	items := make([]types.MemorialEventResp, 0, len(list))
	for _, m := range list {
		items = append(items, *toEventResp(&m))
	}
	return &types.MemorialEventListResp{Items: items}
}
