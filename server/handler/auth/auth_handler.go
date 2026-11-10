package auth

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"server/config"
	"server/errcode"
	"server/errors"
	"server/middleware"
	authservice "server/service/auth"
	"server/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/nbutton23/zxcvbn-go"
)

type AuthHandler struct {
	authSvc   *authservice.AuthService
	protector *middleware.LoginProtector
}

func NewAuthHandler(authSvc *authservice.AuthService, protector *middleware.LoginProtector) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, protector: protector}
}

func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, token string) {
	maxAge := config.Global.JWT.RefreshExpire * 24 * 60 * 60
	secure := config.Global.App.Env != "dev"
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("refresh_token", token, maxAge, "/api/v1/auth/", "", secure, true)
}

func (h *AuthHandler) clearRefreshTokenCookie(c *gin.Context) {
	secure := config.Global.App.Env != "dev"
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("refresh_token", "", -1, "/api/v1/auth/", "", secure, true)
}

func (h *AuthHandler) getRefreshToken(c *gin.Context) string {
	if cookieToken, err := c.Cookie("refresh_token"); err == nil && cookieToken != "" {
		return cookieToken
	}
	return ""
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	if !isValidPassword(req.Password) {
		response.FailMessage(c, errcode.PasswordWeak, "密码强度不足，需8位以上、含字母和数字，且不能为常见弱密码")
		return
	}

	if err := h.authSvc.Register(req.Username, req.Email, req.Password); err != nil {
		switch err {
		case errors.ErrUserAlreadyExist:
			response.Fail(c, errcode.UserAlreadyExist)
		case errors.ErrEmailAlreadyExist:
			response.Fail(c, errcode.EmailAlreadyExist)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}

	response.Success(c)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	if h.protector.IsLocked(ip, req.Username) {
		remaining := h.protector.LockRemaining(ip, req.Username)
		response.FailMessage(c, errcode.LoginLocked,
			fmt.Sprintf("登录失败次数过多，请%d分钟%d秒后再试", int(remaining.Minutes()), int(remaining.Seconds())%60))
		return
	}

	resp, err := h.authSvc.Login(req.Username, req.Password, ip, userAgent)
	if err != nil {
		h.protector.RecordFailure(ip, req.Username)
		switch err {
		case errors.ErrUserNotFound:
			response.Fail(c, errcode.UserNotFound)
		case errors.ErrPasswordWrong:
			response.Fail(c, errcode.PasswordWrong)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}

	h.protector.Reset(ip, req.Username)
	h.setRefreshTokenCookie(c, resp.RefreshToken)
	response.SuccessData(c, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken := h.getRefreshToken(c)
	if refreshToken == "" {
		var req LogoutReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Fail(c, errcode.ParamInvalid)
			return
		}
		refreshToken = req.RefreshToken
	}

	if err := h.authSvc.Logout(refreshToken); err != nil {
		switch err {
		case errors.ErrRefreshTokenNotFound, errors.ErrRefreshTokenRevoked:
			response.Fail(c, errcode.TokenInvalid)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}

	h.clearRefreshTokenCookie(c)
	response.Success(c)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken := h.getRefreshToken(c)
	if refreshToken == "" {
		var req RefreshReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Fail(c, errcode.ParamInvalid)
			return
		}
		refreshToken = req.RefreshToken
	}

	resp, err := h.authSvc.Refresh(refreshToken)
	if err != nil {
		switch err {
		case errors.ErrRefreshTokenNotFound, errors.ErrRefreshTokenRevoked:
			response.Fail(c, errcode.TokenInvalid)
		default:
			response.Fail(c, errcode.Unauthorized)
		}
		return
	}

	h.setRefreshTokenCookie(c, resp.RefreshToken)
	response.SuccessData(c, resp)
}

func (h *AuthHandler) WxLogin(c *gin.Context) {
	var req WxLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	resp, err := h.authSvc.WxLogin(req.Code)
	if err != nil {
		response.FailMessage(c, errcode.WechatLoginFailed, "微信登录失败，请重试")
		return
	}

	h.setRefreshTokenCookie(c, resp.RefreshToken)
	response.SuccessData(c, resp)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, err := response.GetUserID(c)
	if err != nil {
		response.Fail(c, errcode.Unauthorized)
		return
	}

	var req ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}

	if !isValidPassword(req.NewPassword) {
		response.FailMessage(c, errcode.PasswordWeak, "新密码必须至少8位，包含字母和数字")
		return
	}

	if err := h.authSvc.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		switch err {
		case errors.ErrPasswordWrong:
			response.Fail(c, errcode.PasswordWrong)
		case errors.ErrUserNotFound:
			response.Fail(c, errcode.UserNotFound)
		default:
			response.Fail(c, errcode.InternalError)
		}
		return
	}

	response.Success(c)
}

var blacklistPasswords = []string{
	"password", "12345678", "123456789", "1234567890",
	"qwertyui", "qwerty12", "qwerty123",
	"admin123", "administrator",
	"abc12345", "abcdefgh",
	"iloveyou", "monkey12", "letmein1", "trustno1",
	"dragon12", "master12", "sunshine",
	"football", "baseball", "welcome",
	"1234qwer", "1q2w3e4r", "zxcvbnm1",
	"passw0rd", "p@ssw0rd", "pa$$word",
}

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString
	hasDigit := regexp.MustCompile(`\d`).MatchString
	if !hasLetter(password) || !hasDigit(password) {
		return false
	}

	lower := strings.ToLower(password)
	for _, bp := range blacklistPasswords {
		if lower == bp {
			return false
		}
	}

	score := zxcvbn.PasswordStrength(password, nil)
	return score.Score >= 3
}
