package jwt

import (
	"errors"
	"server/configmodel"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

// var (
// 	jwtSecret []byte
// 	secOnce   sync.Once
// )

// func getJwtSecret() []byte {
// 	secOnce.Do(func() {
// 		jwtSecret = []byte(config.Global.JWT.Secret)
// 	})
// 	return jwtSecret
// }

type JWTService struct {
	secret        []byte
	accessExpire  int
	refreshExpire int
}

func NewJWTService(cfg configmodel.JwtConfig) *JWTService {
	return &JWTService{
		secret:        []byte(cfg.Secret),
		accessExpire:  cfg.AccessExpire,
		refreshExpire: cfg.RefreshExpire,
	}
}

func (s *JWTService) GenAccessToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.accessExpire) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}
func (s *JWTService) GenRefreshToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.refreshExpire) * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}
func (s *JWTService) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
func (s *JWTService) AccessExpireDuration() time.Duration {
	return time.Duration(s.accessExpire) * time.Minute
}
func (s *JWTService) RefreshExpireDuration() time.Duration {
	return time.Duration(s.refreshExpire) * 24 * time.Hour
}
