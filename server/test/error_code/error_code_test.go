package error_code_test

import (
	"server/errcode"
	"testing"
)

func TestSuccess(t *testing.T) {
	if errcode.Success != 0 {
		t.Errorf("Success want 0, got %d", errcode.Success)
	}
}

func TestGetMessage_Success(t *testing.T) {
	msg := errcode.GetMessage(errcode.Success)
	if msg != "成功" {
		t.Errorf("GetMessage(Success) want '成功', got '%s'", msg)
	}
}

func TestGetMessage_Common(t *testing.T) {
	tests := []struct {
		code    int
		name    string
		wantMsg string
	}{
		{errcode.Success, "Success", "成功"},
		{errcode.Unknown, "Unknown", "未知错误"},
		{errcode.ParamBind, "ParamBind", "参数绑定失败"},
		{errcode.ParamInvalid, "ParamInvalid", "参数校验不通过"},
		{errcode.Unauthorized, "Unauthorized", "未登录或Token无效"},
		{errcode.Forbidden, "Forbidden", "无权限执行此操作"},
		{errcode.NotFound, "NotFound", "资源不存在"},
		{errcode.MethodDenied, "MethodDenied", "请求方法不允许"},
		{errcode.InternalError, "InternalError", "服务器内部错误"},
		{errcode.DatabaseError, "DatabaseError", "数据库操作失败"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errcode.GetMessage(tt.code)
			if got != tt.wantMsg {
				t.Errorf("GetMessage(%d) = %q, want %q", tt.code, got, tt.wantMsg)
			}
		})
	}
}

func TestGetMessage_Auth(t *testing.T) {
	tests := []struct {
		code    int
		name    string
		wantMsg string
	}{
		{errcode.TokenInvalid, "TokenInvalid", "Token无效"},
		{errcode.TokenExpired, "TokenExpired", "Token已过期"},
		{errcode.LoginRequired, "LoginRequired", "请先登录"},
		{errcode.PasswordWrong, "PasswordWrong", "密码错误"},
		{errcode.PasswordWeak, "PasswordWeak", "密码强度不足，需8位以上、含字母和数字，且不能为常见弱密码"},
		{errcode.CaptchaRequired, "CaptchaRequired", "请填写验证码"},
		{errcode.CaptchaWrong, "CaptchaWrong", "验证码错误"},
		{errcode.RegisterClosed, "RegisterClosed", "注册功能已关闭"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errcode.GetMessage(tt.code)
			if got != tt.wantMsg {
				t.Errorf("GetMessage(%d) = %q, want %q", tt.code, got, tt.wantMsg)
			}
		})
	}
}

func TestGetMessage_User(t *testing.T) {
	tests := []struct {
		code    int
		name    string
		wantMsg string
	}{
		{errcode.UserNotFound, "UserNotFound", "用户不存在"},
		{errcode.UserAlreadyExist, "UserAlreadyExist", "用户名已被占用"},
		{errcode.EmailAlreadyExist, "EmailAlreadyExist", "邮箱已被注册"},
		{errcode.UserUpdateFail, "UserUpdateFail", "用户信息更新失败"},
		{errcode.UserDeleteFail, "UserDeleteFail", "用户注销失败"},
		{errcode.UsernameUpdateFreq, "UsernameUpdateFreq", "用户名修改过于频繁，请稍后再试"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errcode.GetMessage(tt.code)
			if got != tt.wantMsg {
				t.Errorf("GetMessage(%d) = %q, want %q", tt.code, got, tt.wantMsg)
			}
		})
	}
}

func TestGetMessage_Grave(t *testing.T) {
	tests := []struct {
		code    int
		name    string
		wantMsg string
	}{
		{errcode.GraveNotFound, "GraveNotFound", "墓地不存在"},
		{errcode.GraveCreateFail, "GraveCreateFail", "创建墓地失败"},
		{errcode.GraveUpdateFail, "GraveUpdateFail", "更新墓地信息失败"},
		{errcode.GraveDeleteFail, "GraveDeleteFail", "删除墓地失败"},
		{errcode.GraveAccessDenied, "GraveAccessDenied", "无权限访问该墓地"},
		{errcode.GraveNameEmpty, "GraveNameEmpty", "墓地名称不能为空"},
		{errcode.GraveMemberLimit, "GraveMemberLimit", "该墓地成员已达上限"},
		{errcode.GraveStatusInvalid, "GraveStatusInvalid", "墓地状态异常"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errcode.GetMessage(tt.code)
			if got != tt.wantMsg {
				t.Errorf("GetMessage(%d) = %q, want %q", tt.code, got, tt.wantMsg)
			}
		})
	}
}

