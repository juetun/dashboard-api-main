/**
* @Author:changjiang
* @Description:
* @File:admin_import
* @Version: 1.0.0
* @Date 2020/12/9 11:08 下午
 */
package models

import (
	"time"
)

type AdminMenuImport struct {
	Id        int        `gorm:"primary_key" json:"id" form:"id"`
	MenuId    int        `json:"menu_id" gorm :"column:menu_id"`
	ImportId  int        `json:"import_id" gorm :"column:import_id"`
	CreatedAt time.Time  `json:"-" gorm :"column:created_at"`
	UpdatedAt time.Time  `json:"-" gorm :"column:updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"column:deleted_at"`
}

func (r *AdminMenuImport) TableName() string {
	return "admin_menu_import"
}
