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
	"strconv"
	"strings"
	"time"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
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
func (r *PermitServiceImpl) MenuDelete(arg *wrappers.ArgMenuDelete) (res *wrappers.ResultMenuDelete, err error) {
	res = &wrappers.ResultMenuDelete{}
	dao := dao_impl.NewDaoPermit(r.Context)

	if err = dao.DeleteMenuByIds(arg.IdValue...); err != nil {
		return
	}
	daoGroup := dao_impl.NewDaoPermitGroupImpl(r.Context)
	if err = daoGroup.DeleteUserGroupPermit(models.PathTypePage, arg.IdValueNumber...); err != nil {
		return
	}

	// 删除菜单下的所有接口权限
	var importList []models.AdminImport
	importList, err = dao.GetImportMenuId(arg.IdValueNumber...)
	iIds := make([]int, 0, len(importList))
	for _, value := range importList {
		iIds = append(iIds, value.Id)
	}
	if err = daoGroup.DeleteUserGroupPermit(models.PathTypeApi, iIds...); err != nil {
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
	if module == "" {
		return
	}
	var dt []models.AdminMenu
	if dt, err = dao.GetMenuByPermitKey("", module); err != nil {
		return
	}
	if len(dt) == 0 {
		err = fmt.Errorf("您操作的系统信息不存在或已删除")
		return
	}
	systemId = dt[0].Id

	return
}
func (r *PermitServiceImpl) AdminMenu(arg *wrappers.ArgAdminMenu) (res *wrappers.ResultAdminMenu, err error) {

	res = wrappers.NewResultAdminMenu()

	dao := dao_impl.NewDaoPermit(r.Context)

	// 根据module获取当前操作的是哪个系统
	if arg.SystemId, err = r.getSystemIdByModule(dao, arg.Module, arg.SystemId); err != nil {
		return
	}


	var list, dt2 []models.AdminMenu

	if list, err = dao.GetAdminMenuList(arg); err != nil {
		return
	}

	if dt2, err = dao.GetAdminMenuList(&wrappers.ArgAdminMenu{}); err != nil {
		return
	}

	r.permitTab(dt2, &res.Menu, arg.SystemId, arg.Module)

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
		if err = NewSrvPermitImport(r.Context).SetApiPermit(dao, arg); err != nil {
			return
		}
	}
	res.Result = true
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

type MenuHandler struct {
	Name    string `json:"name"`
	handler func(arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn)
}

type ArgGetUserGroupIds struct {
	UserId string `json:"user_id"`
	Dao    daos.DaoPermit
}

func (r *PermitServiceImpl) Flag(arg *wrappers.ArgFlag) (res *wrappers.ResultFlag, err error) {
	res = &wrappers.ResultFlag{}
	return
}
