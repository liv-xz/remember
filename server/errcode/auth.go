package errcode

const (
	TokenInvalid      int = 1001 // Token无效
	TokenExpired      int = 1002 // Token已过期
	LoginRequired     int = 1003 // 需要登录
	PasswordWrong     int = 1004 // 密码错误
	PasswordWeak      int = 1005 // 密码强度不足
	CaptchaRequired   int = 1006 // 需要验证码
	CaptchaWrong      int = 1007 // 验证码错误
	RegisterClosed    int = 1008 // 注册已关闭
	WechatLoginFailed int = 1009 // 微信登录失败
	WechatCodeInvalid int = 1010 // 微信code无效
	LoginLocked       int = 1011 // 登录失败次数过多，账号已锁定
)
