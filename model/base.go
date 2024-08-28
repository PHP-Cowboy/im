package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

const (
	PageSize = 10
)

type Base struct {
	Id        int       `gorm:"primaryKey;type:int(11) unsigned AUTO_INCREMENT;comment:id"`
	CreatedAt time.Time `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;comment:更新时间"`
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 9999:
			pageSize = 9999
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

type GormList []string

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), g)
}
