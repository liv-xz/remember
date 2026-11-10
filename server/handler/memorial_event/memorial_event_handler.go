package memorial_event

import (
	"strconv"

	"server/errcode"
	"server/errors"
	eventservice "server/service/memorial_event"
	"server/utils/paginator"
	"server/utils/response"

	"github.com/gin-gonic/gin"
)

type MemorialEventHandler struct {
	svc *eventservice.MemorialEventService
}

func NewMemorialEventHandler(svc *eventservice.MemorialEventService) *MemorialEventHandler {
	return &MemorialEventHandler{svc: svc}
}

func (h *MemorialEventHandler) Create(c *gin.Context) {
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

	var req CreateEventReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	resp, err := h.svc.Create(uint(gid), userID, req.EventType, req.DeceasedID, req.EventDate, req.Title, req.Lunar)
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
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}
	response.SuccessData(c, resp)
}

func (h *MemorialEventHandler) ListByGrave(c *gin.Context) {
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

func (h *MemorialEventHandler) ListByDeceased(c *gin.Context) {
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

func (h *MemorialEventHandler) Update(c *gin.Context) {
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

	var req UpdateEventReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	params := &eventservice.UpdateEventParams{
		DeceasedID: req.DeceasedID,
		EventType:  req.EventType,
		EventDate:  req.EventDate,
		Title:      req.Title,
		Lunar:      req.Lunar,
	}

	resp, err := h.svc.Update(uint(id), userID, params)
	if err != nil {
		switch err {
		case errors.ErrEventNotFound:
			response.Fail(c, errcode.EventNotFound)
		case errors.ErrGraveNotFound:
			response.Fail(c, errcode.GraveNotFound)
		case errors.ErrGraveAccessDenied:
			response.Fail(c, errcode.GraveAccessDenied)
		case errors.ErrDeceasedNotFound:
			response.Fail(c, errcode.DeceasedNotFound)
		case errors.ErrDeceasedNotInGrave:
			response.Fail(c, errcode.DeceasedNotInGrave)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}
	response.SuccessData(c, resp)
}

func (h *MemorialEventHandler) Delete(c *gin.Context) {
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
		case errors.ErrEventNotFound:
			response.Fail(c, errcode.EventNotFound)
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
