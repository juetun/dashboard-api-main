// Package models
// /**
package models

import (
	"github.com/juetun/base-wrapper/lib/base"
)

type AdminUser struct {
	ID        uint             `gorm:"column:id;primary_key" json:"id"`
	UserHid   string           `json:"user_hid" gorm:"column:user_hid;uniqueIndex:idx_user_hid,priority:1;not null;type:varchar(60) COLLATE utf8mb4_bin;default:'';comment:用户id"`
	RealName  string           `json:"real_name" gorm:"column:real_name;not null;type:varchar(30);default:'';comment:姓名"`
	Mobile    string           `json:"mobile" gorm:"column:mobile;not null;type:varchar(40);default:'';comment:电话号"`
	FlagAdmin uint8            `json:"flag_admin" gorm:"column:flag_admin;not null;type:tinyint(2);default:0;comment:是否可用 0-非客服后台管理员 1-客服后台管理员"`
	CanNotUse uint8            `json:"can_not_use" gorm:"column:can_not_use;not null;type:tinyint(2);default:0;comment:是否可用 0-可用 1-不可用"`
	CreatedAt base.TimeNormal  `json:"created_at" gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt base.TimeNormal  `json:"updated_at" gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt *base.TimeNormal `json:"deleted_at" gorm:"column:deleted_at;type:datetime;"`
}

func (r *AdminUser) GetTableComment() (res string) {
	res = "管理员表"
	return
}

func (r *AdminUser) TableName() string {
	return "admin_user"
}
