package api

import (
	"fmt"
	"gorm.io/gorm"
)

type FriendService struct {
	Db *gorm.DB
}
type FriendRes struct {
	InfoList []Info
	Err      error
}

func NewFriendService(db *gorm.DB) *FriendService {
	return &FriendService{Db: db}
}

// GetFriendList 获取用户的好友列表
func (fs *FriendService) GetFriendList(userID string) ([]Info, error) {
	var friendslist []Info
	err := fs.Db.Table("user_friends").
		Select("user_infos.*").
		Joins("JOIN user_infos ON user_infos.id = user_friends.friend").
		Where("user_friends.user_id = ? AND user_friends.is_delete = ?", userID, false).
		Scan(&friendslist).Error

	if err != nil {
		return nil, fmt.Errorf("查询失败: %s", err)
	}
	return friendslist, nil
}

// GetFriendsOfOtherUsers 获取其他用户的好友列表
func (fs *FriendService) GetFriendsOfOtherUsers(userID string) ([]Info, error) {
	var friendsOfOtherUsers []Info
	err := fs.Db.Table("user_friends").
		Select("user_infos.*").
		Joins("JOIN user_infos ON user_infos.id = user_friends.user_id").
		Where("user_friends.friend = ? AND user_friends.is_delete = ?", userID, false).
		Scan(&friendsOfOtherUsers).Error

	if err != nil {
		return nil, fmt.Errorf("数据库查询失败: %v", err)
	}
	return friendsOfOtherUsers, nil
}
