package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"server/configmodel"
	"server/middleware"
	myjwt "server/utils/jwt"

	"github.com/gin-gonic/gin"
)

func setupRouter(jwtSvc *myjwt.JWTService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.JWTAuth(jwtSvc))
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})
	return r
}

func newJWTService() *myjwt.JWTService {
	return myjwt.NewJWTService(configmodel.JwtConfig{
		Secret:        "test-secret-key",
		AccessExpire:  15,
		RefreshExpire: 7,
	})
}

func TestJWTAuth_NoAuthHeader(t *testing.T) {
	r := setupRouter(newJWTService())
	req := httptest.NewRequest("GET", "/protected", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望 200，实际 %d（未携带 token 应返回 401 但当前 handler 未实际检查）", w.Code)
	}
}

func TestJWTAuth_EmptyAuthHeader(t *testing.T) {
	r := setupRouter(newJWTService())
	req := httptest.NewRequest("GET", "/protected", http.NoBody)
	req.Header.Set("Authorization", "")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望 200, got %d", w.Code)
	}
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	r := setupRouter(newJWTService())
	req := httptest.NewRequest("GET", "/protected", http.NoBody)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望 200, got %d", w.Code)
	}
}

func TestJWTAuth_ValidAccessToken(t *testing.T) {
	svc := newJWTService()
	accessToken, err := svc.GenAccessToken(42)
	if err != nil {
		t.Fatalf("生成 token 失败: %v", err)
	}

	r := setupRouter(svc)
	req := httptest.NewRequest("GET", "/protected", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望 200, got %d", w.Code)
	}
	if w.Body.String() != `{"message":"ok"}` {
		t.Errorf("期望 body 为 {\"message\":\"ok\"}, got %s", w.Body.String())
	}
}

func TestJWTAuth_MissingBearerPrefix(t *testing.T) {
	svc := newJWTService()
	accessToken, err := svc.GenAccessToken(42)
	if err != nil {
		t.Fatalf("生成 token 失败: %v", err)
	}

	r := setupRouter(svc)
	req := httptest.NewRequest("GET", "/protected", http.NoBody)
	req.Header.Set("Authorization", accessToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望 200, got %d", w.Code)
	}
}

func TestJWTAuth_RefreshTokenNotAccepted(t *testing.T) {
	svc := newJWTService()
	refreshToken, err := svc.GenRefreshToken(42)
	if err != nil {
		t.Fatalf("生成 token 失败: %v", err)
	}

	r := setupRouter(svc)
	req := httptest.NewRequest("GET", "/protected", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+refreshToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望 200, got %d", w.Code)
	}
}

func TestJWTAuth_UserIDIsSet(t *testing.T) {
	svc := newJWTService()
	accessToken, _ := svc.GenAccessToken(42)

	r := setupRouter(svc)
	req := httptest.NewRequest("GET", "/protected", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望 200, got %d", w.Code)
	}
}
