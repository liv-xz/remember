package errcode

const (
	RecordNotFound     int = 4101 // 祭奠记录不存在
	RecordCreateFail   int = 4102 // 创建祭奠记录失败
	RecordDeleteFail   int = 4103 // 删除祭奠记录失败
	RecordNotOwned     int = 4104 // 只能删除自己的祭奠记录
	RecordFreqLimit    int = 4105 // 操作过于频繁
	DeceasedNotInGrave int = 4106 // 逝者不属于该纪念馆
)
