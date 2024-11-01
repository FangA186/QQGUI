package model

import "gorm.io/gorm"

// GroupMember 定义群聊成员关联模型
type GroupMember struct {
	gorm.Model
	RoomID string // 房间ID
	UserID int    // 房间成员
}
