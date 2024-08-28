package user

import (
	"gorm.io/gorm"
	"im/model"
	"time"
)

type User struct {
	ID                int        `gorm:"column:id;primaryKey;autoIncrement"`
	CellPhone         string     `gorm:"column:cell_phone;type:varchar(16);default:NULL"`
	UserName          string     `gorm:"column:user_name;type:varchar(16);default:NULL"`
	Password          string     `gorm:"column:password;type:varchar(255);default:NULL"`
	Rich              float64    `gorm:"column:rich;type:decimal(16,2);default:0.00"`
	Point             float64    `gorm:"column:point;type:decimal(16,2);default:0.00"`
	NickName          string     `gorm:"column:nick_name;type:varchar(68);default:NULL"`
	BirthDay          string     `gorm:"column:birth_day;type:varchar(50);default:''"`
	AvatarUrl         string     `gorm:"column:avatar_url;type:varchar(1000);default:NULL"`
	FacebookID        string     `gorm:"column:face_book_id;type:varchar(68);default:NULL"`
	FacebookName      string     `gorm:"column:face_book_name;type:varchar(68);default:NULL"`
	GoogleID          string     `gorm:"column:google_id;type:varchar(68);default:NULL"`
	GoogleName        string     `gorm:"column:google_name;type:varchar(68);default:NULL"`
	Gender            int8       `gorm:"column:gender;default:1"`
	InviteCode        string     `gorm:"column:invite_code;type:varchar(32);default:''"`
	Country           string     `gorm:"column:country;type:varchar(32);default:NULL"`
	Province          string     `gorm:"column:province;type:varchar(32);default:NULL"`
	City              string     `gorm:"column:city;type:varchar(32);default:NULL"`
	Channel           string     `gorm:"column:channel;type:varchar(100);default:NULL"`
	IsAgent           int8       `gorm:"column:is_agent;default:0"`
	IsTest            int8       `gorm:"column:is_test;default:0"`
	IsPass            int8       `gorm:"column:is_pass;default:0"`
	PUserID           int        `gorm:"column:p_user_id;default:NULL"`
	DelFlag           int8       `gorm:"column:del_flag;default:1"`
	ActiveTime        *time.Time `gorm:"column:active_time;default:NULL"`
	Env               string     `gorm:"column:env;type:varchar(50);default:NULL"`
	Oth               string     `gorm:"column:oth;type:varchar(4000);default:NULL"`
	PackageName       string     `gorm:"column:package_name;type:varchar(100);default:NULL"`
	IP                string     `gorm:"column:ip;type:varchar(50);default:''"`
	CountryAbbr       string     `gorm:"column:country_abbr;type:varchar(16);default:''"`
	ChannelSign       string     `gorm:"column:channel_sign;type:varchar(100);default:''"`
	PackageNameActive string     `gorm:"column:package_name_active;type:varchar(255);default:''"`
	FansNum           int        `gorm:"column:fans_num;default:0"`
	FollowNum         int        `gorm:"column:follow_num;default:0"`
	CharmNum          int        `gorm:"column:charm_num;default:0"`
	WealthNum         int        `gorm:"column:wealth_num;default:0"`
	PlayTime          float64    `gorm:"column:play_time;type:decimal(10,2);unsigned;default:0.00"`
	IsBan             int8       `gorm:"column:is_ban;default:0"`
	IsHot             int8       `gorm:"column:is_hot;default:0"`
	IsRecommend       int8       `gorm:"column:is_recommend;default:0"`
	Intro             string     `gorm:"column:intro;type:varchar(255);default:''"`
	VipLevel          int        `gorm:"column:vip_level;default:1"`
	VipExpire         string     `gorm:"column:vip_expire;type:varchar(32);default:''"`
	GoodNumber        string     `gorm:"column:good_number;type:varchar(32);default:''"`
	MyMount           string     `gorm:"column:my_mount;type:varchar(32);default:''"`
	Consumption       float64    `gorm:"column:consumption;type:decimal(16,2);unsigned;default:0.00"`
	Votes             int        `gorm:"column:votes;unsigned;default:0"`
	VotesSub          int        `gorm:"column:votes_sub;unsigned;default:0"`
	AuthorLevel       int        `gorm:"column:author_level;default:1"`
	AuthorExp         int        `gorm:"column:author_exp;default:0"`
	IsAuthor          int8       `gorm:"column:is_author;default:0"`
	UserLevel         int        `gorm:"column:user_level;default:1"`
	UserExp           int        `gorm:"column:user_exp;default:0"`
	SmallGiftEffects  int8       `gorm:"column:small_gift_effects;default:1"`
	AllGiftEffects    int8       `gorm:"column:all_gift_effects;default:1"`
	Cgga              int8       `gorm:"column:cgga;default:1"`
	ClanID            int        `gorm:"column:clan_id;default:0"`
	Percent           float64    `gorm:"column:percent;type:decimal(10,2);default:50.00"`
}

func (t *User) GetOneById(db *gorm.DB, id int) (data User, err error) {
	err = db.Model(t).First(&data, &id).Error

	return
}

func (t *User) GetUserListByIds(db *gorm.DB, ids []int) (dataList []*User, err error) {
	err = db.Model(t).Where("id in (?)", ids).Find(&dataList).Error

	return
}

func (t *User) GetUserPageListByIds(db *gorm.DB, ids []int) (dataList []*User, err error) {
	err = db.Model(t).Where("id in (?)", ids).Scopes(model.Paginate(1, model.PageSize)).Find(&dataList).Error

	return
}
