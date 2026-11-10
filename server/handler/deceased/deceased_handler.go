package deceased

import (
	"strconv"

	"server/errcode"
	"server/errors"
	deceasedservice "server/service/deceased"
	"server/utils/paginator"
	"server/utils/response"

	"github.com/gin-gonic/gin"
)

type DeceasedHandler struct {
	svc *deceasedservice.DeceasedService
}

func NewDeceasedHandler(svc *deceasedservice.DeceasedService) *DeceasedHandler {
	return &DeceasedHandler{svc: svc}
}

func (h *DeceasedHandler) Create(c *gin.Context) {
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

	var req CreateDeceasedReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	params := deceasedservice.CreateDeceasedParams{
		Name: req.Name, Gender: req.Gender,
		BirthDate: req.BirthDate, DeathDate: req.DeathDate,
		Biography: req.Biography, PortraitURL: req.PortraitURL,
		Epitaph: req.Epitaph, BurialOrder: req.BurialOrder, SortOrder: req.SortOrder,
		OriginGraveID: req.OriginGraveID,
	}
	if req.UserRelation != nil {
		params.UserRelationType = req.UserRelation.RelationType
		params.UserRelationNote = req.UserRelation.Note
	}

	resp, err := h.svc.Create(uint(gid), userID, params)
	if err != nil {
		switch err {
		case errors.ErrGraveNotFound:
			response.Fail(c, errcode.GraveNotFound)
		case errors.ErrGraveAccessDenied:
			response.Fail(c, errcode.GraveAccessDenied)
		case errors.ErrDeceasedOverLimit:
			response.Fail(c, errcode.DeceasedOverLimit)
		case errors.ErrDeceasedDateInvalid:
			response.Fail(c, errcode.DeceasedDateInvalid)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}
	response.SuccessData(c, resp)
}

func (h *DeceasedHandler) ListByGrave(c *gin.Context) {
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

	page, size := parsePage(c)
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

func parsePage(c *gin.Context) (int, int) {
	page := 1
	size := 10
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			page = v
		}
	}
	if s := c.Query("size"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			size = v
		}
	}
	return paginator.ParsePage(page, size)
}

func (h *DeceasedHandler) Update(c *gin.Context) {
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

	var req UpdateDeceasedReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	params := &deceasedservice.UpdateDeceasedParams{
		Name: req.Name, Gender: req.Gender,
		BirthDate: req.BirthDate, DeathDate: req.DeathDate,
		Biography: req.Biography, PortraitURL: req.PortraitURL,
		Epitaph: req.Epitaph, BurialOrder: req.BurialOrder, SortOrder: req.SortOrder,
		OriginGraveID: req.OriginGraveID,
	}

	resp, err := h.svc.Update(uint(id), userID, params)
	if err != nil {
		switch err {
		case errors.ErrDeceasedNotFound:
			response.Fail(c, errcode.DeceasedNotFound)
		case errors.ErrGraveNotFound:
			response.Fail(c, errcode.GraveNotFound)
		case errors.ErrGraveAccessDenied:
			response.Fail(c, errcode.GraveAccessDenied)
		case errors.ErrDeceasedDateInvalid:
			response.Fail(c, errcode.DeceasedDateInvalid)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}
	response.SuccessData(c, resp)
}

func (h *DeceasedHandler) Delete(c *gin.Context) {
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
	response.Success(c)
}
