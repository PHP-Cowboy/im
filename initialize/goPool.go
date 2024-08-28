package initialize

import (
	"github.com/panjf2000/ants"
	"im/global"
)

func InitGoPool() {

	pool, err := ants.NewPool(global.ServerConfig.PoolSize)
	if err != nil {
		global.Logger["err"].Errorf("ants.NewPool failed,err:%v", err.Error())
		panic(err)
		return
	}

	global.GoPool = pool
}
