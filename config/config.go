package config

type ServerConfig struct {
	IP        string      `json:"ip"`
	Port      int         `json:"port"`
	MysqlInfo MysqlConfig `json:"mysqlInfo"`
	RedisInfo RedisConfig `json:"redisInfo"`
	JwtInfo   JWTConfig   `json:"jwtInfo"`
	WsInfo    WsConfig    `json:"wsInfo"`
}

type MysqlConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Username string `json:"username"`
	Db       int    `json:"db"`
	Expire   int    `json:"expire"`
}

type JWTConfig struct {
	SigningKey  string `json:"key"`
	ExpiresHour int    `json:"expiresHour"`
	AddHour     int    `json:"addHour"`
}

type WsConfig struct {
	MaxConnectionIdle int `json:"maxConnectionIdle"`
}
