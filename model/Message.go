package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Content        string `json:"content"`                    // 消息内容（文字）
	SendUserID     uint   `json:"send_user_id"`               // 发送用户ID
	ReceiverUserID uint   `json:"receiver_user_id"`           // 接收者用户ID
	RoomID         string `json:"room_id"`                    // 群聊ID
	FileURL        string `json:"file_url" gorm:"default:''"` // 文件的URL或路径
}

// 如果有文件上传那么给文件也单独创建一条信息。
