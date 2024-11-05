package api

import (
	"QQGUI/model"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MessageService struct {
	Db *gorm.DB
	Rs *redis.Client
}

func GetMessageListService(db *gorm.DB, rs *redis.Client) *MessageService {
	return &MessageService{Db: db, Rs: rs}
}
func (a *MessageService) GetMessages(roomID string) []model.Message {
	/*
		后续这里做优化。不然每次点进去都一次获取所有的数据对数据库压力过大
		||| 用户点进某一个房间先拿着这个房间号去对应的redis房间号里面查找缓存数据，如果房间号不存在就
		先从mysql查询100条数据存到redis里面去。但是前端只先拿50条。拿完50条立马再从mysql获取后面50条
		的数据存到redis里面。并且给每条数据设置过期时间，但是时间要保持
		散点分布。防止缓存雪崩
	*/
	var messages []model.Message
	res := a.Db.Where("room_id = ?", roomID).Find(&messages)
	fmt.Println(res.Error)
	fmt.Println(messages)
	return messages
}
