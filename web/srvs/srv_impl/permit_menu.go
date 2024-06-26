// Package srv_impl /**
package srv_impl

import (
	"encoding/json"
	"fmt"
	"github.com/juetun/library/common/app_param"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/plugins/rpc"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type SrvPermitMenuImpl struct {
	SrvPermitCommon
	dao      daos.DaoPermit
	daoGroup daos.DaoPermitGroup
}

func (r *SrvPermitMenuImpl) readySuperAdminMenu(arg *wrappers.ArgPermitMenu) (err error) {

	// 判断参数是否正常
	// 组织并判断当前访问的是哪个系统（如果参数没有则指定默认值）
	if err = r.initParentId(arg); err != nil {
		return
	}

	// 获取当前界面的界面基本信息
	if err = r.GetMenuPermitKeyByPath(&arg.ArgGetImportByMenuIdSingle); err != nil {
		return
	}

	if arg.NowMenuId == 0 {
		arg.NowMenuId, err = r.getDefaultAccessPage(arg.GroupId, arg.IsSuperAdmin, &arg.Module)
	}
	return
}

func (r *SrvPermitMenuImpl) readyGeneralMenu(arg *wrappers.ArgPermitMenu) (err error) {

	// 判断参数是否正常
	// 组织并判断当前访问的是哪个系统（如果参数没有则指定默认值）
	if err = r.initParentId(arg); err != nil {
		return
	}

	// 获取当前界面的界面基本信息
	if err = r.GetMenuPermitKeyByPath(&arg.ArgGetImportByMenuIdSingle); err != nil {
		return
	}

	if arg.NowMenuId == 0 {
		arg.NowMenuId, err = r.getDefaultAccessPage(arg.GroupId, arg.IsSuperAdmin, &arg.Module)
	}
	return
}

func (r *SrvPermitMenuImpl) superAdminMenu(arg *wrappers.ArgPermitMenu, handlers map[string]MenuHandler) (err error) {

	//准备参数数据(1、是否是超级管管理员)
	if err = r.readySuperAdminMenu(arg); err != nil {
		return
	}
	handlers["getPermitMenuList"] = r.getGroupMenuAdmin // 获取当前菜单的树形结构
	if arg.NowMenuId != 0 {
		handlers["getNowMenuImports"] = r.getSuperAdminNowMenuImports // 获取指定菜单下的接口ID列表
	}

	return
}

func (r *SrvPermitMenuImpl) generalAdminMenu(arg *wrappers.ArgPermitMenu, handlers map[string]MenuHandler) (err error) {

	//准备参数数据(1、是否是超级管管理员)
	if err = r.readyGeneralMenu(arg); err != nil {
		return
	}

	handlers["getPermitMenuList"] = r.getGeneralPermitMenuListUser // 获取当前菜单的树形结构

	if arg.NowMenuId != 0 {
		handlers["getNowMenuImports"] = r.getGeneralNowMenuImports // 获取指定菜单下的接口ID列表
	}

	return
}

// 超级管理员
// 获取当前菜单的信息
func (r *SrvPermitMenuImpl) getSuperAdminNowMenuImports(arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn) (err error) {

	param := &wrappers.ArgGetImportByMenuId{ArgGetImportByMenuIdSingle: arg.ArgGetImportByMenuIdSingle,}
	param.SuperAdminFlag = arg.IsSuperAdmin

	logContent := map[string]interface{}{
		"arg": arg,
	}

	defer func() {
		if err != nil {
			logContent["err"] = err.Error()
			r.Context.Error(logContent,
				"SrvPermitMenuImplGetSuperAdminNowMenuImports")
		}
	}()

	err = r.orgNowSuperAdminMenu(arg, res, logContent, param)

	return
}

func (r *SrvPermitMenuImpl) getGeneralNowMenuImports(arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn) (err error) {

	param := &wrappers.ArgGetImportByMenuId{ArgGetImportByMenuIdSingle: arg.ArgGetImportByMenuIdSingle,}
	logContent := map[string]interface{}{
		"arg": arg,
	}

	defer func() {
		if err != nil {
			logContent["err"] = err.Error()
			r.Context.Error(logContent, "SrvPermitMenuImplGetGeneralNowMenuImports")
		}
	}()

	err = r.orgNowGeneralAdminMenu(arg, res, logContent, param)

	return
}

func (r *SrvPermitMenuImpl) generalAdminister(handlers map[string]MenuHandler, arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn) (err error) {

	err = r.generalAdminMenu(arg, handlers)

	// 并行执行逻辑
	r.syncOperate(handlers, arg, res)
	return
}

