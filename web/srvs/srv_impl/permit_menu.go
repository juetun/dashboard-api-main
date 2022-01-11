// Package srv_impl /**
package srv_impl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/plugins/rpc"
	"github.com/juetun/dashboard-api-main/pkg/parameters"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type SrvPermitMenuImpl struct {
	base.ServiceBase
}

func (r *SrvPermitMenuImpl) readyMenu(arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn) (err error) {

	// 判断当前用户是否是超级管理员,
	// 如果不是超级管理员 返回当前用户所属用户组
	if arg.GroupId, arg.IsSuperAdmin, err = NewSrvPermitUserImpl(r.Context).
		GetUserAdminGroupIdByUserHid(arg.UserHid); err != nil {
		return
	}
	res.IsSuperAdmin = arg.IsSuperAdmin

	// 判断参数是否正常
	// 组织并判断当前访问的是哪个系统（如果参数没有则指定默认值）
	if err = r.initParentId(arg); err != nil {
		return
	}

	// 获取当前界面的界面基本信息
	if err = r.GetMenuPermitKeyByPath(&arg.ArgGetImportByMenuIdSingle); err != nil {
		return
	}
	return
}

// Menu 每个界面调用获取菜单列表
func (r *SrvPermitMenuImpl) Menu(arg *wrappers.ArgPermitMenu) (res *wrappers.ResultPermitMenuReturn, err error) {

	res = wrappers.NewResultPermitMenuReturn()

	//准备参数数据(1、是否是超级管管理员)
	if err = r.readyMenu(arg, res); err != nil {
		return
	}

	handlers := map[string]MenuHandler{
		"getNotReadMessage": r.getNotReadMessage, // 获取用户未读消息数),
		"getPermitMenuList": r.getPermitMenuList, // 获取当前菜单的树形结构
	}

	if arg.NowMenuId != 0 {
		handlers["getNowMenuImports"] = r.getNowMenuImports // 获取指定菜单下的接口ID列表
	}

	// 并行执行逻辑
	r.syncOperate(handlers, arg, res)
	return

}

