// Package srv_impl
/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 12:20 上午
 */
package srv_impl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/base-wrapper/lib/plugins/rpc"
	"github.com/juetun/dashboard-api-main/basic/const_obj"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"gorm.io/gorm"
)

type PermitServiceImpl struct {
	base.ServiceBase
}

func NewPermitServiceImpl(context ...*base.Context) srvs.PermitService {
	p := &PermitServiceImpl{}
	p.SetContext(context...)
	return p
}

func (r *PermitServiceImpl) UpdateImportValue(arg *wrappers.ArgUpdateImportValue) (res *wrappers.ResultUpdateImportValue, err error) {
	res = &wrappers.ResultUpdateImportValue{}
	var condition = fmt.Sprintf("id IN (%s)", strings.Join(arg.Ids, ","))
	dao := dao_impl.NewPermitImportImpl(r.Context)
	var data = make(map[string]interface{}, 1)
	t := time.Now().Format("2006-01-02 15:04:05")
	data["updated_at"] = t
	switch arg.Column {
	case "is_del":
		data[arg.Column] = t
		if arg.Val == "1" {
			res.Result, err = dao.DeleteByCondition(condition)
		} else {
			res.Result, err = dao.UpdateByCondition(condition, data)
		}

	case "need_login", "default_open", "need_sign":
		if arg.Val != "1" && arg.Val != "2" {
			err = fmt.Errorf("您设置的值不正确")
			return
		}
		data[arg.Column] = arg.Val
		res.Result, err = dao.UpdateByCondition(condition, data)
	case "":
	default:
		err = fmt.Errorf("您选择的数据值(column:%s)不正确", arg.Column)
		return
	}
	if err != nil {
		return
	}
	return
}
func (r *PermitServiceImpl) GetMenu(arg *wrappers.ArgGetMenu) (res wrappers.ResultGetMenu, err error) {
	res = wrappers.ResultGetMenu{}
	if arg.MenuId == 0 {
		return
	}
	var li []models.AdminMenu
	if li, err = dao_impl.NewDaoPermit(r.Context).GetMenu(arg.MenuId); err != nil {
		return
	}
	if len(li) > 0 {
		res.AdminMenu = li[0]
	}

	return
}
func (r *PermitServiceImpl) AdminMenuSearch(arg *wrappers.ArgAdminMenu) (res wrappers.ResAdminMenuSearch, err error) {
	res = wrappers.ResAdminMenuSearch{
		List: []models.AdminMenu{},
	}
	dao := dao_impl.NewDaoPermit(r.Context)
	res.List, err = dao.GetAdminMenuList(arg)
	return
}
func (r *PermitServiceImpl) AdminUserGroupAdd(arg *wrappers.ArgAdminUserGroupAdd) (res wrappers.ResultAdminUserGroupAdd, err error) {
	res = wrappers.ResultAdminUserGroupAdd{}
	dao := dao_impl.NewDaoPermit(r.Context)

	var args = make([]map[string]interface{}, 0)
	for _, userHId := range arg.UserHIds {
		for _, groupId := range arg.GroupIds {
			args = append(args, map[string]interface{}{
				"group_id":   groupId,
				"user_hid":   userHId,
				"updated_at": time.Now().Format("2006-01-02 15:04:05"),
				"deleted_at": nil,
			})
		}
	}
	err = dao.AdminUserGroupAdd(args)
	if err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitServiceImpl) AdminUserGroupRelease(arg *wrappers.ArgAdminUserGroupRelease) (res wrappers.ResultAdminUserGroupRelease, err error) {
	res = wrappers.ResultAdminUserGroupRelease{}
	dao := dao_impl.NewDaoPermit(r.Context)
	err = dao.AdminUserGroupRelease(arg.IdString...)
	if err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitServiceImpl) AdminUserDelete(arg *wrappers.ArgAdminUserDelete) (res wrappers.ResultAdminUserDelete, err error) {
	res = wrappers.ResultAdminUserDelete{}
	dao := dao_impl.NewDaoPermit(r.Context)
	if err = dao.DeleteAdminUser(arg.IdString); err != nil {
		return
	}
	if err = dao.DeleteUserGroupByUserId(arg.IdString...); err != nil {
		return
	}
	res.Result = true
	return
}

func (r *PermitServiceImpl) AdminUserAdd(arg *wrappers.ArgAdminUserAdd) (res wrappers.ResultAdminUserAdd, err error) {
	res = wrappers.ResultAdminUserAdd{}
	if arg.UserHid == "" {
		err = fmt.Errorf("您没有选择要添加的用户")
		return
	}
	var user *models.UserMain
	if user, err = NewUserService(r.Context).
		GetUserById(strings.TrimSpace(arg.UserHid)); err != nil {
		return
	} else if user.UserHid == "" {
		err = fmt.Errorf("您要添加的用户信息不存在")
		return
	}

	if err = dao_impl.NewDaoPermit(r.Context).
		AdminUserAdd(&models.AdminUser{
			UserHid:  arg.UserHid,
			RealName: user.Name,
			Mobile:   user.Mobile,
			Model: gorm.Model{
				ID:        0,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeletedAt: gorm.DeletedAt{},
			},
		}); err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitServiceImpl) AdminGroupDelete(arg *wrappers.ArgAdminGroupDelete) (res wrappers.ResultAdminGroupDelete, err error) {
	res = wrappers.ResultAdminGroupDelete{}
	dao := dao_impl.NewDaoPermit(r.Context)

	if err = dao.DeleteAdminGroupByIds(arg.IdString...); err != nil {
		return
	}
	// 删除用户组权限
	if err = dao.DeleteUserGroupPermitByGroupId(arg.IdString...); err != nil {
		return
	}

	res.Result = true
	return
}