func (r *SrvPermitMenuImpl) superAdminister(handlers map[string]MenuHandler, arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn) (err error) {
	//如果是超级管理员
	err = r.superAdminMenu(arg, handlers)

	// 并行执行逻辑
	r.syncOperate(handlers, arg, res)
	return
}

// Menu 每个界面调用获取菜单列表
func (r *SrvPermitMenuImpl) Menu(arg *wrappers.ArgPermitMenu) (res *wrappers.ResultPermitMenuReturn, err error) {

	res = wrappers.NewResultPermitMenuReturn()

	var (
		handlers = map[string]MenuHandler{
			"getNotReadMessage": r.getNotReadMessage, // 获取用户未读消息数),
		}
		e error
	)

	// 判断当前用户是否是超级管理员,
	var getUserArgument = base.NewArgGetByNumberIds(base.ArgGetByNumberIdsOptionIds(arg.UUserHid))
	if arg.IsSuperAdmin, e = r.GetAdminUserInfo(getUserArgument, arg.UUserHid); e != nil {
		res.ShowError = true
		res.ErrorMsg = e.Error()
		return
	} else {
		res.IsSuperAdmin = arg.IsSuperAdmin
		arg.SuperAdminFlag = arg.IsSuperAdmin
	}

	//如果超级管理员
	if arg.IsSuperAdmin {
		//如果是超级管理员
		if e = r.superAdminister(handlers, arg, res); e != nil {
			res.ShowError = true
			res.ErrorMsg = e.Error()
			return
		}
		return
	}

	// 如果不是超级管理员 返回当前用户所属用户组
	if arg.GroupId, arg.IsSuperAdmin, e = NewSrvPermitUserImpl(r.Context).
		GetUserAdminGroupIdByUserHid(arg.UUserHid); e != nil {
		res.ShowError = true
		res.ErrorMsg = e.Error()
		return
	}
	if len(arg.GroupId) == 0 {
		res.ShowError = true
		res.ErrorMsg = "您没有后台权限"
		return
	}
	res.IsSuperAdmin = arg.IsSuperAdmin
	arg.SuperAdminFlag = arg.IsSuperAdmin

	if !arg.IsSuperAdmin { //普通管理员
		if e = r.generalAdminister(handlers, arg, res); e != nil {
			res.ShowError = true
			res.ErrorMsg = e.Error()
			return
		}
		return
	}

	return

}

func (r *SrvPermitMenuImpl) initUserPermitGroup(arg *wrappers.ArgAdminMenu) (err error) {
	// 判断当前用户是否是超级管理员,
	var getUserArgument = base.NewArgGetByNumberIds(base.ArgGetByNumberIdsOptionIds(arg.UUserHid))
	if arg.OperatorIsSuperAdmin, err = r.GetAdminUserInfo(getUserArgument, arg.UUserHid); err != nil {
		return
	}

	if !arg.OperatorIsSuperAdmin { //如果账号不是超管
		// 判断当前用户是否是超级管理员,如果不是超级管理员，组织所属组权限
		if arg.OperatorGroupId, arg.OperatorIsSuperAdmin, err = NewSrvPermitUserImpl(r.Context).
			GetUserAdminGroupIdByUserHid(arg.UUserHid); err != nil {
			return
		}
	}
	return
}

//获取所有系统

// AdminMenu 管理后台菜单数据组织
func (r *SrvPermitMenuImpl) AdminMenu(arg *wrappers.ArgAdminMenu) (res *wrappers.ResultAdminMenu, err error) {

	res = wrappers.NewResultAdminMenu()
	var (
		havePermitSystemList []*models.AdminMenu
	)

	//初始化判断用户是否为超管，不为超管，则获取用户所在的权限组
	if err = r.initUserPermitGroup(arg); err != nil {
		return
	}

	// 根据module获取当前操作的是哪个系统
	if havePermitSystemList, arg.SystemId, arg.Module, err = r.GetSystemIdByModule(arg.Module, arg.OperatorGroupId, arg.OperatorIsSuperAdmin, false); err != nil {
		return
	}

	var list []*models.AdminMenu
	if list, err = r.dao.GetAdminMenuList(arg); err != nil {
		return
	}

	r.permitTab(havePermitSystemList, &res.Menu, arg.SystemId, arg.Module)

	//r.permitTab(list, &res.Menu, arg.SystemId, arg.Module)
	r.orgTree(list, arg.SystemId, &res.List, nil)
	return
}

