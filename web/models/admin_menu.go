/**
* @Author:changjiang
* @Description:
* @File:admin_menu
* @Version: 1.0.0
* @Date 2020/9/16 10:37 下午
 */
package models

type AdminMenu struct {
	Id         int    `gorm:"primary_key" json:"id"`
	ParentId   int    `json:"parent_id" gorm:"parent_id"`
	AppName    string `json:"app_name" gorm:"app_name"`
	Label      string `json:"label" gorm:"label"`
	Icon       string `json:"icon" gorm:"icon"`
	IsMenuShow int    `json:"is_menu_show" gorm:"is_menu_show"`
	AppVersion string `json:"app_version" gorm:"app_version"`
	UrlPath    string `json:"url_path" gorm:"url_path"`
	PathType   string `json:"path_type" gorm:"path_type"`
	SortValue  int    `json:"sort_value" gorm:"sort_value"`
	IsDel      int    `json:"-" gorm:"is_del"`
}

func (r *AdminMenu) TableName() string {
	return "admin_menu"
}
