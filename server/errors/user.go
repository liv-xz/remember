package errors

import "errors"

var (
	ErrUserNotFound      = errors.New("用户不存在")
	ErrConfigNotFound    = errors.New("用户配置不存在")
	ErrUserAlreadyExist  = errors.New("用户名已存在")
	ErrEmailAlreadyExist = errors.New("邮箱已存在")
)

var ErrPasswordWrong = errors.New("密码错误")
