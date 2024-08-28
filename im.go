package main

import (
	"fmt"
	"im/global"
	"im/handler"
	"im/initialize"
	"im/ws"
)

func main() {
	initialize.InitLogger()

	initialize.InitConfig()

	initialize.InitMysql()

	initialize.InitRedis()

	initialize.InitGoPool()

	info := global.ServerConfig

	s := ws.NewServer(fmt.Sprintf("%s:%d", info.IP, info.Port))

	handler.RegisterHandlers(s)

	s.Start()
}
