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
	"github.com/juetun/base-wrapper/lib/base"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/pkg/parameters"
	"github.com/juetun/dashboard-api-main/pkg/utils"
	"github.com/juetun/library/common/app_param"

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
		ImportIds []ImportSingle `json:"import_ids,omitempty"` //当前界面能够访问的接口列表
		MenuIds   []MenuSingle   `json:"menu_ids,omitempty"`   //当前界面内能够跳转的界面ID列表
	}
	MenuSingle struct {
		MenuId    int64  `json:"id"`
		PermitKey string `json:"permit_key"`
	}
	ImportSingle struct {
		ImportId      int64  `json:"id"`
		MenuId        int64  `json:"-"`
		MenuPermitKey string `json:"-"`
		PermitKey     string `json:"permit_key"`
	}
	ArgGetAppConfig struct {
		app_param.RequestUser
		Module string `json:"module" form:"module"`
		Env    string `json:"env" form:"env"`
	}
	ResultGetAppConfig map[string]string

	AdminGroup struct {
		models.AdminGroup
		ParentName      string `json:"parent_name"`
		UserCount       int    `json:"user_count"`
		UpdatedAtString string `json:"updated_at"`
	}

	ArgAdminMenuWithCheck struct {
		ArgAdminMenu
		GroupId              int64   `json:"group_id" form:"group_id"`
		OperatorGroupId      []int64 `json:"-" form:"-"` // 当前操作用户所属的组
		OperatorIsSuperAdmin bool    `json:"-" form:"-"` // 当前操作用户是否为超级管理员
	}

	ResultMenuWithCheck struct {
		List []AdminMenuObject       `json:"list"`
		Menu []ResultSystemAdminMenu `json:"menu"` // 一级系统权限列表
	}

	AdminMenuObjectCheck struct {
		ResultAdminMenuSingle
		Children []AdminMenuObject `json:"children"`
	}
	AdminImportWithMenu struct {
		models.AdminImport
		MenuId int64 `json:"menu_id"`
	}
	ArgAdminSetPermit struct {
		app_param.RequestUser
		GroupId        int64   `json:"group_id" form:"group_id"`
		MenuId         int64   `json:"menu_id" form:"menu_id"`
		DefaultOpen    uint8   `json:"default_open" form:"default_open"`
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
		response.Pager
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
		Id                 int64  `json:"id" gorm:"column:id;primary_key" `
		Name               string `json:"name" gorm:"column:name"`
		ParentId           int64  `json:"parent_id"  gorm:"column:parent_id"`
		GroupCode          string `json:"group_code" gorm:"column:group_code"`
		LastChildGroupCode string `json:"last_child_group_code" gorm:"column:last_child_group_code"`
		IsSuperAdmin       uint8  `json:"is_super_admin" gorm:"column:is_super_admin"`
		IsAdminGroup       uint8  `json:"is_admin_group" gorm:"column:is_admin_group"`

		GroupId    int64  `json:"group_id" gorm:"column:group_id;not null;uniqueIndex:idx_gid_hid,priority:1;default:0;comment:组ID"`
		UserHid    string `json:"user_hid"  gorm:"column:user_hid;not null;default:'';uniqueIndex:idx_gid_hid,priority:2;comment:用户ID"`
		SuperAdmin uint8  `json:"is_super_admin" gorm:"column:super_admin;not null;default:0;comment:是否是超级管理员 0-否,1-是"`
		AdminGroup uint8  `json:"is_admin_group" gorm:"column:admin_group;not null;default:0;comment:是否是后台管理员组 0-否,1-是"`
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
		response.Pager
	}
	ArgAdminMenuSearch struct {
		app_param.RequestUser
		UserHid string `json:"user_hid" form:"user_hid"`
	}

	ResAdminMenuSearch struct {
		List []*models.AdminMenu `json:"list"`
	}
	ArgAdminUserAdd struct {
		app_param.RequestUser
		UserHid   int64  `json:"user_hid" form:"user_hid"`
		RealName  string `json:"real_name" form:"real_name"`
		Mobile    string `json:"mobile" form:"mobile"`
		Id        int64  `json:"id" form:"id"`
		FlagAdmin uint8  `json:"flag_admin" form:"flag_admin"`
	}
	ResultAdminUserAdd struct {
		Result bool `json:"result"`
	}
	ArgAdminUserDelete struct {
		app_param.RequestUser
		Ids      string  `json:"ids" form:"ids"`
		IdString []int64 `json:"-" form:"-"`
	}

	ArgMenuAdd struct {
		app_param.RequestUser
		models.AdminMenu
	}

	DaoOrderBy struct {
		Column     string `json:"column"`      // 排序字段
		SortFormat string `json:"sort_format"` // 排序方式
	}

	ResultMenuImportItem struct {
		models.AdminImport
		Checked bool `json:"checked,omitempty"` // 是否要查看选中权限情况
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
		Name         string `json:"name" form:"name"`
		Id           int64  `json:"id" form:"id"`
		IsSuperAdmin uint8  `json:"is_super_admin" form:"is_super_admin"`
		IsAdminGroup uint8  `json:"is_admin_group" form:"is_admin_group"`
	}
)

