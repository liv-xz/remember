package jwt

import (
	"testing"
	"time"

	"server/configmodel"
	myjwt "server/utils/jwt"
)

func newJWTService() *myjwt.JWTService {
	return myjwt.NewJWTService(configmodel.JwtConfig{
		Secret:        "test-secret-key",
		AccessExpire:  15,
		RefreshExpire: 7,
	})
}

func TestGenAccessToken_Success(t *testing.T) {
	svc := newJWTService()

	token, err := svc.GenAccessToken(42)
	if err != nil {
		t.Fatalf("GenAccessToken 失败: %v", err)
	}
	if token == "" {
		t.Fatal("token 为空")
	}

	claims, err := svc.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken 失败: %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("UserID = %d, want 42", claims.UserID)
	}
	if claims.Type != "access" {
		t.Errorf("Type = %s, want access", claims.Type)
	}
}

func TestGenRefreshToken_Success(t *testing.T) {
	svc := newJWTService()

	token, err := svc.GenRefreshToken(99)
	if err != nil {
		t.Fatalf("GenRefreshToken 失败: %v", err)
	}

	claims, err := svc.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken 失败: %v", err)
	}
	if claims.UserID != 99 {
		t.Errorf("UserID = %d, want 99", claims.UserID)
	}
	if claims.Type != "refresh" {
		t.Errorf("Type = %s, want refresh", claims.Type)
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	svc := newJWTService()

	_, err := svc.ParseToken("not.a.valid.jwt")
	if err == nil {
		t.Fatal("期望解析失败，但成功了")
	}
}

func TestParseToken_EmptyToken(t *testing.T) {
	svc := newJWTService()

	_, err := svc.ParseToken("")
	if err == nil {
		t.Fatal("期望解析失败，但成功了")
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	svc1 := newJWTService()
	svc2 := myjwt.NewJWTService(configmodel.JwtConfig{
		Secret:        "different-secret",
		AccessExpire:  15,
		RefreshExpire: 7,
	})

	token, err := svc1.GenAccessToken(1)
	if err != nil {
		t.Fatalf("GenAccessToken 失败: %v", err)
	}

	_, err = svc2.ParseToken(token)
	if err == nil {
		t.Fatal("用错误的 secret 解析应失败")
	}
}

func TestParseToken_RefreshTokenUsedAsAccess(t *testing.T) {
	svc := newJWTService()

	token, err := svc.GenRefreshToken(1)
	if err != nil {
		t.Fatalf("GenRefreshToken 失败: %v", err)
	}

	claims, err := svc.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken 失败: %v", err)
	}
	if claims.Type != "refresh" {
		t.Errorf("token 类型应为 refresh，实际为 %s", claims.Type)
	}
}

func TestAccessExpireDuration(t *testing.T) {
	svc := newJWTService()
	d := svc.AccessExpireDuration()
	if d.Minutes() != 15 {
		t.Errorf("AccessExpireDuration = %v, want 15 minutes", d)
	}
}

func TestRefreshExpireDuration(t *testing.T) {
	svc := newJWTService()
	d := svc.RefreshExpireDuration()
	if d.Hours() != 7*24 {
		t.Errorf("RefreshExpireDuration = %v, want 7 days", d)
	}
}

func TestGenTokens_DifferentUserIDs(t *testing.T) {
	svc := newJWTService()

	token1, _ := svc.GenAccessToken(1)
	token2, _ := svc.GenAccessToken(2)

	claims1, _ := svc.ParseToken(token1)
	claims2, _ := svc.ParseToken(token2)

	if claims1.UserID == claims2.UserID {
		t.Error("不同用户的 token 应包含不同的 UserID")
	}
}

func TestGenTokens_UniqueTokens(t *testing.T) {
	svc := newJWTService()

	token1, _ := svc.GenAccessToken(1)
	time.Sleep(time.Second)
	token2, _ := svc.GenAccessToken(1)

	if token1 == token2 {
		t.Error("相同用户多次生成的 token 应不同（IssuedAt 不同）")
	}
}
