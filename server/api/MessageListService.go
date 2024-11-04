package api

import (
	"QQGUI/model"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

type MessageService struct {
	Db *gorm.DB
}

func GetMessageListService(db *gorm.DB) *MessageService {
	return &MessageService{Db: db}
}
func (a *MessageService) GetMessages(userID, friendID uint, roomID string) []model.Message {
	var messages []model.Message
	fmt.Println(userID, friendID)
	userid := strconv.Itoa(int(userID))
	friendid := strconv.Itoa(int(friendID))
	a.Db.Where("send_user_id = ? AND receiver_user_id = ? AND room_id = ?", userid, friendid, roomID).Find(&messages)
	return messages
}