func (r *PermitServiceImpl) AdminGroupEdit(arg *wrappers.ArgAdminGroupEdit) (res *wrappers.ResultAdminGroupEdit, err error) {
	res = &wrappers.ResultAdminGroupEdit{}
	dao := dao_impl.NewDaoPermit(r.Context)
	var g models.AdminGroup
	var dta models.AdminGroup
	if g, err = dao.FetchByName(arg.Name); err != nil {
		return
	}

	if arg.Id == 0 {
		if g.Name != "" {
			err = fmt.Errorf("您输入的组名已存在")
			return
		}
		dta = models.AdminGroup{
			Name: arg.Name,
		}
		if err = dao.InsertAdminGroup(&dta); err != nil {
			return
		}
		res.Result = true
		return
	}

	if g.Name != "" && g.Id != arg.Id {
		err = fmt.Errorf("您输入的组名已存在")
		return
	}
	var dt []models.AdminGroup
	if dt, err = dao.GetAdminGroupByIds([]int{arg.Id}...); err != nil {
		return
	}
	if len(dt) == 0 {
		err = fmt.Errorf("您要编辑的组不存在或已删除")
		return
	}
	dta = dt[0]
	dta.Name = arg.Name
	dta.Id = arg.Id
	dta.UpdatedAt = base.TimeNormal{Time: time.Now()}
	if err = dao.UpdateAdminGroup(&dta); err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitServiceImpl) MenuAdd(arg *wrappers.ArgMenuAdd) (res *wrappers.ResultMenuAdd, err error) {
	res = &wrappers.ResultMenuAdd{}
	dao := dao_impl.NewDaoPermit(r.Context)
	t := time.Now()
	var list []models.AdminMenu
	if list, err = dao.GetByCondition(map[string]interface{}{
		"label":     arg.Label,
		"module":    arg.Module,
		"parent_id": arg.ParentId,
	}, nil, 0); err != nil {
		return
	} else if len(list) > 0 {
		err = fmt.Errorf("您输入的菜单名已存在")
		return
	}
	data := models.AdminMenu{
		PermitKey:          arg.PermitKey,
		Module:             arg.Module,
		ParentId:           arg.ParentId,
		Label:              arg.Label,
		Icon:               arg.Icon,
		ManageImportPermit: arg.ManageImportPermit,
		HideInMenu:         arg.HideInMenu,
		UrlPath:            arg.UrlPath,
		SortValue:          arg.SortValue,
		OtherValue:         arg.OtherValue,
		CreatedAt:          t,
		UpdatedAt:          t,
	}
	if err = dao.Add(&data); err != nil {
		return
	}
	// 如果是添加系统
	if arg.Module == wrappers.DefaultPermitModule {
		if err = r.addSystemDefaultMenu(dao, &data); err != nil {
			return
		}
	}

	res.Result = true
	return
}
func (r *PermitServiceImpl) addSystemDefaultMenu(dao daos.DaoPermit, data *models.AdminMenu) (err error) {
	if err = dao.Add(&models.AdminMenu{
		Module:             data.PermitKey,
		ParentId:           data.Id,
		Label:              "首页",
		Icon:               "md-home",
		ManageImportPermit: 1,
		HideInMenu:         0,
		UrlPath:            "",
		SortValue:          90000001,
		OtherValue:         "",
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
	}); err != nil {
		return
	}
	if err = dao.Add(&models.AdminMenu{
		Module:             data.PermitKey,
		ParentId:           data.Id,
		Label:              "公共接口",
		Icon:               "",
		ManageImportPermit: 1,
		HideInMenu:         1,
		UrlPath:            "",
		SortValue:          90000000,
		OtherValue:         "",
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
	}); err != nil {
		return
	}
	return
}
func (r *PermitServiceImpl) createImport(dao daos.DaoPermit, arg *wrappers.ArgEditImport) (res bool, err error) {
	t := time.Now()
	data := models.AdminImport{
		AppName:       arg.AppName,
		AppVersion:    arg.AppVersion,
		UrlPath:       arg.UrlPath,
		RequestMethod: arg.RequestMethod,
		SortValue:     arg.SortValue,
		UpdatedAt:     t,
		CreatedAt:     t,
	}
	if _, err = dao.CreateImport(&data); err != nil {
		return
	}
	res = true
	return
}
func (r *PermitServiceImpl) DeleteImport(arg *wrappers.ArgDeleteImport) (res *wrappers.ResultDeleteImport, err error) {
	res = &wrappers.ResultDeleteImport{}
	dao := dao_impl.NewDaoPermit(r.Context)
	if err = dao.DeleteImportByIds([]int{arg.ID}...); err != nil {
		return
	}
	if err = dao.DeleteUserGroupPermit(models.PathTypeApi, []int{arg.ID}...); err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitServiceImpl) editImportParam(arg *wrappers.ArgEditImport, value *models.AdminImport) (err error) {
	if arg.Id == 0 {
		if value.AppName == arg.AppName && value.UrlPath == arg.UrlPath {
			err = fmt.Errorf("您输入的接口信息已存在")
			return
		}
		return
	}
	if value.AppName == arg.AppName && value.UrlPath == arg.UrlPath {
		if arg.Id != value.Id {
			err = fmt.Errorf("您输入的接口信息已存在")
			return
		}
	}

	return
}

func (r *PermitServiceImpl) ImportList(arg *wrappers.ArgImportList) (res *wrappers.ResultImportList, err error) {
	var db *gorm.DB
	if arg.Order == "" {
		arg.Order = "id desc"
	}
	res = &wrappers.ResultImportList{Pager: response.NewPagerAndDefault(&arg.PageQuery)}
	dao := dao_impl.NewDaoPermit(r.Context)

	// 获取分页数据
	if err = res.Pager.CallGetPagerData(func(pager *response.Pager) (err error) {
		pager.TotalCount, db, err = dao.GetImportListCount(db, arg)
		return
	}, func(pager *response.Pager) (err error) {
		var list []models.AdminImport
		list, err = dao.GetImportListData(db, arg, pager)
		pager.List, err = r.orgImportList(dao, list)
		return
	}); err != nil {
		return
	}
	return
}
func (r *PermitServiceImpl) getImportId(l int, list []models.AdminImport) (importId []int) {
	importId = make([]int, 0, l)
	for _, value := range list {
		importId = append(importId, value.Id)
	}
	return
}
func (r *PermitServiceImpl) getImportMenuGroupMap(dao daos.DaoPermit, list []models.AdminMenuImport) (mapAdminMenuModule map[string]models.AdminMenu, mapAdminMenu map[int]models.AdminMenu, err error) {
	ll := len(list)
	menuIds := make([]int, 0, ll)
	for _, value := range list {
		menuIds = append(menuIds, value.MenuId)
	}
	var adminMenu []models.AdminMenu
	if adminMenu, err = dao.GetMenu(menuIds...); err != nil {
		return
	}
	mapAdminMenu = make(map[int]models.AdminMenu, len(adminMenu))
	var (
		m          = make(map[int]int, len(adminMenu))
		modules    = make([]string, 0, len(adminMenu))
		modulesMap = make(map[string]string, len(adminMenu))
	)
	for _, value := range adminMenu {
		if _, ok := modulesMap[value.Module]; !ok {
			modules = append(modules, value.Module)
		}
	}

	var dta []models.AdminMenu
	if dta, err = dao.GetMenuByPermitKey(modules...); err != nil {
		return
	}

	mapAdminMenuModule = make(map[string]models.AdminMenu, len(dta))
	for _, value := range dta {
		mapAdminMenuModule[value.PermitKey] = value
	}

	for _, value := range adminMenu {
		if _, ok := m[value.Id]; !ok {
			mapAdminMenu[value.Id] = value
			m[value.Id] = value.Id
		}
	}

	return
}
func (r *PermitServiceImpl) getImportMenuGroup(dao daos.DaoPermit, l int, data []models.AdminImport) (res map[int][]wrappers.AdminImportListMenu, err error) {
	importId := r.getImportId(l, data)
	daoImportMenu := dao_impl.NewPermitImportImpl(r.Context)

	var list []models.AdminMenuImport
	var mapAdminMenu map[int]models.AdminMenu
	var mapAdminMenuGroup map[string]models.AdminMenu

	if list, err = daoImportMenu.GetImportMenuByImportIds(importId...); err != nil {
		return
	} else if mapAdminMenuGroup, mapAdminMenu, err = r.getImportMenuGroupMap(dao, list); err != nil {
		return
	}
	res = make(map[int][]wrappers.AdminImportListMenu, l)

	var (
		dt  wrappers.AdminImportListMenu
		dtm models.AdminMenu
		ok  bool
		ll  = len(list)
	)
	var dtt models.AdminMenu
	for _, value := range list {
		if _, ok := res[value.ImportId]; !ok {
			res[value.ImportId] = make([]wrappers.AdminImportListMenu, 0, ll)
		}
		dt = wrappers.AdminImportListMenu{
			ImportId: value.Id,
			MenuId:   value.MenuId,
		}

		if dtm, ok = mapAdminMenu[value.MenuId]; ok {
			dt.MenuName = dtm.Label
			dt.Id = dtm.Id
			if dtt, ok = mapAdminMenuGroup[dtm.Module]; ok {
				dt.SystemName = dtt.Label
				dt.SystemModuleId = dtt.Id
				dt.SystemMenuKey = dtt.PermitKey
				dt.SystemIcon = dtt.Icon
			}
		}
		res[value.ImportId] = append(res[value.ImportId], dt)
	}
	return
}
func (r *PermitServiceImpl) orgImportList(dao daos.DaoPermit, list []models.AdminImport) (res []wrappers.AdminImportList, err error) {
	l := len(list)
	res = make([]wrappers.AdminImportList, 0, l)
	var dt wrappers.AdminImportList
	var dta map[int][]wrappers.AdminImportListMenu
	if dta, err = r.getImportMenuGroup(dao, l, list); err != nil {
		return
	}
	for _, value := range list {
		dt = wrappers.AdminImportList{AdminImport: value, Menu: []wrappers.AdminImportListMenu{}}
		if _, ok := dta[value.Id]; ok {
			dt.Menu = dta[value.Id]
		}
		res = append(res, dt)
	}
	return
}

