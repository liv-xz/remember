package errcode

const (
	Success       int = 0  // 成功
	Unknown       int = 1  // 未知错误
	ParamBind     int = 2  // 参数绑定失败
	ParamInvalid  int = 3  // 参数校验不通过
	Unauthorized  int = 4  // 未登录/Token无效
	Forbidden     int = 5  // 无权限
	NotFound      int = 6  // 资源不存在
	MethodDenied  int = 7  // 请求方法不允许
	InternalError int = 8  // 服务器内部错误
	DatabaseError int = 9  // 数据库操作失败
	RateLimited   int = 10 // 请求过于频繁
)
