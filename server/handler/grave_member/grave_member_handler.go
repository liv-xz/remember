package grave_member

import (
	"strconv"

	"server/errcode"
	"server/errors"
	memberservice "server/service/grave_member"
	"server/utils/paginator"
	"server/utils/response"

	"github.com/gin-gonic/gin"
)

type GraveMemberHandler struct {
	svc *memberservice.GraveMemberService
}

func NewGraveMemberHandler(svc *memberservice.GraveMemberService) *GraveMemberHandler {
	return &GraveMemberHandler{svc: svc}
}

func (h *GraveMemberHandler) Add(c *gin.Context) {
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	gid, err := strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	var req AddMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	resp, err := h.svc.Add(uint(gid), userID, req.UserID, req.Role, req.RelationDesc)
	if err != nil {
		switch err {
		case errors.ErrMemberRoleInsufficient:
			response.Fail(c, errcode.MemberRoleInsufficient)
		case errors.ErrMemberAlreadyExist:
			response.Fail(c, errcode.MemberAlreadyExist)
		case errors.ErrMemberNotFound:
			response.Fail(c, errcode.MemberNotFound)
		case errors.ErrUserNotFound:
			response.Fail(c, errcode.UserNotFound)
		default:
			response.Fail(c, errcode.MemberCreateFail)
		}
		return
	}

	response.SuccessData(c, resp)
}

func (h *GraveMemberHandler) ListByGrave(c *gin.Context) {
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	gid, err := strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	page, size := paginator.BindPage(c)
	resp, total, err := h.svc.ListByGrave(uint(gid), userID, page, size)
	if err != nil {
		switch err {
		case errors.ErrMemberNotFound:
			response.Fail(c, errcode.MemberNotFound)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}

	response.SuccessPage(c, resp.Items, total, page, size)
}

func (h *GraveMemberHandler) UpdateRole(c *gin.Context) {
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	var req UpdateMemberRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	resp, err := h.svc.UpdateRole(userID, uint(id), req.Role)
	if err != nil {
		switch err {
		case errors.ErrMemberNotFound:
			response.Fail(c, errcode.MemberNotFound)
		case errors.ErrMemberRoleInsufficient:
			response.Fail(c, errcode.MemberRoleInsufficient)
		case errors.ErrMemberCannotDemoteOwner:
			response.Fail(c, errcode.MemberRoleInsufficient)
		default:
			response.Fail(c, errcode.MemberUpdateFail)
		}
		return
	}

	response.SuccessData(c, resp)
}

func (h *GraveMemberHandler) UpdateRelationDesc(c *gin.Context) {
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	var req UpdateRelationDescReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	resp, err := h.svc.UpdateRelationDesc(userID, uint(id), req.RelationDesc)
	if err != nil {
		switch err {
		case errors.ErrMemberNotFound:
			response.Fail(c, errcode.MemberNotFound)
		case errors.ErrMemberRoleInsufficient:
			response.Fail(c, errcode.MemberRoleInsufficient)
		default:
			response.Fail(c, errcode.MemberUpdateFail)
		}
		return
	}

	response.SuccessData(c, resp)
}

func (h *GraveMemberHandler) Remove(c *gin.Context) {
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	err = h.svc.Remove(userID, uint(id))
	if err != nil {
		switch err {
		case errors.ErrMemberNotFound:
			response.Fail(c, errcode.MemberNotFound)
		case errors.ErrMemberCannotRemoveOwner:
			response.Fail(c, errcode.MemberOwnerTransfer)
		case errors.ErrMemberRoleInsufficient:
			response.Fail(c, errcode.MemberRoleInsufficient)
		default:
			response.Fail(c, errcode.MemberDeleteFail)
		}
		return
	}

	response.Success(c)
}

func (h *GraveMemberHandler) TransferOwnership(c *gin.Context) {
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	var req TransferOwnerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	oldOwner, newOwner, err := h.svc.TransferOwnership(userID, req.TargetMemberID)
	if err != nil {
		switch err {
		case errors.ErrMemberNotFound:
			response.Fail(c, errcode.MemberNotFound)
		case errors.ErrMemberRoleInsufficient:
			response.Fail(c, errcode.MemberRoleInsufficient)
		case errors.ErrMemberTransferTargetInvalid:
			response.Fail(c, errcode.MemberRoleInsufficient)
		default:
			response.Fail(c, errcode.MemberUpdateFail)
		}
		return
	}

	response.SuccessData(c, TransferOwnerResp{
		OldOwner: oldOwner,
		NewOwner: newOwner,
	})
}
