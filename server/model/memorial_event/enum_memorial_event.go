package model

const (
	EventTypeBirth  int8 = 1 // 诞辰
	EventTypeDeath  int8 = 2 // 忌日
	EventTypeInter  int8 = 3 // 合葬日
	EventTypeCustom int8 = 4 // 自定义
)

const (
	LunarSolar int8 = 0 // 公历
	LunarLunar int8 = 1 // 农历
)