func (r *SrvPermitMenuImpl) permitTab(list []*models.AdminMenu, menu *[]wrappers.ResultSystemAdminMenu, systemId int64, moduleSrc string) (sid int64, module string) {
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

// GetSystemIdByModule 获取当前系统权限和默认选中系统
// module 系统模块的key
// OperatorGroupId 权限组ID
// isSuperAdmin是否是超级管理员
// needHaveSystemData 是否需要获取系统组列表
func (r *SrvPermitMenuImpl) GetSystemIdByModule(module string, OperatorGroupId []int64, isSuperAdmin, needHaveSystemData bool) (havePermitSystemList []*models.AdminMenu, systemId int64, moduleResult string, err error) {

	defer func() {
		if err == nil && systemId == 0 {
			err = fmt.Errorf("您操作的系统信息不存在或已删除")
			return
		}
	}()
	// 根据Module获取参数SystemId的值
	if module != "" {
		havePermitSystemList, systemId, moduleResult, err = r.getAdminMenuByUserModule(module, OperatorGroupId, isSuperAdmin, needHaveSystemData)
		return
	}
	havePermitSystemList, systemId, moduleResult, err = r.getDefaultSelectModule(OperatorGroupId, isSuperAdmin)
	return
}

// 如果是超级管理员
func (r *SrvPermitMenuImpl) getAdminMenuByUserModuleSupperAdmin(module string, needHaveSystemData bool) (havePermitSystemList []*models.AdminMenu, systemId int64, moduleResult string, err error) {
	var dt *models.AdminMenu
	daoPermitMenu := dao_impl.NewDaoPermitMenu(r.Context)
	if havePermitSystemList, err = daoPermitMenu.GetAllSystemList(); err != nil {
		return
	}
	if module != "" {
		for _, menu := range havePermitSystemList {
			if menu.PermitKey == module {
				// 默认第一个做为默认选项
				systemId = menu.Id
				moduleResult = menu.PermitKey
				return
			}
		}
	}
	if !needHaveSystemData {
		if dt, err = daoPermitMenu.GetAdminMenuByModule(module); err != nil {
			return
		}
		if dt == nil {
			err = fmt.Errorf("你没有使用该系统(%s)的权限或系统不存在", module)
			return
		}
		moduleResult = module
		systemId = dt.Id
		return
	}
	return

}

func (r *SrvPermitMenuImpl) getAdminMenuByUserModuleGeneral(module string, OperatorGroupId []int64, needHaveSystemData bool) (havePermitSystemList []*models.AdminMenu, systemId int64, moduleResult string, err error) {
	if !needHaveSystemData {
		if havePermitSystemList, err = r.daoGroup.GetGroupAdminMenuByGroupIds(module, OperatorGroupId...); err != nil {
			return
		}
		for _, item := range havePermitSystemList {
			if item.Module != module {
				continue
			}
			systemId = item.Id
			moduleResult = item.Module
		}
	}
	var data []*models.AdminMenu
	if data, err = r.daoGroup.GetGroupAdminMenuByGroupIds(module); err != nil {
		return
	}

	systemId = data[0].Id
	moduleResult = data[0].Module
	return
}

func (r *SrvPermitMenuImpl) getAdminMenuByUserModule(module string, OperatorGroupId []int64, isSuperAdmin, needHaveSystemData bool) (havePermitSystemList []*models.AdminMenu, systemId int64, moduleResult string, err error) {
	if isSuperAdmin { // 如果是超级管理员
		havePermitSystemList, systemId, moduleResult, err = r.getAdminMenuByUserModuleSupperAdmin(module, needHaveSystemData)
		return
	}

	// 如果不是超级管理员
	havePermitSystemList, systemId, moduleResult, err = r.getAdminMenuByUserModuleGeneral(module, OperatorGroupId, needHaveSystemData)
	return
}

// 菜单数据参数准备
// 1、 判断是否为超级管理员
// 2、 判断选中模块（如果未指定，则默认选中一个）
func (r *SrvPermitMenuImpl) adminMenuWithCheckReadyParameterArg(arg *wrappers.ArgAdminMenuWithCheck) (havePermitSystemList []*models.AdminMenu, err error) {

	// 判断当前用户是否是超级管理员,
	var getUserArgument = base.NewArgGetByNumberIds(base.ArgGetByNumberIdsOptionIds(arg.UUserHid))
	if arg.OperatorIsSuperAdmin, err = r.GetAdminUserInfo(getUserArgument, arg.UUserHid); err != nil {
		return
	}
	if !arg.OperatorIsSuperAdmin {
		// 判断当前用户是否是超级管理员,如果不是超级管理员，组织所属组权限
		if arg.OperatorGroupId, arg.OperatorIsSuperAdmin, err = NewSrvPermitUserImpl(r.Context).
			GetUserAdminGroupIdByUserHid(arg.UUserHid); err != nil {
			return
		}
	}

	if arg.SystemId == 0 {
		// 当前应该渲染的系统ID
		if havePermitSystemList, arg.SystemId, arg.Module, err = r.GetSystemIdByModule(arg.Module, arg.OperatorGroupId, arg.OperatorIsSuperAdmin, false); err != nil {
			return
		}
	}
	if arg.SystemId == 0 {
		err = fmt.Errorf("您要操作的系统可能不存在、已删除或没有权限")
	}

	return
}

func (r *SrvPermitMenuImpl) AdminMenuWithCheck(arg *wrappers.ArgAdminMenuWithCheck) (res *wrappers.ResultMenuWithCheck, err error) {
	res = &wrappers.ResultMenuWithCheck{
		List: make([]wrappers.AdminMenuObject, 0, 20),
		Menu: make([]wrappers.ResultSystemAdminMenu, 0, 30),
	}

	var (
		havePermitSystemList []*models.AdminMenu
		list                 []*models.AdminMenu
	)
	if havePermitSystemList, err = r.adminMenuWithCheckReadyParameterArg(arg); err != nil {
		return
	}

	if list, err = r.dao.GetAdminMenuList(&arg.ArgAdminMenu); err != nil {
		return
	}

	r.permitTab(list, &res.Menu, arg.SystemId, arg.Module)

	res.SetSystemList(havePermitSystemList, arg.SystemId)

	// 获取当前已有的权限菜单
	var mapPermitMenu map[int64]int64
	if mapPermitMenu, err = r.getGroupPermitMenu(arg.Module, arg.GroupId); err != nil {
		return
	}

	r.orgTree(list, arg.SystemId, &res.List, mapPermitMenu)

	return
}

//
func (r *SrvPermitMenuImpl) orgTree(list []*models.AdminMenu, parentId int64, res *[]wrappers.AdminMenuObject, mapPermitMenu map[int64]int64) {
	var (
		tmp wrappers.AdminMenuObject
	)

	for _, value := range list {
		// 剔除默认数据那条
		if value.Id == parentId || value.ParentId != parentId {
			continue
		}
		tmp = r.orgAdminMenuObject(value)

		// 组织子菜单数据
		r.orgTree(list, value.Id, &tmp.Children, mapPermitMenu)

		// 判断当前菜单是否被勾选
		if _, ok := mapPermitMenu[value.Id]; len(tmp.Children) == 0 && ok {
			// 只有树的叶子节点才标记成选中
			tmp.ResultAdminMenuSingle.Checked = true
		}
		*res = append(*res, tmp)
	}
}

func (r *SrvPermitMenuImpl) orgAdminMenuObject(value *models.AdminMenu) (res wrappers.AdminMenuObject) {
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
		BadgeKey:           value.BadgeKey,
		Domain:             value.Domain,
		ManageImportPermit: value.ManageImportPermit,
	}

	if value.OtherValue != "" {
		_ = json.Unmarshal([]byte(value.OtherValue), &res.ResultAdminMenuSingle.ResultAdminMenuOtherValue)
	} else {
		res.ResultAdminMenuSingle.ResultAdminMenuOtherValue = wrappers.ResultAdminMenuOtherValue{
			Expand: true,
		}
	}

	return
}

