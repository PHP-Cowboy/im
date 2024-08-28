package handler

import (
	"im/global"
	"im/ws"
)

func Login() ws.HandleFunc {
	return func(s *ws.Server, conn *ws.Conn, msg *ws.Message) {
		//userId := conn.GetUserId()
		err := s.Send(ws.NewMessage(0, 0, "ok"), conn)

		if err != nil {
			global.Logger["err"].Error(err.Error())
			return
		}
	}
}
