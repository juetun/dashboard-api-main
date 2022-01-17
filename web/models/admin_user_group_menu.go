// Package models /**
package models

import (
	"time"
)

const (
	PathTypePage = "page"
	PathTypeApi  = "api"

	SetPermitAdd    = "add"
	SetPermitCancel = "cancel"
)

type AdminUserGroupMenu struct {
	Id            int64      `gorm:"column:id;primary_key" json:"id"`
	GroupId       int64      `json:"group_id" gorm:"column:group_id;uniqueIndex:idx_uid,priority:1;not null;default:0;comment:管理员组ID"`
	Module        string     `json:"module"  gorm:"column:module;uniqueIndex:idx_uid,priority:2;not null;default:'';comment:菜单所属系统"`
	MenuId        int64      `json:"menu_id" gorm:"column:menu_id;uniqueIndex:idx_uid,priority:3;not null;default:0;comment:菜单ID"`
	MenuPermitKey string     `json:"menu_permit_key" gorm:"column:menu_permit_key;not null;default:'';type:varchar(100) COLLATE utf8mb4_bin;comment:菜单唯一Key" form:"menu_permit_key"`
	CreatedAt     time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
	UpdatedAt     time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (r *AdminUserGroupMenu) GetTableComment() (res string) {
	return "用户组所具备的权限"
}

func (r *AdminUserGroupMenu) TableName() string {
	return "admin_user_group_menu"
}