func (r *SrvPermitMenuImpl) getGroupPermitMenu(module string, groupId int64) (mapPermitMenu map[int64]int64, err error) {
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

func (r *SrvPermitMenuImpl) GetMenuPermitKeyByPath(arg *wrappers.ArgGetImportByMenuIdSingle) (err error) {

	var adminMenu []*models.AdminMenu
	if arg.NowModule != "" && arg.NowPermitKey == "" {
		if adminMenu, err = dao_impl.NewDaoPermitMenu(r.Context).
			GetMenuByCondition(map[string]interface{}{
				"module": arg.NowModule,
				"label":  models.CommonMenuDefaultHomePage,
			}); err != nil {
			return
		}
	} else if adminMenu, err = dao_impl.NewDaoPermitMenu(r.Context).
		GetMenuByPermitKey(arg.NowModule, arg.NowPermitKey); err != nil {
		return
	}

	if len(adminMenu) == 0 {
		if !arg.SuperAdminFlag {
			err = fmt.Errorf("您没有权限访问当前界面(%s)", arg.NowRoutePath)
		}
		return
	}
	arg.NowMenuId = adminMenu[0].Id

	return
}

//超级管理员获取默认访问地址
func (r *SrvPermitMenuImpl) getDefaultAccessPageSuperAdmin(adminMenu []*models.AdminMenu, module *string) (res int64, err error) {

	if *module != "" {
		res, err = r.getDefaultAccessPageSuperAdminModuleNull(adminMenu, module)
		return
	}
	res = adminMenu[0].Id

	return
}

func (r *SrvPermitMenuImpl) getDefaultAccessPageSuperAdminModuleNull(adminMenu []*models.AdminMenu, module *string) (res int64, err error) {
	var f bool
	for _, item := range adminMenu {
		if item.Module == *module {
			f = true
			res = item.Id
			return
		}
	}
	if !f {
		*module = adminMenu[0].Module
		res = adminMenu[0].Id
	}
	return
}

func (r *SrvPermitMenuImpl) getDefaultAccessPage(adminGroupId []int64, isSuperAdmin bool, module *string) (res int64, err error) {
	var adminMenu []*models.AdminMenu
	logContent := map[string]interface{}{
		"adminGroupId": adminGroupId,
		"isSuperAdmin": isSuperAdmin,
		"module":       module,
	}
	defer func() {
		if err == nil {
			return
		}
		logContent["err"] = err.Error()
		r.Context.Error(logContent, "SrvPermitMenuImplGetDefaultAccessPage")
	}()

	if adminMenu, err = dao_impl.NewDaoPermitMenu(r.Context).
		GetByCondition(map[string]interface{}{
			"label":        models.CommonMenuDefaultHomePage,
			"is_home_page": models.AdminMenuIsHomePageYes,
		}, nil, 0); err != nil {
		logContent["desc"] = "GetByCondition"
		return
	}
	if len(adminMenu) == 0 {
		logContent["desc"] = "adminMenu is null"
		err = fmt.Errorf("当前没有可使用的系统")
		return
	}
	if isSuperAdmin { //超级管理员
		res, err = r.getDefaultAccessPageSuperAdmin(adminMenu, module)
		if err != nil {
			logContent["err"] = "getDefaultAccessPageSuperAdmin"
		}
		return
	}
	menuIds, _ := r.pickMenuIds(adminMenu)
	var listGroupMenu []*models.AdminUserGroupMenu
	if listGroupMenu, err = r.daoGroup.
		GetAdminMenuByGidAndMenuIds(menuIds, adminGroupId); err != nil {
		logContent["err"] = "GetAdminMenuByGroupIdAndMenuIds"
		return
	}
	if len(listGroupMenu) == 0 {
		err = fmt.Errorf("您没有系统访问权限")
		logContent["err"] = "listGroupMenu is null"
		return
	}

	res = listGroupMenu[0].MenuId
	return
}
func (r *SrvPermitMenuImpl) pickMenuIds(adminMenu []*models.AdminMenu) (menuIds []int64, mapMenu map[int64]*models.AdminMenu) {
	l := len(adminMenu)
	menuIds = make([]int64, 0, l)
	mapMenu = make(map[int64]*models.AdminMenu, l)
	for _, value := range adminMenu {
		menuIds = append(menuIds, value.Id)
		mapMenu[value.Id] = value
	}
	return
}

func (r *SrvPermitMenuImpl) orgNowGeneralAdminMenu(arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn, logContent map[string]interface{}, param *wrappers.ArgGetImportByMenuId) (err error) {
	var childMenus []*models.AdminUserGroupMenu
	if childMenus, err = r.daoGroup.
		GetAdminMenuByGidAndMenuIds([]int64{
			arg.NowMenuId,
		},
			arg.GroupId,
		); err != nil {
		return
	}
	for _, item := range childMenus {
		res.NowImportAndMenu.MenuIds = append(res.NowImportAndMenu.MenuIds, wrappers.MenuSingle{
			MenuId:    item.MenuId,
			PermitKey: item.MenuPermitKey,
		})
	}
	var imports []*models.AdminUserGroupImport
	if imports, err = dao_impl.NewDaoPermitGroupImport(r.Context).
		GetMenuIdsByPermitByGroupIds([]int64{arg.NowMenuId,}, arg.GroupId); err != nil {
		return
	}
	for _, item := range imports {
		res.NowImportAndMenu.ImportIds = append(res.NowImportAndMenu.ImportIds, wrappers.ImportSingle{
			ImportId:  item.ImportId,
			PermitKey: item.ImportPermitKey,
		}, )
	}
	return
}

//获取当前菜单下超级管理员具备的
func (r *SrvPermitMenuImpl) orgNowSuperAdminMenu(arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn, logContent map[string]interface{}, param *wrappers.ArgGetImportByMenuId) (err error) {
	srvImport := NewSrvPermitImport(r.Context)
	if res.NowImportAndMenu.MenuIds, err = srvImport.GetChildMenu(param.NowMenuId); err != nil {
		logContent["err1"] = err.Error()
		return
	}

	if res.NowImportAndMenu.ImportIds, err = srvImport.GetChildImport(param.NowMenuId); err != nil {
		logContent["err2"] = err.Error()
		return
	}
	return
}

// 非超级管理员
func (r *SrvPermitMenuImpl) getGeneralPermitMenuListUser(arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn) (err error) {

	// 获取接口权限列表
	//if res.OpList, err = NewSrvPermitImport(r.Context).
	//	GetOpList(dao, arg); err != nil {
	//	return
	//}

	var menuIds []int64
	if menuIds, err = r.getPermitByGroupIds(arg.Module, arg.GroupId...); err != nil { // 普通管理员
		return
	}

	if err = r.getGroupMenu(r.dao, arg, res, menuIds...); err != nil {
		return
	}
	return
}

func (r *SrvPermitMenuImpl) getPermitByGroupIds(module string, groupIds ...int64) (menuIds []int64, err error) {

	res, err := dao_impl.NewDaoPermitGroupMenu(r.Context).
		GetMenuIdsByPermitByGroupIds(module, groupIds...)

	if err != nil {
		return
	}

	menuIds = make([]int64, 0, len(res))
	for _, value := range res {
		menuIds = append(menuIds, value.MenuId)
	}
	return
}

// 获取用户未读消息数
func (r *SrvPermitMenuImpl) getNotReadMessage(arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn) (err error) {
	res.NotReadMsgCount, _ = r.getMessageCount(arg.UUserHid)
	return
}

// 获取当前用户能够打开的系统列表
func (r *SrvPermitMenuImpl) getNowUserSystemList(tmpList []*models.AdminMenu) (systemList []wrappers.ResultSystemMenu, err error) {
	var (
		data wrappers.ResultSystemMenu
	)
	for _, item := range tmpList {
		data = wrappers.ResultSystemMenu{}
		data.Id = item.Id
		data.PermitKey = item.PermitKey
		data.Label = item.Label
		data.Icon = item.Icon
		data.SortValue = item.SortValue
		data.Module = item.Module
		data.Domain = item.Domain
		systemList = append(systemList, data)
	}

	return
}

// 超级管理员
func (r *SrvPermitMenuImpl) getGroupMenuAdmin(arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn) (err error) {

	var (
		groupPermit          map[int64][]wrappers.ResultPermitMenu
		permitMap            map[int64]wrappers.ResultPermitMenu
		list, systemMenuList []*models.AdminMenu
		current              []string // 获取当前选中的系统
	)

	if arg.Module == "" {
		if _, _, arg.Module, err = r.getDefaultSelectModule([]int64{}, true); err != nil {
			return
		}
	}

	current = []string{arg.Module}

	if list, err = r.dao.GetPermitMenuByIds(current); err != nil {
		return
	}
	if groupPermit, permitMap, err = r.groupPermit(list); err != nil {
		return
	}

	if systemMenuList, err = r.dao.GetPermitMenuByIds([]string{wrappers.DefaultPermitModule}); err != nil {
		return
	}

	// 获取当前用户能够打开的系统列表
	if res.Menu, err = r.getNowUserSystemList(systemMenuList); err != nil {
		return
	}

	if len(res.Menu) <= 0 {
		err = fmt.Errorf("您没有操作权限,请刷新或联系管理员(1)")
		return
	}
	parentId := res.SetActive(arg.ParentId)

	res.ResultPermitMenu = wrappers.ResultPermitMenu{
		Id:       parentId,
		Children: []wrappers.ResultPermitMenu{},
	}
	var (
		badgeKeys    = make([]string, 0, 200)
		mapBadgeKeys map[string]float64
	)
	r.orgPermit(groupPermit, parentId, &res.ResultPermitMenu, permitMap, &badgeKeys)

	if mapBadgeKeys, err = r.getTagCount(badgeKeys); err != nil {
		return
	}
	res.ResultPermitMenu = r.setBadgeCount(mapBadgeKeys, res.ResultPermitMenu)
	res.RoutParentMap = r.orgRoutParentMap(parentId, list, &res.ResultPermitMenu)
	return
}

func (r *SrvPermitMenuImpl) setBadgeCount(mapBadgeKeys map[string]float64, resultPermitMenu wrappers.ResultPermitMenu) (result wrappers.ResultPermitMenu) {
	var (
		ok    bool
		count float64
	)
	result = resultPermitMenu

	if resultPermitMenu.BadgeKey != "" {
		if count, ok = mapBadgeKeys[resultPermitMenu.BadgeKey]; ok {
			result.BadgeCount = count
		}
	}
	result.Children = make([]wrappers.ResultPermitMenu, 0, len(resultPermitMenu.Children))
	for _, item := range resultPermitMenu.Children {
		item = r.setBadgeCount(mapBadgeKeys, item)
		result.Children = append(result.Children, item)
	}
	return
}

func (r *SrvPermitMenuImpl) getGroupMenu(dao daos.DaoPermit, arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn, menuIds ...int64) (err error) {
	var (
		groupPermit          map[int64][]wrappers.ResultPermitMenu
		permitMap            map[int64]wrappers.ResultPermitMenu
		list, systemMenuList []*models.AdminMenu
	)

	// 获取当前选中的系统
	if systemMenuList, err = r.getUserCanUseSystemWithGroupIds(arg.GroupId); err != nil {
		return
	}

	// 获取当前用户能够打开的系统列表
	if res.Menu, err = r.getNowUserSystemList(systemMenuList); err != nil {
		return
	}
	if len(res.Menu) <= 0 {
		err = fmt.Errorf("您没有操作权限,请刷新或联系管理员(1)")
		return
	}

	if list, err = dao.GetPermitMenuByIds([]string{arg.Module}, menuIds...); err != nil {
		return
	} else {
		if groupPermit, permitMap, err = r.groupPermit(list); err != nil {
			return
		}
	}

	parentId := r.orgActiveMenu(arg, res)

	res.ResultPermitMenu = wrappers.ResultPermitMenu{
		Id:       parentId,
		Children: []wrappers.ResultPermitMenu{},
	}
	var (
		//获取徽标的徽标KEY
		badgeKeys = make([]string, 0, 100)
	)
	r.orgPermit(groupPermit, parentId, &res.ResultPermitMenu, permitMap, &badgeKeys)

	res.RoutParentMap = r.orgRoutParentMap(parentId, list, &res.ResultPermitMenu)
	return
}
func (r *SrvPermitMenuImpl) orgActiveMenu(arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn) (parentId int64) {

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
	return
}

func (r *SrvPermitMenuImpl) groupPermit(list []*models.AdminMenu) (groupPermit map[int64][]wrappers.ResultPermitMenu, permitMap map[int64]wrappers.ResultPermitMenu, err error) {
	l := len(list)
	groupPermit = make(map[int64][]wrappers.ResultPermitMenu, l)
	permitMap = make(map[int64]wrappers.ResultPermitMenu, l)

	var (
		data  wrappers.ResultPermitMenu
		ok    bool
		count float64
	)
	for _, item := range list {
		if _, ok = groupPermit[item.ParentId]; !ok {
			groupPermit[item.ParentId] = make([]wrappers.ResultPermitMenu, 0, l)
		}

		data = r.orgResultPermitMenu(item, count)
		data.BadgeKey = item.BadgeKey
		permitMap[item.Id] = data
		groupPermit[item.ParentId] = append(groupPermit[item.ParentId], data)
	}
	return
}

func (r *SrvPermitMenuImpl) orgResultPermitMenu(item *models.AdminMenu, badgeCount float64) (res wrappers.ResultPermitMenu) {
	var f = false
	if item.HideInMenu == models.AdminMenuHideInMenuTrue {
		f = true
	}
	res = wrappers.ResultPermitMenu{
		Id:         item.Id,
		Path:       item.UrlPath,
		Name:       item.PermitKey,
		Module:     item.Module,
		Label:      item.Label,
		BadgeCount: badgeCount,
		Meta: &wrappers.PermitMeta{
			Icon:       item.Icon,
			Title:      item.Label,
			PermitKey:  item.PermitKey,
			HideInMenu: f,
		},
		Children: []wrappers.ResultPermitMenu{},
	}
	return
}

func (r *SrvPermitMenuImpl) orgPermit(group map[int64][]wrappers.ResultPermitMenu, pid int64, res *wrappers.ResultPermitMenu, permitMap map[int64]wrappers.ResultPermitMenu, badgeKeys *[]string) {

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
			if item.BadgeKey != "" {
				*badgeKeys = append(*badgeKeys, item.BadgeKey)
			}
			r.orgPermit(group, item.Id, &(res.Children[key]), permitMap, badgeKeys)
		}
		return
	}
	if len(res.Children) == 0 {
		res.Children = []wrappers.ResultPermitMenu{}
	}
}

