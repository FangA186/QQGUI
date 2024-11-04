package sqlconfig

import (
	"QQGUI/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// Redis 客户端实例
var rdb *redis.Client

// InitRedis 初始化 Redis 客户端
func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 地址
		Password: "",               // 没有密码
		DB:       0,                // 默认数据库
	})
}

// PublishMessage 发布消息到 Redis
func PublishMessage(channel string, message model.Message) {
	marshal, err := json.Marshal(message)
	if err != nil {
		return
	}
	err1 := rdb.Publish(context.Background(), channel, marshal).Err() // 发布消息
	if err1 != nil {
		fmt.Println(err1)
		// 处理错误
		return // 如果出错则返回
	}
}

// SubscribeMessage 订阅 Redis 消息
func SubscribeMessage(channel string) *redis.PubSub {
	pubsub := rdb.Subscribe(context.Background(), channel) // 订阅指定频道
	return pubsub                                          // 返回 PubSub 实例
}
