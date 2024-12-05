package main

import (
	"QQGUI/config/sqlconfig"
	"QQGUI/server/api"
	"context"
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	ctx := context.Background()
	//初始化 Redis 客户端
	rs := sqlconfig.InitRedis()
	db := sqlconfig.MysqlConnect()
	userService := api.NewUserService(db)     // 创建用户服务
	friendService := api.NewFriendService(db) // 创建好友列表服务
	consentFriendService := api.ConsentService(db)
	getApplicationService := api.GetApplicationService(db)
	consentListService := api.ConsentListService(db)
	handleAddFriendService := api.HandleAddFriendService(db)
	searchService := api.SearchService(db)
	speakService := api.SpeakService(db)
	getIsSpeakUserInfoService := api.GetIsSpeakUserService(db)
	getMessageListService := api.GetMessageListService(db, rs)
	createRoomService := api.CreateRoomFactory(db)
	// 启动 HTTP 服务器

	app := &api.MyApp{
		Ctx:                    ctx,
		Db:                     db, // 将数据库连接赋值给 MyApp
		Rs:                     rs,
		UserService:            userService, // 将用户服务赋值给 MyApp
		FriendService:          friendService,
		ConsentService:         consentFriendService,
		GetApplicationService:  getApplicationService,
		ConsentListService:     consentListService,
		HandleAddFriendService: handleAddFriendService,
		SearchService:          searchService,
		SpeakService:           speakService,
		IsSpeakUserInfoService: getIsSpeakUserInfoService,
		GetMessageListService:  getMessageListService,
		CreateRoomService:      createRoomService,
	}

	go app.StartHTTPServer()

	err := wails.Run(&options.App{
		Title:  "",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Bind: []interface{}{
			app,
		},
		OnStartup: func(ctx context.Context) {
			runtime.EventsOn(ctx, "event_name", func(optionalData ...interface{}) {
				// 处理接收到的数据
				log.Println("Received event data:", optionalData)
			})
			runtime.WindowSetMinSize(ctx, 800, 600)
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

//
//package main
//
//import (
//	"fmt"
//	"sync"
//	"time"
//)
//
//// 全局变量
//var (
//	clients     = make(map[string]string)
//	clientsLock sync.RWMutex // 读写锁
//)
//
//func main() {
//	// 使用不加锁的方式
//	go withoutLock("user1", "Alice")
//	//go withoutLock("user2", "Bob")
//	//
//	//time.Sleep(1 * time.Second) // 等待 goroutine 完成
//
//	// 使用加锁的方式
//	//go withLock("user2", "Bob")
//
//	time.Sleep(1 * time.Second) // 等待 goroutine 完成
//
//	fmt.Println("Final clients map:", clients)
//}
//
//// 不加锁的写操作
//func withoutLock(friendID, name string) {
//	for i := 0; i < 10; i++ {
//		clients[friendID] = name
//	}
//}
//
//// 加锁的写操作
//func withLock(friendID, name string) {
//	for i := 0; i < 10; i++ {
//		clientsLock.Lock() // 加锁
//		clients[friendID] = name
//		clientsLock.Unlock() // 解锁
//	}
//}