func (r *ArgDeleteImport) Default(ctx *base.Context) (err error) {
	return
}

func (r *ResultMenuWithCheck) SetSystemList(list []*models.AdminMenu, systemId int64) {
	var data ResultSystemAdminMenu
	for _, item := range list {
		data = ResultSystemAdminMenu{
			Id:        item.Id,
			PermitKey: item.PermitKey,
			Label:     item.Label,
			Icon:      item.Icon,
			SortValue: item.SortValue,
			Module:    item.Module,
			Domain:    item.Domain,
		}
		if item.Id == systemId {
			data.Active = true
		}
	}
	return
}

func (r *ArgGetImportByMenuId) Default(c *base.Context) (err error) {
	_ = c
	if r.NowMenuId == 0 && r.NowRoutePath == "" {
		err = fmt.Errorf("请选择菜单界面")
		return
	}
	return
}

func (r *ArgGetAppConfig) Default(c *base.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	return
}

func (r *ArgAdminMenuWithCheck) Default(c *base.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	r.PageQuery.DefaultPage()
	return
}

func (r *ArgAdminSetPermit) Default(c *base.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	if r.GroupId == 0 {
		err = fmt.Errorf("您没有选择要设置权限的管理组")
		return
	}
	permitIds := strings.Split(r.PermitIdString, ",")
	r.PermitIds = make([]int64, 0, len(permitIds))
	if r.PermitIdString != "" && r.PermitIdString != "0" {
		var id int64
		for _, value := range permitIds {
			if value == "" {
				continue
			}
			if id, err = strconv.ParseInt(value, 10, 64); err != nil {
				return
			}
			if id == 0 {
				err = fmt.Errorf("参数异常,请刷新重试")
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

func (r *ArgUpdateImportValue) Default(c *base.Context) (err error) {
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

func (r *ArgImportList) Default(c *base.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	r.PageQuery.DefaultPage()
	return
}

func (r *ArgEditImport) Default(c *base.Context) (err error) {
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

func (r *ArgMenuImportSet) Default(c *base.Context) (err error) {
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

func (r *ArgGetImport) Default(c *base.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	r.PageQuery.DefaultPage()
	return
}

func (r *ArgAdminUserAdd) Default(c *base.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	return
}
func (r *ArgGetImportByMenuIdSingle) Default(c *base.Context) (err error) {
	divString := "/"
	pathSlice := strings.Split(strings.TrimLeft(r.NowRoutePath, "/"), divString)
	switch len(pathSlice) {
	case 0:
		pathSlice = append(pathSlice, []string{"", ""}...)
	case 1:
		pathSlice = append(pathSlice, "")
	}
	r.NowModule = pathSlice[0]
	r.NowPermitKey = strings.Join(pathSlice[1:], divString)
	switch r.NowPermitKey {
	case "home":
		r.NowPermitKey = ""
	}
	return
}
func (r *ArgAdminUserDelete) Default(c *base.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	if r.Ids != "" {
		IdString := strings.Split(r.Ids, ",")
		var it int64
		for _, i2 := range IdString {
			if it, err = strconv.ParseInt(i2, 10, 64); err != nil {
				return
			}
			r.IdString = append(r.IdString, it)
		}

	}
	if len(r.IdString) == 0 {
		err = fmt.Errorf("您没有选择要删除的用户")
		return
	}
	return
}

func (r *ArgAdminUserGroupRelease) Default(c *base.Context) (err error) {
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
	GroupId      int64  `json:"group_id" form:"group_id"`
	UserHid      int64  `json:"user_hid" form:"user_hid"`
	GroupIdBatch string `json:"group_id_batch" form:"group_id_batch"`
	UserHidBatch string `json:"user_hid_batch" form:"user_hid_batch"`
	GroupIds     []int64
	UserHIds     []int64
}

func (r *ArgAdminUserGroupAdd) Default(c *base.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	if r.UserHid == 0 {
		err = fmt.Errorf("您没有选择要添加权限组的用户ID")
		return
	}
	r.GroupIds = []int64{}
	r.UserHIds = []int64{}

	if r.GroupIdBatch != "" {
		var gid int64
		tmp := strings.Split(r.GroupIdBatch, ",")
		for _, s := range tmp {
			if s == "" {
				continue
			}
			if gid, err = strconv.ParseInt(s, 10, 64); err != nil {
				return
			}
			r.GroupIds = append(r.GroupIds, gid)
		}
	}
	if r.UserHidBatch != "" {
		uidString := strings.Split(r.UserHidBatch, ",")
		r.UserHIds = make([]int64, 0, len(uidString))
		var uidNum int64
		for _, s := range uidString {
			if s == "" {
				continue
			}
			if uidNum, err = strconv.ParseInt(s, 10, 64); err != nil {
				return
			}
			r.UserHIds = append(r.UserHIds, uidNum)
		}
	}
	if r.GroupId != 0 {
		r.GroupIds = append(r.GroupIds, r.GroupId)
	}
	if r.UserHid != 0 {
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

func (r *ArgAdminGroupDelete) Default(c *base.Context) (err error) {
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

func (r *ArgAdminGroupEdit) Default(c *base.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	if utf8.RuneCountInString(r.Name) > models.MAXGroupNameLength {
		err = fmt.Errorf("组名长度不能超过%d个字符", models.MAXGroupNameLength)
		return
	}
	if r.Name == "" {
		err = fmt.Errorf("请输入管理组名称")
		return
	}
	return
}

func (r *ArgMenuAdd) Default(c *base.Context) (err error) {
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

func (r *ArgMenuSave) Default(c *base.Context) (err error) {
	r.PermitKey = strings.TrimSpace(r.PermitKey)//去掉首尾空格
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
	if r.SortValue >= models.AdminMenuMaxSortValue {
		err = fmt.Errorf("排序值不能超过%d", models.AdminMenuMaxSortValue)
		return
	}
	if r.PermitKey == "" {
		err = fmt.Errorf("请输入菜单的KEY")
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

func (r *ArgMenuDelete) Default(c *base.Context) (err error) {
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

func (r *ArgAdminGroup) Default(c *base.Context) (err error) {
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

	OperatorGroupId      []int64 `json:"-" form:"-"` // 当前操作用户所属的组
	OperatorIsSuperAdmin bool    `json:"-" form:"-"` // 当前操作用户是否为超级管理员
}

func (r *ArgAdminMenu) Default(c *base.Context) (err error) {
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
	Name         string  `json:"name" form:"name"`
	UserHIds     string  `json:"user_ids" form:"user_ids"`
	Mobile       string  `json:"mobile" form:"mobile"`
	GroupId      string  `json:"group_id" form:"group_id"`
	GroupIds     []int64 `json:"-" form:"-"`
	UserHIdArray []int64 `json:"-" form:"-"`
	CannotUse    int8    `json:"can_not_use" form:"can_not_use"`
}

func (r *ArgAdminUser) Default(c *base.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	r.UserHIdArray = []int64{}

	if r.UserHIds != "" {
		var uid int64
		tmp := strings.Split(r.UserHIds, ",")
		for _, item := range tmp {
			if item == "" {
				continue
			}
			uid, err = strconv.ParseInt(item, 10, 64)
			r.UserHIdArray = append(r.UserHIdArray, uid)
		}
	}
	r.GroupIds = []int64{}
	var idN int64
	if r.GroupId != "" {
		gIds := strings.Split(r.GroupId, ",")
		for _, id := range gIds {
			if idN, err = strconv.ParseInt(id, 10, 64); err != nil {
				err = fmt.Errorf("group_id格式不正确")
				return
			}
			r.GroupIds = append(r.GroupIds, idN)
		}
	}
	if r.CannotUse > 0 {
		if r.CannotUse > 100 {
			err = fmt.Errorf("can_not_use数据格式错误")
			return
		}
		mapCanUse, _ := models.SliceAdminUserCanNotUse.GetMapAsKeyUint8()
		if _, ok := mapCanUse[uint8(r.CannotUse)]; !ok {
			err = fmt.Errorf("can_not_use数据格式错误")
			return
		}
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
	CreatedAtString string               `json:"created_at"`
	Group           []AdminUserGroupName `json:"group"`
}

type ArgPermitMenu struct {
	app_param.RequestUser
	ArgGetImportByMenuIdSingle          // 通用参数逻辑处理 用于获取当前菜单下的接口列表
	ParentId                   int64    `json:"parent_id"` // 上级菜单ID
	PathType                   string   `json:"path_type" form:"path_type"`
	Module                     string   `json:"module" form:"module"` // 系统ID
	IsSuperAdmin               bool     `json:"-" form:"-"`           // 是否为超级管理员
	GroupId                    []int64  `json:"-" form:"-"`
	PathTypes                  []string `json:"-" form:"-"`
}

// Default 初始化默认值
func (r *ArgPermitMenu) Default(c *base.Context) (err error) {

	if err = r.InitRequestUser(c); err != nil {
		return
	}

	if err = r.ArgGetImportByMenuIdSingle.Default(c); err != nil {
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
	ResultPermitMenu
	IsSuperAdmin    bool   `json:"is_super_admin,omitempty"`     //是否超级管理员
	GoTo301         string `json:"goto301,omitempty"`            //登录成功后，如果当前没有指定访问界面，默认跳转的界面
	NotReadMsgCount int    `json:"not_read_msg_count,omitempty"` // 未读消息数

	// 当前选中的权限
	RoutParentMap map[string][]string `json:"routParentMap"`  // 当前菜单列表
	Menu          []ResultSystemMenu  `json:"menu,omitempty"` // 一级系统权限列表 用户从当前系统跳转到其他管理系统
	//OpList           map[string][]OpOne      `json:"-"`         // 获取接口权限列表
	NowImportAndMenu ResultGetImportByMenuId `json:"import_and_menu"` // 当前菜单下有的菜单和接口列表
}

func NewResultPermitMenuReturn() (res *ResultPermitMenuReturn) {
	res = &ResultPermitMenuReturn{
		ResultPermitMenu: ResultPermitMenu{
			Children: []ResultPermitMenu{},
		},
		RoutParentMap: map[string][]string{},
		Menu:          []ResultSystemMenu{},
		//OpList:          map[string][]OpOne{},
		NotReadMsgCount: 0,
		NowImportAndMenu: ResultGetImportByMenuId{
			ImportIds: []ImportSingle{},
			MenuIds:   []MenuSingle{},
		},
	}
	return
}

func (r *ArgFlag) Default(c *base.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	return
}

func (r *ArgAdminMenuSearch) Default() {

}

func (r *ArgGetMenu) Default(c *base.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	return
}
