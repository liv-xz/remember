package errcode

var messageMap = map[int]string{
	// 通用
	Success:       "成功",
	Unknown:       "未知错误",
	ParamBind:     "参数绑定失败",
	ParamInvalid:  "参数校验不通过",
	Unauthorized:  "未登录或Token无效",
	Forbidden:     "无权限执行此操作",
	NotFound:      "资源不存在",
	MethodDenied:  "请求方法不允许",
	InternalError: "服务器内部错误",
	DatabaseError: "数据库操作失败",
	RateLimited:   "请求过于频繁，请稍后再试",

	// 鉴权
	TokenInvalid:      "Token无效",
	TokenExpired:      "Token已过期",
	LoginRequired:     "请先登录",
	PasswordWrong:     "密码错误",
	PasswordWeak:      "密码强度不足，需8位以上、含字母和数字，且不能为常见弱密码",
	CaptchaRequired:   "请填写验证码",
	CaptchaWrong:      "验证码错误",
	RegisterClosed:    "注册功能已关闭",
	WechatLoginFailed: "微信登录失败，请重试",
	WechatCodeInvalid: "微信code无效",
	LoginLocked:       "登录失败次数过多，请15分钟后再试",

	// 用户
	UserNotFound:       "用户不存在",
	UserAlreadyExist:   "用户名已被占用",
	EmailAlreadyExist:  "邮箱已被注册",
	UserUpdateFail:     "用户信息更新失败",
	UserDeleteFail:     "用户注销失败",
	UsernameUpdateFreq: "用户名修改过于频繁，请稍后再试",

	// 墓地
	GraveNotFound:      "墓地不存在",
	GraveCreateFail:    "创建墓地失败",
	GraveUpdateFail:    "更新墓地信息失败",
	GraveDeleteFail:    "删除墓地失败",
	GraveAccessDenied:  "无权限访问该墓地",
	GraveNameEmpty:     "墓地名称不能为空",
	GraveMemberLimit:   "该墓地成员已达上限",
	GraveStatusInvalid: "墓地状态异常",

	// 分享链接
	ShareLinkNotFound:   "分享链接不存在",
	ShareLinkExpired:    "分享链接已过期或失效",
	ShareLinkCreateFail: "创建分享链接失败",

	// 逝者
	DeceasedNotFound:    "逝者不存在",
	DeceasedCreateFail:  "添加逝者失败",
	DeceasedUpdateFail:  "更新逝者信息失败",
	DeceasedDeleteFail:  "删除逝者失败",
	DeceasedOverLimit:   "该墓地类型无法容纳更多逝者",
	DeceasedDateInvalid: "日期不合法，请检查出生日期和逝世日期",

	// 逝者关系
	RelationNotFound:         "关系不存在",
	RelationCreateFail:       "创建关系失败",
	RelationDeleteFail:       "删除关系失败",
	RelationAlreadyExist:     "这两个逝者之间已存在关系",
	RelationSelfInvalid:      "不能与自己建立关系",
	RelationCrossGrave:       "两个逝者不属于同一墓地",
	RelationDirectionInvalid: "关系方向不合法，夫妻/兄弟应为0，父子/母子/祖孙应为1或2",

	// 墓位成员
	MemberNotFound:         "成员不存在",
	MemberCreateFail:       "添加成员失败",
	MemberUpdateFail:       "更新成员角色失败",
	MemberDeleteFail:       "移除成员失败",
	MemberAlreadyExist:     "该用户已是此墓地成员",
	MemberRoleInsufficient: "权限不足，无法执行此操作",
	MemberOwnerTransfer:    "所有者不能直接退出，请先转让所有权",

	// 纪念事件
	EventNotFound:   "纪念事件不存在",
	EventCreateFail: "创建纪念事件失败",
	EventUpdateFail: "更新纪念事件失败",
	EventDeleteFail: "删除纪念事件失败",

	// 祭奠记录
	RecordNotFound:     "祭奠记录不存在",
	RecordCreateFail:   "创建祭奠记录失败",
	RecordDeleteFail:   "删除祭奠记录失败",
	RecordNotOwned:     "只能删除自己的祭奠记录",
	RecordFreqLimit:    "操作过于频繁，请稍后再试",
	DeceasedNotInGrave: "逝者不属于该纪念馆",

	// 相册
	AlbumNotFound:   "相册不存在",
	AlbumCreateFail: "创建相册失败",
	AlbumDeleteFail: "删除相册失败",
	PhotoNotFound:   "照片不存在",
	PhotoUploadFail: "照片上传失败",
	PhotoDeleteFail: "照片删除失败",
	PhotoOverLimit:  "照片数量已达上限",

	// 文件上传
	FileTooLarge:   "文件大小超出限制（最大10MB）",
	FileTypeDenied: "文件类型不允许，仅支持jpg/png/gif",
	FileUploadFail: "文件上传失败，请稍后重试",
	FileDeleteFail: "文件删除失败",
}

func GetMessage(code int) string {
	msg, ok := messageMap[code]
	if !ok {
		return "未知错误码"
	}
	return msg
}
