package handler

import (
	"im/global"
	"im/model/user"
	"im/ws"
)

type EnterOrLeaveRoomReq struct {
	RoomId string `json:"roomId"`
	//NickName string `json:"nickName"`
}

type EnterOrLeaveRoomRsp struct {
	FromClientUid  int       `json:"fromClientUid"`
	FromClientName int       `json:"fromClientName"`
	Code           string    `json:"code"`
	HotValue       int       `json:"HotValue"` //热度
	ViewData       ViewData  `json:"viewData"`
	UserLevel      UserLevel `json:"userLevel"`
}

type ViewData struct {
	Count int    `json:"count"`
	List  []User `json:"list"`
}

type UserLevel struct {
	VipLevel    int  `json:"vipLevel"`
	AuthorLevel int  `json:"authorLevel"`
	UserLevel   int  `json:"userLevel"`
	AuthorExp   int  `json:"authorExp"`
	UserExp     int  `json:"userExp"`
	IsAuthor    int8 `json:"isAuthor"`
}

type User struct {
	Uid       int    `json:"uid"`
	NickName  string `json:"nickName"`
	AvatarUrl string `json:"avatarUrl"`
	BirthDay  string `json:"birthDay"`
	UserLevel int    `json:"userLevel"`
	Gender    int8   `json:"gender"`
	Intro     string `json:"intro"`
}

// 进入直播间
func EnterRoom() ws.HandleFunc {
	return func(s *ws.Server, conn *ws.Conn, msg *ws.Message) {

		var (
			req        EnterOrLeaveRoomReq
			rsp        EnterOrLeaveRoomRsp
			viewData   ViewData
			list       = make([]User, 0)
			userIdList = make([]int, 0)
		)

		err := ParseData(msg, &req)

		if err != nil {
			global.Logger["err"].Errorf("ParseData failed, err:%v", err.Error())
			return
		}

		userId := conn.GetUserId()

		roomUser, ok := s.RoomUserConnMp[req.RoomId]

		if !ok {
			//构造 用户 连接 map
			roomUser = make(map[int]*ws.Conn)
		}

		roomUser[userId] = conn

		s.RoomUserConnMp[req.RoomId] = roomUser

		//最少有一个，是用户自己
		for id, _ := range roomUser {
			userIdList = append(userIdList, id)
		}

		viewData.Count = len(roomUser)

		userInfo, err := GetUserInfoById(userId)
		if err != nil {
			global.Logger["err"].Errorf("GetUserInfoById failed,err:%v", err.Error())
			return
		}

		level := UserLevel{
			VipLevel:    userInfo.VipLevel,
			AuthorLevel: userInfo.AuthorLevel,
			UserLevel:   userInfo.UserLevel,
			AuthorExp:   userInfo.AuthorExp,
			UserExp:     userInfo.UserExp,
			IsAuthor:    userInfo.IsAuthor,
		}

		rsp.UserLevel = level

		if len(userIdList) > 0 {
			var userList []*user.User

			userList, err = GetUserPageListByIds(userIdList)

			for _, u := range userList {
				list = append(list, User{
					Uid:       u.ID,
					NickName:  u.NickName,
					AvatarUrl: u.AvatarUrl,
					BirthDay:  u.BirthDay,
					UserLevel: u.UserLevel,
					Gender:    u.Gender,
					Intro:     u.Intro, //签名
				})
			}
		}

		viewData.List = list

		rsp.ViewData = viewData

		//发送到直播间
		err = s.SendToRoom(ws.NewMessage(0, 0, rsp), req.RoomId)

		if err != nil {
			global.Logger["err"].Errorf("s.Send failed,err:%v", err.Error())
			return
		}
	}
}

func GetUserInfoById(id int) (*user.User, error) {
	uObj := new(user.User)

	db := global.DB

	u, err := uObj.GetOneById(db, id)

	if err != nil {
		global.Logger["err"].Errorf("user.User.GetOneById() failed,err:[%v]", err.Error())
		return nil, err
	}

	return &u, nil

}

func GetUserPageListByIds(ids []int) ([]*user.User, error) {
	uObj := new(user.User)

	db := global.DB

	userList, err := uObj.GetUserPageListByIds(db, ids)

	if err != nil {
		global.Logger["err"].Errorf("user.User.GetOneById() failed,err:[%v]", err.Error())
		return nil, err
	}

	return userList, nil
}
