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
	MenuId     int       `json:"menu_id" gorm:"column:menu_id" form:"menu_id"`
	AppName    string    `json:"app_name" gorm:"column:app_name" form:"app_name"`
	Label      string    `json:"label" gorm:"column:label" form:"label"`
	AppVersion string    `json:"app_version" gorm:"column:app_version" form:"app_version"`
	UrlPath    string    `json:"url_path" gorm:"column:url_path" form:"url_path"`
	SortValue  int       `json:"sort_value" gorm:"scolumn:ort_value" form:"sort_value"`
	OtherValue string    `json:"other_value" gorm :"column:other_value" form:"other_value"`
	CreatedAt  time.Time `json:"-" gorm :"column:created_at" `
	UpdatedAt  time.Time `json:"-" gorm :"column:updated_at" `
	IsDel      int       `json:"-" gorm:"column:is_del" form:"is_del"`
}

func (r *AdminImport) TableName() string {
	return "admin_import"
}
