package routers

//
//import (
//	"QQGUI/server/api"
//	"QQGUI/utils/createKey"
//	"QQGUI/websocket"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"gorm.io/gorm"
//	"net/http"
//)
//
//func HandlerRun(db *gorm.DB) {
//	r := gin.Default()
//	r.Static("/assets", "./assets")
//	r.StaticFS("/static", http.Dir("./static/avatar")) // 将路由/static 映射到./static/avatar上
//	r.StaticFile("/favicon.ico", "./favicon.ico")
//	r.POST("/login", func(c *gin.Context) {
//		privateKey, _ := createKey.LoadPrivateKey("private.pem")
//		api.LoginFunc(c, db, privateKey)
//	})
//	r.POST("/register", func(c *gin.Context) {
//		privateKey, _ := createKey.LoadPrivateKey("private.pem")
//		api.RegisterFunc(c, db, privateKey)
//	})
//	r.GET("/ws", func(c *gin.Context) {
//		websocket.HandleWebSocket(c, db)
//	}) // WebSocket 连接路由
//	r.POST("/send", func(c *gin.Context) {
//		api.SendMessage(c, db)
//	}) // 发送消息路由
//	r.POST("/upload", func(c *gin.Context) {
//		api.UploadFile(c, db)
//	}) // 上传文件路由
//	r.GET("/FriendList", func(c *gin.Context) {
//		api.FriendList(c, db)
//	})
//	r.POST("/CreateRoom", func(c *gin.Context) {
//		api.CreateRoom(c, db)
//	})
//	r.POST("/addfriend", func(c *gin.Context) {
//		api.HandleAddFriendRequest(c, db)
//	})
//	r.GET("/sse", func(c *gin.Context) {
//		api.HandleSSE(c)
//	})
//	r.GET("/searchuser", func(c *gin.Context) {
//		api.SearchUser(c, db)
//	})
//	r.PUT("/consentFriend", func(c *gin.Context) {
//		api.ConsentFriend(c, db)
//	})
//	r.GET("/getConsentList", func(c *gin.Context) {
//		api.GetConsentList(c, db)
//	})
//	r.GET("/GetApplicationRecord", func(c *gin.Context) {
//		api.GetApplicationRecord(c, db)
//	})
//	err := r.Run(":8000")
//	if err != nil {
//		fmt.Println(err)
//		return
//	} // 监听并在 0.0.0.0:8080 上启动服务
//}
