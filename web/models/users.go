package models

import (
	"time"

	"github.com/juetun/base-wrapper/lib/base"
)

type ZUsers struct {
	UserHid         string           `gorm:"column:user_hid;" json:"user_hid"`
	Name            string           `gorm:"column:name;" jjson:"name"`
	Portrait        string           `gorm:"column:portrait;" json:"portrait"`
	Email           string           `gorm:"column:email;" jjson:"email"`
	Status          int              `gorm:"column:status;" jjson:"status"`
	EmailVerifiedAt time.Time        `gorm:"column:email_verified_at;" jjson:"email_verified_at"`
	Password        string           `gorm:"column:password;" jjson:"password"`
	RememberToken   string           `gorm:"column:remember_token;" jjson:"remember_token"`
	CreatedAt       base.TimeNormal  `json:"created_at"`
	UpdatedAt       base.TimeNormal  `json:"updated_at"`
	DeletedAt       *base.TimeNormal `sql:"index" json:"-"`
}

func (r *ZUsers) TableName() string {
	return "z_users"
}
