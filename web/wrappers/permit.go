/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 12:22 上午
 */
package wrappers

import (
	"strconv"
	"strings"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/common/response"

	"github.com/juetun/dashboard-api-main/web/models"
)

type ArgAdminMenuSearch struct {
	app_obj.JwtUserMessage
	UserHid string `json:"user_hid" form:"user_hid"`
}

func (r *ArgAdminMenuSearch) Default() {

}

type ResAdminMenuSearch struct {
	List []models.AdminMenu `json:"list"`
}
type ArgAdminUserAdd struct {
	app_obj.JwtUserMessage
	UserHid string `json:"user_hid" form:"user_hid"`
}

func (r *ArgAdminUserAdd) Default() {}

type ResultAdminUserAdd struct {
	Result bool `json:"result"`
}
type ArgAdminUserDelete struct {
	app_obj.JwtUserMessage
	Ids      string   `json:"ids" form:"ids"`
	IdString []string `json:"-" form:"-"`
}

func (r *ArgAdminUserDelete) Default() {
	if r.Ids != "" {
		r.IdString = strings.Split(r.Ids, ",")
	}
}

type ResultAdminUserDelete struct {
	Result bool `json:"result"`
}

type ArgAdminUserGroupRelease struct {
	app_obj.JwtUserMessage
	Ids      string   `json:"ids" form:"ids"`
	IdString []string `json:"-" form:"-"`
}

func (r *ArgAdminUserGroupRelease) Default() {
	if r.Ids != "" {
		r.IdString = strings.Split(r.Ids, ",")
	}
}

type ResultAdminUserGroupRelease struct {
	Result bool `json:"result"`
}

type ArgAdminUserGroupAdd struct {
	app_obj.JwtUserMessage
	GroupId      int    `json:"group_id" form:"group_id"`
	UserHid      string `json:"user_hid" form:"user_hid"`
	GroupIdBatch string `json:"group_id_batch" form:"group_id_batch"`
	UserHidBatch string `json:"user_hid_batch" form:"user_hid_batch"`
	GroupIds     []string
	UserHIds     []string
}

func (r *ArgAdminUserGroupAdd) Default() {
	r.GroupIds = []string{}
	r.UserHIds = []string{}
	if r.GroupIdBatch != "" {
		r.GroupIds = strings.Split(r.GroupIdBatch, ",")
	}
	if r.UserHidBatch != "" {
		r.UserHIds = strings.Split(r.UserHidBatch, ",")
	}
	if r.GroupId != 0 {
		r.GroupIds = append(r.GroupIds, strconv.Itoa(r.GroupId))
	}
	if r.UserHid != "" {
		r.UserHIds = append(r.UserHIds, r.UserHid)
	}
}

type ResultAdminUserGroupAdd struct {
	Result bool `json:"result"`
}

type ArgAdminGroupDelete struct {
	app_obj.JwtUserMessage
	Ids      string   `json:"ids" form:"ids"`
	IdString []string `json:"-" form:"-"`
}

func (r *ArgAdminGroupDelete) Default() {
	idString := strings.Split(r.Ids, ",")
	r.IdString = []string{}
	for _, v := range idString {
		if v == "" {
			continue
		}
		r.IdString = append(r.IdString, v)
	}
}

type ResultAdminGroupDelete struct {
	Result bool `json:"result"`
}

type ResultAdminGroupEdit struct {
	Result bool `json:"result"`
}
type ArgAdminGroupEdit struct {
	app_obj.JwtUserMessage
	Name string `json:"name" form:"name"`
	Id   int    `json:"id" form:"id"`
}

func (r *ArgAdminGroupEdit) Default() {

}

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
type AdminUserGroupName struct {
	models.AdminUserGroup
	GroupName string `json:"group_name"`
}
type ResultAdminUserList struct {
	models.AdminUser
	Group []AdminUserGroupName `json:"group"`
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

type ArgGetMenu struct {
	app_obj.JwtUserMessage
	MenuId int `json:"menu_id" form:"menu_id"`
}

type ResultGetMenu struct {
	models.AdminMenu
}

func (r *ArgGetMenu) Default() {

}
