package errors

import "errors"

var ErrGraveNotFound = errors.New("墓地不存在")
var ErrGraveAccessDenied = errors.New("无权限访问该墓地")
