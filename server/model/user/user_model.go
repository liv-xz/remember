package model

import (
	"server/model"
	"time"
)

type User struct {
	model.Model
	Username string  `gorm:"size:32;not null;unique" json:"username"`
	Nickname string  `gorm:"size:32" json:"nickname"`
	Avatar   string  `gorm:"size:255" json:"avatar"`
	Email    string  `gorm:"size:255;not null;uniqueIndex" json:"email"`
	Abstract string  `gorm:"size:255" json:"abstract"`
	Password string  `gorm:"size:255;not null" json:"-"`
	OpenID   *string `gorm:"type:varchar(64);uniqueIndex;comment:微信OpenID" json:"open_id"`
	Status   int8    `gorm:"default:0;not null;comment:0=正常,1=冻结,2=已删除" json:"status"`
}

type UserConfig struct {
	UserID             uint       `gorm:"unique" json:"userID"`
	User               User       `gorm:"foreignKey:UserID" json:"-"`
	LikeTags           []string   `gorm:"type:longtext;serializer:json" json:"LikeTags"`
	UpdateUsernameDate *time.Time `json:"updateUsernameDate"`              // 上次修改用户名的时间
	OpenCollect        bool       `gorm:"default:true" json:"openCollect"` // 是否公开收藏列表
	OpenFans           bool       `gorm:"default:true" json:"openFans"`    // 是否公开粉丝列表
	OpenFollow         bool       `gorm:"default:true" json:"openFollow"`  // 是否公开关注列表
	HomeStyleID        uint       `gorm:"default:1" json:"homeStyleID"`    // 主页风格
}
