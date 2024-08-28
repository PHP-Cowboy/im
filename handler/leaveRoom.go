package handler

import (
	"im/global"
	"im/model/user"
	"im/ws"
)

func LeaveRoom() ws.HandleFunc {
	return func(s *ws.Server, conn *ws.Conn, msg *ws.Message) {
		var (
			req        EnterOrLeaveRoomReq
			rsp        EnterOrLeaveRoomRsp
			viewData   ViewData
			list       = make([]User, 0)
			userIdList = make([]int, 0)
		)

		userId := conn.GetUserId()

		roomUser, ok := s.RoomUserConnMp[req.RoomId]

		if !ok || len(roomUser) <= 1 {
			//房间没人或者那个人是将要离开的自己，不发消息
			return
		}

		delete(s.RoomUserConnMp[req.RoomId], userId)
		delete(roomUser, userId)

		viewData.Count = len(roomUser)

		//有其他用户在房间
		for id, _ := range roomUser {
			userIdList = append(userIdList, id)
		}

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

		err = s.SendToRoom(ws.NewMessage(0, 0, rsp), req.RoomId)

		if err != nil {
			global.Logger["err"].Error(err.Error())
			return
		}
	}
}