func (r *PermitServiceImpl) EditImport(arg *wrappers.ArgEditImport) (res *wrappers.ResultEditImport, err error) {
	res = &wrappers.ResultEditImport{Result: false}
	var (
		dao        = dao_impl.NewDaoPermit(r.Context)
		daoService = dao_impl.NewDaoServiceImpl(r.Context)
		listImport []models.AdminImport
		apps       []models.AdminApp
	)

	if apps, err = daoService.GetByKeys(strings.TrimSpace(arg.AppName)); err != nil {
		return
	}
	if len(apps) == 0 {
		err = fmt.Errorf("您输入的应用(%s)不存在或已删除", arg.AppName)
		return
	}
	if listImport, err = dao.GetImportByCondition(map[string]interface{}{"app_name": arg.AppName, "url_path": arg.UrlPath, "request_method": arg.RequestMethod}); err != nil {
		return
	} else {
		// 验证数据是否重复
		for _, value := range listImport {
			if err = r.editImportParam(arg, &value); err != nil {
				return
			}
		}
	}
	data := map[string]interface{}{
		`app_name`:       arg.AppName,
		`app_version`:    arg.AppVersion,
		`url_path`:       arg.UrlPath,
		`request_method`: arg.RequestMethod,
		`sort_value`:     arg.SortValue,
		`updated_at`:     arg.RequestTime,
		"need_login":     arg.NeedLogin,
		"need_sign":      arg.NeedSign,
	}
	if arg.Id == 0 { // 如果是添加接口
		res.Result, err = r.createImport(dao, arg)
		return
	}

	var m = models.AdminImport{Id: arg.Id}
	var dt []models.AdminImport
	if dt, err = dao.GetAdminImportById(arg.Id); err != nil {
		return
	}

	if len(dt) == 0 {
		err = fmt.Errorf("您编辑的接口信息不存在或已删除")
		return
	}

	if dt[0].PermitKey == "" {
		data["permit_key"] = m.GetPathName()
	}
	// 如果更新了app_name
	if dt[0].AppName != arg.AppName {
		if err = dao_impl.NewPermitImportImpl(r.Context).
			UpdateMenuImport(fmt.Sprintf("import_id = %d ", dt[0].Id),
				map[string]interface{}{"import_app_name": arg.AppName}); err != nil {
			return
		}

	}
	if _, err = dao.UpdateAdminImport(map[string]interface{}{"id": arg.Id}, data); err != nil {
		return
	}
	res.Result = true
	return
}

func (r *PermitServiceImpl) MenuImport(arg *wrappers.ArgMenuImport) (res *wrappers.ResultMenuImport, err error) {
	res = &wrappers.ResultMenuImport{
		Pager: response.NewPager(response.PagerBaseQuery(arg.PageQuery)),
	}
	dao := dao_impl.NewDaoPermit(r.Context)
	var db *gorm.DB
	if db, err = dao.MenuImportCount(arg, &res.TotalCount); err != nil {
		return
	}
	if res.TotalCount == 0 {
		return
	}
	if res.List, err = dao.MenuImportList(db, arg); err != nil {
		return
	}
	return
}
func (r *PermitServiceImpl) GetImport(arg *wrappers.ArgGetImport) (res *wrappers.ResultGetImport, err error) {
	res = &wrappers.ResultGetImport{
		Pager: response.NewPager(response.PagerBaseQuery(arg.PageQuery)),
	}
	dao := dao_impl.NewDaoPermit(r.Context)
	var db *gorm.DB
	if db, err = dao.GetImportCount(arg, &res.TotalCount); err != nil {
		return
	}
	if res.TotalCount == 0 {
		return
	}
	var list []models.AdminImport
	if list, err = dao.GetImportList(db, arg); err != nil {
		return
	}
	if !arg.Checked {
		res.List = list
		return
	}
	res.List, err = r.joinChecked(dao, arg, list)

	// []models.AdminImport{}
	return
}
func (r *PermitServiceImpl) joinChecked(dao daos.DaoPermit, arg *wrappers.ArgGetImport, data []models.AdminImport) (res []wrappers.AdminImport, err error) {
	res = make([]wrappers.AdminImport, 0, len(data))
	var dt wrappers.AdminImport
	var importId = make([]int, 0, len(data))
	for _, value := range data {
		importId = append(importId, value.Id)
	}
	var li []models.AdminUserGroupPermit
	if li, err = dao.GetSelectImportByImportId(arg.GroupId, importId...); err != nil {
		return
	}
	var m = make(map[int]int, len(li))
	for _, it := range li {
		m[it.MenuId] = it.MenuId
	}
	for _, value := range data {
		dt = wrappers.AdminImport{
			AdminImport: value,
		}
		if _, ok := m[value.Id]; ok {
			dt.Checked = true
		}
		res = append(res, dt)
	}
	return
}
func (r *PermitServiceImpl) MenuDelete(arg *wrappers.ArgMenuDelete) (res *wrappers.ResultMenuDelete, err error) {
	res = &wrappers.ResultMenuDelete{}
	dao := dao_impl.NewDaoPermit(r.Context)

	if err = dao.DeleteMenuByIds(arg.IdValue...); err != nil {
		return
	}
	if err = dao.DeleteUserGroupPermit(models.PathTypePage, arg.IdValueNumber...); err != nil {
		return
	}

	// 删除菜单下的所有接口权限
	var importList []models.AdminImport
	importList, err = dao.GetImportMenuId(arg.IdValueNumber...)
	iIds := make([]int, 0, len(importList))
	for _, value := range importList {
		iIds = append(iIds, value.Id)
	}
	if err = dao.DeleteUserGroupPermit(models.PathTypeApi, iIds...); err != nil {
		return
	}
	res.Result = true
	return
}

func (r *PermitServiceImpl) MenuSave(arg *wrappers.ArgMenuSave) (res *wrappers.ResultMenuSave, err error) {
	res = &wrappers.ResultMenuSave{}
	dao := dao_impl.NewDaoPermit(r.Context)

	var list []models.AdminMenu
	if list, err = dao.GetByCondition(map[string]interface{}{
		"label":     arg.Label,
		"module":    arg.Module,
		"parent_id": arg.ParentId,
	}, nil, 0); err != nil {
		return
	} else if len(list) > 0 && arg.Id != list[0].Id {
		err = fmt.Errorf("您输入的菜单名已存在")
		return
	}
	if arg.Id > 0 {
		var menu models.AdminMenu
		var menus []models.AdminMenu
		if menus, err = dao.GetMenu(arg.Id); err != nil {
			return
		} else if len(menus) > 0 {
			menu = menus[0]
		}

		if arg.PermitKey != menu.PermitKey {
			if err = dao.UpdateMenuByCondition(map[string]interface{}{"module": menu.PermitKey}, map[string]interface{}{"module": arg.PermitKey}); err != nil {
				return
			}
		} else if arg.Module != menu.Module {
			// 更新子菜单的 module
			if err = r.updateChildModule(dao, menu.Id, arg.Module); err != nil {
				return
			}

		}

	}
	t := time.Now()
	var m = models.AdminMenu{
		Module:             arg.Module,
		PermitKey:          arg.PermitKey,
		ParentId:           arg.ParentId,
		Label:              arg.Label,
		Icon:               arg.Icon,
		HideInMenu:         arg.HideInMenu,
		Domain:             arg.Domain,
		ManageImportPermit: arg.ManageImportPermit,
		UrlPath:            arg.UrlPath,
		SortValue:          arg.SortValue,
		OtherValue:         arg.OtherValue,
		UpdatedAt:          t,
	}
	err = dao.Save(arg.Id, m.ToMapStringInterface())
	res.Result = true
	return
}

