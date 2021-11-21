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
	Id        int64      `gorm:"column:id;primary_key" json:"id"`
	GroupId   int64      `json:"group_id" gorm:"column:group_id"`
	Module    string     `json:"module"  gorm:"column:module"`
	MenuId    int64      `json:"menu_id" gorm:"column:menu_id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"-" `
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"-" `
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (r *AdminUserGroupMenu) GetTableComment() (res string) {
	return "用户组所具备的权限"
}

func (r *AdminUserGroupMenu) TableName() string {
	return "admin_user_group_menu"
}
