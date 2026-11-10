package grave_share_link

import (
	"strconv"

	"server/errcode"
	"server/errors"
	gslservice "server/service/grave_share_link"
	"server/types"
	"server/utils/paginator"
	"server/utils/response"

	"github.com/gin-gonic/gin"
)

// GraveShareLinkHandler 墓地分享链接处理器
// 负责处理墓地分享链接的 HTTP 请求
type GraveShareLinkHandler struct {
	svc *gslservice.GraveShareLinkService
}

// CreateShareLinkReq 创建分享链接请求结构
type CreateShareLinkReq struct {
	ExpireDays int `json:"expireDays" binding:"omitempty,min=1,max=365"`
}

// NewGraveShareLinkHandler 创建墓地分享链接处理器实例
// 参数:
//   - svc: 墓地分享链接服务
//
// 返回:
//   - GraveShareLinkHandler: 处理器实例
func NewGraveShareLinkHandler(svc *gslservice.GraveShareLinkService) *GraveShareLinkHandler {
	return &GraveShareLinkHandler{svc: svc}
}

// Create 创建分享链接
// 只有墓地创建者可以创建分享链接
// 分享链接用于邀请家人访问墓园并建立与逝者的关系
// 路由: POST /api/v1/graves/:gid/share-links
// 认证: 需要登录，且必须是墓地创建者
func (h *GraveShareLinkHandler) Create(c *gin.Context) {
	// 1. 从上下文中获取用户ID（由 JWT 中间件设置）
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	// 2. 解析路径参数中的墓地ID
	gid, err := strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	// 3. 绑定和验证请求体
	var req CreateShareLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	// 4. 调用服务层创建分享链接
	resp, err := h.svc.Create(uint(gid), userID, req.ExpireDays)
	if err != nil {
		// 5. 根据错误类型返回相应的错误码
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

	// 6. 生成完整的前端分享 URL（前端会根据 shareCode 生成）
	// 这里保留 shareURL 为空，前端根据域名和 shareCode 拼接
	// 例如: https://yourdomain.com/share/{shareCode}

	// 7. 返回成功响应
	response.SuccessData(c, resp)
}

// Validate 验证分享链接
// 这是公开接口，用户通过分享链接访问时调用
// 验证链接是否有效（存在、未过期、未达到使用上限）
// 路由: GET /api/v1/share-links/:code
// 认证: 无需登录（公开接口）
func (h *GraveShareLinkHandler) Validate(c *gin.Context) {
	// 1. 从路径参数中获取分享码
	shareCode := c.Param("code")
	if shareCode == "" {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	// 2. 调用服务层验证分享链接
	resp, err := h.svc.Validate(shareCode)
	if err != nil {
		// 3. 根据错误类型返回相应的错误码
		switch err {
		case errors.ErrShareLinkNotFound:
			response.Fail(c, errcode.ShareLinkNotFound)
		case errors.ErrShareLinkExpired:
			response.Fail(c, errcode.ShareLinkExpired)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}

	// 4. 返回成功响应（包含墓园基本信息）
	response.SuccessData(c, resp)
}

// Deactivate 停用分享链接
// 只有墓地创建者可以停用分享链接
// 停用后，该链接将无法再被使用
// 路由: PUT /api/v1/graves/:gid/share-links/:id/deactivate
// 认证: 需要登录，且必须是墓地创建者
func (h *GraveShareLinkHandler) Deactivate(c *gin.Context) {
	// 1. 从上下文中获取用户ID
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	// 2. 解析路径参数中的链接ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	// 3. 调用服务层停用链接
	err = h.svc.Deactivate(uint(id), userID)
	if err != nil {
		// 4. 根据错误类型返回相应的错误码
		switch err {
		case errors.ErrShareLinkNotFound:
			response.Fail(c, errcode.ShareLinkNotFound)
		case errors.ErrGraveAccessDenied:
			response.Fail(c, errcode.GraveAccessDenied)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}

	// 5. 返回成功响应
	response.Success(c)
}

// ListByGrave 获取墓园的所有分享链接
// 只有墓地创建者可以查看
// 用于在管理页面显示所有有效的分享链接
// 路由: GET /api/v1/graves/:gid/share-links
// 认证: 需要登录，且必须是墓地创建者
func (h *GraveShareLinkHandler) ListByGrave(c *gin.Context) {
	// 1. 从上下文中获取用户ID
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	// 2. 解析路径参数中的墓地ID
	gid, err := strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	// 3. 调用服务层获取分享链接列表
	page, size := paginator.BindPage(c)
	resp, total, err := h.svc.ListByGrave(uint(gid), userID, page, size)
	if err != nil {
		// 4. 根据错误类型返回相应的错误码
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

	// 5. 返回成功响应
	response.SuccessPage(c, resp.Items, total, page, size)
}

// Join 通过分享链接加入墓园
// 分享链接即墓园 Owner 对该用户的提前背书，加入后直接成为成员无需审批
// 路由: POST /api/v1/share-links/:code/join
// 认证: 需要登录
func (h *GraveShareLinkHandler) Join(c *gin.Context) {
	shareCode := c.Param("code")
	if shareCode == "" {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	var req types.JoinShareLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	resp, err := h.svc.Join(shareCode, userID, req)
	if err != nil {
		switch err {
		case errors.ErrShareLinkNotFound:
			response.Fail(c, errcode.ShareLinkNotFound)
		case errors.ErrShareLinkExpired:
			response.Fail(c, errcode.ShareLinkExpired)
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
