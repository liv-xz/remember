package errors

import "errors"

var ErrShareLinkNotFound = errors.New("分享链接不存在")
var ErrShareLinkExpired = errors.New("分享链接已过期或失效")
