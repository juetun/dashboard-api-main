/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 12:22 上午
 */
package pojos

import (
	"strings"

	"github.com/juetun/base-wrapper/lib/app_obj"
	"github.com/juetun/base-wrapper/lib/common/response"

	"github.com/juetun/dashboard-api-main/web/models"
)

type ArgMenuAdd struct {
	app_obj.JwtUserMessage
	models.AdminMenu
}

func (r *ArgMenuAdd) Default() {

}

type ResultMenuAdd struct {
}
type ArgMenuSave struct {
	app_obj.JwtUserMessage
	models.AdminMenu
}

func (r *ArgMenuSave) Default() {

}

type ResultMenuSave struct {
	Result bool `json:"result"`
}
type ArgMenuDelete struct {
	app_obj.JwtUserMessage
	Ids     string   `json:"ids" form:"ids"`
	IdValue []string `json:"-"`
}

func (r *ArgMenuDelete) Default() {
	r.IdValue = make([]string, 0, 5)
	if r.Ids != "" {
		idValue := strings.Split(r.Ids, ",")
		for _, value := range idValue {
			if value != "" {
				r.IdValue = append(r.IdValue, value)
			}
		}
	}
}

type ResultMenuDelete struct {
	Result bool `json:"result"`
}

type ArgAdminGroup struct {
	app_obj.JwtUserMessage
	response.BaseQuery
	Name    string `json:"name" form:"name"`
	GroupId int    `json:"group_id" form:"group_id"`
}

func (r *ArgAdminGroup) Default() {

}

type ResultAdminGroup struct {
	response.Pager
}

type ArgAdminMenu struct {
	app_obj.JwtUserMessage
	response.BaseQuery
	Label      string `json:"label" form:"label"`
	AppName    string `json:"app_name" form:"app_name"`
	UserHId    string `json:"user_hid" form:"user_hid"`
	ParentId   int    `json:"parent_id" form:"parent_id"`
	IsMenuShow int    `json:"is_menu_show" form:"is_menu_show"`
	IsDel      int    `json:"is_del" form:"is_del"`
}

func (r *ArgAdminMenu) Default() {

}

type AdminMenuObject struct {
	ResultAdminMenuSingle
	Children []AdminMenuObject `json:"children"`
}

type ResultAdminMenuOtherValue struct {
	Expand          bool `json:"expand"`
	Disabled        bool `json:"disabled"`
	DisableCheckbox bool `json:"disableCheckbox"`
	Checked         bool `json:"checked"`
}
type ResultAdminMenuSingle struct {
	Id         int    `json:"id"`
	ParentId   int    `json:"parent_id"`
	AppName    string `json:"app_name"`
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	IsMenuShow int    `json:"is_menu_show"`
	AppVersion string `json:"app_version"`
	UrlPath    string `json:"url_path"`
	PathType   string `json:"path_type"`
	SortValue  int    `json:"sort_value"`
	ResultAdminMenuOtherValue
	IsDel int `json:"is_del"`
}
type ResultAdminMenu struct {
	List []AdminMenuObject `json:"list"`
}

type ArgAdminUser struct {
	app_obj.JwtUserMessage
	response.BaseQuery
	Name    string `json:"name" form:"name"`
	UserHId string `json:"user_hid" form:"user_hid"`
}

func (r *ArgAdminUser) Default() {

}

type ResultAdminUser struct {
	response.Pager
}

type ArgPermitMenu struct {
	app_obj.JwtUserMessage
	PathType string `json:"path_type" form:"path_type"`
}

// 初始化默认值
func (r *ArgPermitMenu) Default() {
	if r.PathType == "" {
		r.PathType = "page"
	}
}

type ResultPermitMenu struct {
	Menu []models.AdminMenu `json:"menu"`
}

type ArgFlag struct {
	app_obj.JwtUserMessage
}

type ResultFlag struct {
}
