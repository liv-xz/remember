package errors

import "errors"

var ErrRecordNotFound = errors.New("祭奠记录不存在")
var ErrDeceasedNotInGrave = errors.New("逝者不属于该纪念馆")
