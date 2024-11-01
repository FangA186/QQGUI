package api

import (
	"QQGUI/model"
	"QQGUI/utils/createKey"
	"QQGUI/websocket"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"sync"
)

type Info struct {
	Username  string `json:"username" ts_type:"string"`
	Password  string `json:"password" ts_type:"string"`
	Signature string `json:"signature" ts_type:"string"`
	Avatar    string `json:"avatar" ts_type:"string"`
	UUID      string `json:"uuid" ts_type:"string"`
	ID        string `json:"id" ts-type:"string"`
}

type UserSearchResult struct {
	model.UserInfo     // 嵌入原始 UserInfo 结构体
	IsFriend       int `json:"is_friend"` // 新增的字段
}

type ApplyForResponse struct {
	model.ApplyFor           // 申请记录
	ApplicantID       uint   `json:"applicant_id"`       // 申请者ID
	ApplicantUsername string `json:"applicant_username"` // 申请者用户名
	ApplicantAvatar   string `json:"applicant_avatar"`   // 申请者头像
	FriendID          uint   `json:"friend_id"`          // 被添加者ID
	FriendUsername    string `json:"friend_username"`    // 被添加者用户名
	FriendAvatar      string `json:"friend_avatar"`      // 被添加者头像
	FriendUUID        string `json:"friend_uuid"`        // 被添加者的UUID
}
type MyApp struct {
	Ctx                    context.Context
	Db                     *gorm.DB
	UserService            *UserService // 添加 UserService 字段
	FriendService          *FriendService
	ConsentService         *ConsentFriendService
	GetApplicationService  *GetApplicationRecordService
	ConsentListService     *GetConsentListService
	HandleAddFriendService *HandleAddFriendRequestService
	SearchService          *SearchUserService
	SpeakService           *IsSpeakService
	IsSpeakUserInfoService *GetIsSpeakUserInfoService
}

func NewMyApp(db *gorm.DB) *MyApp {
	return &MyApp{Db: db}
}

// // 客户端映射表，用于管理活跃的 SSE 连接
var clients = make(map[string]chan string)
var clientsLock sync.RWMutex

type LoginResponse struct {
	Message    string `json:"message"`
	UserID     uint   `json:"userID"`
	UUID       string `json:"UUID"`
	UserAvatar string `json:"useravatar"`
	Status     int    `json:"status"`
}

type IsSpeakInfo struct {
	model.IsSpeak
	UserName  string `json:"user_name"`
	Avatar    string `json:"avatar"`
	UserID    string `json:"user_id"`
	Signature string `json:"signature"`
}

// Login 登录方法
func (a *MyApp) Login(username, password string) LoginResponse {
	var de createKey.De
	privateKey, err := createKey.LoadPrivateKey("private.pem")
	if err != nil {
		return LoginResponse{Message: "加载私钥错误", Status: 500}
	}

	de = createKey.DecryptPassword(privateKey, password) // 解密密码
	if de.Err != nil {
		return LoginResponse{Message: "密码解密错误", Status: 400}
	}

	// 使用 UserService 进行用户查找和验证
	user, err := a.UserService.GetUserByUsername(username)
	if err != nil {
		return LoginResponse{Message: "无效的账户", Status: 400}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(de.Str)); err != nil {
		return LoginResponse{Message: "无效的密码", Status: 400}
	}

	return LoginResponse{
		Message:    "登录成功",
		UserID:     user.ID,
		UUID:       user.UUID,
		UserAvatar: user.Avatar,
		Status:     200,
	}
}

