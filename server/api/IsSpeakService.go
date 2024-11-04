package api

import (
	"QQGUI/model"
	"gorm.io/gorm"
	"strconv"
)

type IsSpeakService struct {
	Db *gorm.DB
}

func SpeakService(db *gorm.DB) *IsSpeakService {
	return &IsSpeakService{Db: db}
}

func (i *IsSpeakService) Speak(userID, friendID string) (message string, status int) {
	var speakInfo model.IsSpeak
	i.Db.Where("user_id = ? AND friend_id = ?", userID, friendID).First(&speakInfo)
	if speakInfo.ID != 0 {
		return "已聊天", 200
	}
	speakInfo.Speak = true
	speakInfo.UserID, _ = strconv.Atoi(userID)
	speakInfo.FriendID, _ = strconv.Atoi(friendID)
	res := i.Db.Create(&speakInfo)
	if res.RowsAffected == 0 {
		return "创建聊天状态失败", 400
	}
	return "创建聊天状态成功", 200
}
