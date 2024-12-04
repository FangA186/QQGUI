package model

import "gorm.io/gorm"

// GroupMember 定义群聊成员关联模型
type GroupMember struct {
	gorm.Model
	RoomID     string // 房间ID
	UserID     int    // 房间成员
	CreateUser int    `gorm:"default:0"`     // 创建者ID
	Speak      bool   `gorm:"default:false"` // 默认一开始没有在这个群里面聊天，不出现在主页聊天记录
}