// 获取默认打开的系统
func (r *SrvPermitMenuImpl) getDefaultSelectModule(groupId []int64, isSupperAdmin bool) (havePermitSystemList []*models.AdminMenu, systemId int64, moduleResult string, err error) {
	if isSupperAdmin { // 如果是超级管理员
		dao := dao_impl.NewDaoPermitMenu(r.Context)
		if havePermitSystemList, err = dao.GetAllSystemList(); err != nil {
			return
		}
		for k, menu := range havePermitSystemList {
			if k == 0 { // 默认第一个做为默认选项
				systemId = menu.Id
				moduleResult = menu.PermitKey

			}

		}
		if systemId == 0 {
			err = fmt.Errorf("")
			return
		}
	}
	return
}

func (r *SrvPermitMenuImpl) getMessageCount(userHid int64) (count int, err error) {
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
		Context:     r.Context,
		Method:      http.MethodGet,
		AppName:     app_param.AppNameNotice,
		Header:      httpHeader,
		URI:         "/announce/has_not_msg",
		PathVersion: app_obj.App.AppRouterPrefix.Intranet,
		Value:       url.Values{},
	}
	request.Value.Set("user_hid", fmt.Sprintf("%d", userHid))
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

func (r *SrvPermitMenuImpl) getTagCount(badgeKeys []string) (mapBadgeKeys map[string]float64, err error) {
	var (
		l = len(badgeKeys)
	)
	mapBadgeKeys = make(map[string]float64, l)
	if l == 0 {
		return
	}
	arg := url.Values{}
	params := rpc.RequestOptions{
		Context:     r.Context,
		Method:      http.MethodPost,
		AppName:     app_param.AppNameTag,
		URI:         "/permit/get_plat_tag_count",
		Value:       arg,
		PathVersion: app_obj.App.AppRouterPrefix.Intranet,
		Header:      http.Header{},
	}
	arg.Set("badge_keys", strings.Join(badgeKeys, ","))
	req := rpc.NewHttpRpc(&params).
		Send().GetBody()
	if err = req.Error; err != nil {
		return
	}
	var body []byte
	if body = req.Body; len(body) == 0 {
		return
	}

	var resResult struct {
		Code int                `json:"code"`
		Data map[string]float64 `json:"data"`
		Msg  string             `json:"message"`
	}
	if err = json.Unmarshal(body, &resResult); err != nil {
		return
	}
	if resResult.Code > 0 {
		err = fmt.Errorf(resResult.Msg)
		return
	}
	mapBadgeKeys = resResult.Data

	return
}
func (r *SrvPermitMenuImpl) orgRoutParentMap(parentId int64, list []*models.AdminMenu, dataChild *wrappers.ResultPermitMenu) (returnData map[string][]string) {
	_ = dataChild
	l := len(list)
	returnData = make(map[string][]string, l)
	res := make(map[int64][]int64, l)
	mapIdToParent := make(map[int64]int64, l)
	mapIdConfig := make(map[int64]*models.AdminMenu, l)
	for _, value := range list {
		mapIdConfig[value.Id] = value
		if value.Id == parentId {
			continue
		}
		if _, ok := res[value.Id]; !ok {
			res[value.Id] = make([]int64, 0, l)
		}
		mapIdToParent[value.Id] = value.ParentId
	}

	for key, value := range res {
		r.getParentId(mapIdToParent, key, parentId, &value)
		res[key] = value
	}
	for key, value := range res {
		pk := mapIdConfig[key].PermitKey
		if _, ok := returnData[pk]; !ok {
			returnData[pk] = make([]string, 0, l)
		}
		for _, id := range value {
			if dtm, ok := mapIdConfig[id]; ok {
				returnData[pk] = append(returnData[pk], dtm.PermitKey)
			}
		}
	}
	return
}

