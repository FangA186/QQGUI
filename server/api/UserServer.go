package api

import (
	"QQGUI/model"
	"gorm.io/gorm"
)

type UserService struct {
	Db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{Db: db}
}

func (s *UserService) GetUserByUsername(username string) (model.UserInfo, error) {
	var user model.UserInfo
	err := s.Db.Where("username = ?", username).First(&user).Error
	return user, err
}

func (s *UserService) CreateUser(info *model.UserInfo) error {
	err := s.Db.Create(info).Error
	if err != nil {
		return err
	}
	return err
}
