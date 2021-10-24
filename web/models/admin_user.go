// Package models
// /**
package models

import (
	"github.com/juetun/base-wrapper/lib/base"
)

type AdminUser struct {
	ID        uint             `gorm:"column:id;primary_key" json:"id"`
	UserHid   string           `json:"user_hid" gorm:"column:user_hid"`
	RealName  string           `json:"real_name" gorm:"column:real_name"`
	Mobile    string           `json:"mobile" gorm:"column:mobile"`
	FlagAdmin int              `json:"flag_admin" gorm:"column:flag_admin"`
	CreatedAt base.TimeNormal  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt base.TimeNormal  `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt *base.TimeNormal `json:"deleted_at" gorm:"column:deleted_at"`
}

func (r *AdminUser) GetTableComment() (res string) {
	res = "管理员库"
	return
}

func (r *AdminUser) TableName() string {
	return "admin_user"
}
