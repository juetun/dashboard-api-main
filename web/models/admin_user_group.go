// Package models /**
package models

import (
	"fmt"
	"time"

	"github.com/juetun/base-wrapper/lib/base"
)

const (
	IsSuperAdminYes = 1 // 是超级管理员
	IsSuperAdminNo  = 0 // 不是超级管理员

	IsAdminGroupYes = 1 // 是管理员
	IsAdminGroupNo  = 0 // 不是管理员
)
const (
	GatewayNotHavePermit = iota + 11000
)

var (
	GatewayErrMap = map[int]string{
		GatewayNotHavePermit: "您不是系统管理员",
	}
)

type AdminUserGroup struct {
	Id           int             `json:"id" gorm:"column:id;primary_key" `
	GroupId      int64           `json:"group_id" gorm:"column:group_id;not null;uniqueIndex:idx_gid_hid,priority:1;default:0;comment:组ID"`
	UserHid      int64          `json:"user_hid"  gorm:"column:user_hid;not null;default:0;type:bigint(20);uniqueIndex:idx_gid_hid,priority:2;comment:用户ID"`
	IsSuperAdmin uint8           `json:"is_super_admin" gorm:"column:is_super_admin;not null;default:0;comment:是否是超级管理员 0-否,1-是"`
	IsAdminGroup uint8           `json:"is_admin_group" gorm:"column:is_admin_group;not null;default:0;comment:是否是后台管理员组 0-否,1-是"`
	CreatedAt    base.TimeNormal `json:"-" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" `
	UpdatedAt    base.TimeNormal `json:"updated_at" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" `
	DeletedAt    *time.Time      `json:"-" gorm:"column:deleted_at" `
}

func (r *AdminUserGroup) GetTableComment() (res string) {
	return "用户组用户列表"
}

func (r *AdminUserGroup) TableName() string {
 	return fmt.Sprintf("%suser_group", TablePrefix)
}
