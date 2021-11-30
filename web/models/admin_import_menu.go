// Package models /**
package models

import (
	"time"
)

type AdminMenuImport struct {
	Id            int64      `gorm:"column:id;primary_key" json:"id"`
	MenuId        int64      `json:"menu_id" gorm:"column:menu_id"`
	MenuModule    string     `json:"menu_module" gorm:"column:menu_module"`
	ImportAppName string     `json:"import_app_name" gorm:"column:import_app_name"`
	ImportId      int64      `json:"import_id" gorm:"column:import_id"`
	DefaultOpen   uint8      `json:"default_open" gorm:"column:default_open"` // 默认开启菜单就开启本接口权限 1-开启 不为1-不开启（ 默认值2不开启）
	CreatedAt     time.Time  `json:"-" gorm:"column:created_at"`
	UpdatedAt     time.Time  `json:"-" gorm:"column:updated_at"`
	DeletedAt     *time.Time `json:"-" gorm:"column:deleted_at"`
}

func (r *AdminMenuImport) GetTableComment() (res string) {
	return "菜单所有的接口"
}

func (r *AdminMenuImport) TableName() string {
	return "admin_menu_import"
}