func (r *PermitServiceImpl) getChildIds(dao daos.DaoPermit, parentId []string, ids *[]string) (err error) {
	if len(parentId) == 0 {
		return
	}
	var li []models.AdminMenu
	if li, err = dao.GetMenuByCondition(fmt.Sprintf("parent_id IN (%s)", strings.Join(parentId, ","))); err != nil {
		return
	}
	pIds := make([]string, 0, len(li))
	for _, value := range li {
		pIds = append(pIds, strconv.Itoa(value.Id))
		*ids = append(*ids, strconv.Itoa(value.Id))
	}

	if err = r.getChildIds(dao, pIds, ids); err != nil {
		return
	}
	return
}
func (r *PermitServiceImpl) updateChildModule(dao daos.DaoPermit, parentId int, module string) (err error) {

	ids := make([]string, 0, 20)
	ids = append(ids, strconv.Itoa(parentId))
	if err = r.getChildIds(dao, []string{strconv.Itoa(parentId)}, &ids); err != nil {
		return
	}
	if err = dao.UpdateMenuByCondition(
		fmt.Sprintf("id IN(%s)", strings.Join(ids, ",")),
		map[string]interface{}{"module": module},
	); err != nil {
		return
	}

	ids = append(ids, strconv.Itoa(parentId))
	// 更新菜单接口关系表的menu_module
	if err = dao_impl.NewPermitImportImpl(r.Context).
		UpdateMenuImport(fmt.Sprintf("menu_id IN('%s')", strings.Join(ids, "','")),
			map[string]interface{}{"module": module}); err != nil {
		return
	}

	return
}

func (r *PermitServiceImpl) AdminMenuWithCheck(arg *wrappers.ArgAdminMenuWithCheck) (res *wrappers.ResultMenuWithCheck, err error) {
	res = &wrappers.ResultMenuWithCheck{
		List: make([]wrappers.AdminMenuObject, 0, 20),
		Menu: make([]wrappers.ResultSystemAdminMenu, 0, 30),
	}
	dao := dao_impl.NewDaoPermit(r.Context)
	var list []models.AdminMenu
	if list, err = dao.GetAdminMenuList(&arg.ArgAdminMenu); err != nil {
		return
	}
	if arg.SystemId, err = r.getSystemIdByModule(dao, arg.Module, arg.SystemId); err != nil {
		return
	}
	arg.SystemId, arg.Module = r.permitTab(list, &res.Menu, arg.SystemId, arg.Module)
	var mapPermitMenu map[int]int
	if mapPermitMenu, err = r.getGroupPermitMenu(dao, arg.Module, arg.GroupId); err != nil {
		return
	}

	r.orgTree(list, arg.SystemId, &res.List, mapPermitMenu)

	return
}

func (r *PermitServiceImpl) getGroupPermitMenu(dao daos.DaoPermit, module string, groupId int) (mapPermitMenu map[int]int, err error) {
	var groupPermit []models.AdminUserGroupPermit
	if groupPermit, err = dao.GetMenuIdsByPermitByGroupIds(module, []string{models.PathTypePage}, []int{groupId}...); err != nil {
		return
	} else {
		mapPermitMenu = make(map[int]int, len(groupPermit))
		for _, it := range groupPermit {
			mapPermitMenu[it.MenuId] = it.MenuId
		}
	}
	return
}

func (r *PermitServiceImpl) getSystemIdByModule(dao daos.DaoPermit, module string, sysId int) (systemId int, err error) {
	systemId = sysId
	// 如果选择了系统模块则，直接初始化systemId
	if module != "" {
		var dt []models.AdminMenu
		if dt, err = dao.GetMenuByPermitKey(module); err != nil {
			return
		}
		if len(dt) == 0 {
			err = fmt.Errorf("您操作的系统信息不存在或已删除")
			return
		}
		systemId = dt[0].Id
	}
	return
}
func (r *PermitServiceImpl) AdminMenu(arg *wrappers.ArgAdminMenu) (res *wrappers.ResultAdminMenu, err error) {

	res = &wrappers.ResultAdminMenu{
		List: make([]wrappers.AdminMenuObject, 0, 20),
		Menu: make([]wrappers.ResultSystemAdminMenu, 0, 30),
	}
	dao := dao_impl.NewDaoPermit(r.Context)
	if arg.SystemId, err = r.getSystemIdByModule(dao, arg.Module, arg.SystemId); err != nil {
		return
	}
	var list, dt2 []models.AdminMenu
	if list, err = dao.GetAdminMenuList(arg); err != nil {
		return
	}

	if dt2, err = dao.GetAdminMenuList(&wrappers.ArgAdminMenu{}); err != nil {
		return
	} else {
		r.permitTab(dt2, &res.Menu, arg.SystemId, arg.Module)
	}

	arg.SystemId, arg.Module = r.permitTab(list, &res.Menu, arg.SystemId, arg.Module)
	r.orgTree(list, arg.SystemId, &res.List, nil)
	return
}
func (r *PermitServiceImpl) permitTab(list []models.AdminMenu, menu *[]wrappers.ResultSystemAdminMenu, systemId int, moduleSrc string) (sid int, module string) {
	module = moduleSrc
	var data wrappers.ResultSystemAdminMenu
	var ind int
	for _, item := range list {
		if item.ParentId != wrappers.DefaultPermitParentId {
			continue
		}
		data = wrappers.ResultSystemAdminMenu{
			Id:        item.Id,
			PermitKey: item.PermitKey,
			Label:     item.Label,
			Icon:      item.Icon,
			SortValue: item.SortValue,
			Module:    item.Module,
			Domain:    item.Domain,
		}
		if systemId == 0 && ind == 0 {
			data.Active = true
			systemId = item.Id
			if moduleSrc == "" {
				module = item.PermitKey
			}
		} else if systemId > 0 && item.Id == systemId {
			data.Active = true
			if moduleSrc == "" {
				module = item.PermitKey
			}
		}
		*menu = append(*menu, data)
		ind++
	}
	sid = systemId

	return
}

//
func (r *PermitServiceImpl) orgTree(list []models.AdminMenu, parentId int, res *[]wrappers.AdminMenuObject, mapPermitMenu map[int]int) () {
	var tmp wrappers.AdminMenuObject
	for _, value := range list {
		// 剔除默认数据那条
		if value.Id == parentId || value.ParentId != parentId {
			continue
		}
		tmp = r.orgAdminMenuObject(&value, mapPermitMenu)
		r.orgTree(list, value.Id, &tmp.Children, mapPermitMenu)
		*res = append(*res, tmp)
	}
}

