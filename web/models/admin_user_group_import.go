// Package models /**
package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type AdminUserGroupImport struct {
	Id              int64      `gorm:"column:id;primary_key" json:"id"`
	GroupId         int64      `json:"group_id" gorm:"column:group_id;uniqueIndex:idx_uid,priority:1;type:bigint(20);not null;default:0;comment:管理员组ID"`
	AppName         string     `json:"app_name" gorm:"column:app_name;not null;default:'';type:varchar(100) COLLATE utf8mb4_bin;comment:应用名"`
	Module          string     `json:"module"  gorm:"column:module;not null;default:'';type:varchar(100) COLLATE utf8mb4_bin;comment:接口所属菜单所属系统"`
	MenuId          int64      `json:"menu_id" gorm:"column:menu_id;uniqueIndex:idx_uid,priority:2;type:bigint(20);not null;default:0;comment:接口所属菜单界面"`
	ImportPermitKey string     `json:"import_permit_key" gorm:"column:import_permit_key;not null;default:'';type:varchar(100) COLLATE utf8mb4_bin;comment:接口唯一Key" form:"import_permit_key"`
	DefaultOpen     uint8      `json:"default_open" gorm:"column:default_open;not null;type:TINYINT(2);default:0;comment:0-默认不随菜单开启 1-随菜单开启权限"`
	ImportId        int64      `json:"import_id" gorm:"column:import_id;uniqueIndex:idx_uid,priority:3;type:bigint(20);not null;default:0;comment:接口ID"`
	MenuPermitKey   string     `json:"menu_permit_key" gorm:"column:menu_permit_key;not null;default:'';type:varchar(100) COLLATE utf8mb4_bin;comment:菜单唯一Key" form:"menu_permit_key"`
	CreatedAt       time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
	UpdatedAt       time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
	DeletedAt       *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (r *AdminUserGroupImport) GetTableComment() (res string) {
	return "用户组所具备的权限"
}

func (r *AdminUserGroupImport) TableName() string {
	return fmt.Sprintf("%suser_group_import", TablePrefix)
}
func (r *AdminUserGroupImport) UnmarshalBinary(data []byte) (err error) {
	if r == nil {
		r = &AdminUserGroupImport{}
	}
	err = json.Unmarshal(data, r)
	return
}

func (r *AdminUserGroupImport) MarshalBinary() (data []byte, err error) {
	data, err = json.Marshal(r)
	return
}