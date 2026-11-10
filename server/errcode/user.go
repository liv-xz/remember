package errcode

const (
	UserNotFound       int = 1101 // 用户不存在
	UserAlreadyExist   int = 1102 // 用户名已存在
	EmailAlreadyExist  int = 1103 // 邮箱已被注册
	UserUpdateFail     int = 1104 // 用户信息更新失败
	UserDeleteFail     int = 1105 // 用户注销失败
	UsernameUpdateFreq int = 1106 // 用户名修改过于频繁
)
