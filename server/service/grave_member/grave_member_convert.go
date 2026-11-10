package grave_member

import (
	"time"

	membermodel "server/model/grave_member"
	usermodel "server/model/user"
	"server/types"
)

func toMemberResp(m *membermodel.GraveMember, user *usermodel.User) *types.GraveMemberResp {
	username := ""
	nickname := ""
	avatar := ""
	if user != nil {
		username = user.Username
		nickname = user.Nickname
		avatar = user.Avatar
	}
	if nickname == "" {
		nickname = username
	}
	return &types.GraveMemberResp{
		ID:           m.ID,
		GraveID:      m.GraveID,
		UserID:       m.UserID,
		Username:     username,
		Nickname:     nickname,
		Avatar:       avatar,
		Role:         m.Role,
		RelationDesc: m.RelationDesc,
		CreatedAt:    m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    m.UpdatedAt.Format(time.RFC3339),
	}
}

func toMemberListResp(list []membermodel.GraveMember, userMap map[uint]*usermodel.User) *types.GraveMemberListResp {
	items := make([]types.GraveMemberResp, 0, len(list))
	for _, m := range list {
		user, ok := userMap[m.UserID]
		if !ok {
			continue
		}
		items = append(items, *toMemberResp(&m, user))
	}
	return &types.GraveMemberListResp{Items: items}
}
