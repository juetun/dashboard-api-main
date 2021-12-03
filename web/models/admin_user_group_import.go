// Package models /**
package models

import (
	"time"
)

type AdminUserGroupImport struct {
	Id          int64      `gorm:"column:id;primary_key" json:"id"`
	GroupId     int64      `json:"group_id" gorm:"column:group_id;uniqueIndex:idx_uid,priority:1;not null;default:0;comment:管理员组ID"`
	AppName     string     `json:"app_name" gorm:"column:app_name;not null;default:'';type:varchar(100) COLLATE utf8mb4_bin;comment:应用名"`
	Module      string     `json:"module"  gorm:"column:module;not null;default:'';type:varchar(100) COLLATE utf8mb4_bin;comment:接口所属菜单所属系统"`
	MenuId      int64      `json:"menu_id" gorm:"column:menu_id;uniqueIndex:idx_uid,priority:2;not null;default:0;comment:接口所属菜单界面"`
	DefaultOpen uint8      `json:"default_open" gorm:"column:default_open;not null;type:TINYINT(2);default:0;comment:0-默认不随菜单开启 1-随菜单开启权限"`
	ImportId    int64      `json:"import_id" gorm:"column:import_id;uniqueIndex:idx_uid,priority:3;not null;default:0;comment:接口ID"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
	UpdatedAt   time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (r *AdminUserGroupImport) GetTableComment() (res string) {
	return "用户组所具备的权限"
}

func (r *AdminUserGroupImport) TableName() string {
	return "admin_user_group_import"
}