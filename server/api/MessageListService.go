package api

import (
	"QQGUI/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MessageService struct {
	Db *gorm.DB
	Rs *redis.Client
}
type GroupMessageList struct {
	model.Message
	model.UserInfo
}

func GetMessageListService(db *gorm.DB, rs *redis.Client) *MessageService {
	return &MessageService{Db: db, Rs: rs}
}

//func (a *MessageService) GetMessages(roomID string) []GroupMessageList {
//	/*
//		后续这里做优化。不然每次点进去都一次获取所有的数据对数据库压力过大
//		||| 用户点进某一个房间先拿着这个房间号去对应的redis房间号里面查找缓存数据，如果房间号不存在就
//		先从mysql查询100条数据存到redis里面去。但是前端只先拿50条。拿完50条立马再从mysql获取后面50条
//		的数据存到redis里面。并且给每条数据设置过期时间，但是时间要保持
//		散点分布。防止缓存雪崩
//	*/
//	//var messages []model.Message
//	//res := a.Db.Debug().Where("room_id = ?", roomID).Find(&messages)
//	//fmt.Println(res.Error)
//	//return messages
//	var info []GroupMessageList
//	a.Db.Debug().Raw("SELECT  messages.*,ui.id,ui.avatar,ui.username,ui.uuid,ui.signature FROM `messages` left join user_infos ui on ui.id = messages.send_user_id where room_id = ?", roomID).Scan(&info)
//	fmt.Println(info)
//	return info
//}

func (a *MessageService) GetGroupMessages(roomID string) []GroupMessageList {
	var info []GroupMessageList
	a.Db.Debug().Raw("SELECT  messages.*,ui.id,ui.avatar,ui.username,ui.uuid,ui.signature FROM `messages` left join user_infos ui on ui.id = messages.send_user_id where room_id = ?", roomID).Scan(&info)
	return info
}
