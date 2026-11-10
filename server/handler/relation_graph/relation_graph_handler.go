package relation_graph

import (
	"strconv"

	"server/errcode"
	"server/errors"
	rgservice "server/service/relation_graph"
	"server/utils/response"

	"github.com/gin-gonic/gin"
)

// RelationGraphHandler 关系图处理器
// 负责处理关系图数据的 HTTP 请求
type RelationGraphHandler struct {
	svc *rgservice.RelationGraphService
}

// NewRelationGraphHandler 创建关系图处理器实例
// 参数:
//   - svc: 关系图服务
//
// 返回:
//   - RelationGraphHandler: 处理器实例
func NewRelationGraphHandler(svc *rgservice.RelationGraphService) *RelationGraphHandler {
	return &RelationGraphHandler{svc: svc}
}

// GetGraph 获取关系图数据
// 关系图用于展示家族谱系，包括：
//   - 节点（Node）：逝者节点和用户节点
//   - 边（Edge）：逝者间关系、用户-逝者关系
//
// 前端可以使用这些数据来绘制家谱图
//
// 路由: GET /api/v1/graves/:gid/relation-graph
// 认证: 需要登录（墓地成员或创建者）
//
// 请求示例:
//
//	GET /api/v1/graves/1/relation-graph
//	Authorization: Bearer <token>
//
// 响应示例:
//
//	{
//	  "code": 0,
//	  "message": "成功",
//	  "data": {
//	    "graveID": 1,
//	    "graveName": "张家墓园",
//	    "nodes": [
//	      {
//	        "id": "deceased_1",
//	        "type": "deceased",
//	        "name": "张三",
//	        "avatar": "",
//	        "gender": 1,
//	        "generation": 1,
//	        "isDead": true
//	      },
//	      {
//	        "id": "user_10",
//	        "type": "user",
//	        "name": "张小明",
//	        "avatar": "",
//	        "gender": 1,
//	        "generation": 0,
//	        "isDead": false
//	      }
//	    ],
//	    "edges": [
//	      {
//	        "from": "user_10",
//	        "to": "deceased_1",
//	        "relationType": "son",
//	        "confirmed": true,
//	        "confirmationStatus": 1,
//	        "confirmedBy": 10
//	      }
//	    ]
//	  }
//	}
func (h *RelationGraphHandler) GetGraph(c *gin.Context) {
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

	// 3. 调用服务层获取关系图数据
	resp, err := h.svc.GetGraph(uint(gid), userID)
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
	response.SuccessData(c, resp)
}