func TestGetMessage_Deceased(t *testing.T) {
	tests := []struct {
		code    int
		name    string
		wantMsg string
	}{
		{errcode.DeceasedNotFound, "DeceasedNotFound", "逝者不存在"},
		{errcode.DeceasedCreateFail, "DeceasedCreateFail", "添加逝者失败"},
		{errcode.DeceasedUpdateFail, "DeceasedUpdateFail", "更新逝者信息失败"},
		{errcode.DeceasedDeleteFail, "DeceasedDeleteFail", "删除逝者失败"},
		{errcode.DeceasedOverLimit, "DeceasedOverLimit", "该墓地类型无法容纳更多逝者"},
		{errcode.DeceasedDateInvalid, "DeceasedDateInvalid", "日期不合法，请检查出生日期和逝世日期"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errcode.GetMessage(tt.code)
			if got != tt.wantMsg {
				t.Errorf("GetMessage(%d) = %q, want %q", tt.code, got, tt.wantMsg)
			}
		})
	}
}

func TestGetMessage_Relation(t *testing.T) {
	tests := []struct {
		code    int
		name    string
		wantMsg string
	}{
		{errcode.RelationNotFound, "RelationNotFound", "关系不存在"},
		{errcode.RelationCreateFail, "RelationCreateFail", "创建关系失败"},
		{errcode.RelationDeleteFail, "RelationDeleteFail", "删除关系失败"},
		{errcode.RelationAlreadyExist, "RelationAlreadyExist", "这两个逝者之间已存在关系"},
		{errcode.RelationSelfInvalid, "RelationSelfInvalid", "不能与自己建立关系"},
		{errcode.RelationCrossGrave, "RelationCrossGrave", "两个逝者不属于同一墓地"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errcode.GetMessage(tt.code)
			if got != tt.wantMsg {
				t.Errorf("GetMessage(%d) = %q, want %q", tt.code, got, tt.wantMsg)
			}
		})
	}
}

func TestGetMessage_Member(t *testing.T) {
	tests := []struct {
		code    int
		name    string
		wantMsg string
	}{
		{errcode.MemberNotFound, "MemberNotFound", "成员不存在"},
		{errcode.MemberCreateFail, "MemberCreateFail", "添加成员失败"},
		{errcode.MemberUpdateFail, "MemberUpdateFail", "更新成员角色失败"},
		{errcode.MemberDeleteFail, "MemberDeleteFail", "移除成员失败"},
		{errcode.MemberAlreadyExist, "MemberAlreadyExist", "该用户已是此墓地成员"},
		{errcode.MemberRoleInsufficient, "MemberRoleInsufficient", "权限不足，无法执行此操作"},
		{errcode.MemberOwnerTransfer, "MemberOwnerTransfer", "所有者不能直接退出，请先转让所有权"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errcode.GetMessage(tt.code)
			if got != tt.wantMsg {
				t.Errorf("GetMessage(%d) = %q, want %q", tt.code, got, tt.wantMsg)
			}
		})
	}
}

func TestGetMessage_Event(t *testing.T) {
	tests := []struct {
		code    int
		name    string
		wantMsg string
	}{
		{errcode.EventNotFound, "EventNotFound", "纪念事件不存在"},
		{errcode.EventCreateFail, "EventCreateFail", "创建纪念事件失败"},
		{errcode.EventUpdateFail, "EventUpdateFail", "更新纪念事件失败"},
		{errcode.EventDeleteFail, "EventDeleteFail", "删除纪念事件失败"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errcode.GetMessage(tt.code)
			if got != tt.wantMsg {
				t.Errorf("GetMessage(%d) = %q, want %q", tt.code, got, tt.wantMsg)
			}
		})
	}
}

