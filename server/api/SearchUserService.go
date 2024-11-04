package api

import (
	"QQGUI/model"
	"gorm.io/gorm"
)

// SearchUserService 用户服务
type SearchUserService struct {
	Db *gorm.DB
}

func SearchService(db *gorm.DB) *SearchUserService {
	return &SearchUserService{Db: db}
}

// SearchUser 搜索用户并添加 is_friend 字段
func (s *UserService) SearchUser(currentUserID, username string) ([]UserSearchResult, error) {
	var friendList []model.UserFriend
	// 查询当前用户的好友列表
	if err := s.Db.Where("user_id = ? AND is_delete = ?", currentUserID, false).Find(&friendList).Error; err != nil {
		return nil, err
	}

	// 将好友列表转换为映射
	friendMap := make(map[int]bool)
	for _, friend := range friendList {
		friendMap[friend.Friend] = true
	}

	var searchResults []model.UserInfo
	// 根据用户名搜索用户
	if err := s.Db.Select("id", "username", "signature", "avatar", "uuid").
		Where("username LIKE ?", "%"+username+"%").Find(&searchResults).Error; err != nil {
		return nil, err
	}

	var results []UserSearchResult
	for _, user := range searchResults {
		result := UserSearchResult{
			UserInfo: user,
			IsFriend: 0,
		}
		// 检查好友关系
		if friendMap[int(user.ID)] {
			result.IsFriend = 1
		} else {
			var reverseFriend model.UserFriend
			if err := s.Db.Where("user_id = ? AND friend = ? AND is_delete = ?", user.ID, currentUserID, false).First(&reverseFriend).Error; err == nil {
				result.IsFriend = 1
			}
		}
		results = append(results, result)
	}

	return results, nil
}
