package model

const (
	GraveTypeSingle int8 = 1 // 单人墓
	GraveTypeCouple int8 = 2 // 双人墓（夫妻）
	GraveTypeFamily int8 = 3 // 家族多人墓
)

const (
	LayoutStyleSingle   int8 = 1 // 单碑
	LayoutStyleDouble   int8 = 2 // 双碑并立
	LayoutStyleMultiple int8 = 3 // 多碑排列
)

const (
	GraveStatusActive   int8 = 1 // 启用
	GraveStatusInactive int8 = 0 // 停用
)

const (
	GraveModeMemorial   int8 = 1 // 纪念模式
	GraveModeFamilyTree int8 = 2 // 族谱模式
)
