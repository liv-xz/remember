package errcode

const (
	RelationNotFound         int = 3101 // 关系不存在
	RelationCreateFail       int = 3102 // 创建关系失败
	RelationDeleteFail       int = 3103 // 删除关系失败
	RelationAlreadyExist     int = 3104 // 关系已存在
	RelationSelfInvalid      int = 3105 // 不能与自己建立关系
	RelationCrossGrave       int = 3106 // 两个逝者不属于同一墓地
	RelationDirectionInvalid int = 3107 // 关系方向不合法
)
