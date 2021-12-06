// Package wrappers
/**
* @Author:ChangJiang
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
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/common/app_param"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/pkg/parameters"
	"github.com/juetun/dashboard-api-main/pkg/utils"

	"github.com/juetun/dashboard-api-main/web/models"
)

const (
	DefaultPermitParentId = -1
	DefaultPermitModule   = "platform"
)

type (
	ArgGetImportByMenuIdSingle struct {
		NowMenuId       int64  `json:"now_menu_id" form:"now_menu_id"`
		NowModule       string `json:"now_module" form:"now_module"`
		NowPermitKey    string `json:"now_permit_key" form:"now_permit_key"`
		NowRoutePath    string `json:"now_route_path" form:"now_route_path"`
		NowRouterOrigin string `json:"now_router_origin" form:"now_router_origin"`
		SuperAdminFlag  bool   `json:"super_admin_flag" form:"super_admin_flag"` // 是否为超级管理员
	}
	ArgGetImportByMenuId struct {
		app_param.RequestUser
		ArgGetImportByMenuIdSingle
	}
	ResultGetImportByMenuId struct {
		ImportIds []int64 `json:"import_ids"`
		MenuIds   []int64 `json:"menu_ids"`
	}
	ArgGetAppConfig struct {
		app_param.RequestUser
		Module string `json:"module" form:"module"`
		Env    string `json:"env" form:"env"`
	}
	ResultGetAppConfig map[string]string

	AdminGroup struct {
		models.AdminGroup
		ParentName string `json:"parent_name"`
	}

	ArgAdminMenuWithCheck struct {
		ArgAdminMenu
		GroupId int64 `json:"group_id" form:"group_id"`
	}

	ResultMenuWithCheck struct {
		List []AdminMenuObject       `json:"list"`
		Menu []ResultSystemAdminMenu `json:"menu"` // 一级系统权限列表
	}

	AdminMenuObjectCheck struct {
		ResultAdminMenuSingle
		Children []AdminMenuObject `json:"children"`
	}

	ArgAdminSetPermit struct {
		app_param.RequestUser
		GroupId        int64   `json:"group_id" form:"group_id"`
		Type           string  `json:"type" form:"type"`
		PermitIdString string  `json:"permit_ids" form:"permit_ids"`
		PermitIds      []int64 `json:"-" form:"-"`
		Act            string  `json:"act" form:"act"`
		Module         string  `json:"-"`
	}

	ResultAdminSetPermit struct {
		Result bool `json:"result"`
	}
	ArgDeleteImport struct {
		ID int64 `uri:"id" binding:"required"`
	}
	ResultDeleteImport struct {
		Result bool `json:"result"`
	}
	ArgImportList struct {
		app_param.RequestUser
		response.PageQuery
		PermitKey   string `json:"permit_key" form:"permit_key"`
		AppName     string `json:"app_name" form:"app_name"`
		DefaultOpen uint8  `json:"default_open" form:"default_open"`
		NeedLogin   uint8  `json:"need_login" form:"need_login"`
		NeedSign    uint8  `json:"need_sign" form:"need_sign"`
		UrlPath     string `json:"url_path" form:"url_path"`
	}
	ArgUpdateImportValue struct {
		app_param.RequestUser
		Id     string   `json:"id" form:"id"`
		Ids    []string `json:"-" form:"-"`
		Column string   `json:"column" form:"column"`
		Val    string   `json:"val" form:"val"`
	}

	ResultImportList struct {
		*response.Pager
	}
	ArgEditImport struct {
		Id            int64    `json:"id" form:"id"`
		AppName       string   `json:"app_name" form:"app_name"`
		AppVersion    string   `json:"app_version" form:"app_version"`
		SortValue     int      `json:"sort_value" form:"sort_value"`
		NeedLogin     uint8    `json:"need_login" form:"need_login"`
		NeedSign      uint8    `json:"need_sign" form:"need_sign"`
		UrlPath       string   `json:"url_path" form:"url_path"`
		UrlPaths      []string `json:"-" form:"-"`
		RequestMethod []string `json:"request_methods" form:"request_methods"`
		DefaultOpen   uint8    `json:"default_open" gorm:"column:default_open" form:"default_open"`
		RequestTime   string   `json:"request_time" form:"-"`
	}

	ResultEditImport struct {
		Result bool `json:"result"`
	}
	ResultUpdateImportValue struct {
		Result bool `json:"result"`
	}
	AdminImportList struct {
		models.AdminImport
		RequestMethods []string              `json:"request_methods"`
		Menu           []AdminImportListMenu `json:"menu"`
	}
	AdminImportListMenu struct {
		SystemModuleId int64  `json:"system_module_id"`
		MenuId         int64  `json:"menu_id"`
		Id             int64  `json:"id"` // 接口
		ImportId       int64  `json:"import_id"`
		MenuName       string `json:"menu_name"`
		SystemMenuKey  string `json:"system_menu_key"`
		SystemIcon     string `json:"system_icon"`
		SystemName     string `json:"system_name"`
	}

	OpOne string
	Op    struct {
		PermitKey     string `json:"pk"  gorm:"column:permit_key"`
		MenuPermitKey string `json:"-" gorm:"column:menu_permit_key"`
	}

	ResultSystemMenu struct {
		Id        int64  `gorm:"primary_key" json:"id" form:"id"`
		PermitKey string `json:"permit_key" gorm:"column:permit_key"`
		Label     string `json:"label" gorm:"column:label" form:"label"`
		Icon      string `json:"icon" gorm:"column:icon" form:"icon"`
		SortValue int    `json:"sort_value,omitempty" gorm:"column:sort_value" form:"sort_value"`
		Module    string `json:"module" gorm:"column:module" form:"module"`
		Domain    string `json:"domain" gorm:"column:domain" form:"domain"`
		Active    bool   `json:"active"`
	}
	ResultPermitMenu struct {
		Id        int64              `json:"-"`
		Path      string             `json:"path,omitempty"`
		Module    string             `json:"-"`
		Name      string             `json:"name,omitempty"`
		Label     string             `json:"label,omitempty"`
		Meta      PermitMeta         `json:"meta,omitempty"`
		Children  []ResultPermitMenu `json:"children"`
		Component interface{}        `json:"component,omitempty"`
	}
	AdminMenu struct {
		ID         int    `json:"id"`
		PathName   string `json:"path_name" form:"path_name"`
		ParentId   int    `json:"parent_id" gorm:"parent_id" form:"parent_id"`
		Label      string `json:"label" gorm:"label" form:"label"`
		Icon       string `json:"icon" gorm:"icon" form:"icon"`
		IsMenuShow int    `json:"is_menu_show" gorm:"is_menu_show" form:"is_menu_show"`
		UrlPath    string `json:"url_path" gorm:"url_path" form:"url_path"`
		SortValue  int    `json:"sort_value" gorm:"sort_value" form:"sort_value"`
		OtherValue string `json:"other_value" gorm:"other_value" form:"other_value"`
	}
	ArgFlag struct {
		app_param.RequestUser
	}
	ResultFlag struct {
	}

	AdminGroupUserStruct struct {
		models.AdminUserGroup
		models.AdminGroup
	}
	ArgGetMenu struct {
		app_param.RequestUser
		MenuId int64 `json:"menu_id" form:"menu_id"`
	}

	ResultGetMenu struct {
		models.AdminMenu
	}

	ArgGetImport struct {
		app_param.RequestUser
		response.PageQuery
		Select      string `form:"select" json:"select,omitempty"`
		MenuId      int64  `json:"menu_id,omitempty" form:"menu_id"`
		Checked     bool   `json:"checked,omitempty" form:"checked"` // 是否要查看选中权限情况
		GroupId     int64  `json:"group_id,omitempty" form:"group_id"`
		PermitKey   string `json:"permit_key" form:"permit_key"`
		AppName     string `json:"app_name" form:"app_name"`
		DefaultOpen uint8  `json:"default_open" form:"default_open"`
		NeedLogin   uint8  `json:"need_login" form:"need_login"`
		NeedSign    uint8  `json:"need_sign" form:"need_sign"`
		UrlPath     string `json:"url_path" form:"url_path"`
	}
	AdminImport struct {
		models.AdminImport
		Checked bool `json:"checked"`
	}

	ResultGetImport struct {
		*response.Pager
	}
	ArgAdminMenuSearch struct {
		app_param.RequestUser
		UserHid string `json:"user_hid" form:"user_hid"`
	}

	ResAdminMenuSearch struct {
		List []models.AdminMenu `json:"list"`
	}
	ArgAdminUserAdd struct {
		app_param.RequestUser
		UserHid string `json:"user_hid" form:"user_hid"`
	}
	ResultAdminUserAdd struct {
		Result bool `json:"result"`
	}
	ArgAdminUserDelete struct {
		app_param.RequestUser
		Ids      string   `json:"ids" form:"ids"`
		IdString []string `json:"-" form:"-"`
	}

	ArgMenuAdd struct {
		app_param.RequestUser
		models.AdminMenu
	}

	DaoOrderBy struct {
		Column     string `json:"column"`      // 排序字段
		SortFormat string `json:"sort_format"` // 排序方式
	}

	ResultMenuImport struct {
		*response.Pager
	}
	ResultMenuImportItem struct {
		models.AdminImport
		Checked bool `json:"checked,omitempty"` // 是否要查看选中权限情况
	}
	ArgMenuImport struct {
		app_param.RequestUser
		response.PageQuery
		MenuId  int    `json:"menu_id" form:"menu_id"`
		AppName string `json:"app_name" form:"app_name"`
		UrlPath string `json:"url_path" form:"url_path"`
	}

	ArgMenuImportSet struct {
		app_param.RequestUser
		MenuId    int64   `json:"menu_id" form:"menu_id"`
		ImportId  string  `json:"import_id" form:"import_id"`
		Column    string  `json:"column" form:"column"`
		Value     string  `json:"value" form:"value"`
		ImportIds []int64 `json:"-" form:"-"`
		Type      string  `json:"type" form:"type"`
	}

	ResultMenuImportSet struct {
		Result bool `json:"result"`
	}

	ResultAdminUserDelete struct {
		Result bool `json:"result"`
	}

	ArgAdminUserGroupRelease struct {
		app_param.RequestUser
		Ids      string   `json:"ids" form:"ids"`
		IdString []string `json:"-" form:"-"`
	}

	ResultAdminUserGroupRelease struct {
		Result bool `json:"result"`
	}

	ResultAdminGroupDelete struct {
		Result bool `json:"result"`
	}

	ResultAdminGroupEdit struct {
		Result bool `json:"result"`
	}
	ArgAdminGroupEdit struct {
		app_param.RequestUser
		Name string `json:"name" form:"name"`
		Id   int64  `json:"id" form:"id"`
	}
)

func (r *ArgGetImportByMenuId) Default(c *gin.Context) (err error) {
	_ = c
	if r.NowMenuId == 0 && r.NowRoutePath == "" {
		err = fmt.Errorf("请选择菜单界面")
		return
	}
	return
}

func (r *ArgGetAppConfig) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	return
}

func (r *ArgAdminMenuWithCheck) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	r.PageQuery.DefaultPage()
	return
}

func (r *ArgAdminSetPermit) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	if r.GroupId == 0 {
		err = fmt.Errorf("您没有选择要设置权限的管理组")
		return
	}
	permitIds := strings.Split(r.PermitIdString, ",")
	r.PermitIds = make([]int64, 0, len(permitIds))
	if r.PermitIdString != "" {
		var id int64
		for _, value := range permitIds {
			if value == "" {
				continue
			}
			if id, err = strconv.ParseInt(value, 10, 64); err != nil {
				return
			}
			if id == 0 {
				err = fmt.Errorf("参数异常,请联系管理员")
				return
			}
			r.PermitIds = append(r.PermitIds, id)
		}
	}
	if r.Type == "" {
		err = fmt.Errorf("type is null")
		return
	}
	switch r.Type {
	case models.PathTypePage:
	case models.PathTypeApi:
		if r.Act != models.SetPermitAdd && r.Act != models.SetPermitCancel {
			err = fmt.Errorf("act格式不正确")
			return
		}
	default:
		err = fmt.Errorf("type格式错误")
		return
	}

	return
}

func (r *ArgUpdateImportValue) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	if r.Id == "" {
		err = fmt.Errorf("请选择要修改的数据")
		return
	}
	ids := strings.Split(r.Id, ",")
	for _, value := range ids {
		if value == "" {
			continue
		}
		r.Ids = append(r.Ids, value)
	}
	if len(r.Ids) == 0 {
		err = fmt.Errorf("请选择要修改的数据")
		return
	}

	return
}

func (r *ArgImportList) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	r.PageQuery.DefaultPage()
	return
}

func (r *ArgEditImport) Default(c *gin.Context) (err error) {
	_ = c
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
	if len(r.RequestMethod) == 0 {
		err = fmt.Errorf("请选择请求方法")
		return
	}
	if r.UrlPath != "" {
		r.UrlPaths = strings.Split(r.UrlPath, ",")
	}
	return
}

func (r *ArgMenuImport) Default(c *gin.Context) (err error) {
	_ = c
	return
}
func (r *ArgMenuImportSet) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	if r.Type != "delete" && r.Type != "add" && r.Type != "update" {
		err = fmt.Errorf("type must be delete or add")
		return
	}
	iIds := strings.Split(r.ImportId, ",")
	r.ImportIds = make([]int64, 0, len(iIds))
	for _, value := range iIds {
		if value == "" {
			continue
		}
		id, _ := strconv.ParseInt(value, 10, 64)
		if id > 0 {
			r.ImportIds = append(r.ImportIds, id)
		}
	}
	return
}

func (r *ArgGetImport) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	r.PageQuery.DefaultPage()
	return
}

func (r *ArgAdminUserAdd) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	return
}

func (r *ArgAdminUserDelete) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	if r.Ids != "" {
		r.IdString = strings.Split(r.Ids, ",")
	}
	if len(r.IdString) == 0 {
		err = fmt.Errorf("您没有选择要删除的用户")
		return
	}
	return
}

func (r *ArgAdminUserGroupRelease) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	if r.Ids != "" {
		r.IdString = strings.Split(r.Ids, ",")
	}
	return
}

type ArgAdminUserGroupAdd struct {
	app_param.RequestUser
	GroupId      int    `json:"group_id" form:"group_id"`
	UserHid      string `json:"user_hid" form:"user_hid"`
	GroupIdBatch string `json:"group_id_batch" form:"group_id_batch"`
	UserHidBatch string `json:"user_hid_batch" form:"user_hid_batch"`
	GroupIds     []string
	UserHIds     []string
}

func (r *ArgAdminUserGroupAdd) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}

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
	return
}

type ResultAdminUserGroupAdd struct {
	Result bool `json:"result"`
}

type ArgAdminGroupDelete struct {
	app_param.RequestUser
	Ids      string  `json:"ids" form:"ids"`
	IdString []int64 `json:"-" form:"-"`
}

func (r *ArgAdminGroupDelete) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	idString := strings.Split(r.Ids, ",")
	r.IdString = []int64{}
	var id int64
	for _, v := range idString {
		if v == "" {
			continue
		}
		if id, err = strconv.ParseInt(v, 10, 64); err != nil {
			return
		}
		r.IdString = append(r.IdString, id)
	}
	return
}

func (r *ArgAdminGroupEdit) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	if utf8.RuneCountInString(r.Name) > models.MAXGroupNameLength {
		err = fmt.Errorf("组名长度不能超过%d个字符", models.MAXGroupNameLength)
		return
	}
	return
}

func (r *ArgMenuAdd) Default(c *gin.Context) (err error) {
	_ = c
	if r.ParentId == 0 {
		r.ParentId = DefaultPermitParentId
	}
	if r.Label != "" {
		r.Label = strings.TrimSpace(r.Label)
	}

	if r.Label == "" {
		err = fmt.Errorf("请输入菜单名")
		return
	}
	return
}

type ResultMenuAdd struct {
	Result bool `json:"result"`
}
type ArgMenuSave struct {
	app_param.RequestUser
	models.AdminMenu
}

func (r *ArgMenuSave) Default(c *gin.Context) (err error) {
	_ = c
	if r.ParentId == 0 {
		r.ParentId = DefaultPermitParentId
	}
	if r.PermitKey == DefaultPermitModule {
		err = fmt.Errorf("permit_key不能设置为%s", DefaultPermitModule)
	}
	if r.Domain != "" {
		if strings.TrimPrefix(r.Domain, "https") != r.Domain || strings.TrimPrefix(r.Domain, "https") != r.Domain {
			if strings.TrimSuffix(r.Domain, "/") != r.Domain {
				err = fmt.Errorf("domain格式不正确")
				return
			}
			return
		}
		err = fmt.Errorf("domain格式不正确")
	}
	if utf8.RuneCountInString(r.AdminMenu.Label) > models.AdminMenuNameMaxLength {
		err = fmt.Errorf("菜单名长度不能超过%d个字符", models.AdminMenuNameMaxLength)
		return
	}
	return
}

type ResultMenuSave struct {
	Result bool `json:"result"`
}
type ArgMenuDelete struct {
	app_param.RequestUser
	Ids           string   `json:"ids" form:"ids"`
	IdValue       []string `json:"-"`
	IdValueNumber []int64  `json:"-"`
}

func (r *ArgMenuDelete) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	r.IdValue = make([]string, 0, 5)
	r.IdValueNumber = make([]int64, 0, 5)
	var v int64
	if r.Ids != "" {
		idValue := strings.Split(r.Ids, ",")
		for _, value := range idValue {
			if value != "" {
				r.IdValue = append(r.IdValue, value)
				if v, err = strconv.ParseInt(value, 10, 64); err != nil {
					return
				}
				r.IdValueNumber = append(r.IdValueNumber, v)
			}
		}
	}
	return
}

type ResultMenuDelete struct {
	Result bool `json:"result"`
}

type ArgAdminGroup struct {
	app_param.RequestUser
	response.PageQuery
	Name    string `json:"name" form:"name"`
	GroupId string `json:"group_id" form:"group_id"`
}

func (r *ArgAdminGroup) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	r.PageType = response.DefaultPageTypeList
	return
}

type ResultAdminGroup struct {
	response.Pager
}

type ArgAdminMenu struct {
	app_param.RequestUser
	response.PageQuery
	Id         int64  `json:"id" form:"id"`
	Label      string `json:"label" form:"label"`
	AppName    string `json:"app_name" form:"app_name"`
	UserHId    string `json:"user_hid" form:"user_hid"`
	ParentId   int64  `json:"parent_id" form:"parent_id"`
	IsMenuShow int    `json:"is_menu_show" form:"is_menu_show"`
	IsDel      int    `json:"is_del" form:"is_del"`
	Module     string `json:"module" form:"module"`
	SystemId   int64  `json:"system_id" form:"system_id"`
}

func (r *ArgAdminMenu) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	if r.Module == "" {
		r.Module = parameters.DefaultSystem
	}
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
	Id                 int64  `json:"id"`
	ParentId           int64  `json:"parent_id"`
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
	Domain             string `json:"domain"`
	PermitKey          string `json:"permit_key"`
	ManageImportPermit uint8  `json:"manage_import_permit"`
	ResultAdminMenuOtherValue
	IsDel int `json:"is_del"`
}
type ResultAdminMenu struct {
	List []AdminMenuObject       `json:"list"`
	Menu []ResultSystemAdminMenu `json:"menu"` // 一级系统权限列表
}

func NewResultAdminMenu() *ResultAdminMenu {
	return &ResultAdminMenu{
		List: make([]AdminMenuObject, 0, 20),
		Menu: make([]ResultSystemAdminMenu, 0, 30),
	}
}

type ResultSystemAdminMenu struct {
	Id                 int64  `gorm:"primary_key" json:"id" form:"id"`
	PermitKey          string `json:"permit_key" gorm:"column:permit_key"`
	ManageImportPermit uint8  `json:"manage_import_permit" gorm:"column:manage_import_permit" form:"manage_import_permit"`
	Label              string `json:"label" gorm:"column:label" form:"label"`
	Icon               string `json:"icon" gorm:"column:icon" form:"icon"`
	SortValue          int    `json:"sort_value" gorm:"column:sort_value" form:"sort_value"`
	Module             string `json:"module" gorm:"column:module" form:"module"`
	Domain             string `json:"domain" gorm:"column:domain" form:"domain"`
	Active             bool   `json:"active"`
}
type ArgAdminUser struct {
	app_param.RequestUser
	response.PageQuery
	Name    string `json:"name" form:"name"`
	UserHId string `json:"user_hid" form:"user_hid"`
}

func (r *ArgAdminUser) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	return
}

type ResultAdminUser struct {
	response.Pager
}
type AdminUserGroupName struct {
	models.AdminUserGroup
	GroupName    string `json:"group_name"`
	IsSuperAdmin uint8  `json:"is_super_admin"`
	IsAdminGroup uint8  `json:"is_admin_group"`
}
type ResultAdminUserList struct {
	models.AdminUser
	Group []AdminUserGroupName `json:"group"`
}

type ArgPermitMenu struct {
	app_param.RequestUser
	ArgGetImportByMenuIdSingle          // 通用参数逻辑处理 用于获取当前菜单下的接口列表
	ParentId                   int64    `json:"parent_id"` // 上级菜单ID
	PathType                   string   `json:"path_type" form:"path_type"`
	PathTypes                  []string `json:"-" form:"-"`
	Module                     string   `json:"module" form:"module"` // 系统ID
	IsSuperAdmin               bool     `json:"-" form:"-"`           // 是否为超级管理员
	GroupId                    []int64  `json:"-" form:"-"`
}

// Default 初始化默认值
func (r *ArgPermitMenu) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	r.PathTypes = []string{}

	if r.PathType != "" {
		pType := strings.Split(r.PathType, ",")
		for _, value := range pType {
			if value == "" {
				continue
			}
			r.PathTypes = append(r.PathTypes, value)
		}
		// r.PathTypes = []string{"pages", "system"}
	}
	return
}

type PermitMeta struct {
	PermitKey  string `json:"permitKey"` // 控制权限结构的参数
	Icon       string `json:"icon"`
	Title      string `json:"title"`
	HideInMenu bool   `json:"hideInMenu"`
}
type ResultPermitMenuReturn struct {
	ResultPermitMenu                         // 当前选中的权限
	RoutParentMap    map[string][]string     `json:"routParentMap"`      // 当前菜单列表
	Menu             []ResultSystemMenu      `json:"menu"`               // 一级系统权限列表 用户从当前系统跳转到其他管理系统
	OpList           map[string][]OpOne      `json:"op_list"`            // 获取接口权限列表
	NotReadMsgCount  int                     `json:"not_read_msg_count"` // 未读消息数量
	NowImportAndMenu ResultGetImportByMenuId `json:"import_and_menu"`    // 当前菜单下有的菜单和接口列表
}

func NewResultPermitMenuReturn() (res *ResultPermitMenuReturn) {
	res = &ResultPermitMenuReturn{
		ResultPermitMenu: ResultPermitMenu{
			Children: []ResultPermitMenu{},
		},
		RoutParentMap:   map[string][]string{},
		Menu:            []ResultSystemMenu{},
		OpList:          map[string][]OpOne{},
		NotReadMsgCount: 0,
		NowImportAndMenu: ResultGetImportByMenuId{
			ImportIds: []int64{},
			MenuIds:   []int64{},
		},
	}
	return
}

func (r *ArgFlag) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	return
}

func (r *ArgAdminMenuSearch) Default() {

}
func (r *ArgGetMenu) Default(c *gin.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	return
}
