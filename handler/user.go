package handler

import (
	"im/global"
	"im/ws"
)

func Online() ws.HandleFunc {
	return func(s *ws.Server, conn *ws.Conn, msg *ws.Message) {
		uids := s.GetUsers()
		u := s.GetUsers(conn)
		err := s.Send(ws.NewMessage(u[0], uids), conn)

		if err != nil {
			global.Logger["err"].Error(err.Error())
			return
		}
	}
}
