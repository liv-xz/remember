package response

import (
	"fmt"

	"server/applog"
	"server/errcode"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Page struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}

type PageResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
	HasMore bool        `json:"hasMore"`
}

func GetUserID(c *gin.Context) (uint, error) {
	val, exists := c.Get("userID")
	if !exists {
		return 0, fmt.Errorf("userID not found in context")
	}
	userID, ok := val.(uint)
	if !ok {
		return 0, fmt.Errorf("userID type assertion failed")
	}
	return userID, nil
}

func Success(c *gin.Context) {
	c.JSON(200, Response{
		Code:    errcode.Success,
		Message: errcode.GetMessage(errcode.Success),
	})
}

func SuccessData(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    errcode.Success,
		Message: errcode.GetMessage(errcode.Success),
		Data:    data,
	})
}

func SuccessPage(c *gin.Context, data interface{}, total int64, page int, size int) {
	c.JSON(200, PageResponse{
		Code:    errcode.Success,
		Message: errcode.GetMessage(errcode.Success),
		Data:    data,
		Total:   total,
		Page:    page,
		Size:    size,
		HasMore: int64(page*size) < total,
	})
}

func Fail(c *gin.Context, code int) {
	c.JSON(200, Response{
		Code:    code,
		Message: errcode.GetMessage(code),
	})
}

func FailMessage(c *gin.Context, code int, msg string) {
	c.JSON(200, Response{
		Code:    code,
		Message: msg,
	})
}

func FailWithError(c *gin.Context, code int, err error) {
	applog.WithFields(map[string]interface{}{
		"code":  code,
		"error": err.Error(),
	}).Warn("请求失败")

	c.JSON(200, Response{
		Code:    code,
		Message: errcode.GetMessage(code),
	})
}
