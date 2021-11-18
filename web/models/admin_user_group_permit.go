// Package models /**
package models

import (
	"time"
)

const (
	PathTypePage = "pages"
	PathTypeApi  = "api"

	SetPermitAdd    = "add"
	SetPermitCancel = "cancel"
)

type AdminUserGroupPermit struct {
	Id        int        `gorm:"column:id;primary_key" json:"id"`
	GroupId   int64      `json:"group_id" gorm:"column:group_id"`
	AppName   string     `json:"app_name" gorm:"column:app_name"`
	PathType  string     `json:"path_type" gorm:"column:path_type"`
	Module    string     `json:"module"  gorm:"column:module"`
	MenuId    int        `json:"menu_id" gorm:"column:menu_id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"-" `
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"-" `
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (r *AdminUserGroupPermit) TableName() string {
	return "admin_user_group_permit"
}
