package handler

import (
	"im/global"
	"im/ws"
)

type SendGiftReq struct {
	RoomId string `json:"roomId"`
}

func SendGift() ws.HandleFunc {
	return func(s *ws.Server, conn *ws.Conn, msg *ws.Message) {
		var req SendGiftReq

		err := ParseData(msg, &req)

		if err != nil {
			global.Logger["err"].Errorf("ParseData failed, err:%v", err.Error())
			return
		}

		err = s.SendToRoom(ws.NewMessage(0, 0, msg), req.RoomId)

		if err != nil {
			global.Logger["err"].Error(err.Error())
			return
		}
	}
}
