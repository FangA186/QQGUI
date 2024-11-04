package model

import (
	"database/sql"
	"gorm.io/gorm"
)

type DeletedAt sql.NullTime
type UserInfo struct {
	gorm.Model
	Username  string
	Password  string
	Signature string `gorm:"default:''"`                                      // 个性签名
	Avatar    string `gorm:"default:'http://127.0.0.1:8080/static/lxy.jpg/'"` // 头像URL
	UUID      string
}
type UserFriend struct {
	gorm.Model
	UserID   int
	Friend   int
	IsDelete bool
}

// TagUserFriend 用户好友标签分类信息
type TagUserFriend struct {
	gorm.Model
	UserID   int
	FriendID int
	TagID    int
}

type Tag struct {
	gorm.Model
	TagName string
	UserID  int
}
