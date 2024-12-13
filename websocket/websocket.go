package websocket

import (
	"QQGUI/config/sqlconfig"
	"QQGUI/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"sync"
)

/* 整个流程 HandleWebSocket(创建连接) -> handleClientMessages(从客户端读取消息，如果读取到了就直接广播到redis)
-> listenForMessages(监听到了有消息发不到了redis就通过sendToRoom) -> sendToRoom(发送到客户端)
*/

// GroupWebSocket 结构体，用于存储特定房间的 WebSocket 连接信息。
// - WsRoomInfo: 存储房间中每个用户的 WebSocket 连接列表。
// - mu: 确保对 WsRoomInfo 的并发访问是安全的。
type GroupWebSocket struct {
	WsRoomInfo map[string][]*websocket.Conn // 用户ID与其对应WebSocket连接的映射。
	mu         sync.Mutex                   // 保护 WsRoomInfo 的互斥锁。
}

type DebugGroupWebSocket struct {
	WsRoomInfo map[string][]string `json:"WsRoomInfo"`
}

// 全局变量
var (
	clients = make(map[string]*GroupWebSocket) // 存储所有房间的 WebSocket 信息。
	mu      sync.Mutex                         // 保护 clients 的全局互斥锁。
)

// upgrader 配置 WebSocket 的升级器。
// - CheckOrigin: 允许所有来源连接。
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// DebugClients 返回一个可序列化的调试版本的 clients。
func DebugClients() map[string]*DebugGroupWebSocket {
	mu.Lock()
	defer mu.Unlock()

	debugData := make(map[string]*DebugGroupWebSocket)
	for roomID, group := range clients {
		debugGroup := &DebugGroupWebSocket{WsRoomInfo: make(map[string][]string)}
		group.mu.Lock()
		for userID, conns := range group.WsRoomInfo {
			for _, conn := range conns {
				// 用 conn 的内存地址表示 WebSocket 连接
				debugGroup.WsRoomInfo[userID] = append(debugGroup.WsRoomInfo[userID], fmt.Sprintf("%p", conn))
			}
		}
		group.mu.Unlock()
		debugData[roomID] = debugGroup
	}
	return debugData
}
func PrintDebugClients() {
	data, err := json.MarshalIndent(DebugClients(), "", "  ")
	if err != nil {
		fmt.Println("JSON 序列化失败:", err)
		return
	}
	fmt.Println(string(data))
}

// HandleWebSocket 处理 WebSocket 连接。
// - 参数：
//   - w: HTTP 响应写入器。
//   - r: HTTP 请求。
//   - db: 数据库连接实例。
func HandleWebSocket(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// 从请求中获取用户和房间相关参数。
	userID := r.URL.Query().Get("userIDUUID")
	isGroup := r.URL.Query().Get("IsGroup")
	roomID := r.URL.Query().Get("roomid")

	// 参数验证。
	if userID == "" || isGroup == "" || roomID == "" {
		http.Error(w, "缺少必要的参数", http.StatusBadRequest)
		return
	}

	// 升级 HTTP 连接到 WebSocket。
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket 升级失败:", err)
		http.Error(w, "无法升级连接", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// 初始化或获取指定房间的 GroupWebSocket。
	mu.Lock()
	group, exists := clients[roomID]
	if !exists {
		group = &GroupWebSocket{WsRoomInfo: make(map[string][]*websocket.Conn)}
		clients[roomID] = group
	}
	mu.Unlock()

	// 将当前用户的连接加入房间。
	group.mu.Lock()
	group.WsRoomInfo[userID] = append(group.WsRoomInfo[userID], conn)
	group.mu.Unlock()
	PrintDebugClients()
	// 确定 Redis 频道名称。
	channel := ""
	if isGroup == "1" {
		channel = "group:" + userID + roomID
	} else {
		channel = "user:" + roomID + ":" + userID
	}

	// 订阅 Redis 消息频道。
	pubsub := sqlconfig.SubscribeMessage(channel)
	defer pubsub.Close()

	// 启动协程监听 Redis 频道的消息。
	go listenForMessages(pubsub, roomID)

	// 处理来自客户端的 WebSocket 消息。
	handleClientMessages(conn, db, channel, userID, roomID)

}

// listenForMessages 持续监听指定 Redis 频道的消息并广播到房间。
// - 参数：
//   - pubsub: Redis 订阅实例。
//   - roomID: 房间 ID。
func listenForMessages(pubsub *redis.PubSub, roomID string) {
	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			fmt.Println("订阅消息接收失败:", err)
			return
		}
		sendToRoom(msg.Payload, roomID)
	}
}

// handleClientMessages 处理来自 WebSocket 客户端的消息。
// - 参数：
//   - conn: 客户端的 WebSocket 连接。
//   - db: 数据库实例。
//   - channel: Redis 频道名称。
//   - userID: 用户 ID。
//   - roomID: 房间 ID。

func handleClientMessages(conn *websocket.Conn, db *gorm.DB, channel, userID, roomID string) {
	var message model.Messages1
	for {
		// 从客户端读取 JSON 格式的消息。
		err := conn.ReadJSON(&message)
		if err != nil {
			fmt.Println("读取消息失败:", err)
			removeClient(conn, roomID, userID)
			return
		}

		// 重置消息 ID，确保插入新记录。
		message.ID = 0
		if err := db.Create(&message.Message).Error; err != nil {
			fmt.Println("消息保存失败:", err)
			continue
		}

		// 将消息发布到 Redis 频道。
		sqlconfig.PublishMessage(channel, message)
	}
}

// sendToRoom 将消息广播到房间内的所有 WebSocket 连接。
// - 参数：
//   - message: 要发送的消息。
//   - roomID: 房间
//     .3+0ID。
func sendToRoom(message, roomID string) {
	mu.Lock()
	group, exists := clients[roomID]
	mu.Unlock()
	if !exists {
		return
	}

	group.mu.Lock()
	defer group.mu.Unlock()
	for _, conns := range group.WsRoomInfo {
		for _, conn := range conns {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				fmt.Println("发送消息失败，关闭连接:", err)
				_ = conn.Close()
			}
		}
	}
}

// removeClient 移除房间中指定用户的 WebSocket 连接。
// - 参数：
//   - conn: 要移除的连接。
//   - roomID: 房间 ID。
//   - userID: 用户 ID。
func removeClient(conn *websocket.Conn, roomID, userID string) {
	mu.Lock()
	group, exists := clients[roomID]
	mu.Unlock()
	if !exists {
		return
	}

	group.mu.Lock()
	defer group.mu.Unlock()

	// 移除指定用户的连接。
	if conns, ok := group.WsRoomInfo[userID]; ok {
		for i, c := range conns {
			if c == conn {
				group.WsRoomInfo[userID] = append(conns[:i], conns[i+1:]...)
				break
			}
		}
		if len(group.WsRoomInfo[userID]) == 0 {
			delete(group.WsRoomInfo, userID)
		}
	}

	// 如果房间中没有任何用户，删除房间。
	if len(group.WsRoomInfo) == 0 {
		mu.Lock()
		delete(clients, roomID)
		mu.Unlock()
	}
}