func (r *PermitServiceImpl) orgAdminMenuObject(value *models.AdminMenu, mapPermitMenu map[int]int) (res wrappers.AdminMenuObject) {
	res = wrappers.AdminMenuObject{Children: make([]wrappers.AdminMenuObject, 0, 20)}
	res.ResultAdminMenuSingle = wrappers.ResultAdminMenuSingle{
		Id:                 value.Id,
		ParentId:           value.ParentId,
		Title:              value.Label,
		Label:              value.Label,
		Icon:               value.Icon,
		HideInMenu:         value.HideInMenu,
		UrlPath:            value.UrlPath,
		SortValue:          value.SortValue,
		Module:             value.Module,
		PermitKey:          value.PermitKey,
		Domain:             value.Domain,
		ManageImportPermit: value.ManageImportPermit,
	}
	if value.OtherValue != "" {
		_ = json.Unmarshal([]byte(value.OtherValue), &res.ResultAdminMenuSingle.ResultAdminMenuOtherValue)
	} else {
		res.ResultAdminMenuSingle.ResultAdminMenuOtherValue = wrappers.ResultAdminMenuOtherValue{Expand: true}
	}
	if mapPermitMenu != nil {
		if _, ok := mapPermitMenu[value.Id]; ok {
			res.ResultAdminMenuSingle.Checked = true
		}
	}
	return
}
func (r *PermitServiceImpl) getGroupHavePermit(arg *wrappers.ArgAdminSetPermit, dao daos.DaoPermit) (newPermit []int, notPermitId []int, err error) {
	listSelectMenu, err := dao.GetMenuIdsByPermitByGroupIds(arg.Module, []string{arg.Type}, arg.GroupId)
	if err != nil {
		return
	}
	commonPermit := make([]int, 0, len(listSelectMenu)) // 当前已经选中的菜单
	notPermitId = make([]int, 0, len(listSelectMenu))   // 当前取消权限的菜单
	newPermit = make([]int, 0, len(listSelectMenu))     // 新增的菜单
	for _, item := range listSelectMenu {
		if r.inSlice(item.MenuId, arg.PermitIds) {
			commonPermit = append(commonPermit, item.MenuId)
		} else {
			notPermitId = append(notPermitId, item.MenuId)
		}
	}
	for _, it := range arg.PermitIds {
		if !r.inSlice(it, commonPermit) {
			newPermit = append(newPermit, it)
		}
	}
	return
}
func (r *PermitServiceImpl) inSlice(id int, slice []int) (res bool) {
	for _, value := range slice {
		if id == value {
			res = true
			return
		}
	}
	return
}

func (r *PermitServiceImpl) isHomePage(dao daos.DaoPermit, permitIds []int) (res bool, homePageId int, module string, err error) {
	var li []models.AdminMenu
	if li, err = dao.GetMenu(permitIds...); err != nil {
		return
	}
	for _, it := range li {
		if module == "" {
			module = it.Module
		}
		if it.Label == "首页" {
			res = true
			module = it.Module
			homePageId = it.Id
			return
		}
	}
	return
}
func (r *PermitServiceImpl) commonImport(dao daos.DaoPermit, module string) (res []int, err error) {
	res = make([]int, 0, 2)
	var da []models.AdminMenu
	da, err = dao.GetMenuByCondition(map[string]interface{}{
		"module": module,
		"label":  "公共接口",
	})
	for _, value := range da {
		res = append(res, value.Id)
	}
	return
}
func (r *PermitServiceImpl) orgNeedMenu(dao daos.DaoPermit, arg *wrappers.ArgAdminSetPermit) (permitIds []int, err error) {
	var (
		isHomePage bool
		homePageId int
	)
	permitIds = make([]int, 0, 2)
	// 判断是否为系统首页，如果是首页，则自动绑定公共接口隐藏界面
	if isHomePage, homePageId, arg.Module, err = r.isHomePage(dao, arg.PermitIds); err != nil {
		return
	}
	if permitIds, err = r.commonImport(dao, arg.Module); err != nil {
		return
	}
	if isHomePage { // 如果设置的有首页权限，则有公共权限
		arg.PermitIds = append(arg.PermitIds, permitIds...)
		return
	}

	if homePageId == 0 {
		return
	}
	// 如果没有首页设置，则删除公共接口隐藏界面（隐藏界面只是用于存储公共接口权限使用，实际界面不存在）
	permitIds = append(permitIds, homePageId)
	return
}
func (r *PermitServiceImpl) deleteGroupPermitByGroupId(dao daos.DaoPermit, arg *wrappers.ArgAdminSetPermit) (err error) {
	if len(arg.PermitIds) > 0 {
		return
	}
	// 没有权限ID，则说明清除所有的权限
	if err = dao.DeleteGroupPermitByGroupId(arg.GroupId); err != nil {
		return
	}
	return
}
func (r *PermitServiceImpl) setMenuPermit(dao daos.DaoPermit, arg *wrappers.ArgAdminSetPermit) (err error) {
	var (
		newPermit   []int
		notPermitId []int
		permitIds   []int
	)

	// permitId 为空处理逻辑
	if err = r.deleteGroupPermitByGroupId(dao, arg); err != nil {
		return
	}

	if permitIds, err = r.orgNeedMenu(dao, arg); err != nil {
		return
	}

	if newPermit, notPermitId, err = r.getGroupHavePermit(arg, dao); err != nil {
		return
	}
	if err = r.deleteNotMenuPermitId(dao, append(notPermitId, permitIds...), arg); err != nil {
		return
	}
	if err = r.addNewMenuPermit(dao, newPermit, arg); err != nil {
		return
	}
	return
}
func (r *PermitServiceImpl) AdminSetPermit(arg *wrappers.ArgAdminSetPermit) (res *wrappers.ResultAdminSetPermit, err error) {
	res = &wrappers.ResultAdminSetPermit{}

	dao := dao_impl.NewDaoPermit(r.Context)
	switch arg.Type {
	case models.PathTypePage: // 设置菜单权限
		if err = r.setMenuPermit(dao, arg); err != nil {
			return
		}
	case models.PathTypeApi: // 设置API权限
		if err = r.setApiPermit(dao, arg); err != nil {
			return
		}
	}
	res.Result = true
	return
}
func (r *PermitServiceImpl) setApiPermit(dao daos.DaoPermit, arg *wrappers.ArgAdminSetPermit) (err error) {
	switch arg.Act {
	case models.SetPermitAdd:
		var dt models.AdminUserGroupPermit
		var list []models.AdminUserGroupPermit
		var t = time.Now()
		for _, pid := range arg.PermitIds {
			dt = models.AdminUserGroupPermit{
				GroupId:   arg.GroupId,
				MenuId:    pid,
				PathType:  models.PathTypeApi,
				CreatedAt: t,
				UpdatedAt: t,
				DeletedAt: nil,
			}
			list = append(list, dt)
		}
		var m models.AdminUserGroupPermit
		if err = dao.BatchGroupPermit(m.TableName(), list); err != nil {
			err = fmt.Errorf("操作异常")
			return
		}

	case models.SetPermitCancel:
		if err = dao.DeleteGroupPermit(arg.GroupId, models.PathTypeApi, arg.PermitIds...); err != nil {
			return
		}
	default:
		err = fmt.Errorf("act格式错误")
		return
	}
	return
}
func (r *PermitServiceImpl) setApiPermitOld(dao daos.DaoPermit, arg *wrappers.ArgAdminSetPermit) (err error) {
	var (
		newPermit   []int
		notPermitId []int
	)
	if newPermit, notPermitId, err = r.getGroupHavePermit(arg, dao); err != nil {
		return
	}
	if err = r.deleteNotApiPermitId(dao, notPermitId, arg); err != nil {
		return
	}
	if err = r.addNewApiPermit(dao, newPermit, arg.GroupId); err != nil {
		return
	}
	return
}

