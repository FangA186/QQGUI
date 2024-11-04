package api

import (
	"QQGUI/model"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

type HandleAddFriendRequestService struct {
	Db *gorm.DB
}

func HandleAddFriendService(db *gorm.DB) *HandleAddFriendRequestService {
	return &HandleAddFriendRequestService{Db: db}
}

// HandleAddFriendService 处理添加好友请求
func (fs *HandleAddFriendRequestService) HandleAddFriendService(userID, friendID string, isGroup bool) (string, error) {
	userid, _ := strconv.Atoi(userID)
	friendid, _ := strconv.Atoi(friendID) // 被申请ID
	// 检查是否已存在未处理的好友请求
	var existingRequest model.ApplyFor
	if err := fs.Db.Where("user_id = ? AND friend_id = ? AND is_consent = 0", userid, friendid).First(&existingRequest).Error; err == nil {
		return "已存在好友请求，请勿重复发送", nil
	}

	// 创建 ApplyFor 实例
	applyFor := model.ApplyFor{
		UserID:   userid,
		FriendID: friendid,
		IsGroup:  isGroup,
	}

	// 插入新请求
	if err := fs.Db.Create(&applyFor).Error; err != nil {
		return "好友请求发送失败", err
	}

	// 发送通知给被添加的好友
	var info model.UserInfo
	if err := fs.Db.Where("id = ?", friendid).First(&info).Error; err == nil {
		NotifyFriend(friendID, "icon-xiaolian1", "请求添加你为好友", info)
	}
	fmt.Println("adasdasdasdasdas")
	return "好友请求已发送", nil
}
