package api

import (
	"QQGUI/model"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

type RoomService struct {
	Db *gorm.DB
}

func CreateRoomFactory(db *gorm.DB) *RoomService {
	return &RoomService{Db: db}
}

func (cr *RoomService) CreateRoom(group map[string]map[string]string, UserID string, roomID string, userName string) (error, error, string) {
	roomName := ""
	//fmt.Printf("userName : %s", userName)
	for _, user := range group {
		roomName += user["userName"] + "--"
	}
	roomName = roomName + userName
	UserID1, _ := strconv.Atoi(UserID)

	//fmt.Println(UserID1)
	gp := model.Group{RoomID: roomID, Name: roomName, IsGroup: true, CreateUser: UserID1}
	num := cr.Db.Create(&gp)
	gm := model.GroupMember{RoomID: roomID, UserID: UserID1, CreateUser: UserID1}
	num1 := cr.Db.Create(&gm)
	if num.RowsAffected != 0 && num1.RowsAffected != 0 {
		return nil, nil, roomName
	} else {
		return num.Error, num1.Error, ""
	}
}

func (cr *RoomService) CreateRoomMember(roomID string, userID string) {
	UserID1, _ := strconv.Atoi(userID)
	gm := model.GroupMember{RoomID: roomID, UserID: UserID1}
	cr.Db.Create(&gm)
}

func (cr *RoomService) QueryRoomService(userID int) (*Response, error) {
	var result []struct {
		RoomID     string `json:"room_id"`
		UserID     int    `json:"user_id"`
		UserName   string `json:"user_name"`
		Avatar     string `json:"avatar"`
		Name       string `json:"name"`
		CreateUser int    `json:"create_user"`
	}

	// 执行查询，获取该用户所在房间的信息
	err := cr.Db.Raw("SELECT "+
		"gm.room_id, gm.user_id, ui1.username AS user_name, ui1.avatar AS avatar,"+
		" g.name AS name, g.create_user AS create_user FROM group_members gm "+
		"JOIN user_infos ui1 ON ui1.id = gm.user_id JOIN `groups` g "+
		"ON g.room_id = gm.room_id WHERE gm.room_id IN "+
		"(SELECT room_id FROM group_members WHERE user_id = ?)", userID).Scan(&result).Error

	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	//fmt.Printf("数据在这呢:%v", result)
	// 合并房间信息
	roomsMap := make(map[string]*Room)

	for _, item := range result {
		// 如果该房间还没有记录，创建一个新的
		if _, exists := roomsMap[item.RoomID]; !exists {
			roomsMap[item.RoomID] = &Room{
				RoomID:     item.RoomID,
				RoomName:   item.Name,
				CreateUser: item.CreateUser,
				UserInfo:   []UserInfo{},
			}
		}

		// 将用户信息添加到该房间的 userInfo 中
		roomsMap[item.RoomID].UserInfo = append(roomsMap[item.RoomID].UserInfo, UserInfo{
			UserName: item.UserName,
			UserID:   item.UserID,
			Avatar:   item.Avatar,
		})
	}
	//fmt.Printf("shuju:%v", roomsMap)
	// 将 map 转换成 slice
	var rooms []Room
	for _, room := range roomsMap {
		rooms = append(rooms, *room)
	}

	// 返回结果
	response := &Response{
		Room: rooms,
	}

	return response, nil
}

func (cr *RoomService) DeleteRoomMemberService() {

}
func (cr *RoomService) EditRoomInfoService() {

}
