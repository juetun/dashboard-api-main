// Package models
// /**
package models

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"time"
)

const (
	AdminUserCanNotUseYes uint8 = iota // 可用
	AdminUserCanNotUseNo               // 不可用
)

const (
	AdminUserFlagAdminNo  uint8 = iota //不是超级管理员
	AdminUserFlagAdminYes              //是超级管理员
)

var (
	SliceAdminUserCanNotUse = base.ModelItemOptions{
		{
			Value: AdminUserCanNotUseYes,
			Label: "可用",
		},
		{
			Value: AdminUserCanNotUseNo,
			Label: "不可用",
		},
	}
	SliceAdminUserFlagAdmin = base.ModelItemOptions{
		{
			Value: AdminUserFlagAdminNo,
			Label: "非超管",
		},
		{
			Value: AdminUserFlagAdminYes,
			Label: "超管",
		},
	}
)

type (
	AdminUser struct {
		ID        uint   `gorm:"column:id;primary_key" json:"id"`
		UserHid   int64  `json:"user_hid" gorm:"column:user_hid;uniqueIndex:idx_user_hid,priority:1;not null;type:bigint(20);default:0;comment:用户id"`
		RealName  string `json:"real_name" gorm:"column:real_name;not null;type:varchar(30);default:'';comment:姓名"`
		Mobile    string `json:"mobile" gorm:"column:mobile;not null;type:varchar(40);default:'';comment:电话号"`
		Email     string `json:"email" gorm:"column:email;not null;type:varchar(100);default:'';comment:邮箱"`
		FlagAdmin uint8  `json:"flag_admin" gorm:"column:flag_admin;not null;type:tinyint(2);default:0;comment:是否可用 0-非客服后台管理员 1-客服后台管理员"`
		CanNotUse uint8  `json:"can_not_use" gorm:"column:can_not_use;not null;type:tinyint(2);default:0;comment:是否可用 0-可用 1-不可用"`

		CreatedAt time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
		UpdatedAt time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
		DeletedAt *time.Time `gorm:"column:deleted_at" json:"-"`
	}

	AdminUserMap map[int64]*AdminUser
)

func (r *AdminUser) GetTableComment() (res string) {
	res = "管理员表"
	return
}
func (r *AdminUser) ParseFlagAdmin() (res string) {
	adminUserFlagAdmin, _ := SliceAdminUserFlagAdmin.GetMapAsKeyUint8()
	if tp, ok := adminUserFlagAdmin[r.FlagAdmin]; ok {
		res = tp
	}
	return
}

func (r *AdminUser) ParseCanNotUse() (res string) {
	AdminUserCanNotUseMap, _ := SliceAdminUserCanNotUse.GetMapAsKeyUint8()
	if tp, ok := AdminUserCanNotUseMap[r.CanNotUse]; ok {
		res = tp
	}
	return
}

func (r *AdminUser) TableName() string {
	return fmt.Sprintf("%suser", TablePrefix)
}
