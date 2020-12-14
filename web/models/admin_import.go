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

type AdminImport struct {
	Id         int       `gorm:"primary_key" json:"id" form:"id"`
	MenuId     int       `json:"menu_id" gorm:"menu_id" form:"menu_id"`
	AppName    string    `json:"app_name" gorm:"app_name" form:"app_name"`
	Label      string    `json:"label" gorm:"label" form:"label"`
	AppVersion string    `json:"app_version" gorm:"app_version" form:"app_version"`
	UrlPath    string    `json:"url_path" gorm:"url_path" form:"url_path"`
	SortValue  int       `json:"sort_value" gorm:"sort_value" form:"sort_value"`
	OtherValue string    `json:"other_value" gorm :"other_value" form:"other_value"`
	CreatedAt  time.Time `json:"-" gorm :"created_at" `
	UpdatedAt  time.Time `json:"-" gorm :"updated_at" `
	IsDel      int       `json:"-" gorm:"is_del" form:"is_del"`
}

func (r *AdminImport) TableName() string {
	return "admin_import"
}
