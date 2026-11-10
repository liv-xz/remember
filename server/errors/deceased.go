package errors

import "errors"

var ErrDeceasedNotFound = errors.New("逝者不存在")
var ErrDeceasedOverLimit = errors.New("超过该类型墓地可容纳逝者数量")
var ErrDeceasedDateInvalid = errors.New("逝者日期不合法")
