package handler

import (
	"im/global"
	"im/ws"
)

// 广播消息，一般不是由客户端发起，应该是后端服务发起 目前是预定义状态，后期优化
func Broadcast() ws.HandleFunc {
	return func(s *ws.Server, conn *ws.Conn, msg *ws.Message) {

		err := s.Broadcast(ws.NewMessage(0, 0, ""))

		if err != nil {
			global.Logger["err"].Error(err.Error())
			return
		}
	}
}