// AdminMenu 管理后台菜单数据组织
func (r *SrvPermitMenuImpl) AdminMenu(arg *wrappers.ArgAdminMenu) (res *wrappers.ResultAdminMenu, err error) {

	res = wrappers.NewResultAdminMenu()

	dao := dao_impl.NewDaoPermit(r.Context)

	// 判断当前用户是否是超级管理员,如果不是超级管理员，组织所属组权限
	if arg.OperatorGroupId, arg.OperatorIsSuperAdmin, err = NewSrvPermitUserImpl(r.Context).
		GetUserAdminGroupIdByUserHid(arg.UserHid); err != nil {
		return
	}

	// 根据module获取当前操作的是哪个系统
	if _, arg.SystemId, arg.Module, err = r.GetSystemIdByModule(arg.Module, arg.OperatorGroupId, arg.OperatorIsSuperAdmin, false); err != nil {
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

	r.permitTab(list, &res.Menu, arg.SystemId, arg.Module)
	r.orgTree(list, arg.SystemId, &res.List, nil)
	return
}

func (r *SrvPermitMenuImpl) permitTab(list []models.AdminMenu, menu *[]wrappers.ResultSystemAdminMenu, systemId int64, moduleSrc string) (sid int64, module string) {
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
	var dt models.AdminMenu
	if !needHaveSystemData {
		if dt, err = dao_impl.NewDaoPermitMenu(r.Context).
			GetAdminMenuByModule(module); err != nil {
			return
		}
		systemId = dt.Id
		return
	}

	dao := dao_impl.NewDaoPermitMenu(r.Context)
	if havePermitSystemList, err = dao.GetAllSystemList(); err != nil {
		return
	}

	for _, menu := range havePermitSystemList {
		if menu.Module != module {
			continue
		} // 默认第一个做为默认选项
		systemId = menu.Id
		moduleResult = menu.PermitKey
	}
	return

}

func (r *SrvPermitMenuImpl) getAdminMenuByUserModuleGeneral(module string, OperatorGroupId []int64, needHaveSystemData bool) (havePermitSystemList []*models.AdminMenu, systemId int64, moduleResult string, err error) {
	if !needHaveSystemData {
		if havePermitSystemList, err = dao_impl.NewDaoPermitGroup(r.Context).
			GetGroupAdminMenuByGroupIds(module, OperatorGroupId...); err != nil {
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
	if data, err = dao_impl.NewDaoPermitGroup(r.Context).
		GetGroupAdminMenuByGroupIds(module); err != nil {
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

	// 判断当前用户是否是超级管理员,如果不是超级管理员，组织所属组权限
	if arg.OperatorGroupId, arg.OperatorIsSuperAdmin, err = NewSrvPermitUserImpl(r.Context).
		GetUserAdminGroupIdByUserHid(arg.UserHid); err != nil {
		return
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
	dao := dao_impl.NewDaoPermit(r.Context)

	var (
		havePermitSystemList []*models.AdminMenu
		list                 []models.AdminMenu
	)
	if havePermitSystemList, err = r.adminMenuWithCheckReadyParameterArg(arg); err != nil {
		return
	}

	if list, err = dao.GetAdminMenuList(&arg.ArgAdminMenu); err != nil {
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
func (r *SrvPermitMenuImpl) orgTree(list []models.AdminMenu, parentId int64, res *[]wrappers.AdminMenuObject, mapPermitMenu map[int64]int64) {
	var (
		tmp wrappers.AdminMenuObject
	)

	for _, value := range list {
		// 剔除默认数据那条
		if value.Id == parentId || value.ParentId != parentId {
			continue
		}
		tmp = r.orgAdminMenuObject(&value)

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
		Domain:             value.Domain,
		ManageImportPermit: value.ManageImportPermit,
	}

	if value.OtherValue != "" {
		_ = json.Unmarshal([]byte(value.OtherValue), &res.ResultAdminMenuSingle.ResultAdminMenuOtherValue)
	} else {
		res.ResultAdminMenuSingle.ResultAdminMenuOtherValue = wrappers.ResultAdminMenuOtherValue{Expand: true}
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
	divString := "/"
	pathSlice := strings.Split(strings.TrimLeft(arg.NowRoutePath, "/"), divString)
	switch len(pathSlice) {
	case 0:
		pathSlice = append(pathSlice, []string{"", ""}...)
	case 1:
		pathSlice = append(pathSlice, "")
	}
	arg.NowModule = pathSlice[0]
	arg.NowPermitKey = strings.Join(pathSlice[1:], divString)

	var adminMenu []models.AdminMenu
	if adminMenu, err = dao_impl.NewDaoPermitMenu(r.Context).
		GetMenuByPermitKey(arg.NowModule, arg.NowPermitKey); err != nil {
		return
	}
	if len(adminMenu) == 0 {
		return
	}
	arg.NowMenuId = adminMenu[0].Id

	return
}

// 获取当前菜单的信息
func (r *SrvPermitMenuImpl) getNowMenuImports(arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn) (err error) {

	param := &wrappers.ArgGetImportByMenuId{ArgGetImportByMenuIdSingle: arg.ArgGetImportByMenuIdSingle,}
	param.SuperAdminFlag = arg.IsSuperAdmin

	logContent := map[string]interface{}{
		"arg": arg,
	}

	defer func() {
		if err != nil {
			logContent["err"] = err.Error()
			r.Context.Error(logContent, "SrvPermitMenuImplGetNowMenuImports")
		}
	}()

	//如果是超级管理员
	if param.SuperAdminFlag {
		err = r.orgNowSuperAdminMenu(arg, res, logContent, param)
	} else {

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
func (r *SrvPermitMenuImpl) getPermitMenuListGeneralUser(arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn) (err error) {
	dao := dao_impl.NewDaoPermit(r.Context)

	// 获取接口权限列表
	//if res.OpList, err = NewSrvPermitImport(r.Context).
	//	GetOpList(dao, arg); err != nil {
	//	return
	//}

	var menuIds []int64
	if menuIds, err = r.getPermitByGroupIds(arg.Module, arg.GroupId...); err != nil { // 普通管理员
		return
	}

	if err = r.getGroupMenu(dao, arg, res, menuIds...); err != nil {
		return
	}
	return
}

func (r *SrvPermitMenuImpl) getPermitMenuList(arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn) (err error) {

	if arg.IsSuperAdmin { // 如果是超级管理员
		err = r.getGroupMenuAdmin(arg, res)
		return
	}

	if err = r.getPermitMenuListGeneralUser(arg, res); err != nil {
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
	res.NotReadMsgCount, _ = r.getMessageCount(arg.UserHid)
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
	dao := dao_impl.NewDaoPermit(r.Context)
	var (
		groupPermit          map[int64][]wrappers.ResultPermitMenu
		permitMap            map[int64]wrappers.ResultPermitMenu
		list, systemMenuList []*models.AdminMenu
	)

	// 获取当前选中的系统
	var (
		current []string
	)

	if arg.Module == "" {
		if _, _, arg.Module, err = r.getDefaultSelectModule([]int64{}, true); err != nil {
			return
		}
	}

	current = []string{arg.Module}

	if list, err = dao.GetPermitMenuByIds(current); err != nil {
		return
	} else {
		groupPermit, permitMap = r.groupPermit(list)
	}

	if systemMenuList, err = dao.GetPermitMenuByIds([]string{wrappers.DefaultPermitModule}); err != nil {
		return
	} else {
		// 获取当前用户能够打开的系统列表
		if res.Menu, err = r.getNowUserSystemList(systemMenuList); err != nil {
			return
		}
	}

	if len(res.Menu) <= 0 {
		err = fmt.Errorf("您没有操作权限,请刷新或联系管理员(1)")
		return
	}

	var parentId int64
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
		groupPermit, permitMap = r.groupPermit(list)
	}

	parentId := r.orgActiveMenu(arg, res)

	res.ResultPermitMenu = wrappers.ResultPermitMenu{
		Id:       parentId,
		Children: []wrappers.ResultPermitMenu{},
	}

	r.orgPermit(groupPermit, parentId, &res.ResultPermitMenu, permitMap)

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
func (r *SrvPermitMenuImpl) groupPermit(list []*models.AdminMenu) (groupPermit map[int64][]wrappers.ResultPermitMenu, permitMap map[int64]wrappers.ResultPermitMenu) {
	l := len(list)
	groupPermit = make(map[int64][]wrappers.ResultPermitMenu, l)
	permitMap = make(map[int64]wrappers.ResultPermitMenu, l)

	var data wrappers.ResultPermitMenu

	for _, item := range list {
		if _, ok := groupPermit[item.ParentId]; !ok {
			groupPermit[item.ParentId] = make([]wrappers.ResultPermitMenu, 0, l)
		}
		data = r.orgResultPermitMenu(item)
		permitMap[item.Id] = data
		groupPermit[item.ParentId] = append(groupPermit[item.ParentId], data)
	}
	return
}

func (r *SrvPermitMenuImpl) orgResultPermitMenu(item *models.AdminMenu) (res wrappers.ResultPermitMenu) {
	var f = false
	if item.HideInMenu == models.AdminMenuHideInMenuTrue {
		f = true
	}
	res = wrappers.ResultPermitMenu{
		Id:     item.Id,
		Path:   item.UrlPath,
		Name:   item.PermitKey,
		Module: item.Module,
		Label:  item.Label,
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

func (r *SrvPermitMenuImpl) orgPermit(group map[int64][]wrappers.ResultPermitMenu, pid int64, res *wrappers.ResultPermitMenu, permitMap map[int64]wrappers.ResultPermitMenu) {

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
		Context: r.Context,
		Method:  "GET",
		AppName: parameters.MicroUser,
		Header:  httpHeader,
		URI:     "/in/user/has_not_msg",
		Value:   url.Values{},
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
	var dt []models.AdminMenu
	if dt, err = dao_impl.NewDaoPermitMenu(r.Context).
		GetMenuByCondition(map[string]interface{}{"permit_key": arg.Module}); err != nil {
		return
	}
	if len(dt) == 0 {
		err = fmt.Errorf("您查看的系统(%s)不存在或已删除", arg.Module)
		return
	}
	arg.ParentId = dt[0].Id
	if arg.Module == "" {
		arg.Module = dt[0].PermitKey
	}
	if arg.ArgGetImportByMenuIdSingle.NowModule == "" {
		arg.ArgGetImportByMenuIdSingle.NowModule = dt[0].PermitKey
	}
	return
}

// 获取当前用户组能够使用的系统列表
func (r *SrvPermitMenuImpl) getUserCanUseSystemWithGroupIds(groupIds []int64) (res []*models.AdminMenu, err error) {
	dao := dao_impl.NewDaoPermitGroup(r.Context)

	res, err = dao.GetGroupSystemByGroupId(wrappers.DefaultPermitModule, groupIds...)
	return
}

func NewSrvPermitMenu(ctx ...*base.Context) (res srvs.SrvPermitMenu) {
	p := &SrvPermitMenuImpl{}
	p.SetContext(ctx...)
	return p
}
