package kinship_deceased_relation

import (
	"strconv"

	"server/errcode"
	"server/errors"
	kdrservice "server/service/kinship_deceased_relation"
	"server/utils/paginator"
	"server/utils/response"

	"github.com/gin-gonic/gin"
)

type DeceasedRelationHandler struct {
	svc *kdrservice.DeceasedRelationService
}

func NewDeceasedRelationHandler(svc *kdrservice.DeceasedRelationService) *DeceasedRelationHandler {
	return &DeceasedRelationHandler{svc: svc}
}

func (h *DeceasedRelationHandler) Create(c *gin.Context) {
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

	var req CreateRelationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	resp, err := h.svc.Create(uint(gid), userID, req.DeceasedID1, req.DeceasedID2, req.RelationType, req.Direction)
	if err != nil {
		switch err {
		case errors.ErrGraveNotFound:
			response.Fail(c, errcode.GraveNotFound)
		case errors.ErrGraveAccessDenied:
			response.Fail(c, errcode.GraveAccessDenied)
		case errors.ErrDeceasedNotFound:
			response.Fail(c, errcode.DeceasedNotFound)
		case errors.ErrDeceasedNotInGrave:
			response.Fail(c, errcode.DeceasedNotInGrave)
		case errors.ErrRelationSelfInvalid:
			response.Fail(c, errcode.RelationSelfInvalid)
		case errors.ErrRelationAlreadyExist:
			response.Fail(c, errcode.RelationAlreadyExist)
		case errors.ErrRelationDirectionInvalid:
			response.Fail(c, errcode.RelationDirectionInvalid)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}
	response.SuccessData(c, resp)
}

func (h *DeceasedRelationHandler) ListByGrave(c *gin.Context) {
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

func (h *DeceasedRelationHandler) ListByDeceased(c *gin.Context) {
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	did, err := strconv.ParseUint(c.Param("did"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	page, size := paginator.BindPage(c)
	resp, total, err := h.svc.ListByDeceased(uint(did), userID, page, size)
	if err != nil {
		switch err {
		case errors.ErrDeceasedNotFound:
			response.Fail(c, errcode.DeceasedNotFound)
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

func (h *DeceasedRelationHandler) Delete(c *gin.Context) {
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

	err = h.svc.Delete(uint(id), userID)
	if err != nil {
		switch err {
		case errors.ErrRelationNotFound:
			response.Fail(c, errcode.RelationNotFound)
		case errors.ErrGraveNotFound:
			response.Fail(c, errcode.GraveNotFound)
		case errors.ErrGraveAccessDenied:
			response.Fail(c, errcode.GraveAccessDenied)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}
	response.Success(c)
}
