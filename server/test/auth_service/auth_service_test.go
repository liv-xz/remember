package auth_service

import (
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	"server/configmodel"
	"server/errors"
	loginLogRepo "server/repository/login_log"
	tokenRepo "server/repository/token"
	userRepo "server/repository/user"
	authservice "server/service/auth"
	"server/test/testutil"
	myjwt "server/utils/jwt"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	reUserByUsername  = `SELECT .+ FROM \x60users\x60 WHERE username = \? .+`
	reUserByEmail     = `SELECT .+ FROM \x60users\x60 WHERE email = \? .+`
	reUserByID        = `SELECT .+ FROM \x60users\x60 WHERE \x60users\x60.\x60id\x60 = \? .+`
	reInsUser         = `INSERT INTO \x60users\x60`
	reInsUserConfig   = `INSERT INTO \x60user_configs\x60`
	reInsLoginLog     = `INSERT INTO \x60login_logs\x60`
	reInsToken        = `INSERT INTO \x60refresh_tokens\x60`
	reDelInvalidToken = `DELETE FROM \x60refresh_tokens\x60`
	reTokenByHash     = `SELECT .+ FROM \x60refresh_tokens\x60 WHERE token_hash = \? .+`
	reRevokeToken     = `UPDATE \x60refresh_tokens\x60 SET .+ WHERE id = \?`
	reRevokeHash      = `UPDATE \x60refresh_tokens\x60 SET .+ WHERE token_hash = \? AND revoked = \?`
	reUpdUser         = `UPDATE \x60users\x60 SET .+ WHERE .+`
	reRevokeAllTokens = `UPDATE \x60refresh_tokens\x60 SET .+ WHERE user_id = \? AND revoked = \?`
)

func newJWTService() *myjwt.JWTService {
	return myjwt.NewJWTService(configmodel.JwtConfig{
		Secret:        "test-secret-key",
		AccessExpire:  15,
		RefreshExpire: 7,
	})
}

func hashPassword(p string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(h)
}

func tokenHash(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

func TestRegister_Success(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	ur := userRepo.NewUserRepository(mockDB.DB)
	tr := tokenRepo.NewRefreshTokenRepository(mockDB.DB)
	lr := loginLogRepo.NewLoginLogRepository(mockDB.DB)
	svc := authservice.NewAuthService(ur, tr, lr, newJWTService())

	mockDB.Mock.ExpectQuery(reUserByUsername).WithArgs("testuser", sqlmock.AnyArg()).
		WillReturnError(gorm.ErrRecordNotFound)
	mockDB.Mock.ExpectQuery(reUserByEmail).WithArgs("test@example.com", sqlmock.AnyArg()).
		WillReturnError(gorm.ErrRecordNotFound)
	mockDB.Mock.ExpectExec(reInsUser).WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.Mock.ExpectExec(reInsUserConfig).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := svc.Register("testuser", "test@example.com", "StrongPass1"); err != nil {
		t.Fatalf("Register 失败: %v", err)
	}
	if err := mockDB.Mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未满足的 mock 期望: %v", err)
	}
}

