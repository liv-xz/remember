package admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	adminhandler "server/handler/admin"
	"server/test/testutil"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

func setupAdminTest(t *testing.T) (*gin.Engine, sqlmock.Sqlmock, func()) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	mockDB := testutil.NewMockDB(t)
	handler := adminhandler.NewAdminHandler(mockDB.DB, time.Now())

	r := gin.New()
	r.GET("/admin/stats/overview", handler.OverviewHandler)
	r.GET("/admin/trend", handler.TrendHandler)
	r.GET("/admin/hot", handler.HotHandler)
	r.GET("/admin/users", handler.ListUsersHandler)
	r.GET("/admin/graves", handler.ListGravesHandler)
	r.PUT("/admin/users/:id/status", handler.UpdateUserStatusHandler)
	r.DELETE("/admin/users/:id", handler.DeleteUserHandler)

	return r, mockDB.Mock, func() { mockDB.Close() }
}

func parseBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("JSON 解析失败: %v", err)
	}
	return body
}

func TestAllAdminRoutesRegistered(t *testing.T) {
	r, _, cleanup := setupAdminTest(t)
	defer cleanup()

	routes := r.Routes()
	paths := make(map[string]bool)
	for _, rt := range routes {
		paths[rt.Method+" "+rt.Path] = true
	}

	expected := []string{
		"GET /admin/stats/overview",
		"GET /admin/trend",
		"GET /admin/hot",
		"GET /admin/users",
		"GET /admin/graves",
		"PUT /admin/users/:id/status",
		"DELETE /admin/users/:id",
	}
	for _, p := range expected {
		if !paths[p] {
			t.Errorf("路由未注册: %s", p)
		}
	}
}