func (r *PermitServiceImpl) deleteNotMenuPermitId(dao daos.DaoPermit, notPermitId []int, args *wrappers.ArgAdminSetPermit) (err error) {
	if len(notPermitId) == 0 {
		return
	}
	var listImport []models.AdminImport

	// notPermitId = append(notPermitId, args.PermitIds...)
	// 获取菜单下的接口
	var notPermitImportId []int
	if listImport, err = dao.GetDefaultImportByMenuIds(args.Type, args.Module, notPermitId...); err != nil {
		return
	}
	if len(listImport) > 0 {
		notPermitImportId = make([]int, 0, len(listImport))
		for _, it := range listImport {
			notPermitImportId = append(notPermitImportId, it.Id)
		}
	}

	if err = dao.DeleteGroupPermitByMenuIds(args.GroupId, args.Module, notPermitId, notPermitImportId); err != nil {
		return
	}
	return
}
func (r *PermitServiceImpl) deleteNotApiPermitId(dao daos.DaoPermit, notPermitId []int, arg *wrappers.ArgAdminSetPermit) (err error) {
	var listImport []models.AdminImport

	// 获取菜单下的接口
	var notPermitImportId []int
	if listImport, err = dao.GetDefaultImportByMenuIds(arg.Type, arg.Module, notPermitId...); err != nil {
		return
	} else if len(listImport) > 0 {
		notPermitImportId = make([]int, 0, len(listImport))
		for _, it := range listImport {
			notPermitImportId = append(notPermitImportId, it.Id)
		}
	}

	if err = dao.DeleteGroupPermitByMenuIds(arg.GroupId, arg.Module, notPermitId, notPermitImportId); err != nil {
		return
	}
	return
}
func (r *PermitServiceImpl) addNewMenuPermit(dao daos.DaoPermit, newPermit []int, arg *wrappers.ArgAdminSetPermit) (err error) {
	l := len(newPermit)
	if l == 0 {
		return
	}
	list := make([]models.AdminUserGroupPermit, 0, l)

	var (
		dt         models.AdminUserGroupPermit
		t          = time.Now()
		listImport []models.AdminImport
	)

	if listImport, err = dao.GetDefaultOpenImportByMenuIds(newPermit...); err != nil {
		return
	}

	for _, importData := range listImport {
		dt = models.AdminUserGroupPermit{
			Module:    arg.Module,
			GroupId:   arg.GroupId,
			MenuId:    importData.Id,
			PathType:  models.PathTypeApi,
			CreatedAt: t,
			UpdatedAt: t,
		}
		list = append(list, dt)
	}

	for _, pid := range newPermit {
		dt = models.AdminUserGroupPermit{
			Module:    arg.Module,
			GroupId:   arg.GroupId,
			MenuId:    pid,
			PathType:  models.PathTypePage,
			CreatedAt: t,
			UpdatedAt: t,
		}
		list = append(list, dt)
	}
	var m models.AdminUserGroupPermit
	if err = dao.BatchGroupPermit(m.TableName(), list); err != nil {
		err = fmt.Errorf("操作异常")
		return
	}
	return
}
func (r *PermitServiceImpl) addNewApiPermit(dao daos.DaoPermit, newPermit []int, groupId int) (err error) {
	if len(newPermit) == 0 {
		return
	}
	list := make([]models.AdminUserGroupPermit, 0, len(newPermit))
	var dt models.AdminUserGroupPermit
	var t = time.Now()
	for _, importId := range newPermit {
		dt = models.AdminUserGroupPermit{
			GroupId:   groupId,
			MenuId:    importId,
			PathType:  models.PathTypeApi,
			CreatedAt: t,
			UpdatedAt: t,
		}
		list = append(list, dt)
	}

	var m models.AdminUserGroupPermit
	if err = dao.BatchGroupPermit(m.TableName(), list); err != nil {
		err = fmt.Errorf("操作异常")
		return
	}
	return
}

func (r *PermitServiceImpl) AdminGroup(arg *wrappers.ArgAdminGroup) (res *wrappers.ResultAdminGroup, err error) {

	res = &wrappers.ResultAdminGroup{Pager: *response.NewPagerAndDefault(&arg.PageQuery)}

	var db *gorm.DB
	dao := dao_impl.NewDaoPermit(r.Context)
	// 获取分页数据
	if err = res.Pager.CallGetPagerData(func(pagerObject *response.Pager) (err error) {
		pagerObject.TotalCount, db, err = dao.GetAdminGroupCount(db, arg)
		return
	}, func(pagerObject *response.Pager) (err error) {
		var list []models.AdminGroup
		list, err = dao.GetAdminGroupList(db, arg, pagerObject)
		pagerObject.List, err = r.orgGroupList(dao, list)
		return
	}); err != nil {
		return
	}
	return
}

func (r *PermitServiceImpl) orgGroupList(dao daos.DaoPermit, list []models.AdminGroup) (res []wrappers.AdminGroup, err error) {
	l := len(list)
	res = make([]wrappers.AdminGroup, 0, l)
	ids := make([]int, 0, l)
	for _, value := range list {
		ids = append(ids, value.ParentId)
	}
	var listG []models.AdminGroup
	if listG, err = dao.GetAdminGroupByIds(ids...); err != nil {
		return
	}
	var m = make(map[int]models.AdminGroup, l)
	for _, value := range listG {
		m[value.Id] = value
	}

	var dt wrappers.AdminGroup
	for _, it := range list {
		dt = wrappers.AdminGroup{
			AdminGroup: it,
		}
		if dta, ok := m[it.ParentId]; ok {
			dt.ParentName = dta.Name
		}
		res = append(res, dt)
	}

	return
}
func (r *PermitServiceImpl) AdminUser(arg *wrappers.ArgAdminUser) (res *wrappers.ResultAdminUser, err1 error) {

	res = &wrappers.ResultAdminUser{Pager: *response.NewPagerAndDefault(&arg.PageQuery)}

	dao := dao_impl.NewDaoPermit(r.Context)

	var db *gorm.DB

	// 分页获取数据
	err1 = res.Pager.CallGetPagerData(func(pagerObject *response.Pager) (err error) {
		pagerObject.TotalCount, db, err = dao.GetAdminUserCount(db, arg)
		return
	}, func(pagerObject *response.Pager) (err error) {
		var list []models.AdminUser
		if list, err = dao.GetAdminUserList(db, arg, pagerObject); err != nil {
			return
		}
		pagerObject.List, err = r.leftAdminUser(list, dao)
		return
	})
	return
}

func (r *PermitServiceImpl) getUserGroup(list []models.AdminUser, dao daos.DaoPermit) (res map[string][]wrappers.AdminUserGroupName, err error) {
	res = map[string][]wrappers.AdminUserGroupName{}
	uIds := make([]string, 0, len(list))
	for _, value := range list {
		uIds = append(uIds, value.UserHid)
	}
	listResult, err := dao.GetUserGroupByUIds(uIds...)
	if err != nil {
		return
	}
	var tmp wrappers.AdminUserGroupName
	gIds := make([]int, 0, len(listResult))
	gIdsMap := make(map[int]int, len(listResult))
	for _, item := range listResult {
		if _, ok := gIdsMap[item.GroupId]; !ok {
			gIdsMap[item.GroupId] = item.GroupId
			gIds = append(gIds, item.GroupId)
		}
		if _, ok := res[item.UserHid]; !ok {
			res[item.UserHid] = []wrappers.AdminUserGroupName{}
		}
	}

	var groupList []models.AdminGroup
	if groupList, err = dao.GetAdminGroupByIds(gIds...); err != nil {
		return
	}

	groupMap := make(map[int]models.AdminGroup, len(groupList))
	for _, value := range groupList {
		groupMap[value.Id] = value
	}

	for _, item := range listResult {
		tmp = wrappers.AdminUserGroupName{
			AdminUserGroup: item,
		}
		if _, ok := groupMap[item.GroupId]; ok {
			tmp.GroupName = groupMap[item.GroupId].Name
			tmp.IsSuperAdmin = groupMap[item.GroupId].IsSuperAdmin
			tmp.IsAdminGroup = groupMap[item.GroupId].IsAdminGroup
		}
		res[item.UserHid] = append(res[item.UserHid], tmp)
	}
	return
}