// RegisterFunc 注册
func (a *MyApp) RegisterFunc(username, pwd string) map[string]interface{} {
	var regInfo model.UserInfo
	var de createKey.De
	privateKey, _ := createKey.LoadPrivateKey("private.pem")
	de = createKey.DecryptPassword(privateKey, pwd)
	if de.Err != nil {
		return map[string]interface{}{
			"err": "密码解密失败",
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(de.Str), bcrypt.DefaultCost)
	if err != nil {
		return map[string]interface{}{
			"err": "密码错误",
		}
	}
	regInfo.Username = username
	regInfo.Password = string(hashedPassword)
	v4, _ := uuid.NewRandom()
	regInfo.UUID = v4.String()
	if err := a.UserService.CreateUser(&regInfo).Error; err != nil {
		return map[string]interface{}{
			"err": "用户已存在",
		}
	}

	return map[string]interface{}{
		"success": "用户创建成功",
	}
}

// SendMessage 处理发送消息请求
//func (a *MyApp) SendMessage(messageContent string) (string, error) {
//	// 保存消息到数据库
//	message := model.Message{}
//	if err := db.Create(&message).Error; err != nil {
//		return "", fmt.Errorf("保存消息失败: %v", err)
//	}
//
//	// 发布消息到 Redis
//	sqlconfig.PublishMessage("chat", message.Content)
//
//	return "消息已发送", nil
//}

// StartHTTPServer // UploadFile 上传文件
//
//	func (a *MyApp) UploadFile(file interface{}) (string, error) {
//		// 保存文件逻辑（例如保存到本地或云存储）
//		fileURL := "文件的URL" // 替换为实际的文件 URL
//
//		return fileURL, nil
//	}
//
// FriendList 获取好友列表
func (a *MyApp) StartHTTPServer() {
	http.HandleFunc("/sse", HandleSSE)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Jinlaile ")
		websocket.HandleWebSocket(w, r, a.Db)
	})
	// 提供静态文件服务
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/avatar"))))
	http.Handle("/favicon.ico", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
}
func (a *MyApp) FriendList(userID string) FriendRes {
	// 创建 FriendService 实例
	friendService := &FriendService{Db: a.Db}
	// 获取用户的好友
	friendslist, err := friendService.GetFriendList(userID)
	if err != nil {
		return FriendRes{
			InfoList: nil,
			Err:      err,
		}
	}
	// 获取其他用户的好友
	friendsOfOtherUsers, err := friendService.GetFriendsOfOtherUsers(userID)
	if err != nil {
		return FriendRes{
			InfoList: nil,
			Err:      err,
		}
	}

	// 合并好友列表
	infos := append(friendslist, friendsOfOtherUsers...)

	return FriendRes{InfoList: infos, Err: nil}
}

// HandleAddFriendRequest 处理好友请求
func (a *MyApp) HandleAddFriendRequest(userID, friendID string, isGroup bool) string {
	// 创建 FriendService 实例
	friendService := &HandleAddFriendRequestService{Db: a.Db}
	// 调用服务层处理好友请求
	message, err := friendService.HandleAddFriendService(userID, friendID, isGroup)
	if err != nil {
		return "处理好友请求时发生错误"
	}
	return message
}

// NotifyFriend 发送通知
func NotifyFriend(friendID, icon string, message string, friendInfo interface{}) {
	clientsLock.RLock()
	defer clientsLock.RUnlock()
	if messageChan, ok := clients[friendID]; ok {
		// 直接将对象发送到好友的 SSE 通道
		message1 := map[string]interface{}{
			"icon":    icon,
			"message": message,
			"friend":  friendInfo,
			"i":       0,
		}
		jsonStr, _ := json.Marshal(message1)
		messageChan <- string(jsonStr)
	} else {
		// 当未找到好友的 SSE 通道时执行的逻辑
		fmt.Printf("未找到好友的 SSE 通道，friendID: %s\n", friendID)
		// 这里可以选择记录日志、发送错误消息，或进行其他处理
	}
}

/*
Receive
用户A点击发送按钮，发送信息给用户B
就要在那个用户B那里判断一下如果他处于和用户A的聊天界面就不显示通知如果在就显示
*/
func (a *MyApp) Receive(userID, FriendID, Avatar, Username string) {
	clientsLock.RLock()
	defer clientsLock.RUnlock()
	if messageChan, ok := clients[FriendID]; ok {
		// 直接将对象发送到好友的 SSE 通道
		message1 := map[string]interface{}{
			//"icon":    icon,
			//"message": "",
			"userID":   userID,
			"friend":   FriendID,
			"avatar":   Avatar,
			"username": Username,
			"i":        1,
		}
		jsonStr, _ := json.Marshal(message1)
		messageChan <- string(jsonStr)
	} else {
		// 当未找到好友的 SSE 通道时执行的逻辑
		fmt.Printf("未找到好友的 SSE 通道，friendID: %s\n", FriendID)
		// 这里可以选择记录日志、发送错误消息，或进行其他处理
	}
}

