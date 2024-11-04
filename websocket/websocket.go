package websocket

import (
	"QQGUI/config/sqlconfig"
	"QQGUI/model"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
	"sync"
)

type websocketInfo struct {
	WS     *websocket.Conn
	UserID string
}

// 定义全局变量
var (
	clients []websocketInfo
	mu      sync.Mutex // 用于保护对 clients 的并发访问
)

// 定义 WebSocket 升级器，并自定义 CheckOrigin
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有来源，防止跨域问题
		return true
	},
}

// GenerateRoomID 生成房间 ID，通过字典顺序比较两个 UUID
func GenerateRoomID(uuid1, uuid2 string) string {
	// 对两个 UUID 进行字典顺序比较
	if uuid1 < uuid2 {
		return uuid1 + "_" + uuid2 // 返回 UUID1 + UUID2 形式的房间 ID
	}
	return uuid2 + "_" + uuid1 // 返回 UUID2 + UUID1 形式的房间 ID
}

// HandleWebSocket 处理 WebSocket 连接
func HandleWebSocket(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// 获取请求中的查询参数
	var wsinfo websocketInfo
	_ = r.URL.Query().Get("user_id")            // 获取 user_id
	useruuid := r.URL.Query().Get("userIDUUID") // 获取用户 UUID
	//frienduuid := r.URL.Query().Get("friendIDUUID") // 获取好友 UUID
	IsGroup := r.URL.Query().Get("IsGroup") // 获取是否为群聊的标识
	roomid := r.URL.Query().Get("roomid")   // 获取房间 ID
	var channel string                      // 定义频道变量
	// 根据 IsGroup 参数决定频道类型
	if IsGroup == "1" {
		channel = "group" + "1" // 设置为群聊频道
	} else if IsGroup == "0" {
		channel = "user:" + roomid // 设置为用户频道
	}
	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)                               // 打印错误信息
		http.Error(w, "无法升级连接", http.StatusBadRequest) // 返回错误响应
		return
	}
	defer conn.Close() // 函数退出时关闭连接
	wsinfo.UserID = useruuid
	wsinfo.WS = conn
	mu.Lock() // 锁定以保护 clients
	//clients[conn] = true // 将新连接加入到客户端列表
	clients = append(clients, wsinfo)
	mu.Unlock() // 解锁

	pubsub := sqlconfig.SubscribeMessage(channel) // 订阅消息频道
	defer pubsub.Close()                          // 函数退出时关闭订阅

	// 启动一个 Goroutine 监听频道消息
	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(context.Background()) // 接收消息
			if err != nil {
				fmt.Println(err)
				return // 出现错误则退出
			}
			sendToClients(msg.Payload, useruuid, channel) // 发送消息给所有连接的客户端
		}
	}()

	// 处理客户端发送的消息
	for {
		var message model.Message      // 定义消息结构体
		err := conn.ReadJSON(&message) // 读取客户端发送的 JSON 消息
		if err != nil {
			fmt.Println(err) // 打印错误信息
			mu.Lock()        // 锁定以保护 clients
			// 从切片中删除关闭的客户端
			for i, client := range clients {
				if client.WS == conn {
					clients = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			mu.Unlock() // 解锁
			return      // 退出循环
		}
		// 将消息保存到数据库
		result := db.Create(&message)
		if result.Error != nil {
			fmt.Println(err) // 打印保存失败的错误信息
			continue         // 跳过当前循环
		}
		// 发布消息到 Redis
		sqlconfig.PublishMessage(channel, message)
	}
}

// sendToClients 将消息发送给所有连接的客户端
func sendToClients(message string, senderID string, channel string) {
	mu.Lock()         // 锁定以保护 clients
	defer mu.Unlock() // 函数退出时解锁

	for _, client := range clients {
		if client.UserID == senderID {
			continue // 跳过发送者
		}
		if err := client.WS.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			err := client.WS.Close()
			if err != nil {
				return
			}
			// 从切片中删除失效的客户端
			for i, c := range clients {
				if c.WS == client.WS {
					clients = append(clients[:i], clients[i+1:]...)
					break
				}
			}
		}
	}
}
