package handler

import (
	"im/global"
	"im/ws"
	"time"
)

type PrivateChatReq struct {
	ToUserId int    `json:"toUserId"` //消息接收方
	Msg      string `json:"msg"`
}

type PrivateChatRsp struct {
	UserId        int    `json:"userId"`
	UserName      string `json:"userName"`
	UserAvatarUrl string `json:"userAvatarUrl"`
	Msg           string `json:"msg"`
	MsgId         int64  `json:"msgId"`
}

func PrivateChat() ws.HandleFunc {
	return func(s *ws.Server, conn *ws.Conn, msg *ws.Message) {
		var (
			req PrivateChatReq
		)

		userId := conn.GetUserId()

		u, err := GetUserInfoById(userId)
		if err != nil {
			global.Logger["err"].Errorf("GetUserInfoById failed,err:%v", err.Error())
			return
		}

		rsp := PrivateChatRsp{
			UserId:        userId,
			UserName:      u.UserName,
			UserAvatarUrl: u.AvatarUrl,
			Msg:           req.Msg,
			MsgId:         time.Now().UnixNano(),
		}

		userConnMp, ok := s.UserToConnMp[req.ToUserId]

		if !ok {
			//存入离线消息

			return
		}

		err = s.Send(ws.NewMessage(userId, req.ToUserId, rsp), userConnMp)

		if err != nil {
			global.Logger["err"].Error(err.Error())
			return
		}
	}
}
