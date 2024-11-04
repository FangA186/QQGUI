package model

import "gorm.io/gorm"

type IsSpeak struct {
	gorm.Model
	UserID   int
	FriendID int
	Speak    bool
}
