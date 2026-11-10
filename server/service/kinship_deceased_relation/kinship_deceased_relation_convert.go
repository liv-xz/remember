package kinship_deceased_relation

import (
	"time"

	kdr "server/model/kinship_deceased_relation"
	"server/types"
)

func toRelationResp(m *kdr.DeceasedRelation) *types.DeceasedRelationResp {
	return &types.DeceasedRelationResp{
		ID:           m.ID,
		GraveID:      m.GraveID,
		DeceasedID1:  m.DeceasedID1,
		DeceasedID2:  m.DeceasedID2,
		RelationType: m.RelationType,
		Direction:    m.Direction,
		CreatedAt:    m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    m.UpdatedAt.Format(time.RFC3339),
	}
}

func toRelationListResp(list []kdr.DeceasedRelation) *types.DeceasedRelationListResp {
	items := make([]types.DeceasedRelationResp, 0, len(list))
	for _, m := range list {
		items = append(items, *toRelationResp(&m))
	}
	return &types.DeceasedRelationListResp{Items: items}
}
