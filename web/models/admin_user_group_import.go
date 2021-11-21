// Package models /**
package models

import (
	"time"
)



type AdminUserGroupImport struct {
	Id        int64      `gorm:"column:id;primary_key" json:"id"`
	GroupId   int64      `json:"group_id" gorm:"column:group_id"`
	AppName   string     `json:"app_name" gorm:"column:app_name"`
 	Module    string     `json:"module"  gorm:"column:module"`
	MenuId    int64      `json:"menu_id" gorm:"column:menu_id"`
	ImportId  int64      `json:"import_id" gorm:"column:import_id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"-" `
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"-" `
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (r *AdminUserGroupImport) GetTableComment() (res string) {
	return "用户组所具备的权限"
}

func (r *AdminUserGroupImport) TableName() string {
	return "admin_user_group_import"
}
