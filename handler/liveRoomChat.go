package handler

import (
	"im/global"
	"im/ws"
	"time"
)

type LiveRoomChatReq struct {
	ToRoomId string `json:"toRoomId"` //发送消息房间
	Msg      string `json:"msg"`
}

type LiveRoomChatRsp struct {
	UserId        int    `json:"userId"`
	UserName      string `json:"userName"`
	UserAvatarUrl string `json:"userAvatarUrl"`
	Msg           string `json:"msg"`
	MsgId         int64  `json:"msgId"`
}

func LiveRoomChat() ws.HandleFunc {
	return func(s *ws.Server, conn *ws.Conn, msg *ws.Message) {
		var (
			req LiveRoomChatReq
		)

		userId := conn.GetUserId()

		u, err := GetUserInfoById(userId)
		if err != nil {
			global.Logger["err"].Errorf("GetUserInfoById failed,err:%v", err.Error())
			return
		}

		rsp := LiveRoomChatRsp{
			UserId:        userId,
			UserName:      u.UserName,
			UserAvatarUrl: u.AvatarUrl,
			Msg:           req.Msg,
			MsgId:         time.Now().UnixNano(),
		}

		err = s.SendToRoom(ws.NewMessage(userId, 0, rsp), req.ToRoomId)

		if err != nil {
			global.Logger["err"].Error(err.Error())
			return
		}

		//热度 聊天 + 1
	}
}
