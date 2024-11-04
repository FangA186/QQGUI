package api

import (
	"gorm.io/gorm"
)

type GetApplicationRecordService struct {
	Db *gorm.DB
}

func GetApplicationService(db *gorm.DB) *GetApplicationRecordService {
	return &GetApplicationRecordService{Db: db}
}

// GetApplicationRecord 获取该用户的申请记录
func (s *GetApplicationRecordService) GetApplicationRecord(userID string) []ApplyForResponse {
	var applyForRequests []ApplyForResponse

	err := s.Db.Table("apply_fors").
		Select("apply_fors.*, applicant.id AS applicant_id, applicant.username AS applicant_username, applicant.avatar AS applicant_avatar, friend.id AS friend_id, friend.username AS friend_username, friend.avatar AS friend_avatar,friend.uuid AS friend_uuid").
		Joins("JOIN user_infos AS applicant ON applicant.id = apply_fors.user_id").
		Joins("JOIN user_infos AS friend ON friend.id = apply_fors.friend_id").
		Where("apply_fors.friend_id = ?", userID).
		Scan(&applyForRequests).Error

	if err != nil {
		return nil
	}

	return applyForRequests
}
