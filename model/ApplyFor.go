package model

import "gorm.io/gorm"

type ApplyFor struct {
	gorm.Model
	UserID    int  `gorm:"comment:申请添加用户ID"`                           // 申请添加用户ID
	FriendID  int  `gorm:"comment:被添加用户ID"`                            // 被添加用户ID
	IsGroup   bool `gorm:"comment:群聊邀请(true)or用户添加(false)"`            // 群聊邀请(true)or用户添加(false)
	IsConsent int  `gorm:"comment:是否同意(0代表等待,1代表同意,2代表拒绝)?;default:0"` // 是否同意
}