func TestRegister_UsernameAlreadyExists(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	ur := userRepo.NewUserRepository(mockDB.DB)
	tr := tokenRepo.NewRefreshTokenRepository(mockDB.DB)
	lr := loginLogRepo.NewLoginLogRepository(mockDB.DB)
	svc := authservice.NewAuthService(ur, tr, lr, newJWTService())

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
		AddRow(1, "testuser", "test@example.com", hashPassword("StrongPass1"), time.Now(), time.Now())
	mockDB.Mock.ExpectQuery(reUserByUsername).WithArgs("testuser", sqlmock.AnyArg()).WillReturnRows(rows)

	if err := svc.Register("testuser", "test@example.com", "StrongPass1"); err != errors.ErrUserAlreadyExist {
		t.Errorf("期望 ErrUserAlreadyExist, got %v", err)
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	ur := userRepo.NewUserRepository(mockDB.DB)
	tr := tokenRepo.NewRefreshTokenRepository(mockDB.DB)
	lr := loginLogRepo.NewLoginLogRepository(mockDB.DB)
	svc := authservice.NewAuthService(ur, tr, lr, newJWTService())

	mockDB.Mock.ExpectQuery(reUserByUsername).WithArgs("testuser", sqlmock.AnyArg()).
		WillReturnError(gorm.ErrRecordNotFound)
	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
		AddRow(1, "other", "test@example.com", hashPassword("StrongPass1"), time.Now(), time.Now())
	mockDB.Mock.ExpectQuery(reUserByEmail).WithArgs("test@example.com", sqlmock.AnyArg()).WillReturnRows(rows)

	if err := svc.Register("testuser", "test@example.com", "StrongPass1"); err != errors.ErrEmailAlreadyExist {
		t.Errorf("期望 ErrEmailAlreadyExist, got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	ur := userRepo.NewUserRepository(mockDB.DB)
	tr := tokenRepo.NewRefreshTokenRepository(mockDB.DB)
	lr := loginLogRepo.NewLoginLogRepository(mockDB.DB)
	jwtSvc := newJWTService()
	svc := authservice.NewAuthService(ur, tr, lr, jwtSvc)

	hashedPw := hashPassword("StrongPass1")
	rows := sqlmock.NewRows([]string{"id", "username", "nickname", "avatar", "email", "password", "created_at", "updated_at"}).
		AddRow(1, "testuser", "Test", "", "test@example.com", hashedPw, time.Now(), time.Now())
	mockDB.Mock.ExpectQuery(reUserByUsername).WithArgs("testuser", sqlmock.AnyArg()).WillReturnRows(rows)
	mockDB.Mock.ExpectExec(reInsLoginLog).WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.Mock.ExpectExec(reDelInvalidToken).WillReturnResult(sqlmock.NewResult(0, 0))
	mockDB.Mock.ExpectExec(reInsToken).WillReturnResult(sqlmock.NewResult(1, 1))

	resp, err := svc.Login("testuser", "StrongPass1", "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Login 失败: %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("AccessToken 为空")
	}
	if resp.RefreshToken == "" {
		t.Error("RefreshToken 为空")
	}
	if resp.User.Username != "testuser" {
		t.Errorf("Username = %s, want testuser", resp.User.Username)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	ur := userRepo.NewUserRepository(mockDB.DB)
	tr := tokenRepo.NewRefreshTokenRepository(mockDB.DB)
	lr := loginLogRepo.NewLoginLogRepository(mockDB.DB)
	svc := authservice.NewAuthService(ur, tr, lr, newJWTService())

	mockDB.Mock.ExpectQuery(reUserByUsername).WithArgs("nobody", sqlmock.AnyArg()).
		WillReturnError(gorm.ErrRecordNotFound)
	mockDB.Mock.ExpectExec(reInsLoginLog).WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := svc.Login("nobody", "anything", "127.0.0.1", "test-agent")
	if err != errors.ErrUserNotFound {
		t.Errorf("期望 ErrUserNotFound, got %v", err)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	ur := userRepo.NewUserRepository(mockDB.DB)
	tr := tokenRepo.NewRefreshTokenRepository(mockDB.DB)
	lr := loginLogRepo.NewLoginLogRepository(mockDB.DB)
	svc := authservice.NewAuthService(ur, tr, lr, newJWTService())

	rows := sqlmock.NewRows([]string{"id", "username", "nickname", "avatar", "email", "password", "created_at", "updated_at"}).
		AddRow(1, "testuser", "Test", "", "test@example.com", hashPassword("RightPass"), time.Now(), time.Now())
	mockDB.Mock.ExpectQuery(reUserByUsername).WithArgs("testuser", sqlmock.AnyArg()).WillReturnRows(rows)
	mockDB.Mock.ExpectExec(reInsLoginLog).WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := svc.Login("testuser", "WrongPass", "127.0.0.1", "test-agent")
	if err != errors.ErrPasswordWrong {
		t.Errorf("期望 ErrPasswordWrong, got %v", err)
	}
}

func TestLogout_Success(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	ur := userRepo.NewUserRepository(mockDB.DB)
	tr := tokenRepo.NewRefreshTokenRepository(mockDB.DB)
	lr := loginLogRepo.NewLoginLogRepository(mockDB.DB)
	jwtSvc := newJWTService()
	svc := authservice.NewAuthService(ur, tr, lr, jwtSvc)

	refreshToken, _ := jwtSvc.GenRefreshToken(1)

	rows := sqlmock.NewRows([]string{"id", "user_id", "token_hash", "revoked", "expires_at", "created_at", "updated_at"}).
		AddRow(1, 1, tokenHash(refreshToken), false, time.Now().Add(7*24*time.Hour), time.Now(), time.Now())
	mockDB.Mock.ExpectQuery(reTokenByHash).WithArgs(tokenHash(refreshToken), sqlmock.AnyArg()).WillReturnRows(rows)
	mockDB.Mock.ExpectExec(reRevokeToken).WillReturnResult(sqlmock.NewResult(0, 1))

	if err := svc.Logout(refreshToken); err != nil {
		t.Fatalf("Logout 失败: %v", err)
	}
}

func TestRefresh_Success(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	ur := userRepo.NewUserRepository(mockDB.DB)
	tr := tokenRepo.NewRefreshTokenRepository(mockDB.DB)
	lr := loginLogRepo.NewLoginLogRepository(mockDB.DB)
	jwtSvc := newJWTService()
	svc := authservice.NewAuthService(ur, tr, lr, jwtSvc)

	refreshToken, _ := jwtSvc.GenRefreshToken(1)

	mockDB.Mock.ExpectExec(reRevokeHash).WillReturnResult(sqlmock.NewResult(0, 1))
	userRows := sqlmock.NewRows([]string{"id", "username", "nickname", "avatar", "email", "password", "created_at", "updated_at"}).
		AddRow(1, "testuser", "Test", "", "test@example.com", hashPassword("any"), time.Now(), time.Now())
	mockDB.Mock.ExpectQuery(reUserByID).WithArgs(uint(1), sqlmock.AnyArg()).WillReturnRows(userRows)
	mockDB.Mock.ExpectExec(reDelInvalidToken).WillReturnResult(sqlmock.NewResult(0, 0))
	mockDB.Mock.ExpectExec(reInsToken).WillReturnResult(sqlmock.NewResult(1, 1))

	resp, err := svc.Refresh(refreshToken)
	if err != nil {
		t.Fatalf("Refresh 失败: %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("AccessToken 为空")
	}
}

func TestRefresh_InvalidToken(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	ur := userRepo.NewUserRepository(mockDB.DB)
	tr := tokenRepo.NewRefreshTokenRepository(mockDB.DB)
	lr := loginLogRepo.NewLoginLogRepository(mockDB.DB)
	svc := authservice.NewAuthService(ur, tr, lr, newJWTService())

	_, err := svc.Refresh("not-a-valid-token")
	if err == nil {
		t.Fatal("期望 Refresh 失败，但成功了")
	}
}

func TestChangePassword_Success(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	ur := userRepo.NewUserRepository(mockDB.DB)
	tr := tokenRepo.NewRefreshTokenRepository(mockDB.DB)
	lr := loginLogRepo.NewLoginLogRepository(mockDB.DB)
	svc := authservice.NewAuthService(ur, tr, lr, newJWTService())

	userRows := sqlmock.NewRows([]string{"id", "username", "nickname", "avatar", "email", "password", "created_at", "updated_at"}).
		AddRow(1, "testuser", "Test", "", "test@example.com", hashPassword("oldPass1"), time.Now(), time.Now())
	mockDB.Mock.ExpectQuery(reUserByID).WithArgs(uint(1), sqlmock.AnyArg()).WillReturnRows(userRows)
	mockDB.Mock.ExpectExec(reUpdUser).WillReturnResult(sqlmock.NewResult(0, 1))
	mockDB.Mock.ExpectExec(reRevokeAllTokens).WillReturnResult(sqlmock.NewResult(0, 0))

	if err := svc.ChangePassword(1, "oldPass1", "NewPass123"); err != nil {
		t.Fatalf("ChangePassword 失败: %v", err)
	}
}

func TestChangePassword_WrongOldPassword(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	ur := userRepo.NewUserRepository(mockDB.DB)
	tr := tokenRepo.NewRefreshTokenRepository(mockDB.DB)
	lr := loginLogRepo.NewLoginLogRepository(mockDB.DB)
	svc := authservice.NewAuthService(ur, tr, lr, newJWTService())

	userRows := sqlmock.NewRows([]string{"id", "username", "nickname", "avatar", "email", "password", "created_at", "updated_at"}).
		AddRow(1, "testuser", "Test", "", "test@example.com", hashPassword("realOldPass"), time.Now(), time.Now())
	mockDB.Mock.ExpectQuery(reUserByID).WithArgs(uint(1), sqlmock.AnyArg()).WillReturnRows(userRows)

	if err := svc.ChangePassword(1, "wrongOldPass", "NewPass123"); err != errors.ErrPasswordWrong {
		t.Errorf("期望 ErrPasswordWrong, got %v", err)
	}
}
