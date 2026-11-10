package memorial_relation

import (
	"strconv"

	"server/errcode"
	"server/errors"
	mrservice "server/service/memorial_relation"
	"server/utils/paginator"
	"server/utils/response"

	"github.com/gin-gonic/gin"
)

type MemorialRelationHandler struct {
	svc *mrservice.MemorialRelationService
}

func NewMemorialRelationHandler(svc *mrservice.MemorialRelationService) *MemorialRelationHandler {
	return &MemorialRelationHandler{svc: svc}
}

func (h *MemorialRelationHandler) Create(c *gin.Context) {
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

	var req CreateMemorialReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	resp, err := h.svc.Create(uint(gid), userID, req.DeceasedID, req.RelationType, req.Note)
	if err != nil {
		switch err {
		case errors.ErrDeceasedNotFound:
			response.Fail(c, errcode.DeceasedNotFound)
		case errors.ErrDeceasedNotInGrave:
			response.Fail(c, errcode.DeceasedNotInGrave)
		case errors.ErrRelationAlreadyExist:
			response.Fail(c, errcode.RelationAlreadyExist)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}

	response.SuccessData(c, resp)
}

func (h *MemorialRelationHandler) CreatePublic(c *gin.Context) {
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

	var req CreateMemorialReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	resp, err := h.svc.Create(uint(gid), userID, req.DeceasedID, req.RelationType, req.Note)
	if err != nil {
		switch err {
		case errors.ErrDeceasedNotFound:
			response.Fail(c, errcode.DeceasedNotFound)
		case errors.ErrDeceasedNotInGrave:
			response.Fail(c, errcode.DeceasedNotInGrave)
		case errors.ErrRelationAlreadyExist:
			response.Fail(c, errcode.RelationAlreadyExist)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}

	response.SuccessData(c, resp)
}

func (h *MemorialRelationHandler) Confirm(c *gin.Context) {
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

	err = h.svc.Confirm(uint(id), userID)
	if err != nil {
		switch err {
		case errors.ErrRelationNotFound:
			response.Fail(c, errcode.RelationNotFound)
		case errors.ErrGraveAccessDenied:
			response.Fail(c, errcode.GraveAccessDenied)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}

	response.Success(c)
}

func (h *MemorialRelationHandler) Reject(c *gin.Context) {
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

	err = h.svc.Reject(uint(id), userID)
	if err != nil {
		switch err {
		case errors.ErrRelationNotFound:
			response.Fail(c, errcode.RelationNotFound)
		case errors.ErrGraveAccessDenied:
			response.Fail(c, errcode.GraveAccessDenied)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}

	response.Success(c)
}

func (h *MemorialRelationHandler) ListByGrave(c *gin.Context) {
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
		case errors.ErrGraveNotFound:
			response.Fail(c, errcode.GraveNotFound)
		case errors.ErrGraveAccessDenied:
			response.Fail(c, errcode.GraveAccessDenied)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}

	response.SuccessPage(c, resp.Items, total, page, size)
}
