package main

import (
	"im/handler"
	"im/initialize"
	"im/ws"
)

func main() {
	initialize.InitLogger()

	initialize.InitConfig()

	//initialize.InitMysql()

	s := ws.NewServer("127.0.0.1:1234")

	handler.RegisterHandlers(s)

	s.Start()

	//select {}
}
