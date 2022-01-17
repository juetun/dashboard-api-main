// Package models /**
package models

import (
	"time"
)

type AdminMenuImport struct {
	Id              int64  `gorm:"column:id;primary_key" json:"id"`
	MenuId          int64  `json:"menu_id" gorm:"column:menu_id;type:bigint(20);uniqueIndex:idx_uid,priority:1;not null;default:0;comment:菜单ID"`
	MenuModule      string `json:"menu_module" gorm:"column:menu_module;not null;default:'';type:varchar(100) COLLATE utf8mb4_bin;comment:菜单所属模块"`
	ImportAppName   string `json:"import_app_name" gorm:"column:import_app_name;not null;default:'';type:varchar(100) COLLATE utf8mb4_bin;comment:接口应用名"`
	ImportId        int64  `json:"import_id" gorm:"column:import_id;type:bigint(20);uniqueIndex:idx_uid,priority:2;not null;default:0;comment:接口ID"`
	ImportPermitKey string `json:"import_permit_key" gorm:"column:import_permit_key;not null;default:'';type:varchar(100) COLLATE utf8mb4_bin;comment:接口唯一Key" form:"import_permit_key"`
	DefaultOpen uint8      `json:"default_open" gorm:"column:default_open;type:TINYINT(2);not null;default:0;comment:默认开启菜单就开启本接口权限 1-开启 不为1-不开启（ 默认值2不开启）"` // 默认开启菜单就开启本接口权限 1-开启 不为1-不开启（ 默认值2不开启）
	CreatedAt   time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
	UpdatedAt   time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (r *AdminMenuImport) GetTableComment() (res string) {
	return "菜单所有的接口"
}

func (r *AdminMenuImport) TableName() string {
	return "admin_menu_import"
}
