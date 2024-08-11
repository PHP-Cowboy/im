package main

import (
	"im/handler"
	"im/ws"
)

func main() {
	s := ws.NewServer("127.0.0.1:1234")

	handler.RegisterHandlers(s)

	s.Start()

	//select {}
}