func (r *PermitServiceImpl) leftAdminUser(list []models.AdminUser, dao daos.DaoPermit) (res []wrappers.ResultAdminUserList, err error) {
	res = []wrappers.ResultAdminUserList{}
	mapGroupPermit, err := r.getUserGroup(list, dao)
	if err != nil {
		return
	}
	var tmp wrappers.ResultAdminUserList
	for _, item := range list {
		tmp = wrappers.ResultAdminUserList{
			AdminUser: item,
		}
		if _, ok := mapGroupPermit[item.UserHid]; ok {
			tmp.Group = mapGroupPermit[item.UserHid]

		}
		res = append(res, tmp)
	}

	return

}
func (r *PermitServiceImpl) Menu(arg *wrappers.ArgPermitMenu) (res *wrappers.ResultPermitMenuReturn, err error) {
	res = &wrappers.ResultPermitMenuReturn{
		ResultPermitMenu: wrappers.ResultPermitMenu{
			Children: []wrappers.ResultPermitMenu{},
		},
		RoutParentMap:   map[string][]string{},
		Menu:            []wrappers.ResultSystemMenu{},
		OpList:          map[string][]wrappers.OpOne{},
		NotReadMsgCount: 0,
	}

	var syncG sync.WaitGroup
	syncG.Add(2)
	go func() { // 获取用户未读消息数
		defer syncG.Done()
		res.NotReadMsgCount, _ = r.getMessageCount(arg.UserId)
	}()
	go func() {
		defer syncG.Done()
		dao := dao_impl.NewDaoPermit(r.Context)
		if arg.GroupId, arg.IsSuperAdmin, err = r.getUserGroupIds(&ArgGetUserGroupIds{Dao: dao, UserId: arg.UserId}); err != nil {
			return
		}

		if arg.Module != "" {
			var dt []models.AdminMenu
			if dt, err = dao.GetMenuByCondition(map[string]interface{}{"permit_key": arg.Module}); err != nil {
				return
			}
			if len(dt) == 0 {
				err = fmt.Errorf("您查看的系统(%s)不存在或已删除", arg.Module)
				return
			}
			arg.ParentId = dt[0].Id
		}
		// 获取接口权限列表
		if res.OpList, err = r.getOpList(dao, arg); err != nil {
			return
		}
		if arg.IsSuperAdmin { // 如果是超级管理员
			err = r.getGroupMenu(dao, arg, res)
			return
		}
		var menuIds []int
		if menuIds, err = r.getPermitByGroupIds(dao, arg.Module, arg.PathTypes, arg.GroupId...); err != nil { // 普通管理员
			return
		}
		if err = r.getGroupMenu(dao, arg, res, menuIds...); err != nil {
			return
		}
	}()

	syncG.Wait()
	return
}
func (r *PermitServiceImpl) getMessageCount(userHid string) (count int, err error) {
	var httpHeader = http.Header{}
	logContent := map[string]interface{}{
		"userHid": userHid,
	}
	defer func() {
		if err != nil {
			logContent["err"] = err.Error()
			r.Context.Error(logContent, "PermitServiceImplGetMessageCount")
		} else {
			r.Context.Info(logContent, "PermitServiceImplGetMessageCount")
		}
	}()
	httpHeader.Set(app_obj.HttpUserToken, r.Context.GinContext.GetHeader(app_obj.HttpUserToken))
	httpHeader.Set(app_obj.HttpUserHid, r.Context.GinContext.GetHeader(app_obj.HttpUserHid))
	request := &rpc.RequestOptions{
		Context: r.Context,
		Method:  "GET",
		AppName: const_obj.MicroUser,
		Header:  httpHeader,
		URI:     "/in/user/has_not_msg",
		Value:   url.Values{},
	}
	request.Value.Set("user_hid", userHid)
	logContent["request"] = request
	var body string
	action := rpc.NewHttpRpc(request).Send().GetBody()
	if err = action.Error; action.Error != nil {
		return
	}

	body = action.GetBodyAsString()
	logContent["body"] = body

	type DataCount struct {
		Count int `json:"count"`
	}
	type AutoGenerated struct {
		Code    int       `json:"code"`
		Data    DataCount `json:"data"`
		Message string    `json:"message"`
	}
	var res AutoGenerated
	if err = json.Unmarshal(action.Body, &res); err != nil {
		return
	}
	count = res.Data.Count
	return
}
func (r *PermitServiceImpl) getOpList(dao daos.DaoPermit, arg *wrappers.ArgPermitMenu) (opList map[string][]wrappers.OpOne, err error) {
	var list []wrappers.Op
	if list, err = dao.GetPermitImportByModule(arg); err != nil {
		return
	}
	l := len(list)
	opList = make(map[string][]wrappers.OpOne, l)
	var t wrappers.OpOne
	for _, value := range list {
		if _, ok := opList[value.MenuPermitKey]; !ok {
			opList[value.MenuPermitKey] = make([]wrappers.OpOne, 0, l)
		}
		t = wrappers.OpOne(value.PermitKey)
		opList[value.MenuPermitKey] = append(opList[value.MenuPermitKey], t)
	}
	return
}

func (r *PermitServiceImpl) getCurrentSystem(dao daos.DaoPermit, arg *wrappers.ArgPermitMenu, ) (res []string, err error) {
	if arg.Module != "" {
		res = []string{arg.Module, wrappers.DefaultPermitModule}
		return
	}
	res = make([]string, 0, 2)
	res = append(res, wrappers.DefaultPermitModule)

	var list []models.AdminMenu
	if list, err = dao.GetByCondition(map[string]interface{}{"module": wrappers.DefaultPermitModule}, []wrappers.DaoOrderBy{
		{
			Column:     "sort_value",
			SortFormat: "asc",
		},
	}, 1); err != nil {
		return
	}
	if len(list) == 0 {
		err = fmt.Errorf("您的账号当前没有权限访问本系统")
		return
	}

	res = append(res, list[0].PermitKey)

	return
}

