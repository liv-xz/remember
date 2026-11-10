package errors

import "errors"

var ErrMemberNotFound = errors.New("成员不存在")
var ErrMemberAlreadyExist = errors.New("该用户已是此墓地成员")
var ErrMemberRoleInsufficient = errors.New("权限不足，无法执行此操作")
var ErrMemberOwnerTransfer = errors.New("所有者不能直接退出，请先转让所有权")
var ErrMemberCannotDemoteOwner = errors.New("不能降级墓地所有者")
var ErrMemberCannotRemoveOwner = errors.New("不能移除墓地所有者")
var ErrMemberTransferTargetInvalid = errors.New("转让目标必须是管理员")
var ErrMemberNotInGrave = errors.New("成员不属于该墓地")
