package errors

import "errors"

var (
	ErrRefreshTokenNotFound = errors.New("refresh token 不存在")
	ErrRefreshTokenRevoked  = errors.New("refresh token 已作废")
)
