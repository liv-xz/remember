package errors

import "errors"

var ErrRelationNotFound = errors.New("关系不存在")
var ErrRelationAlreadyExist = errors.New("这两个逝者之间已存在关系")
var ErrRelationSelfInvalid = errors.New("不能与自己建立关系")
var ErrRelationCrossGrave = errors.New("两个逝者不属于同一墓地")
var ErrRelationDirectionInvalid = errors.New("关系方向不合法")