// ConsentFriend 同意好友的申请
func (a *MyApp) ConsentFriend(userID, friendID, consent string) string {
	// 创建 FriendService 实例
	friendService := &ConsentFriendService{Db: a.Db}

	// 调用服务层处理好友申请
	message, err := friendService.ConsentFriend1(userID, friendID, consent)
	if err != nil {
		return "处理好友申请时发生错误"
	}

	return message
}

// HandleSSE 处理 SSE 连接
func HandleSSE(w http.ResponseWriter, r *http.Request) {
	// 获取当前用户的ID
	userID := r.URL.Query().Get("user_id")
	ctx := r.Context()

	if userID == "" {
		http.Error(w, "缺少用户ID", http.StatusBadRequest)
		return
	}

	// 确保连接支持流式响应
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "当前连接不支持流式传输", http.StatusInternalServerError)
		return
	}

	// 设置 HTTP 头信息，告知客户端这是 SSE 通道
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// 为当前用户创建一个消息通道
	messageChan := make(chan string)

	// 将此用户的 SSE 连接加入到映射表中
	clientsLock.Lock()
	clients[userID] = messageChan
	clientsLock.Unlock()
	// 退出函数时，移除该用户的 SSE 连接，防止内存泄漏
	defer func() {
		clientsLock.Lock()
		delete(clients, userID)
		clientsLock.Unlock()
		close(messageChan) // 关闭通道，避免泄漏
	}()

	// 启动一个 goroutine 来异步处理消息接收和推送
	go func() {
		for {
			select {
			case message := <-messageChan:
				// 当有消息时，将消息发送给客户端
				fmt.Fprintf(w, "data: %s\n\n", message)
				flusher.Flush() // 刷新缓冲区
			case <-r.Context().Done():
				return
			}
		}
	}()

	// 保持连接
	<-ctx.Done() // 等待客户端断开连接
}

// SearchUser 搜索用户并添加 is_friend 字段
func (a *MyApp) SearchUser(currentUserID, username string) []UserSearchResult {
	// 创建 UserService 实例
	userService := &UserService{Db: a.Db}

	// 调用服务层进行用户搜索
	results, err := userService.SearchUser(currentUserID, username)
	if err != nil {
		fmt.Println("搜索用户时发生错误:", err)
		return nil
	}

	fmt.Println("搜索结果:", results)
	return results
}

// GetConsentList 获取该用户被他人申请添加的记录
func (a *MyApp) GetConsentList(userID string) []ApplyForResponse {
	// 创建 ApplyForService 实例
	applyForService := &GetConsentListService{Db: a.Db}

	// 调用服务层获取申请记录
	applyForRequests, err := applyForService.GetConsentList(userID)
	if err != nil {
		fmt.Println("获取申请记录时发生错误:", err)
		return nil
	}

	fmt.Println("申请记录:", applyForRequests)
	return applyForRequests
}

// GetApplicationRecord 获取该用户的申请记录
func (a *MyApp) GetApplicationRecord(userID string) []ApplyForResponse {
	// 创建 ApplyForService 实例
	applyForService := &GetApplicationRecordService{Db: a.Db}
	// 调用服务层获取申请记录
	applyForRequests := applyForService.GetApplicationRecord(userID)
	return applyForRequests
}

// IsSpeak 将好友标记为对话状态
func (a *MyApp) IsSpeak(userID, friendID string) map[string]interface{} {
	/*
		return
		{
		message:"xxxxx",
		status:000000
		}

	*/
	speakService := &IsSpeakService{Db: a.Db}
	message, status := speakService.Speak(userID, friendID)
	return map[string]interface{}{
		"message": message,
		"status":  status,
	}
}

func (a *MyApp) GetIsSpeakUserInfo(userID string) []ApplyForResponse {
	fmt.Println("suerId", userID)
	getSpeak := &GetIsSpeakUserInfoService{Db: a.Db}
	res := getSpeak.GetIsSpeakUser(userID)
	return res
}
