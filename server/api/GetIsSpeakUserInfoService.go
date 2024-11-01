package api

import (
	"fmt"
	"gorm.io/gorm"
)

type GetIsSpeakUserInfoService struct {
	Db *gorm.DB
}

func GetIsSpeakUserService(db *gorm.DB) *GetIsSpeakUserInfoService {
	return &GetIsSpeakUserInfoService{Db: db}
}
func (g *GetIsSpeakUserInfoService) GetIsSpeakUser(userID string) []ApplyForResponse {
	var speakInfo []ApplyForResponse

	err := g.Db.Table("is_speaks").
		Select("is_speaks.*, applicant.id AS applicant_id, applicant.username AS applicant_username, applicant.avatar AS applicant_avatar, friend.id AS friend_id, friend.username AS friend_username, friend.avatar AS friend_avatar,friend.uuid AS friend_uuid").
		Joins("JOIN user_infos AS applicant ON applicant.id = is_speaks.user_id").
		Joins("JOIN user_infos AS friend ON friend.id = is_speaks.friend_id").
		Where("is_speaks.user_id = ?", userID).
		Scan(&speakInfo).Error

	if err != nil {
		fmt.Println(err)
		return nil
	}
	return speakInfo
}