func (r *SrvPermitMenuImpl) getParentId(mapIdToParent map[int64]int64, nodeId, parentId int64, value *[]int64) {
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

func (r *SrvPermitMenuImpl) syncOperate(handlers map[string]MenuHandler, arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn) {
	var syncG sync.WaitGroup
	syncG.Add(len(handlers))
	for _, funcHandler := range handlers {
		go func(h MenuHandler) {
			defer syncG.Done()
			_ = h(arg, res)
		}(funcHandler)
	}
	syncG.Wait()
}

func (r *SrvPermitMenuImpl) initParentId(arg *wrappers.ArgPermitMenu) (err error) {

	if arg.Module == "" {
		return
	}
	var dt []*models.AdminMenu
	if dt, err = dao_impl.NewDaoPermitMenu(r.Context).
		GetMenuByCondition(map[string]interface{}{"permit_key": arg.Module}); err != nil {
		return
	}
	if len(dt) == 0 {
		err = fmt.Errorf("您查看的系统(%s)不存在或已删除", arg.Module)
		return
	}
	arg.ParentId = dt[0].Id

	if arg.ArgGetImportByMenuIdSingle.NowModule == "" {
		arg.ArgGetImportByMenuIdSingle.NowModule = dt[0].PermitKey
	}
	return
}

// 获取当前用户组能够使用的系统列表
func (r *SrvPermitMenuImpl) getUserCanUseSystemWithGroupIds(groupIds []int64) (res []*models.AdminMenu, err error) {

	res, err = r.daoGroup.GetGroupSystemByGroupId(wrappers.DefaultPermitModule, groupIds...)
	return
}

func NewSrvPermitMenu(ctx ...*base.Context) (res srvs.SrvPermitMenu) {
	p := &SrvPermitMenuImpl{}
	p.SetContext(ctx...)
	p.dao = dao_impl.NewDaoPermit(p.Context)
	p.daoGroup = dao_impl.NewDaoPermitGroup(p.Context)
	return p
}
