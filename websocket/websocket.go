package websocket

import (
	"QQGUI/config/sqlconfig"
	"QQGUI/model"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"sync"
)

type websocketInfo struct {
	WS       *websocket.Conn
	UserID   string
	RoomInfo map[string][]string
	mu       sync.Mutex // 添加互斥锁以保护 RoomInfo
}

// 定义全局变量
var (
	clients []*websocketInfo
	mu      sync.Mutex // 用于保护对 clients 的并发访问
)

// 定义 WebSocket 升级器，并自定义 CheckOrigin
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有来源，防止跨域问题
		return true
	},
}

// HandleWebSocket 处理 WebSocket 连接
func HandleWebSocket(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// 获取请求中的查询参数
	var wsinfo websocketInfo
	var message model.Message // 定义消息结构体
	//_ = r.URL.Query().Get("user_id")            // 获取 user_id
	useruuid := r.URL.Query().Get("userIDUUID") // 获取用户 UUID
	IsGroup := r.URL.Query().Get("IsGroup")     // 获取是否为群聊的标识
	roomid := r.URL.Query().Get("roomid")       // 获取房间 ID
	_ = r.URL.Query().Get("userid")
	var channel string // 定义频道变量
	// 根据 IsGroup 参数决定频道类型
	if IsGroup == "1" {
		channel = "group:" + useruuid + roomid // 设置为群聊频道
		//
		// 确保 wsinfo.RoomInfo 已初始化
		if wsinfo.RoomInfo == nil {
			wsinfo.RoomInfo = make(map[string][]string)
		}

		// 将用户 ID 添加到指定房间
		wsinfo.mu.Lock()         // 加锁保护 RoomInfo
		defer wsinfo.mu.Unlock() // 确保函数退出时解锁

		// 如果房间不存在，则初始化
		if _, exists := wsinfo.RoomInfo[roomid]; !exists {
			wsinfo.RoomInfo[roomid] = []string{}
		}

		// 检查用户是否已存在于房间，避免重复添加
		userExists := false
		for _, id := range wsinfo.RoomInfo[roomid] {
			if id == useruuid {
				userExists = true
				break
			}
		}

		if !userExists {
			wsinfo.RoomInfo[roomid] = append(wsinfo.RoomInfo[roomid], useruuid) // 添加用户到房间
			fmt.Printf("用户 %s 加入房间 %s\n", useruuid, roomid)
		} else {
			fmt.Printf("用户 %s 已在房间 %s 中\n", useruuid, roomid)
		}
		fmt.Println(wsinfo.RoomInfo[roomid])
	} else if IsGroup == "0" {
		channel = "user:" + roomid + ":" + useruuid // 设置为用户频道
		wsinfo.UserID = useruuid

	}
	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)                               // 打印错误信息
		http.Error(w, "无法升级连接", http.StatusBadRequest) // 返回错误响应
		return
	}
	defer conn.Close() // 函数退出时关闭连接

	wsinfo.WS = conn
	mu.Lock() // 锁定以保护 clients
	//clients[conn] = true // 将新连接加入到客户端列表
	clients = append(clients, &wsinfo)
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
			sendToClients(msg.Payload, roomid, IsGroup) // 发送消息到该房间内的用户
		}
	}()
	// 处理客户端发送的消息
	for {
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

		message.ID = 0
		// 将消息保存到数据库
		mu.Lock()
		result := db.Create(&message)
		mu.Unlock()
		if result.Error != nil {
			fmt.Println(result.Error) // 打印保存失败的错误信息
			continue                  // 跳过当前循环
		}
		// 发布消息到 Redis
		sqlconfig.PublishMessage(channel, message)
	}
}

// sendToClients 将消息发送给所有连接的客户端
func sendToClients(message string, roomID string, IsGroup string) {
	fmt.Println("执行了n次")
	mu.Lock() // 锁定以保护 clients
	defer mu.Unlock()
	if IsGroup == "0" {
		// 解析房间内的用户
		userIDs := ParseRoomID(roomID)
		// 遍历 clients，仅向属于该房间的用户发送消息
		for _, client := range clients {
			for _, userID := range userIDs {
				if client.UserID == userID { // 判断是否是房间内的用户
					if err := client.WS.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
						fmt.Printf("发送失败，关闭连接 (用户: %s): %v\n", client.UserID, err)
						_ = client.WS.Close()
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
		}
	} else {
	loop:
		for _, c := range clients {
			for key, _ := range c.RoomInfo {
				if key == roomID {
					for _, c2 := range clients {
						err := c2.WS.WriteMessage(websocket.TextMessage, []byte(message))
						if err != nil {
							fmt.Printf("发送失败，关闭连接 (用户: %s): %v\n", c2.UserID, err)
							_ = c2.WS.Close()
							// 从切片中删除失效的客户端
							for i, c := range clients {
								if c.WS == c2.WS {
									clients = append(clients[:i], clients[i+1:]...)
									break
								}
							}
						}
					}
					break loop
				}
			}
		}
	}

}

func ParseRoomID(roomid string) []string {
	return strings.Split(roomid, "_") // 根据 "_" 分割返回用户 ID 列表
}
