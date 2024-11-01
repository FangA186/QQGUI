package api

import (
	"QQGUI/model"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

type ConsentFriendService struct {
	Db *gorm.DB
}

func ConsentService(db *gorm.DB) *ConsentFriendService {
	return &ConsentFriendService{Db: db}
}

// ConsentFriend1 同意或拒绝好友申请
func (fs *ConsentFriendService) ConsentFriend1(userID, friendID, consent string) (string, error) {
	userid, _ := strconv.Atoi(userID)
	friendid, _ := strconv.Atoi(friendID)
	var info model.UserInfo
	if err := fs.Db.Where("id = ?", userid).First(&info).Error; err != nil {
		return "", err
	}

	var num int
	var str string
	var message string
	fmt.Println("consent", consent)
	if consent == "2" {
		num = 2
		str = "icon-cuowu1" // 拒绝
		message = "已拒绝您的好友申请"
	} else if consent == "1" {
		num = 1
		str = "icon-duigou" // 同意
		message = "好友已通过您的好友申请"

		// 创建标签
		tag := model.Tag{
			TagName: "我的好友",
			UserID:  userid,
		}
		if err := fs.Db.Create(&tag).Error; err != nil {
			return "", err
		}

		taguserfriend := model.TagUserFriend{
			FriendID: friendid,
			UserID:   userid,
			TagID:    int(tag.ID),
		}
		if err := fs.Db.Create(&taguserfriend).Error; err != nil {
			return "", err
		}
	}
	// 检查是否已存在相同的好友关系
	var existingFriend model.UserFriend
	fs.Db.Where("(user_id = ? AND friend = ?) OR (user_id = ? AND friend = ?)", userid, friendid, friendid, userid).First(&existingFriend)

	// 如果已存在相同的好友关系，则不创建
	if existingFriend.ID != 0 {
		return "好友关系已存在", nil
	}

	// 创建好友关系
	userfriend := model.UserFriend{
		UserID:   userid,
		Friend:   friendid,
		IsDelete: false,
	}
	if err := fs.Db.Create(&userfriend).Error; err != nil {
		return "", err
	}

	// 更新申请状态
	result := fs.Db.Model(&model.ApplyFor{}).Where("user_id = ? AND friend_id = ?", friendid, userid).Update("is_consent", num)
	if result.RowsAffected > 0 {
		fmt.Println("通知了？")
		// 通知好友
		NotifyFriend(strconv.Itoa(friendid), str, message, info)
	}

	return "处理成功", nil
}
