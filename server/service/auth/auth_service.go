package auth

import (
	"crypto/sha256"
	"fmt"
	"time"

	"server/database/redis"
	"server/errors"
	loginlogmodel "server/model/login_log"
	rtmodel "server/model/token"
	usermodel "server/model/user"
	loginLogRepo "server/repository/login_log"
	tokenRepo "server/repository/token"
	userRepo "server/repository/user"
	"server/types"

	myjwt "server/utils/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     *userRepo.UserRepository
	rtRepo       *tokenRepo.RefreshTokenRepository
	loginLogRepo *loginLogRepo.LoginLogRepository
	jwtSvc       *myjwt.JWTService
}

func NewAuthService(userRepo *userRepo.UserRepository, rtRepo *tokenRepo.RefreshTokenRepository, loginLogRepo *loginLogRepo.LoginLogRepository, jwtSvc *myjwt.JWTService) *AuthService {
	return &AuthService{userRepo: userRepo, rtRepo: rtRepo, loginLogRepo: loginLogRepo, jwtSvc: jwtSvc}
}

func (s *AuthService) Register(username, email, password string) error {
	if _, err := s.userRepo.GetByUsername(username); err == nil {
		return errors.ErrUserAlreadyExist
	} else if err != errors.ErrUserNotFound {
		return err
	}

	if _, err := s.userRepo.GetByEmail(email); err == nil {
		return errors.ErrEmailAlreadyExist
	} else if err != errors.ErrUserNotFound {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &usermodel.User{
		Username: username,
		Email:    email,
		Password: string(hashed),
	}
	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	cfg := &usermodel.UserConfig{UserID: user.ID}
	if err := s.userRepo.CreateConfig(cfg); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Logout(refreshToken string) error {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(refreshToken)))
	stored, err := s.rtRepo.GetByHash(hash)
	if err != nil {
		return err
	}
	if err := s.rtRepo.Revoke(stored.ID); err != nil {
		return err
	}
	ttl := time.Until(stored.ExpiresAt)
	if ttl > 0 {
		redis.BlacklistToken(hash, ttl)
	}
	return nil
}

func (s *AuthService) Login(username, password, ip, userAgent string) (*types.LoginResp, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		_ = s.loginLogRepo.Create(&loginlogmodel.LoginLog{
			UserID: 0, IP: ip, UserAgent: userAgent, Status: "fail", Reason: "用户不存在",
		})
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		_ = s.loginLogRepo.Create(&loginlogmodel.LoginLog{
			UserID: user.ID, IP: ip, UserAgent: userAgent, Status: "fail", Reason: "密码错误",
		})
		return nil, errors.ErrPasswordWrong
	}

	_ = s.loginLogRepo.Create(&loginlogmodel.LoginLog{
		UserID: user.ID, IP: ip, UserAgent: userAgent, Status: "success", Reason: "",
	})

	return s.generateTokens(user)
}

func (s *AuthService) Refresh(refreshToken string) (*types.LoginResp, error) {
	claims, err := s.jwtSvc.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}
	if claims.Type != "refresh" {
		return nil, fmt.Errorf("invalid token type")
	}

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(refreshToken)))

	if redis.IsTokenBlacklisted(hash) {
		return nil, fmt.Errorf("token has been revoked")
	}

	if err := s.rtRepo.RevokeByHash(hash); err != nil {
		return nil, err
	}

	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining > 0 {
		redis.BlacklistToken(hash, remaining)
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	return s.generateTokens(user)
}

func (s *AuthService) generateTokens(user *usermodel.User) (*types.LoginResp, error) {
	_ = s.rtRepo.DeleteInvalidByUserID(user.ID)

	accessToken, err := s.jwtSvc.GenAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtSvc.GenRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(refreshToken)))
	rt := &rtmodel.RefreshToken{
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(s.jwtSvc.RefreshExpireDuration()),
	}
	if err := s.rtRepo.Create(rt); err != nil {
		return nil, err
	}

	return &types.LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: types.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Email:    user.Email,
		},
	}, nil
}

func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.ErrPasswordWrong
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashed)
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	_ = s.rtRepo.RevokeAllByUserID(userID)

	return nil
}
