package user

import (
	"server/errcode"
	"server/errors"
	userservice "server/service/user"
	"server/utils/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *userservice.UserService
}

func NewUserHandler(svc *userservice.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	resp, err := h.svc.GetProfile(userID)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			response.Fail(c, errcode.UserNotFound)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}

	response.SuccessData(c, resp)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	var req UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	resp, err := h.svc.UpdateProfile(userID, req.Nickname, req.Avatar, req.Abstract)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			response.Fail(c, errcode.UserNotFound)
		default:
			response.Fail(c, errcode.UserUpdateFail)
		}
		return
	}

	response.SuccessData(c, resp)
}
