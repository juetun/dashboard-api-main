/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 12:22 上午
 */
package wrappers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/pkg/utils"

	"github.com/juetun/dashboard-api-main/web/models"
)

const DefaultPermitParentId = 1

type ArgDeleteImport struct {
	ID int `uri:"id" binding:"required"`
}
type ResultDeleteImport struct {
	Result bool `json:"result"`
}
type ArgEditImport struct {
	AppName       string   `json:"app_name" form:"app_name"`
	AppVersion    string   `json:"app_version" form:"app_version"`
	Id            int      `json:"id" form:"id"`
	MenuId        int      `json:"menu_id" form:"menu_id"`
	SortValue     int      `json:"sort_value" form:"sort_value"`
	UrlPath       string   `json:"url_path" form:"url_path"`
	RequestMethod []string `json:"request_method" form:"request_method"`
	RequestTime   string   `json:"request_time" form:"-"`
}

type ResultEditImport struct {
	Result bool `json:"result"`
}

func (r *ArgEditImport) Default(c *gin.Context) (err error) {
	if r.MenuId == 0 {
		err = fmt.Errorf("您没有选择要添加接口权限的菜单")
		return
	}
	if r.AppName == "" {
		err = fmt.Errorf("请输入接口所属应用KEY")
		return
	}
	if r.UrlPath == "" {
		err = fmt.Errorf("请输入接口path路径")
		return
	}
	if r.AppVersion == "" {
		r.AppVersion = "1.0"
	}
	r.RequestTime = utils.DateTime(time.Now())
	return
}

type ArgGetImport struct {
	app_obj.JwtUserMessage
	base.ReqPager
	MenuId int `json:"menu_id" form:"menu_id"`
}

func (r *ArgGetImport) Default(c *gin.Context) {
	r.JwtUserMessage = GetUser(c)
	r.ReqPager.DefaultPager()
}

type ResultGetImport struct {
	base.Pager
}
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
	if r.ParentId == 0 {
		r.ParentId = DefaultPermitParentId
	}
	if r.Label != "" {
		r.Label = strings.TrimSpace(r.Label)
	}
}

type ResultMenuAdd struct {
	Result bool `json:"result"`
}
type ArgMenuSave struct {
	app_obj.JwtUserMessage
	models.AdminMenu
}

func (r *ArgMenuSave) Default() {
	if r.ParentId == 0 {
		r.ParentId = DefaultPermitParentId
	}
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
	Id         int    `json:"id" form:"id"`
	Label      string `json:"label" form:"label"`
	AppName    string `json:"app_name" form:"app_name"`
	UserHId    string `json:"user_hid" form:"user_hid"`
	ParentId   int    `json:"parent_id" form:"parent_id"`
	IsMenuShow int    `json:"is_menu_show" form:"is_menu_show"`
	IsDel      int    `json:"is_del" form:"is_del"`
	Module     string `json:"module" form:"module"`
	SystemId   int    `json:"system_id" form:"system_id"`
}

func (r *ArgAdminMenu) Default(c *gin.Context) (err error) {
	r.JwtUserMessage = GetUser(c)
	return
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
	Id                 int    `json:"id"`
	ParentId           int    `json:"parent_id"`
	AppName            string `json:"app_name"`
	Title              string `json:"title"`
	Label              string `json:"label"`
	Icon               string `json:"icon"`
	HideInMenu         uint8  `json:"hide_in_menu"`
	AppVersion         string `json:"app_version"`
	UrlPath            string `json:"url_path"`
	PathType           string `json:"path_type"`
	SortValue          int    `json:"sort_value"`
	Module             string `json:"module"`
	PermitKey          string `json:"permit_key"`
	ManageImportPermit uint8  `json:"manage_import_permit"`
	ResultAdminMenuOtherValue
	IsDel int `json:"is_del"`
}
type ResultAdminMenu struct {
	List []AdminMenuObject       `json:"list"`
	Menu []ResultSystemAdminMenu `json:"menu"` // 一级系统权限列表
}
type ResultSystemAdminMenu struct {
	Id                 int    `gorm:"primary_key" json:"id" form:"id"`
	PermitKey          string `json:"permit_key" gorm:"permit_key"`
	ManageImportPermit uint8  `json:"manage_import_permit" gorm:"column:manage_import_permit" form:"manage_import_permit"`
	Label              string `json:"label" gorm:"label" form:"label"`
	Icon               string `json:"icon" gorm:"icon" form:"icon"`
	SortValue          int    `json:"sort_value" gorm:"sort_value" form:"sort_value"`
	Module             string `json:"module" gorm:"module" form:"module"`
	Active             bool   `json:"active"`
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
	ParentId  int      `json:"parent_id"`
	PathType  string   `json:"path_type" form:"path_type"`
	PathTypes []string `json:"path_type" form:"path_type"`
}

// 初始化默认值
func (r *ArgPermitMenu) Default() {
	r.PathTypes = []string{}
	if r.PathType == "" {
		r.PathTypes = []string{"page", "system"}
	}
}

type PermitMeta struct {
	PermitKey  string `json:"permitKey"` // 控制权限结构的参数
	Icon       string `json:"icon"`
	Title      string `json:"title"`
	HideInMenu bool   `json:"hideInMenu"`
}
type ResultPermitMenuReturn struct {
	ResultPermitMenu                     // 当前选中的权限
	RoutParentMap    map[string][]string `json:"routParentMap"`
	Menu             []ResultSystemMenu  `json:"menu"` // 一级系统权限列表
}

type ResultSystemMenu struct {
	Id        int    `gorm:"primary_key" json:"id" form:"id"`
	PermitKey string `json:"permit_key" gorm:"permit_key"`
	Label     string `json:"label" gorm:"label" form:"label"`
	Icon      string `json:"icon" gorm:"icon" form:"icon"`
	SortValue int    `json:"sort_value" gorm:"sort_value" form:"sort_value"`
	Module    string `json:"module" gorm:"module" form:"module"`
	Active    bool   `json:"active"`
}
type ResultPermitMenu struct {
	Id        int                `json:"id"`
	Path      string             `json:"path"`
	Module    string             `json:"module"`
	Name      string             `json:"name"`
	Meta      PermitMeta         `json:"meta"`
	Children  []ResultPermitMenu `json:"children"`
	Component interface{}        `json:"component"`
}
type AdminMenu struct {
	ID         int    `json:"id"`
	PathName   string `json:"path_name" form:"path_name"`
	ParentId   int    `json:"parent_id" gorm:"parent_id" form:"parent_id"`
	Label      string `json:"label" gorm:"label" form:"label"`
	Icon       string `json:"icon" gorm:"icon" form:"icon"`
	IsMenuShow int    `json:"is_menu_show" gorm:"is_menu_show" form:"is_menu_show"`
	UrlPath    string `json:"url_path" gorm:"url_path" form:"url_path"`
	SortValue  int    `json:"sort_value" gorm:"sort_value" form:"sort_value"`
	OtherValue string `json:"other_value" gorm :"other_value" form:"other_value"`
}
type ArgFlag struct {
	app_obj.JwtUserMessage
}

type ResultFlag struct {
}

type AdminGroupUserStruct struct {
	models.AdminUserGroup
	models.AdminGroup
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
func GetUser(c *gin.Context) (jwtUser app_obj.JwtUserMessage) {
	jwtUser = app_obj.JwtUserMessage{}
	v, e := c.Get(app_obj.ContextUserObjectKey)
	if e {
		jwtUser = v.(app_obj.JwtUserMessage)
	}

	return jwtUser
}
