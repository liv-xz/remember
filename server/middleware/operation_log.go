package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"

	oplogmodel "server/model/operation_log"
	oplogrepo "server/repository/operation_log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func OperationLog(repo *oplogrepo.OperationLogRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != "POST" && method != "PUT" && method != "DELETE" {
			c.Next()
			return
		}

		userID, exists := c.Get("userID")
		if !exists {
			c.Next()
			return
		}
		uid, ok := userID.(uint)
		if !ok {
			c.Next()
			return
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logrus.Warnf("读取请求体失败: %v", err)
			bodyBytes = []byte("{}")
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		c.Next()

		go func() {
			defer func() {
				if r := recover(); r != nil {
					logrus.Errorf("操作日志写入 panic: %v", r)
				}
			}()

			detail := sanitizeBody(bodyBytes)
			if detail == "" {
				detail = "{}"
			}

			log := &oplogmodel.OperationLog{
				UserID:         uid,
				Action:         actionFromMethod(method),
				TargetType:     targetTypeFromPath(c.FullPath()),
				TargetID:       targetIDFromParam(c),
				ResponseStatus: c.Writer.Status(),
				Detail:         detail,
				IP:             c.ClientIP(),
			}

			if err := repo.Create(log); err != nil {
				logrus.Warnf("操作日志写入失败: %v", err)
			}
		}()
	}
}

func actionFromMethod(method string) string {
	switch method {
	case "POST":
		return "CREATE"
	case "PUT":
		return "UPDATE"
	case "DELETE":
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

var actionSegments = map[string]bool{
	"confirm": true, "reject": true, "deactivate": true,
	"password": true, "transfer": true, "profile": true,
	"role": true, "relation-desc": true, "with-deceased": true,
}

func targetTypeFromPath(path string) string {
	trimmed := strings.TrimPrefix(path, "/api/v1/")
	segments := strings.Split(trimmed, "/")
	for i := len(segments) - 1; i >= 0; i-- {
		seg := segments[i]
		if seg == "" || strings.HasPrefix(seg, ":") {
			continue
		}
		if actionSegments[seg] {
			continue
		}
		return seg
	}
	return "unknown"
}

func targetIDFromParam(c *gin.Context) uint {
	idStr := c.Param("id")
	if idStr == "" {
		return 0
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0
	}
	return uint(id)
}

var sensitiveFields = map[string]bool{
	"password":      true,
	"oldPassword":   true,
	"newPassword":   true,
	"refresh_token": true,
	"token":         true,
}

func sanitizeBody(raw []byte) string {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" {
		return "{}"
	}

	var body map[string]interface{}
	if err := json.Unmarshal(raw, &body); err != nil {
		return trimmed
	}

	for key := range body {
		if sensitiveFields[key] || strings.Contains(strings.ToLower(key), "password") || strings.Contains(strings.ToLower(key), "secret") {
			body[key] = "***FILTERED***"
		}
	}

	sanitized, err := json.Marshal(body)
	if err != nil {
		return trimmed
	}
	return string(sanitized)
}