func TestGetMessage_Record(t *testing.T) {
	tests := []struct {
		code    int
		name    string
		wantMsg string
	}{
		{errcode.RecordNotFound, "RecordNotFound", "祭奠记录不存在"},
		{errcode.RecordCreateFail, "RecordCreateFail", "创建祭奠记录失败"},
		{errcode.RecordDeleteFail, "RecordDeleteFail", "删除祭奠记录失败"},
		{errcode.RecordNotOwned, "RecordNotOwned", "只能删除自己的祭奠记录"},
		{errcode.RecordFreqLimit, "RecordFreqLimit", "操作过于频繁，请稍后再试"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errcode.GetMessage(tt.code)
			if got != tt.wantMsg {
				t.Errorf("GetMessage(%d) = %q, want %q", tt.code, got, tt.wantMsg)
			}
		})
	}
}

func TestGetMessage_Album(t *testing.T) {
	tests := []struct {
		code    int
		name    string
		wantMsg string
	}{
		{errcode.AlbumNotFound, "AlbumNotFound", "相册不存在"},
		{errcode.AlbumCreateFail, "AlbumCreateFail", "创建相册失败"},
		{errcode.AlbumDeleteFail, "AlbumDeleteFail", "删除相册失败"},
		{errcode.PhotoNotFound, "PhotoNotFound", "照片不存在"},
		{errcode.PhotoUploadFail, "PhotoUploadFail", "照片上传失败"},
		{errcode.PhotoDeleteFail, "PhotoDeleteFail", "照片删除失败"},
		{errcode.PhotoOverLimit, "PhotoOverLimit", "照片数量已达上限"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errcode.GetMessage(tt.code)
			if got != tt.wantMsg {
				t.Errorf("GetMessage(%d) = %q, want %q", tt.code, got, tt.wantMsg)
			}
		})
	}
}

func TestGetMessage_File(t *testing.T) {
	tests := []struct {
		code    int
		name    string
		wantMsg string
	}{
		{errcode.FileTooLarge, "FileTooLarge", "文件大小超出限制（最大10MB）"},
		{errcode.FileTypeDenied, "FileTypeDenied", "文件类型不允许，仅支持jpg/png/gif"},
		{errcode.FileUploadFail, "FileUploadFail", "文件上传失败，请稍后重试"},
		{errcode.FileDeleteFail, "FileDeleteFail", "文件删除失败"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errcode.GetMessage(tt.code)
			if got != tt.wantMsg {
				t.Errorf("GetMessage(%d) = %q, want %q", tt.code, got, tt.wantMsg)
			}
		})
	}
}

func TestGetMessage_UnknownCode(t *testing.T) {
	got := errcode.GetMessage(99999)
	if got != "未知错误码" {
		t.Errorf("GetMessage(99999) want '未知错误码', got '%s'", got)
	}
}

