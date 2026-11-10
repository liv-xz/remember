package errcode

const (
	DeceasedNotFound    int = 3001 // 逝者不存在
	DeceasedCreateFail  int = 3002 // 添加逝者失败
	DeceasedUpdateFail  int = 3003 // 更新逝者信息失败
	DeceasedDeleteFail  int = 3004 // 删除逝者失败
	DeceasedOverLimit   int = 3005 // 超过该类型墓地可容纳逝者数量
	DeceasedDateInvalid int = 3006 // 逝者日期不合法
)
