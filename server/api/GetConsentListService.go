package api

import (
	"gorm.io/gorm"
)

// GetConsentListService 申请服务
type GetConsentListService struct {
	Db *gorm.DB
}

func ConsentListService(db *gorm.DB) *GetConsentListService {
	return &GetConsentListService{Db: db}
}

// GetConsentList 获取该用户被他人申请添加的记录
func (s *GetConsentListService) GetConsentList(userID string) ([]ApplyForResponse, error) {
	var applyForRequests []ApplyForResponse

	err := s.Db.Table("apply_fors").
		Select("apply_fors.*, applicant.id AS applicant_id, applicant.username AS applicant_username, applicant.avatar AS applicant_avatar, friend.id AS friend_id, friend.username AS friend_username, friend.avatar AS friend_avatar").
		Joins("JOIN user_infos AS applicant ON applicant.id = apply_fors.user_id").
		Joins("JOIN user_infos AS friend ON friend.id = apply_fors.friend_id").
		Where("apply_fors.user_id = ?", userID).
		Scan(&applyForRequests).Error

	if err != nil {
		return nil, err
	}

	return applyForRequests, nil
}
