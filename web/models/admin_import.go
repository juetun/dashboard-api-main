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
	Id            int        `gorm:"primary_key" json:"id" form:"id"`
	MenuId        int        `json:"menu_id" gorm:"column:menu_id" form:"menu_id"`
	AppName       string     `json:"app_name" gorm:"column:app_name" form:"app_name"`
	AppVersion    string     `json:"app_version" gorm:"column:app_version" form:"app_version"`
	UrlPath       string     `json:"url_path" gorm:"column:url_path" form:"url_path"`
	SortValue     int        `json:"sort_value" gorm:"scolumn:ort_value" form:"sort_value"`
	RequestMethod string     `json:"request_method" gorm :"column:request_method" form:"request_method"`
	CreatedAt     time.Time  `json:"-" gorm :"column:created_at" `
	UpdatedAt     time.Time  `json:"-" gorm :"column:updated_at" `
	DeletedAt     *time.Time `json:"-" gorm:"column:deleted_at"`
}

func (r *AdminImport) TableName() string {
	return "admin_import"
}
