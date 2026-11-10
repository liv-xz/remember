package memorial_record

import (
	"time"

	recordmodel "server/model/memorial_record"
	"server/types"
)

func toRecordResp(m *recordmodel.MemorialRecord) *types.MemorialRecordResp {
	return &types.MemorialRecordResp{
		ID:         m.ID,
		GraveID:    m.GraveID,
		DeceasedID: m.DeceasedID,
		UserID:     m.UserID,
		Type:       m.Type,
		Message:    m.Message,
		CreatedAt:  m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  m.UpdatedAt.Format(time.RFC3339),
	}
}

func toRecordListResp(list []recordmodel.MemorialRecord) *types.MemorialRecordListResp {
	items := make([]types.MemorialRecordResp, 0, len(list))
	for _, m := range list {
		items = append(items, *toRecordResp(&m))
	}
	return &types.MemorialRecordListResp{Items: items}
}
