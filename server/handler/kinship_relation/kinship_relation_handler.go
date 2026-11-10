package kinship_relation

import (
	"strconv"

	"server/errcode"
	"server/errors"
	krservice "server/service/kinship_relation"
	"server/utils/paginator"
	"server/utils/response"

	"github.com/gin-gonic/gin"
)

type KinshipRelationHandler struct {
	svc *krservice.KinshipRelationService
}

type CreateKinshipReq struct {
	DeceasedID   uint   `json:"deceasedID" binding:"required"`
	RelationType int8   `json:"relationType" binding:"required,oneof=1 2 3 4 5 6 7 8 9 10 11 99"`
	Note         string `json:"note" binding:"max=255"`
}

func NewKinshipRelationHandler(svc *krservice.KinshipRelationService) *KinshipRelationHandler {
	return &KinshipRelationHandler{svc: svc}
}

func (h *KinshipRelationHandler) CreatePublic(c *gin.Context) {
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

	var req CreateKinshipReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	resp, err := h.svc.CreatePublic(uint(gid), userID, req.DeceasedID, req.RelationType, req.Note)
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

func (h *KinshipRelationHandler) Confirm(c *gin.Context) {
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

func (h *KinshipRelationHandler) Reject(c *gin.Context) {
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

func (h *KinshipRelationHandler) ListByGrave(c *gin.Context) {
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
