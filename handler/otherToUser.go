package handler

import (
	"im/global"
	"im/ws"
)

type OtherToUserReq struct {
	ToUserId int `json:"toUserId"`
}

func OtherToUser() ws.HandleFunc {
	return func(s *ws.Server, conn *ws.Conn, msg *ws.Message) {

		var req OtherToUserReq

		err := ParseData(msg, &req)

		if err != nil {
			global.Logger["err"].Errorf("ParseData failed, err:%v", err.Error())
			return
		}

		userConn, ok := s.UserToConnMp[req.ToUserId]

		if !ok {
			global.Logger["err"].Errorf("用户不在线")
			return
		}

		err = s.Send(ws.NewMessage(0, 0, msg), userConn)

		if err != nil {
			global.Logger["err"].Error(err.Error())
			return
		}
	}
}
