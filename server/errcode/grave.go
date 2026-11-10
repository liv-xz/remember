package errcode

const (
	GraveNotFound      int = 2001 // 墓地不存在
	GraveCreateFail    int = 2002 // 创建墓地失败
	GraveUpdateFail    int = 2003 // 更新墓地失败
	GraveDeleteFail    int = 2004 // 删除墓地失败
	GraveAccessDenied  int = 2005 // 无权限访问该墓地
	GraveNameEmpty     int = 2006 // 墓地名称不能为空
	GraveMemberLimit   int = 2007 // 墓地成员已达上限
	GraveStatusInvalid int = 2008 // 墓地状态异常
)
