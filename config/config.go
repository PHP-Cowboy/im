package config

type ServerConfig struct {
	Port      int         `json:"port"`
	MysqlInfo MysqlConfig `json:"mysqlInfo"`
	JwtInfo   JWTConfig   `json:"jwtInfo"`
}

type MysqlConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type JWTConfig struct {
	SigningKey  string `json:"key"`
	ExpiresHour int    `json:"expiresHour"`
	AddHour     int    `json:"addHour"`
}
