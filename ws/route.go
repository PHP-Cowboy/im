package ws

import "github.com/gorilla/websocket"

type Route struct {
	Method  string
	Handler HandleFunc
}

type HandleFunc func(s *Server, conn *websocket.Conn, msg *Message)
