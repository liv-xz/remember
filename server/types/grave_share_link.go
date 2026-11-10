package types

import "time"

type GraveShareLinkResp struct {
	ID          uint       `json:"id"`
	GraveID     uint       `json:"graveID"`
	GraveName   string     `json:"graveName"`
	ShareCode   string     `json:"shareCode"`
	ShareURL    string     `json:"shareURL"`
	ExpireAt    *time.Time `json:"expireAt"`
	MaxUseCount int        `json:"maxUseCount"`
	UsedCount   int        `json:"usedCount"`
	IsActive    bool       `json:"isActive"`
	CreatedAt   string     `json:"createdAt"`
}

type GraveShareLinkListResp struct {
	Items []GraveShareLinkResp `json:"items"`
}

type GraveShareLinkValidResp struct {
	GraveID   uint            `json:"graveID"`
	GraveName string          `json:"graveName"`
	GraveDesc string          `json:"graveDesc"`
	CreatedBy uint            `json:"createdBy"`
	Deceased  []DeceasedBrief `json:"deceased"`
}

type JoinShareLinkReq struct {
	DeceasedID   uint   `json:"deceasedID"`
	RelationType int8   `json:"relationType"`
	Note         string `json:"note"`
}

type DeceasedBrief struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Gender int8   `json:"gender"`
}
