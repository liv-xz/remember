package response_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"server/errcode"
	"server/utils/response"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func withRecorder(handler func(c *gin.Context)) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler(c)
	return w
}

func TestSuccess(t *testing.T) {
	w := withRecorder(func(c *gin.Context) {
		response.Success(c)
	})

	var body struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}

	if body.Code != errcode.Success {
		t.Errorf("code = %d, want %d", body.Code, errcode.Success)
	}
	if body.Message != "成功" {
		t.Errorf("message = %q, want %q", body.Message, "成功")
	}
}

func TestSuccessData(t *testing.T) {
	w := withRecorder(func(c *gin.Context) {
		response.SuccessData(c, gin.H{"id": 1, "name": "test"})
	})

	var body struct {
		Code    int                    `json:"code"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}

	if body.Code != errcode.Success {
		t.Errorf("code = %d, want %d", body.Code, errcode.Success)
	}
	if body.Data["id"] != float64(1) {
		t.Errorf("data.id = %v, want 1", body.Data["id"])
	}
	if body.Data["name"] != "test" {
		t.Errorf("data.name = %v, want test", body.Data["name"])
	}
}

func TestSuccessData_Nil(t *testing.T) {
	w := withRecorder(func(c *gin.Context) {
		response.SuccessData(c, nil)
	})

	var body struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}

	if body.Code != errcode.Success {
		t.Errorf("code = %d, want %d", body.Code, errcode.Success)
	}
}

func TestSuccessPage(t *testing.T) {
	w := withRecorder(func(c *gin.Context) {
		response.SuccessPage(c, []string{"a", "b"}, int64(50), 1, 20)
	})

	var body struct {
		Code    int      `json:"code"`
		Message string   `json:"message"`
		Data    []string `json:"data"`
		Total   int64    `json:"total"`
		Page    int      `json:"page"`
		Size    int      `json:"size"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}

	if body.Code != errcode.Success {
		t.Errorf("code = %d, want %d", body.Code, errcode.Success)
	}
	if body.Total != 50 {
		t.Errorf("total = %d, want 50", body.Total)
	}
	if body.Page != 1 {
		t.Errorf("page = %d, want 1", body.Page)
	}
	if body.Size != 20 {
		t.Errorf("size = %d, want 20", body.Size)
	}
	if len(body.Data) != 2 {
		t.Errorf("len(data) = %d, want 2", len(body.Data))
	}
}

func TestFail(t *testing.T) {
	w := withRecorder(func(c *gin.Context) {
		response.Fail(c, errcode.GraveNotFound)
	})

	var body struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}

	if body.Code != errcode.GraveNotFound {
		t.Errorf("code = %d, want %d", body.Code, errcode.GraveNotFound)
	}
	if body.Message != "墓地不存在" {
		t.Errorf("message = %q, want %q", body.Message, "墓地不存在")
	}
}

func TestFail_DataOmitted(t *testing.T) {
	w := withRecorder(func(c *gin.Context) {
		response.Fail(c, errcode.GraveNotFound)
	})

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}

	if _, exists := body["data"]; exists {
		t.Errorf("Fail response should not contain 'data' field, got %v", body)
	}
}

func TestFailMessage(t *testing.T) {
	w := withRecorder(func(c *gin.Context) {
		response.FailMessage(c, errcode.ParamInvalid, "自定义错误信息")
	})

	var body struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}

	if body.Code != errcode.ParamInvalid {
		t.Errorf("code = %d, want %d", body.Code, errcode.ParamInvalid)
	}
	if body.Message != "自定义错误信息" {
		t.Errorf("message = %q, want %q", body.Message, "自定义错误信息")
	}
}

func TestContentType(t *testing.T) {
	w := withRecorder(func(c *gin.Context) {
		response.Success(c)
	})

	ct := w.Header().Get("Content-Type")
	if ct != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q, want 'application/json; charset=utf-8'", ct)
	}
}

func TestHTTPStatus(t *testing.T) {
	w := withRecorder(func(c *gin.Context) {
		response.Fail(c, errcode.Unauthorized)
	})

	if w.Code != 200 {
		t.Errorf("HTTP status = %d, want 200", w.Code)
	}
}
