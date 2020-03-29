package models

import (
	"time"

	"github.com/juetun/app-dashboard/lib/base"
)

type ZUsers struct {
	base.Model
	Name            string    `gorm:"column:name;" jjson:"name"`
	Email           string    `gorm:"column:email;" jjson:"email"`
	Status          int       `gorm:"column:status;" jjson:"status"`
	EmailVerifiedAt time.Time `gorm:"column:email_verified_at;" jjson:"email_verified_at"`
	Password        string    `gorm:"column:password;" jjson:"password"`
	RememberToken   string    `gorm:"column:remember_token;" jjson:"remember_token"`
}

func (r *ZUsers) TableName() string {
	return "z_users"
}
