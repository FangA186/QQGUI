package model

import "gorm.io/gorm"

// Group 定义群聊模型
type Group struct {
	gorm.Model
	RoomID     string // 房间ID
	IsGroup    bool   `gorm:"default:false"` // 是否为群聊
	CreateUser int    `gorm:"default:0"`     //创建群聊的人为群主
}