func TestOverviewHandler(t *testing.T) {
	r, mock, cleanup := setupAdminTest(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT count\(\*\) FROM .?users.? WHERE status`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(42))
	mock.ExpectQuery(`SELECT count\(\*\) FROM .?graves.?`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(30))
	mock.ExpectQuery(`SELECT count\(\*\) FROM .?deceased.?`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(15))
	mock.ExpectQuery(`SELECT count\(\*\) FROM .?memorial_records.?`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(88))
	mock.ExpectQuery(`SELECT count\(\*\) FROM .?memorial_events.?`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(22))
	mock.ExpectQuery(`SELECT count\(\*\) FROM .?grave_members.?`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(55))
	mock.ExpectQuery(`SELECT count\(\*\) FROM .?users.? WHERE DATE\(created_at\)`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	mock.ExpectQuery(`SELECT count\(\*\) FROM .?graves.? WHERE DATE\(created_at\)`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`SELECT count\(\*\) FROM .?deceased.? WHERE DATE\(created_at\)`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	mock.ExpectQuery(`SELECT count\(\*\) FROM .?users.? WHERE updated_at`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

	req := httptest.NewRequest("GET", "/admin/stats/overview", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	body := parseBody(t, w)

	if body["code"] != float64(0) {
		t.Errorf("期望 code=0, 实际 %v", body["code"])
	}
}

func TestTrendHandler(t *testing.T) {
	r, mock, cleanup := setupAdminTest(t)
	defer cleanup()

	startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	mock.ExpectQuery("SELECT DATE").WithArgs(startDate).WillReturnRows(sqlmock.NewRows([]string{"date", "count"}))
	mock.ExpectQuery("SELECT DATE").WithArgs(startDate).WillReturnRows(sqlmock.NewRows([]string{"date", "count"}))
	mock.ExpectQuery("SELECT DATE").WithArgs(startDate).WillReturnRows(sqlmock.NewRows([]string{"date", "count"}))

	req := httptest.NewRequest("GET", "/admin/trend?days=7", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	body := parseBody(t, w)

	if body["code"] != float64(0) {
		t.Errorf("期望 code=0, 实际 %v", body["code"])
	}
}

func TestHotHandler(t *testing.T) {
	r, mock, cleanup := setupAdminTest(t)
	defer cleanup()

	mock.ExpectQuery("SELECT g.id, g.name").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "created_at", "deceased_count", "member_count", "record_count"}).
			AddRow(1, "热门墓园", time.Now(), 10, 5, 20),
	)

	req := httptest.NewRequest("GET", "/admin/hot?limit=5", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	body := parseBody(t, w)

	if body["code"] != float64(0) {
		t.Errorf("期望 code=0, 实际 %v", body["code"])
	}
}

func TestListUsersHandler(t *testing.T) {
	r, mock, cleanup := setupAdminTest(t)
	defer cleanup()

	mock.ExpectQuery("SELECT count").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT .+ FROM .?users.? WHERE status != 2 ORDER BY").WillReturnRows(
		sqlmock.NewRows([]string{"id", "username", "nickname", "email", "avatar", "status", "created_at"}).
			AddRow(1, "admin", "管理员", "admin@test.com", "", 0, time.Now()),
	)

	req := httptest.NewRequest("GET", "/admin/users?page=1&size=10", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	body := parseBody(t, w)

	if body["code"] != float64(0) {
		t.Errorf("期望 code=0, 实际 %v", body["code"])
	}
}

func TestListUsers_Search(t *testing.T) {
	r, mock, cleanup := setupAdminTest(t)
	defer cleanup()

	mock.ExpectQuery("SELECT count").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery("SELECT .+ FROM .?users.? WHERE status != 2 AND .+username LIKE").WillReturnRows(
		sqlmock.NewRows([]string{"id", "username", "nickname", "email", "avatar", "status", "created_at"}),
	)

	req := httptest.NewRequest("GET", "/admin/users?page=1&size=10&search=test", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	body := parseBody(t, w)

	if body["code"] != float64(0) {
		t.Errorf("期望 code=0, 实际 %v", body["code"])
	}
}

func TestListGravesHandler(t *testing.T) {
	r, mock, cleanup := setupAdminTest(t)
	defer cleanup()

	mock.ExpectQuery("SELECT count").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "grave_type", "grave_mode", "status", "created_by", "creator_name", "deceased_count", "member_count", "created_at"}).
			AddRow(1, "测试墓园", 1, 1, 1, 1, "管理员", 3, 5, time.Now()),
	)

	req := httptest.NewRequest("GET", "/admin/graves?page=1&size=10", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	body := parseBody(t, w)

	if body["code"] != float64(0) {
		t.Errorf("期望 code=0, 实际 %v", body["code"])
	}
}

func TestUpdateUserStatus_Success(t *testing.T) {
	r, mock, cleanup := setupAdminTest(t)
	defer cleanup()

	mock.ExpectExec("UPDATE .?users.? SET .+ WHERE id = ?").WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest("PUT", "/admin/users/1/status", strings.NewReader(`{"status":1}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	body := parseBody(t, w)

	if body["code"] != float64(0) {
		t.Errorf("期望 code=0, 实际 %v", body["code"])
	}
}

func TestUpdateUserStatus_InvalidID(t *testing.T) {
	r, _, cleanup := setupAdminTest(t)
	defer cleanup()

	req := httptest.NewRequest("PUT", "/admin/users/abc/status", strings.NewReader(`{"status":1}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	body := parseBody(t, w)

	if body["code"] == float64(0) {
		t.Errorf("期望非零错误码, 实际 %v", body["code"])
	}
}

func TestUpdateUserStatus_MissingBody(t *testing.T) {
	r, _, cleanup := setupAdminTest(t)
	defer cleanup()

	req := httptest.NewRequest("PUT", "/admin/users/1/status", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	body := parseBody(t, w)

	if body["code"] == float64(0) {
		t.Errorf("期望非零错误码, 实际 %v", body["code"])
	}
}

func TestUpdateUserStatus_BadStatus(t *testing.T) {
	r, _, cleanup := setupAdminTest(t)
	defer cleanup()

	req := httptest.NewRequest("PUT", "/admin/users/1/status", strings.NewReader(`{"status":99}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	body := parseBody(t, w)

	if body["code"] == float64(0) {
		t.Errorf("期望非零错误码, 实际 %v", body["code"])
	}
}

func TestDeleteUser_Success(t *testing.T) {
	r, mock, cleanup := setupAdminTest(t)
	defer cleanup()

	mock.ExpectExec("UPDATE .?users.? SET .+ WHERE id = ?").WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest("DELETE", "/admin/users/1", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	body := parseBody(t, w)

	if body["code"] != float64(0) {
		t.Errorf("期望 code=0, 实际 %v", body["code"])
	}
}

func TestDeleteUser_InvalidID(t *testing.T) {
	r, _, cleanup := setupAdminTest(t)
	defer cleanup()

	req := httptest.NewRequest("DELETE", "/admin/users/abc", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	body := parseBody(t, w)

	if body["code"] == float64(0) {
		t.Errorf("期望非零错误码, 实际 %v", body["code"])
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	r, mock, cleanup := setupAdminTest(t)
	defer cleanup()

	mock.ExpectExec("UPDATE .?users.? SET .+ WHERE id = ?").WillReturnResult(sqlmock.NewResult(0, 0))

	req := httptest.NewRequest("DELETE", "/admin/users/999", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	body := parseBody(t, w)

	if body["code"] == float64(0) {
		t.Errorf("期望非零错误码(用户不存在), 实际 %v", body["code"])
	}
}
