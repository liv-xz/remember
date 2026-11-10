package memorial_record

import (
	"strconv"

	"server/errcode"
	"server/errors"
	recordservice "server/service/memorial_record"
	"server/utils/paginator"
	"server/utils/response"

	"github.com/gin-gonic/gin"
)

type MemorialRecordHandler struct {
	svc *recordservice.MemorialRecordService
}

func NewMemorialRecordHandler(svc *recordservice.MemorialRecordService) *MemorialRecordHandler {
	return &MemorialRecordHandler{svc: svc}
}

func (h *MemorialRecordHandler) Create(c *gin.Context) {
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

	var req CreateRecordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	resp, err := h.svc.Create(uint(gid), userID, req.Type, req.DeceasedID, req.Message)
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

func (h *MemorialRecordHandler) ListByGrave(c *gin.Context) {
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

func (h *MemorialRecordHandler) ListByDeceased(c *gin.Context) {
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

	page, size := parsePage(c)
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

func (h *MemorialRecordHandler) Delete(c *gin.Context) {
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
		case errors.ErrRecordNotFound:
			response.Fail(c, errcode.RecordNotFound)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}
	response.Success(c)
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
