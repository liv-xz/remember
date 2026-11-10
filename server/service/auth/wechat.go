package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"server/config"
	"server/errors"
	usermodel "server/model/user"
	"server/types"
	"server/utils/wechat"

	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) WxLogin(code string) (*types.LoginResp, error) {
	wxResp, err := wechat.Code2Session(
		config.Global.App.WxAppID,
		config.Global.App.WxAppSecret,
		code,
	)
	if err != nil {
		return nil, fmt.Errorf("微信code2session失败: %w", err)
	}

	user, err := s.userRepo.GetByOpenID(wxResp.OpenID)
	if err == errors.ErrUserNotFound {
		user, err = s.createWechatUser(wxResp.OpenID)
		if err != nil {
			return nil, fmt.Errorf("创建微信用户失败: %w", err)
		}
	} else if err != nil {
		return nil, err
	}

	return s.generateTokens(user)
}

func (s *AuthService) createWechatUser(openID string) (*usermodel.User, error) {
	if len(openID) < 12 {
		return nil, fmt.Errorf("微信OpenID长度不足: %d", len(openID))
	}

	username := fmt.Sprintf("wx_%s", openID[:12])

	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, fmt.Errorf("生成密码失败: %w", err)
	}
	randomPwd := hex.EncodeToString(randomBytes)
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(randomPwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码哈希失败: %w", err)
	}

	user := &usermodel.User{
		Username: username,
		Nickname: "微信用户",
		Email:    fmt.Sprintf("%s@wechat.local", openID[:12]),
		Password: string(hashedPwd),
		OpenID:   &openID,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	cfg := &usermodel.UserConfig{UserID: user.ID}
	if err := s.userRepo.CreateConfig(cfg); err != nil {
		return nil, err
	}

	return user, nil
}
