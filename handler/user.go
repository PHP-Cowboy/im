package handler

import (
	"github.com/gorilla/websocket"
	"im/ws"
	"log"
)

func Online() ws.HandleFunc {
	return func(s *ws.Server, conn *websocket.Conn, msg *ws.Message) {
		uids := s.GetUsers()
		u := s.GetUsers(conn)
		err := s.Send(ws.NewMessage(u[0], uids), conn)

		if err != nil {
			log.Println(err.Error())
		}
	}
}
