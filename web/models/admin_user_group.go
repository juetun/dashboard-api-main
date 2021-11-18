// Package models /**
package models

import (
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
	GroupId      int64           `json:"group_id" gorm:"column:group_id"`
	UserHid      string          `json:"user_hid"  gorm:"column:user_hid"`
	IsSuperAdmin uint8           `json:"is_super_admin" gorm:"column:is_super_admin"`
	IsAdminGroup uint8           `json:"is_admin_group" gorm:"column:is_admin_group"`
	CreatedAt    base.TimeNormal `json:"-" gorm:"column:created_at" `
	UpdatedAt    base.TimeNormal `json:"updated_at" gorm:"column:updated_at" `
	DeletedAt    *time.Time      `json:"-" gorm:"column:deleted_at" `
}

func (r *AdminUserGroup) TableName() string {
	return "admin_user_group"
}