func (r *PermitServiceImpl) getGroupMenu(
	dao daos.DaoPermit,
	arg *wrappers.ArgPermitMenu,
	res *wrappers.ResultPermitMenuReturn,
	menuIds ...int) (err error) {
	var (
		groupPermit map[int][]wrappers.ResultPermitMenu
		permitMap   map[int]wrappers.ResultPermitMenu
		list        []models.AdminMenu
	)

	// 获取当前选中的系统
	var current []string
	if current, err = r.getCurrentSystem(dao, arg); err != nil {
		return
	}
	if list, err = dao.GetPermitMenuByIds(current, menuIds...); err != nil {
		return
	}

	if res.Menu, groupPermit, permitMap = r.groupPermit(list); len(res.Menu) <= 0 {
		err = fmt.Errorf("您没有操作权限,请刷新或联系管理员(1)")
		return
	}
	var parentId int
	if arg.ParentId > 0 {
		for k, item := range res.Menu {
			if item.Id == arg.ParentId {
				parentId = arg.ParentId
				res.Menu[k].Active = true
			}
		}
	} else {
		parentId = res.Menu[0].Id
		res.Menu[0].Active = true
	}

	res.ResultPermitMenu = wrappers.ResultPermitMenu{
		Id:       parentId,
		Children: []wrappers.ResultPermitMenu{},
	}
	r.orgPermit(groupPermit, parentId, &res.ResultPermitMenu, permitMap)

	res.RoutParentMap = r.orgRoutParentMap(parentId, list, &res.ResultPermitMenu)
	return
}
func (r *PermitServiceImpl) orgRoutParentMap(parentId int, list []models.AdminMenu, dataChild *wrappers.ResultPermitMenu) (returnData map[string][]string) {
	l := len(list)
	returnData = make(map[string][]string, l)
	res := make(map[int][]int, l)
	mapIdToParent := make(map[int]int, l)
	mapIdConfig := make(map[int]models.AdminMenu, l)
	for _, value := range list {
		mapIdConfig[value.Id] = value
		if value.Id == parentId {
			continue
		}
		if _, ok := res[value.Id]; !ok {
			res[value.Id] = make([]int, 0, 6)
		}
		mapIdToParent[value.Id] = value.ParentId
	}

	for key, value := range res {
		r.getParentId(mapIdToParent, key, parentId, &value)
		res[key] = value
	}
	for key, value := range res {
		pk := mapIdConfig[key].PermitKey
		returnData[pk] = make([]string, 0, 6)
		for _, id := range value {
			returnData[pk] = append(returnData[pk], mapIdConfig[id].PermitKey)
		}
	}

	return
}
func (r *PermitServiceImpl) getParentId(mapIdToParent map[int]int, nodeId, parentId int, value *[]int) {
	if parentId <= 0 {
		return
	}
	if _, ok := mapIdToParent[nodeId]; ok {
		if mapIdToParent[nodeId] != parentId {
			*value = append(*value, mapIdToParent[nodeId])
			r.getParentId(mapIdToParent, mapIdToParent[nodeId], parentId, value)
		}
	}
	return
}
func (r *PermitServiceImpl) groupPermit(list []models.AdminMenu) (
	systemList []wrappers.ResultSystemMenu,
	groupPermit map[int][]wrappers.ResultPermitMenu,
	permitMap map[int]wrappers.ResultPermitMenu,
) {
	systemList = make([]wrappers.ResultSystemMenu, 0, 30)
	groupPermit = map[int][]wrappers.ResultPermitMenu{}
	l := len(list)
	permitMap = make(map[int]wrappers.ResultPermitMenu, l)
	var data wrappers.ResultPermitMenu
	for _, item := range list {

		if item.ParentId == wrappers.DefaultPermitParentId {
			systemList = append(systemList, wrappers.ResultSystemMenu{
				Id:        item.Id,
				PermitKey: item.PermitKey,
				Label:     item.Label,
				Icon:      item.Icon,
				SortValue: item.SortValue,
				Module:    item.Module,
				Domain:    item.Domain,
			})
			continue
		}
		if _, ok := groupPermit[item.ParentId]; !ok {
			groupPermit[item.ParentId] = make([]wrappers.ResultPermitMenu, 0, l)
		}
		data = r.orgResultPermitMenu(item)
		permitMap[item.Id] = data
		groupPermit[item.ParentId] = append(groupPermit[item.ParentId], data)
	}
	return
}

func (r *PermitServiceImpl) orgPermit(group map[int][]wrappers.ResultPermitMenu, pid int, res *wrappers.ResultPermitMenu, permitMap map[int]wrappers.ResultPermitMenu) {

	if res == nil {
		res = &wrappers.ResultPermitMenu{}
	}
	if dt, ok := permitMap[pid]; ok {
		*res = dt
	}
	var ok bool
	if res.Children, ok = group[pid]; ok {
		if len(res.Children) == 0 {
			res.Children = []wrappers.ResultPermitMenu{}
			return
		}
		delete(group, pid)
		for key, item := range res.Children {
			r.orgPermit(group, item.Id, &(res.Children[key]), permitMap)
		}
		return
	}
	if len(res.Children) == 0 {
		res.Children = []wrappers.ResultPermitMenu{}
	}
}
func (r *PermitServiceImpl) orgResultPermitMenu(item models.AdminMenu) (res wrappers.ResultPermitMenu) {
	var f = false
	if item.HideInMenu == 1 {
		f = true
	}
	res = wrappers.ResultPermitMenu{
		Id:     item.Id,
		Path:   item.UrlPath,
		Name:   item.PermitKey,
		Module: item.Module,
		Meta: wrappers.PermitMeta{
			Icon:       item.Icon,
			Title:      item.Label,
			PermitKey:  item.PermitKey,
			HideInMenu: f,
		},
		Children: []wrappers.ResultPermitMenu{},
	}
	return
}

func (r *PermitServiceImpl) getPermitByGroupIds(dao daos.DaoPermit, module string, pathType []string, groupIds ...int) (menuIds []int, err error) {
	res, err := dao.GetMenuIdsByPermitByGroupIds(module, pathType, groupIds...)
	if err != nil {
		return
	}
	menuIds = make([]int, 0, len(res))
	for _, value := range res {
		menuIds = append(menuIds, value.MenuId)
	}
	return
}

type ArgGetUserGroupIds struct {
	UserId string `json:"user_id"`
	Dao    daos.DaoPermit
}

func (r *PermitServiceImpl) getUserGroupIds(arg *ArgGetUserGroupIds) (res []int, superAdmin bool, err error) {
	res = []int{}
	groups, err := arg.Dao.GetGroupByUserId(arg.UserId)
	if err != nil {
		return
	}
	res = make([]int, 0, len(groups))
	for _, group := range groups {
		if group.IsSuperAdmin > 0 { // 如果是超级管理员
			superAdmin = true
			return
		}
		res = append(res, group.GroupId)
	}
	return
}

func (r *PermitServiceImpl) Flag(arg *wrappers.ArgFlag) (res *wrappers.ResultFlag, err error) {
	res = &wrappers.ResultFlag{}
	return
}
func (r *PermitServiceImpl) getAppConfigList(dao daos.DaoService, arg *wrappers.ArgGetAppConfig, res *wrappers.ResultGetAppConfig) (err error) {
	var list []models.AdminApp
	if list, err = dao.GetList(nil, nil, nil); err != nil {
		return
	}
	err = r.orgResAppConfig(list, res, arg)
	return
}
func (r *PermitServiceImpl) orgResAppConfig(list []models.AdminApp, res *wrappers.ResultGetAppConfig, arg *wrappers.ArgGetAppConfig) (err error) {
	*res = make(map[string]string, len(list))
	for _, it := range list {
		if err = it.UnmarshalHosts(); err != nil {
			r.Context.Error(map[string]interface{}{
				"err": err.Error(),
				"it":  it,
				"arg": arg,
			}, "permitServiceGetAppConfigList")
			return
		}
		(*res)[it.UniqueKey] = fmt.Sprintf("http://%s:%d", it.HostConfig[arg.Env], it.Port)
	}
	return
}
func (r *PermitServiceImpl) GetAppConfig(arg *wrappers.ArgGetAppConfig) (res *wrappers.ResultGetAppConfig, err error) {
	res = &wrappers.ResultGetAppConfig{}
	dao := dao_impl.NewDaoServiceImpl(r.Context)
	if arg.Module == "" {
		err = r.getAppConfigList(dao, arg, res)
		return
	}
	err = r.getAppConfigListByModule(dao, arg, res)
	return
}
func (r *PermitServiceImpl) getImportMenu(dao daos.DaoService, arg *wrappers.ArgGetAppConfig) (uniqueKeys []string, err error) {
	var importMenus []wrappers.ImportMenu
	if importMenus, err = dao.GetImportMenuByModule(arg.Module); err != nil {
		return
	}
	uniqueKeys = make([]string, 0, len(importMenus))
	for _, it := range importMenus {
		uniqueKeys = append(uniqueKeys, it.AppName)
	}
	return
}
func (r *PermitServiceImpl) getAppConfigListByModule(dao daos.DaoService, arg *wrappers.ArgGetAppConfig, res *wrappers.ResultGetAppConfig) (err error) {
	var (
		list       []models.AdminApp
		uniqueKeys []string
	)
	if uniqueKeys, err = r.getImportMenu(dao, arg); err != nil {
		return
	}

	if list, err = dao.GetList(nil, &wrappers.ArgServiceList{
		UniqueKeys: uniqueKeys,
	}, nil); err != nil {
		return
	}
	if err = r.orgResAppConfig(list, res, arg); err != nil {
		return
	}
	return
}
