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
	daoPermitMenu := dao_impl.NewDaoPermitMenu(r.Context)
	if li, err = daoPermitMenu.GetMenu(arg.MenuId); err != nil {
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
	res.List, err = dao_impl.NewDaoPermit(r.Context).
		GetAdminMenuList(arg)
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
	var adminMenu models.AdminMenu

	listData := adminMenu.InitDefaultSystemMenu(&models.DefaultSystemMenuNeedParams{
		Module:         data.PermitKey,
		UpdateTime:     data.UpdatedAt,
		CreateTime:     data.CreatedAt,
		ParentSystemId: data.ParentId,
	})
	for _, datum := range listData {
		if err = dao.Add(datum); err != nil {
			return
		}
	}

	return
}

func (r *PermitServiceImpl) MenuImport(arg *wrappers.ArgMenuImport) (res *wrappers.ResultMenuImport, err error) {
	res = &wrappers.ResultMenuImport{
		Pager: response.NewPager(response.PagerBaseQuery(&arg.PageQuery)),
	}
	var (
		dao = dao_impl.NewDaoPermit(r.Context)
		db  *gorm.DB
	)

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

	if err = dao_impl.NewDaoPermitGroupImpl(r.Context).
		DeleteUserGroupPermit(arg.IdValueNumber...); err != nil {
		return
	}

	// 删除菜单下的所有接口权限
	var importList []models.AdminImport
	importList, err = dao.GetImportMenuId(arg.IdValueNumber...)
	iIds := make([]int64, 0, len(importList))
	for _, value := range importList {
		iIds = append(iIds, value.Id)
	}
	if err = dao_impl.NewDaoPermitGroupImport(r.Context).
		DeleteGroupImportWithMenuId(iIds...); err != nil {
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
		daoPermitMenu := dao_impl.NewDaoPermitMenu(r.Context)
		if menus, err = daoPermitMenu.GetMenu(arg.Id); err != nil {
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

func (r *PermitServiceImpl) getChildIds(dao daos.DaoPermitMenu, parentId []string, ids *[]string) (err error) {
	if len(parentId) == 0 {
		return
	}
	var li []models.AdminMenu
	if li, err = dao.GetMenuByCondition(fmt.Sprintf("parent_id IN (%s)", strings.Join(parentId, ","))); err != nil {
		return
	}
	pIds := make([]string, 0, len(li))
	for _, value := range li {
		idString := fmt.Sprintf("%d", value.Id)
		pIds = append(pIds, idString)
		*ids = append(*ids, idString)
	}

	if err = r.getChildIds(dao, pIds, ids); err != nil {
		return
	}
	return
}
func (r *PermitServiceImpl) updateChildModule(dao daos.DaoPermit, parentId int64, module string) (err error) {

	ids := make([]string, 0, 20)
	pidString := fmt.Sprintf("%d", parentId)
	ids = append(ids, pidString)

	daoMenu := dao_impl.NewDaoPermitMenu(r.Context)
	if err = r.getChildIds(daoMenu, []string{pidString}, &ids); err != nil {
		return
	}
	if err = dao.UpdateMenuByCondition(
		fmt.Sprintf("id IN(%s)", strings.Join(ids, ",")),
		map[string]interface{}{"module": module},
	); err != nil {
		return
	}

	ids = append(ids, pidString)
	// 更新菜单接口关系表的menu_module
	if err = dao_impl.NewDaoPermitImport(r.Context).
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
	var mapPermitMenu map[int64]int64
	if mapPermitMenu, err = r.getGroupPermitMenu(arg.Module, arg.GroupId); err != nil {
		return
	}

	r.orgTree(list, arg.SystemId, &res.List, mapPermitMenu)

	return
}

func (r *PermitServiceImpl) getGroupPermitMenu(module string, groupId int64) (mapPermitMenu map[int64]int64, err error) {
	var groupPermit []models.AdminUserGroupMenu
	daoGroupMenu := dao_impl.NewDaoPermitGroupMenu(r.Context)
	if groupPermit, err = daoGroupMenu.GetMenuIdsByPermitByGroupIds(module, []int64{groupId}...); err != nil {
		return
	} else {
		mapPermitMenu = make(map[int64]int64, len(groupPermit))
		for _, it := range groupPermit {
			mapPermitMenu[it.MenuId] = it.MenuId
		}
	}
	return
}

func (r *PermitServiceImpl) getSystemIdByModule(dao daos.DaoPermit, module string, sysId int64) (systemId int64, err error) {
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
func (r *PermitServiceImpl) permitTab(list []models.AdminMenu, menu *[]wrappers.ResultSystemAdminMenu, systemId int64, moduleSrc string) (sid int64, module string) {
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
func (r *PermitServiceImpl) orgTree(list []models.AdminMenu, parentId int64, res *[]wrappers.AdminMenuObject, mapPermitMenu map[int64]int64) {
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

func (r *PermitServiceImpl) orgAdminMenuObject(value *models.AdminMenu, mapPermitMenu map[int64]int64) (res wrappers.AdminMenuObject) {
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
func (r *PermitServiceImpl) getGroupHavePermit(arg *wrappers.ArgAdminSetPermit) (newPermit []int64, notPermitId []int64, err error) {

	listSelectMenu, err := dao_impl.NewDaoPermitGroupMenu(r.Context).
		GetMenuIdsByPermitByGroupIds(arg.Module, arg.GroupId)
	if err != nil {
		return
	}
	commonPermit := make([]int64, 0, len(listSelectMenu)) // 当前已经选中的菜单
	notPermitId = make([]int64, 0, len(listSelectMenu))   // 当前取消权限的菜单
	newPermit = make([]int64, 0, len(listSelectMenu))     // 新增的菜单
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
func (r *PermitServiceImpl) inSlice(id int64, slice []int64) (res bool) {
	for _, value := range slice {
		if id == value {
			res = true
			return
		}
	}
	return
}

func (r *PermitServiceImpl) isHomePage(permitIds []int64) (res bool, homePageId int64, module string, err error) {
	var li []models.AdminMenu
	permitMenuDao := dao_impl.NewDaoPermitMenu(r.Context)
	if li, err = permitMenuDao.GetMenu(permitIds...); err != nil {
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
func (r *PermitServiceImpl) commonImport(module string) (res []int64, err error) {
	res = make([]int64, 0, 2)
	var da []models.AdminMenu
	daoPermitMenu := dao_impl.NewDaoPermitMenu(r.Context)

	da, err = daoPermitMenu.GetMenuByCondition(map[string]interface{}{
		"module": module,
		"label":  models.CommonMenuDefaultLabel,
	})
	for _, value := range da {
		res = append(res, value.Id)
	}
	return
}
func (r *PermitServiceImpl) orgNeedMenu(arg *wrappers.ArgAdminSetPermit) (permitIds []int64, err error) {
	var (
		isHomePage bool
		homePageId int64
	)
	permitIds = make([]int64, 0, 2)
	// 判断是否为系统首页，如果是首页，则自动绑定公共接口隐藏界面
	if isHomePage, homePageId, arg.Module, err = r.isHomePage(arg.PermitIds); err != nil {
		return
	}
	if permitIds, err = r.commonImport(arg.Module); err != nil {
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

func (r *PermitServiceImpl) deleteGroupMenuPermitByGroupId(arg *wrappers.ArgAdminSetPermit) (err error) {
	if len(arg.PermitIds) == 0 {
		return
	}
	dao := dao_impl.NewDaoPermitGroupImport(r.Context)
	// 没有权限ID，则说明清除所有的权限
	if err = dao.DeleteGroupMenuPermitByGroupIdAndMenuIds(arg.GroupId, arg.PermitIds...); err != nil {
		return
	}

	daoGroupMenu := dao_impl.NewDaoPermitGroupMenu(r.Context)

	if err = daoGroupMenu.DeleteGroupMenuByGroupIdAndIds(arg.GroupId, arg.PermitIds...); err != nil {
		return
	}

	return
}
func (r *PermitServiceImpl) setMenuPermit(dao daos.DaoPermit, arg *wrappers.ArgAdminSetPermit) (err error) {
	var (
		newPermit   []int64
		notPermitId []int64
		permitIds   []int64
	)

	// permitId 为空处理逻辑
	if err = r.deleteGroupMenuPermitByGroupId(arg); err != nil {
		return
	}

	if permitIds, err = r.orgNeedMenu(arg); err != nil {
		return
	}

	if newPermit, notPermitId, err = r.getGroupHavePermit(arg); err != nil {
		return
	}

	if err = r.deleteNotMenuPermitId(dao, append(notPermitId, permitIds...), arg); err != nil {
		return
	}
	t := time.Now()

	if err = r.addNewMenuPermit(newPermit, arg, t); err != nil {
		return
	}

	if err = r.addNewImportPermit(newPermit, arg, t); err != nil {
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
		if err = NewSrvPermitImport(r.Context).SetApiPermit(arg); err != nil {
			return
		}
	default:
		err = fmt.Errorf("当前不支持你选择的类型(%s)", arg.Type)
		return
	}
	res.Result = true
	return
}

func (r *PermitServiceImpl) setApiPermitOld(dao daos.DaoPermit, arg *wrappers.ArgAdminSetPermit) (err error) {
	var (
		newPermit   []int64
		notPermitId []int64
	)
	if newPermit, notPermitId, err = r.getGroupHavePermit(arg); err != nil {
		return
	}
	if err = r.deleteNotApiPermitId(dao, notPermitId, arg); err != nil {
		return
	}

	if err = r.addNewApiPermit(newPermit, arg.GroupId); err != nil {
		return
	}
	return
}

func (r *PermitServiceImpl) deleteNotMenuPermitId(dao daos.DaoPermit, notPermitId []int64, args *wrappers.ArgAdminSetPermit) (err error) {
	if len(notPermitId) == 0 {
		return
	}
	var listImport []models.AdminImport

	// notPermitId = append(notPermitId, args.PermitIds...)
	// 获取菜单下的接口
	var notPermitImportId []int64

	if listImport, err = dao.GetDefaultImportByMenuIds(args.Type, args.Module, notPermitId...); err != nil {
		return
	}
	if len(listImport) > 0 {
		notPermitImportId = make([]int64, 0, len(listImport))
		for _, it := range listImport {
			notPermitImportId = append(notPermitImportId, it.Id)
		}
	}
	if err = dao_impl.NewDaoPermitGroupMenu(r.Context).
		DeleteGroupPermitByMenuIds(args.GroupId, args.Module, notPermitId); err != nil {
		return
	}

	if err = dao_impl.NewDaoPermitGroupImport(r.Context).
		DeleteGroupMenuPermitByGroupIdAndMenuIds(args.GroupId, notPermitId...); err != nil {
		return
	}
	return
}

func (r *PermitServiceImpl) deleteNotApiPermitId(dao daos.DaoPermit, notPermitId []int64, arg *wrappers.ArgAdminSetPermit) (err error) {
	var listImport []models.AdminImport

	// 获取菜单下的接口
	var notPermitImportId []int64
	if listImport, err = dao.GetDefaultImportByMenuIds(arg.Type, arg.Module, notPermitId...); err != nil {
		return
	} else if len(listImport) > 0 {
		notPermitImportId = make([]int64, 0, len(listImport))
		for _, it := range listImport {
			notPermitImportId = append(notPermitImportId, it.Id)
		}
	}
	if err = dao_impl.NewDaoPermitGroupImport(r.Context).
		DeleteGroupMenuPermitByGroupIdAndMenuIds(arg.GroupId, notPermitImportId...); err != nil {
		return
	}
	return
}

func (r *PermitServiceImpl) addNewImportPermit(newPermit []int64, arg *wrappers.ArgAdminSetPermit, t time.Time) (err error) {
	var (
		listImport []models.AdminImport
		dt         models.AdminUserGroupImport
	)
	daoPermitImport := dao_impl.NewDaoPermitImport(r.Context)
	if listImport, err = daoPermitImport.GetDefaultOpenImportByMenuIds(newPermit...); err != nil {
		return
	}

	list := make([]base.ModelBase, 0, len(listImport))
	for _, importData := range listImport {
		dt = models.AdminUserGroupImport{
			Module:    arg.Module,
			AppName:   importData.AppName,
			GroupId:   arg.GroupId,
			MenuId:    importData.Id,
			ImportId:  importData.Id,
			CreatedAt: t,
			UpdatedAt: t,
		}
		list = append(list, &dt)
	}
	var m models.AdminUserGroupImport
	daoGImport := dao_impl.NewDaoPermitGroupImport(r.Context)
	if err = daoGImport.BatchAddData(m.TableName(), list); err != nil {
		err = fmt.Errorf("操作异常")
		return
	}
	return
}

func (r *PermitServiceImpl) addNewMenuPermit(newPermit []int64, arg *wrappers.ArgAdminSetPermit, t time.Time) (err error) {
	l := len(newPermit)
	if l == 0 {
		return
	}

	var (
		list = make([]base.ModelBase, 0, l)
		dt   *models.AdminUserGroupMenu
	)

	for _, pid := range newPermit {
		dt = &models.AdminUserGroupMenu{
			Module:    arg.Module,
			GroupId:   arg.GroupId,
			MenuId:    pid,
			CreatedAt: t,
			UpdatedAt: t,
		}
		list = append(list, dt)
	}
	daoGroupMenu := dao_impl.NewDaoPermitGroupMenu(r.Context)
	if err = daoGroupMenu.BatchAddData(dt.TableName(), list); err != nil {
		err = fmt.Errorf("操作异常")
		return
	}
	return
}

func (r *PermitServiceImpl) addNewApiPermit(newPermit []int64, groupId int64) (err error) {
	if len(newPermit) == 0 {
		return
	}
	list := make([]base.ModelBase, 0, len(newPermit))
	var dt *models.AdminUserGroupImport
	var t = time.Now()
	for _, importId := range newPermit {
		dt = &models.AdminUserGroupImport{
			GroupId:   groupId,
			MenuId:    importId,
			CreatedAt: t,
			UpdatedAt: t,
		}
		list = append(list, dt)
	}
	dao := dao_impl.NewDaoPermitGroupImport(r.Context)
	var m models.AdminUserGroupImport
	if err = dao.BatchAddData(m.TableName(), list); err != nil {
		return
	}
	return
}

type MenuHandler struct {
	Name    string `json:"name"`
	handler func(arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn) (err error)
}

type ArgGetUserGroupIds struct {
	UserId string `json:"user_id"`
	Dao    daos.DaoPermit
}