func TestConstValues_Unique(t *testing.T) {
	seen := make(map[int]string)
	tests := []struct {
		code int
		name string
	}{
		{errcode.Success, "Success"},
		{errcode.Unknown, "Unknown"},
		{errcode.ParamBind, "ParamBind"},
		{errcode.ParamInvalid, "ParamInvalid"},
		{errcode.Unauthorized, "Unauthorized"},
		{errcode.Forbidden, "Forbidden"},
		{errcode.NotFound, "NotFound"},
		{errcode.MethodDenied, "MethodDenied"},
		{errcode.InternalError, "InternalError"},
		{errcode.DatabaseError, "DatabaseError"},
		{errcode.TokenInvalid, "TokenInvalid"},
		{errcode.TokenExpired, "TokenExpired"},
		{errcode.LoginRequired, "LoginRequired"},
		{errcode.PasswordWrong, "PasswordWrong"},
		{errcode.PasswordWeak, "PasswordWeak"},
		{errcode.CaptchaRequired, "CaptchaRequired"},
		{errcode.CaptchaWrong, "CaptchaWrong"},
		{errcode.RegisterClosed, "RegisterClosed"},
		{errcode.UserNotFound, "UserNotFound"},
		{errcode.UserAlreadyExist, "UserAlreadyExist"},
		{errcode.EmailAlreadyExist, "EmailAlreadyExist"},
		{errcode.UserUpdateFail, "UserUpdateFail"},
		{errcode.UserDeleteFail, "UserDeleteFail"},
		{errcode.UsernameUpdateFreq, "UsernameUpdateFreq"},
		{errcode.GraveNotFound, "GraveNotFound"},
		{errcode.GraveCreateFail, "GraveCreateFail"},
		{errcode.GraveUpdateFail, "GraveUpdateFail"},
		{errcode.GraveDeleteFail, "GraveDeleteFail"},
		{errcode.GraveAccessDenied, "GraveAccessDenied"},
		{errcode.GraveNameEmpty, "GraveNameEmpty"},
		{errcode.GraveMemberLimit, "GraveMemberLimit"},
		{errcode.GraveStatusInvalid, "GraveStatusInvalid"},
		{errcode.DeceasedNotFound, "DeceasedNotFound"},
		{errcode.DeceasedCreateFail, "DeceasedCreateFail"},
		{errcode.DeceasedUpdateFail, "DeceasedUpdateFail"},
		{errcode.DeceasedDeleteFail, "DeceasedDeleteFail"},
		{errcode.DeceasedOverLimit, "DeceasedOverLimit"},
		{errcode.DeceasedDateInvalid, "DeceasedDateInvalid"},
		{errcode.RelationNotFound, "RelationNotFound"},
		{errcode.RelationCreateFail, "RelationCreateFail"},
		{errcode.RelationDeleteFail, "RelationDeleteFail"},
		{errcode.RelationAlreadyExist, "RelationAlreadyExist"},
		{errcode.RelationSelfInvalid, "RelationSelfInvalid"},
		{errcode.RelationCrossGrave, "RelationCrossGrave"},
		{errcode.MemberNotFound, "MemberNotFound"},
		{errcode.MemberCreateFail, "MemberCreateFail"},
		{errcode.MemberUpdateFail, "MemberUpdateFail"},
		{errcode.MemberDeleteFail, "MemberDeleteFail"},
		{errcode.MemberAlreadyExist, "MemberAlreadyExist"},
		{errcode.MemberRoleInsufficient, "MemberRoleInsufficient"},
		{errcode.MemberOwnerTransfer, "MemberOwnerTransfer"},
		{errcode.EventNotFound, "EventNotFound"},
		{errcode.EventCreateFail, "EventCreateFail"},
		{errcode.EventUpdateFail, "EventUpdateFail"},
		{errcode.EventDeleteFail, "EventDeleteFail"},
		{errcode.RecordNotFound, "RecordNotFound"},
		{errcode.RecordCreateFail, "RecordCreateFail"},
		{errcode.RecordDeleteFail, "RecordDeleteFail"},
		{errcode.RecordNotOwned, "RecordNotOwned"},
		{errcode.RecordFreqLimit, "RecordFreqLimit"},
		{errcode.AlbumNotFound, "AlbumNotFound"},
		{errcode.AlbumCreateFail, "AlbumCreateFail"},
		{errcode.AlbumDeleteFail, "AlbumDeleteFail"},
		{errcode.PhotoNotFound, "PhotoNotFound"},
		{errcode.PhotoUploadFail, "PhotoUploadFail"},
		{errcode.PhotoDeleteFail, "PhotoDeleteFail"},
		{errcode.PhotoOverLimit, "PhotoOverLimit"},
		{errcode.FileTooLarge, "FileTooLarge"},
		{errcode.FileTypeDenied, "FileTypeDenied"},
		{errcode.FileUploadFail, "FileUploadFail"},
		{errcode.FileDeleteFail, "FileDeleteFail"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if prev, ok := seen[tt.code]; ok {
				t.Errorf("冲突: %s 和 %s 使用了相同的错误码 %d", tt.name, prev, tt.code)
			}
			seen[tt.code] = tt.name
		})
	}
}
