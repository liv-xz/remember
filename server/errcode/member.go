package errcode

const (
	MemberNotFound         int = 3201 // 成员不存在
	MemberCreateFail       int = 3202 // 添加成员失败
	MemberUpdateFail       int = 3203 // 更新成员角色失败
	MemberDeleteFail       int = 3204 // 移除成员失败
	MemberAlreadyExist     int = 3205 // 该用户已是成员
	MemberRoleInsufficient int = 3206 // 权限不足
	MemberOwnerTransfer    int = 3207 // 所有者不能直接退出，需先转让
)
