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
	ParentId   int       `json:"parent_id" gorm:"column:parent_id" form:"parent_id"`
	Label      string    `json:"label" gorm:"column:label" form:"label"`
	Icon       string    `json:"icon" gorm:"column:icon" form:"icon"`
	IsMenuShow int       `json:"is_menu_show" gorm:"column:is_menu_show" form:"is_menu_show"`
	UrlPath    string    `json:"url_path" gorm:"column:url_path" form:"url_path"`
	SortValue  int       `json:"sort_value" gorm:"column:sort_value" form:"sort_value"`
	OtherValue string    `json:"other_value" gorm :"column:other_value" form:"other_value"`
	CreatedAt  time.Time `json:"-" gorm :"column:created_at" `
	UpdatedAt  time.Time `json:"-" gorm :"column:updated_at" `
	IsDel      int       `json:"-" gorm:"column:is_del" form:"is_del"`
}

func (r *AdminMenu) TableName() string {
	return "admin_menu"
}
