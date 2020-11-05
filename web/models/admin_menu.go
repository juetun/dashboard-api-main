/**
* @Author:changjiang
* @Description:
* @File:admin_menu
* @Version: 1.0.0
* @Date 2020/9/16 10:37 下午
 */
package models

import (
	"time"
)

type AdminMenu struct {
	Id         int       `gorm:"primary_key" json:"id" form:"id"`
	ParentId   int       `json:"parent_id" gorm:"parent_id" form:"parent_id"`
	AppName    string    `json:"app_name" gorm:"app_name" form:"app_name"`
	Label      string    `json:"label" gorm:"label" form:"label"`
	Icon       string    `json:"icon" gorm:"icon" form:"icon"`
	IsMenuShow int       `json:"is_menu_show" gorm:"is_menu_show" form:"is_menu_show"`
	AppVersion string    `json:"app_version" gorm:"app_version" form:"app_version"`
	UrlPath    string    `json:"url_path" gorm:"url_path" form:"url_path"`
	PathType   string    `json:"path_type" gorm:"path_type" form:"path_type"`
	SortValue  int       `json:"sort_value" gorm:"sort_value" form:"sort_value"`
	OtherValue string    `json:"other_value" gorm :"other_value" form:"other_value"`
	CreatedAt  time.Time `json:"-" gorm :"created_at" `
	UpdatedAt  time.Time `json:"-" gorm :"updated_at" `
	IsDel      int       `json:"-" gorm:"is_del" form:"is_del"`
}

func (r *AdminMenu) TableName() string {
	return "admin_menu"
}
