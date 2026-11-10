package grave

import (
	"strconv"

	"server/errcode"
	"server/errors"
	deceasedservice "server/service/deceased"
	graveservice "server/service/grave"
	"server/types"
	"server/utils/paginator"
	"server/utils/response"

	"github.com/gin-gonic/gin"
)

type GraveHandler struct {
	svc         *graveservice.GraveService
	deceasedSvc *deceasedservice.DeceasedService
}

func NewGraveHandler(svc *graveservice.GraveService, deceasedSvc *deceasedservice.DeceasedService) *GraveHandler {
	return &GraveHandler{svc: svc, deceasedSvc: deceasedSvc}
}

func (h *GraveHandler) Create(c *gin.Context) {
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	var req CreateGraveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	resp, err := h.svc.Create(userID, req.Name, req.GraveType, req.GraveMode, req.LayoutStyle, req.Description, req.CoverImage, req.Latitude, req.Longitude, req.Address)
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

	response.SuccessData(c, resp)
}

func (h *GraveHandler) GetByID(c *gin.Context) {
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	graveID, err := strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	resp, err := h.svc.GetByID(userID, uint(graveID))
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

	response.SuccessData(c, resp)
}

func (h *GraveHandler) ListMy(c *gin.Context) {
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	page, size := paginator.BindPage(c)
	resp, total, err := h.svc.ListMy(userID, page, size)
	if err != nil {
		response.Fail(c, errcode.InternalError)
		return
	}

	response.SuccessPage(c, resp.Items, total, page, size)
}

func (h *GraveHandler) Update(c *gin.Context) {
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	graveID, err := strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	var req UpdateGraveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	params := &graveservice.UpdateGraveParams{}
	if req.Name != "" {
		params.Name = &req.Name
	}
	if req.GraveType != 0 {
		params.GraveType = &req.GraveType
	}
	if req.GraveMode != 0 {
		params.GraveMode = &req.GraveMode
	}
	if req.LayoutStyle != 0 {
		params.LayoutStyle = &req.LayoutStyle
	}
	if req.Description != "" {
		params.Description = &req.Description
	}
	if req.CoverImage != "" {
		params.CoverImage = &req.CoverImage
	}
	if req.Status != nil {
		params.Status = req.Status
	}
	if req.Latitude != nil {
		params.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		params.Longitude = req.Longitude
	}
	if req.Address != "" {
		params.Address = &req.Address
	}
	if req.ClearLocation != nil && *req.ClearLocation {
		params.ClearLocation = true
	}

	resp, err := h.svc.Update(userID, uint(graveID), params)
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

	response.SuccessData(c, resp)
}

func (h *GraveHandler) Delete(c *gin.Context) {
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	graveID, err := strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	err = h.svc.Delete(userID, uint(graveID))
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

	response.Success(c)
}

func (h *GraveHandler) CreateWithDeceased(c *gin.Context) {
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	var req CreateGraveWithDeceasedReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailMessage(c, errcode.ParamInvalid, err.Error())
		return
	}

	if req.GraveType == 1 && len(req.Deceased) > 1 {
		response.Fail(c, errcode.ParamInvalid)
		return
	}
	if req.GraveType == 2 && len(req.Deceased) > 2 {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	graveResp, err := h.svc.Create(userID, req.Name, req.GraveType, req.GraveMode, req.LayoutStyle, req.Description, req.CoverImage, req.Latitude, req.Longitude, req.Address)
	if err != nil {
		response.Fail(c, errcode.InternalError)
		return
	}

	deceasedList := make([]types.DeceasedResp, 0, len(req.Deceased))
	for _, item := range req.Deceased {
		params := deceasedservice.CreateDeceasedParams{
			Name:          item.Name,
			Gender:        item.Gender,
			BirthDate:     item.BirthDate,
			DeathDate:     item.DeathDate,
			Biography:     item.Biography,
			PortraitURL:   item.PortraitURL,
			Epitaph:       item.Epitaph,
			BurialOrder:   item.BurialOrder,
			OriginGraveID: item.OriginGraveID,
		}
		deceasedResp, err := h.deceasedSvc.Create(graveResp.ID, userID, params)
		if err != nil {
			response.Fail(c, errcode.InternalError)
			return
		}
		deceasedList = append(deceasedList, *deceasedResp)
	}

	response.SuccessData(c, types.GraveWithDeceasedResp{
		Grave:    *graveResp,
		Deceased: deceasedList,
	})
}
