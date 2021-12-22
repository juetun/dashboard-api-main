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

func (r *SrvPermitMenuImpl) Menu(arg *wrappers.ArgPermitMenu) (res *wrappers.ResultPermitMenuReturn, err error) {

	res = wrappers.NewResultPermitMenuReturn()
	dao := dao_impl.NewDaoPermit(r.Context)

	// 判断参数是否正常
	if err = r.initParentId(arg); err != nil {
		return
	}

	// 判断当前用户是否是超级管理员,如果不是超级管理员，组织所属组权限
	if err = r.initGroupAndIsSuperAdmin(arg, dao); err != nil {
		return
	}

	if err = r.GetMenuPermitKeyByPath(&arg.ArgGetImportByMenuIdSingle, dao); err != nil {
		return
	}

	handlers := []MenuHandler{
		{
			Name:    "getNotReadMessage",
			handler: r.getNotReadMessage, // 获取用户未读消息数)
		},
		{
			Name:    "getPermitMenuList",
			handler: r.getPermitMenuList, // 获取当前菜单的树形结构
		},
	}
	if arg.NowMenuId != 0 {
		handlers = append(handlers, MenuHandler{
			Name:    "getNowMenuImports",
			handler: r.getNowMenuImports, // 获取指定菜单下的接口ID列表
		})
	}

	// 并行执行逻辑
	r.syncOperate(handlers, arg, res)
	return

}
func (r *SrvPermitMenuImpl) AdminMenu(arg *wrappers.ArgAdminMenu) (res *wrappers.ResultAdminMenu, err error) {

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
func (r *SrvPermitMenuImpl) getSystemIdByModule(dao daos.DaoPermit, module string, sysId int64) (systemId int64, err error) {
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
func (r *SrvPermitMenuImpl) AdminMenuWithCheck(arg *wrappers.ArgAdminMenuWithCheck) (res *wrappers.ResultMenuWithCheck, err error) {
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

//
func (r *SrvPermitMenuImpl) orgTree(list []models.AdminMenu, parentId int64, res *[]wrappers.AdminMenuObject, mapPermitMenu map[int64]int64) {
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

func (r *SrvPermitMenuImpl) orgAdminMenuObject(value *models.AdminMenu, mapPermitMenu map[int64]int64) (res wrappers.AdminMenuObject) {
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


func (r *SrvPermitMenuImpl) GetMenuPermitKeyByPath(arg *wrappers.ArgGetImportByMenuIdSingle, dao daos.DaoPermit) (err error) {
	_ = dao
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
	if adminMenu, err = dao_impl.NewDaoPermit(r.Context).
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

	param := &wrappers.ArgGetImportByMenuId{
		ArgGetImportByMenuIdSingle: arg.ArgGetImportByMenuIdSingle,
	}
	param.SuperAdminFlag = arg.IsSuperAdmin

	logContent := map[string]interface{}{
		"arg": arg,
	}
	defer func() {
		r.Context.Error(logContent)
	}()

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

func (r *SrvPermitMenuImpl) getPermitMenuList(arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn) (err error) {
	dao := dao_impl.NewDaoPermit(r.Context)
	if arg.IsSuperAdmin { // 如果是超级管理员
		err = r.getGroupMenu(dao, arg, res)
		return
	}
	// 获取接口权限列表
	if res.OpList, err = NewSrvPermitImport(r.Context).
		GetOpList(dao, arg); err != nil {
		return
	}

	var menuIds []int64
	if menuIds, err = r.getPermitByGroupIds(arg.Module, arg.GroupId...); err != nil { // 普通管理员
		return
	}

	if err = r.getGroupMenu(dao, arg, res, menuIds...); err != nil {
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

func (r *SrvPermitMenuImpl) getGroupMenu(
	dao daos.DaoPermit,
	arg *wrappers.ArgPermitMenu,
	res *wrappers.ResultPermitMenuReturn,
	menuIds ...int64) (err error) {
	var (
		groupPermit map[int64][]wrappers.ResultPermitMenu
		permitMap   map[int64]wrappers.ResultPermitMenu
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
	res.Menu, groupPermit, permitMap = r.groupPermit(list)
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

func (r *SrvPermitMenuImpl) groupPermit(list []models.AdminMenu) (systemList []wrappers.ResultSystemMenu, groupPermit map[int64][]wrappers.ResultPermitMenu, permitMap map[int64]wrappers.ResultPermitMenu) {
	l := len(list)
	systemList = make([]wrappers.ResultSystemMenu, 0, l)
	groupPermit = make(map[int64][]wrappers.ResultPermitMenu, l)
	permitMap = make(map[int64]wrappers.ResultPermitMenu, l)

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

func (r *SrvPermitMenuImpl) orgResultPermitMenu(item models.AdminMenu) (res wrappers.ResultPermitMenu) {
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

func (r *SrvPermitMenuImpl) getCurrentSystem(dao daos.DaoPermit, arg *wrappers.ArgPermitMenu) (res []string, err error) {
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

func (r *SrvPermitMenuImpl) getMessageCount(userHid string) (count int, err error) {
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

// initGroupAndIsSuperAdmin 判断当前用户是否是超级管理员,如果不是超级管理员，组织所属组权限
func (r *SrvPermitMenuImpl) initGroupAndIsSuperAdmin(arg *wrappers.ArgPermitMenu, dao daos.DaoPermit) (err error) {
	if arg.GroupId, arg.IsSuperAdmin, err = r.getUserGroupIds(&ArgGetUserGroupIds{Dao: dao, UserId: arg.UserHid}); err != nil {
		return
	}
	return
}

func (r *SrvPermitMenuImpl) orgRoutParentMap(parentId int64, list []models.AdminMenu, dataChild *wrappers.ResultPermitMenu) (returnData map[string][]string) {
	_ = dataChild
	l := len(list)
	returnData = make(map[string][]string, l)
	res := make(map[int64][]int64, l)
	mapIdToParent := make(map[int64]int64, l)
	mapIdConfig := make(map[int64]models.AdminMenu, l)
	for _, value := range list {
		mapIdConfig[value.Id] = value
		if value.Id == parentId {
			continue
		}
		if _, ok := res[value.Id]; !ok {
			res[value.Id] = make([]int64, 0, 6)
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

func (r *SrvPermitMenuImpl) getUserGroupIds(arg *ArgGetUserGroupIds) (res []int64, superAdmin bool, err error) {
	res = []int64{}
	var groups []wrappers.AdminGroupUserStruct
	if groups, err = arg.Dao.GetGroupByUserId(arg.UserId); err != nil {
		return
	}
	res = make([]int64, 0, len(groups))
	for _, group := range groups {

		if group.IsSuperAdmin == models.IsAdminGroupYes ||
			group.SuperAdmin == models.IsSuperAdminYes { // 如果是超级管理员
			superAdmin = true
			return
		}
		res = append(res, group.GroupId)
	}
	return
}

func (r *SrvPermitMenuImpl) syncOperate(handlers []MenuHandler, arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn) {
	var syncG sync.WaitGroup
	syncG.Add(len(handlers))
	for _, handler := range handlers {
		go func(h MenuHandler) {
			defer syncG.Done()
			_ = h.handler(arg, res)
		}(handler)
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
	return
}

func NewSrvPermitMenu(ctx ...*base.Context) (res srvs.SrvPermitMenu) {
	p := &SrvPermitMenuImpl{}
	p.SetContext(ctx...)
	return p
}
