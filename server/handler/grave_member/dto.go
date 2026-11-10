package grave_member

import "server/types"

type AddMemberReq struct {
	UserID       uint   `json:"userID" binding:"required"`
	Role         int8   `json:"role" binding:"required,oneof=10 20 30"`
	RelationDesc string `json:"relationDesc" binding:"max=50"`
}

type UpdateMemberRoleReq struct {
	Role int8 `json:"role" binding:"required,oneof=20 30"`
}

type UpdateRelationDescReq struct {
	RelationDesc string `json:"relationDesc" binding:"required,max=50"`
}

type TransferOwnerReq struct {
	TargetMemberID uint `json:"targetMemberID" binding:"required"`
}

type TransferOwnerResp struct {
	OldOwner *types.GraveMemberResp `json:"oldOwner"`
	NewOwner *types.GraveMemberResp `json:"newOwner"`
}
