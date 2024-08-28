package handler

import (
	"im/global"
	"im/ws"
)

type OtherToRoomReq struct {
	ToRoomId string `json:"toRoomId"`
}

func OtherToAll() ws.HandleFunc {
	return func(s *ws.Server, conn *ws.Conn, msg *ws.Message) {

		var req OtherToRoomReq

		err := ParseData(msg, &req)

		if err != nil {
			global.Logger["err"].Errorf("ParseData failed, err:%v", err.Error())
			return
		}

		err = s.Broadcast(ws.NewMessage(0, 0, msg))

		if err != nil {
			global.Logger["err"].Error(err.Error())
			return
		}
	}
}
