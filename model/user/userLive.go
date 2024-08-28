package user

import (
	"gorm.io/gorm"
	"time"
)

type UserLive struct {
	Uid         int       `gorm:"column:uid;primaryKey"`
	ShowId      string    `gorm:"column:show_id;type:varchar(50);default:''"`
	ClanId      int       `gorm:"column:clan_id;default:0"`
	IsLive      int       `gorm:"column:is_live;default:0"`
	StartTime   int       `gorm:"column:start_time;default:0"`
	Title       string    `gorm:"column:title;type:varchar(255);default:''"`
	Province    string    `gorm:"column:province;type:varchar(255);default:''"`
	City        string    `gorm:"column:city;type:varchar(255);default:''"`
	Stream      string    `gorm:"column:stream;type:varchar(255);default:''"`
	Thumb       string    `gorm:"column:thumb;type:varchar(255);default:''"`
	Pull        string    `gorm:"column:pull;type:varchar(255);default:''"`
	Push        string    `gorm:"column:push;type:varchar(255);default:''"`
	Lng         string    `gorm:"column:lng;type:varchar(32);default:''"`
	Lat         string    `gorm:"column:lat;type:varchar(32);default:''"`
	LiveClassid int       `gorm:"column:live_classid;default:0"`
	IsHot       bool      `gorm:"column:is_hot;default:0"`
	IsRecommend bool      `gorm:"column:is_recommend;default:0"`
	GameID      int       `gorm:"column:game_id;default:0"`
	CreateTime  time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP"`
	ModifyTime  time.Time `gorm:"column:modify_time;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP"`
	Votes       int       `gorm:"column:votes;default:0"`
	GameMark    string    `gorm:"column:game_mark;type:varchar(32);default:''"`
	CdnSwitch   int8      `gorm:"column:cdn_switch;default:0"`
	CloseReason string    `gorm:"column:close_reason;type:varchar(255);default:''"`
	Snapshot    string    `gorm:"column:snapshot;type:varchar(255);default:''"`
	Type        int8      `gorm:"column:type;default:1"`
	Tag         int       `gorm:"column:tag;default:1"`
	RoomSkin    string    `gorm:"column:room_skin;type:varchar(255);default:''"`
	RoomNotice  string    `gorm:"column:room_notice;type:varchar(255);default:''"`
}

func (t *UserLive) GetOneByShowId(db *gorm.DB, showId string) (data UserLive, err error) {
	err = db.Model(t).First(&data, &showId).Error

	return
}
