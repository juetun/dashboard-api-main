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

	"github.com/jinzhu/gorm"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type PermitService struct {
	base.ServiceBase
}

func NewPermitService(context ...*base.Context) (p *PermitService) {
	p = &PermitService{}
	p.SetContext(context...)
	return
}

func (r *PermitService) GetMenu(arg *wrappers.ArgGetMenu) (res wrappers.ResultGetMenu, err error) {
	res = wrappers.ResultGetMenu{}
	if arg.MenuId == 0 {
		return
	}
	res.AdminMenu, err = dao_impl.NewDaoPermit(r.Context).GetMenu(arg.MenuId)
	if err != nil {
		return
	}
	return
}
func (r *PermitService) AdminMenuSearch(arg *wrappers.ArgAdminMenu) (res wrappers.ResAdminMenuSearch, err error) {
	res = wrappers.ResAdminMenuSearch{
		List: []models.AdminMenu{},
	}
	dao := dao_impl.NewDaoPermit(r.Context)
	res.List, err = dao.GetAdminMenuList(arg)
	return
}
func (r *PermitService) AdminUserGroupAdd(arg *wrappers.ArgAdminUserGroupAdd) (res wrappers.ResultAdminUserGroupAdd, err error) {
	res = wrappers.ResultAdminUserGroupAdd{}
	dao := dao_impl.NewDaoPermit(r.Context)

	var args = make([]map[string]interface{}, 0)
	for _, userHId := range arg.UserHIds {
		for _, groupId := range arg.GroupIds {
			args = append(args, map[string]interface{}{
				"group_id": groupId,
				"user_hid": userHId,
				"is_del":   0,
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
func (r *PermitService) AdminUserGroupRelease(arg *wrappers.ArgAdminUserGroupRelease) (res wrappers.ResultAdminUserGroupRelease, err error) {
	res = wrappers.ResultAdminUserGroupRelease{}
	dao := dao_impl.NewDaoPermit(r.Context)
	err = dao.AdminUserGroupRelease(arg.IdString)
	if err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitService) AdminUserDelete(arg *wrappers.ArgAdminUserDelete) (res wrappers.ResultAdminUserDelete, err error) {
	res = wrappers.ResultAdminUserDelete{}
	dao := dao_impl.NewDaoPermit(r.Context)
	err = dao.AdminUserDelete(arg.IdString)
	if err != nil {
		return
	}
	res.Result = true
	return
}

func (r *PermitService) AdminUserAdd(arg *wrappers.ArgAdminUserAdd) (res wrappers.ResultAdminUserAdd, err error) {
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
				DeletedAt: nil,
			},
		}); err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitService) AdminGroupDelete(arg *wrappers.ArgAdminGroupDelete) (res wrappers.ResultAdminGroupDelete, err error) {
	res = wrappers.ResultAdminGroupDelete{}
	dao := dao_impl.NewDaoPermit(r.Context)
	err = dao.DeleteAdminGroupByIds(arg.IdString)
	if err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitService) AdminGroupEdit(arg *wrappers.ArgAdminGroupEdit) (res *wrappers.ResultAdminGroupEdit, err error) {
	res = &wrappers.ResultAdminGroupEdit{}
	dao := dao_impl.NewDaoPermit(r.Context)
	var g models.AdminGroup
	if g, err = dao.FetchByName(arg.Name); err != nil {
		return
	}
	if arg.Id == 0 {
		if g.Name != "" {
			err = fmt.Errorf("您输入的组名已存在")
			return
		}
		if err = dao.InsertAdminGroup(&models.AdminGroup{
			Name: arg.Name,
		}); err != nil {
			return
		}
		res.Result = true
		return
	}

	if g.Name != "" && g.Id != arg.Id {
		err = fmt.Errorf("您输入的组名已存在")
		return
	}

	if dao.UpdateAdminGroup(&models.AdminGroup{
		Name: arg.Name,
		Id:   arg.Id,
	}); err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitService) MenuAdd(arg *wrappers.ArgMenuAdd) (res *wrappers.ResultMenuAdd, err error) {
	res = &wrappers.ResultMenuAdd{}
	dao := dao_impl.NewDaoPermit(r.Context)
	t := time.Now()
	var list []models.AdminMenu
	if list, err = dao.GetByCondition(map[string]interface{}{
		"label":     arg.Label,
		"module":    arg.Module,
		"parent_id": arg.ParentId,
	}); err != nil {
		return
	} else if len(list) > 0 {
		err = fmt.Errorf("您输入的菜单名已存在")
		return
	}
	err = dao.Add(&models.AdminMenu{
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
		IsDel:              0,
	})
	res.Result = true
	return
}

func (r *PermitService) createImport(dao *dao_impl.DaoPermit, arg *wrappers.ArgEditImport) (res bool, err error) {
	t := time.Now()
	data := models.AdminImport{
		MenuId:        arg.MenuId,
		AppName:       arg.AppName,
		AppVersion:    arg.AppVersion,
		UrlPath:       arg.UrlPath,
		RequestMethod: strings.Join(arg.RequestMethod, ","),
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
func (r *PermitService) DeleteImport(arg *wrappers.ArgDeleteImport) (res *wrappers.ResultDeleteImport, err error) {
	res = &wrappers.ResultDeleteImport{}
	dao := dao_impl.NewDaoPermit(r.Context)
	if err = dao.DeleteImportByIds([]int{arg.ID}...); err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitService) EditImport(arg *wrappers.ArgEditImport) (res *wrappers.ResultEditImport, err error) {
	res = &wrappers.ResultEditImport{Result: false}
	dao := dao_impl.NewDaoPermit(r.Context)

	var listImport []models.AdminImport
	if listImport, err = dao.GetImportMenuId(arg.MenuId); err != nil {
		return
	}
	for _, value := range listImport {
		if arg.Id == 0 {
			if value.AppName == arg.AppName && value.UrlPath == arg.UrlPath {
				err = fmt.Errorf("您输入的接口信息已存在")
				return
			}
		} else {
			if value.AppName == arg.AppName && value.UrlPath == arg.UrlPath {
				if arg.Id != value.Id {
					err = fmt.Errorf("您输入的接口信息已存在")
					return
				}
			}
		}
	}

	if arg.Id == 0 { // 如果是添加接口
		res.Result, err = r.createImport(dao, arg)
		return
	}

	if _, err = dao.UpdateAdminImport(map[string]interface{}{"id": arg.Id}, map[string]interface{}{
		`menu_id`:        arg.MenuId,
		`app_name`:       arg.AppName,
		`app_version`:    arg.AppVersion,
		`url_path`:       arg.UrlPath,
		`request_method`: strings.Join(arg.RequestMethod, ","),
		`sort_value`:     arg.SortValue,
		`updated_at`:     arg.RequestTime,
	}); err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitService) GetImport(arg *wrappers.ArgGetImport) (res *wrappers.ResultGetImport, err error) {
	res = &wrappers.ResultGetImport{Pager: base.Pager{ReqPager: arg.ReqPager, List: []models.AdminImport{}, TotalCount: 0,},}
	dao := dao_impl.NewDaoPermit(r.Context)
	var db *gorm.DB
	if db, err = dao.GetImportCount(arg, &res.TotalCount); err != nil {
		return
	}
	if res.TotalCount > 0 {
		if res.List, err = dao.GetImportList(db, arg); err != nil {
			return
		}
	}

	// []models.AdminImport{}
	return
}
func (r *PermitService) MenuDelete(arg *wrappers.ArgMenuDelete) (res *wrappers.ResultMenuDelete, err error) {
	res = &wrappers.ResultMenuDelete{}
	dao := dao_impl.NewDaoPermit(r.Context)
	err = dao.DeleteByIds(arg.IdValue)
	res.Result = true
	return
}

func (r *PermitService) MenuSave(arg *wrappers.ArgMenuSave) (res *wrappers.ResultMenuSave, err error) {
	res = &wrappers.ResultMenuSave{}
	dao := dao_impl.NewDaoPermit(r.Context)

	var list []models.AdminMenu
	if list, err = dao.GetByCondition(map[string]interface{}{
		"label":     arg.Label,
		"module":    arg.Module,
		"parent_id": arg.ParentId,
	}); err != nil {
		return
	} else if len(list) > 0 && arg.Id != list[0].Id {
		err = fmt.Errorf("您输入的菜单名已存在")
		return
	}

	t := time.Now()
	err = dao.Save(arg.Id, &models.AdminMenu{
		PermitKey:          arg.PermitKey,
		ParentId:           arg.ParentId,
		Label:              arg.Label,
		Icon:               arg.Icon,
		HideInMenu:         arg.HideInMenu,
		ManageImportPermit: arg.ManageImportPermit,
		UrlPath:            arg.UrlPath,
		SortValue:          arg.SortValue,
		OtherValue:         arg.OtherValue,
		UpdatedAt:          t,
		IsDel:              0,
	})
	res.Result = true
	return
}
func (r *PermitService) AdminMenu(arg *wrappers.ArgAdminMenu) (res *wrappers.ResultAdminMenu, err error) {
	res = &wrappers.ResultAdminMenu{
		List: make([]wrappers.AdminMenuObject, 0, 20),
		Menu: make([]wrappers.ResultSystemAdminMenu, 0, 30),
	}
	dao := dao_impl.NewDaoPermit(r.Context)
	var list []models.AdminMenu
	list, err = dao.GetAdminMenuList(arg)

	var data wrappers.ResultSystemAdminMenu
	for _, item := range list {
		if item.ParentId == 0 {
			data = wrappers.ResultSystemAdminMenu{
				Id:        item.Id,
				PermitKey: item.PermitKey,
				Label:     item.Label,
				Icon:      item.Icon,
				SortValue: item.SortValue,
				Module:    item.Module,
			}
			if item.Id == arg.SystemId {
				data.Active = true
			}
			res.Menu = append(res.Menu, data)
		}
	}
	r.orgTree(list, wrappers.DefaultPermitParentId, &res.List)
	return
}

//
func (r *PermitService) orgTree(list []models.AdminMenu, parentId int, res *[]wrappers.AdminMenuObject) () {
	var tmp wrappers.AdminMenuObject
	for _, value := range list {
		// 剔除默认数据那条
		if value.Id == parentId || value.ParentId != parentId {
			continue
		}
		tmp = r.orgAdminMenuObject(&value)
		r.orgTree(list, value.Id, &tmp.Children)
		*res = append(*res, tmp)
	}
}

func (r *PermitService) orgAdminMenuObject(value *models.AdminMenu) (res wrappers.AdminMenuObject) {
	res = wrappers.AdminMenuObject{Children: make([]wrappers.AdminMenuObject, 0, 20),}
	res.ResultAdminMenuSingle = wrappers.ResultAdminMenuSingle{
		Id:                 value.Id,
		ParentId:           value.ParentId,
		Title:              value.Label,
		Label:              value.Label,
		Icon:               value.Icon,
		HideInMenu:         value.HideInMenu,
		UrlPath:            value.UrlPath,
		SortValue:          value.SortValue,
		IsDel:              value.IsDel,
		Module:             value.Module,
		PermitKey:          value.PermitKey,
		ManageImportPermit: value.ManageImportPermit,
	}
	if value.OtherValue != "" {
		json.Unmarshal([]byte(value.OtherValue), &res.ResultAdminMenuSingle.ResultAdminMenuOtherValue)
	} else {
		res.ResultAdminMenuSingle.ResultAdminMenuOtherValue = wrappers.ResultAdminMenuOtherValue{Expand: true,}
	}
	return
}
func (r *PermitService) AdminSetPermit(arg *wrappers.ArgAdminSetPermit) (res *wrappers.ResultAdminSetPermit, err error) {
	res = &wrappers.ResultAdminSetPermit{}

	list := make([]models.AdminUserGroupPermit, 0, len(arg.PermitIds))
	var dt models.AdminUserGroupPermit
	dao := dao_impl.NewDaoPermit(r.Context)
	if err = dao.DeleteGroupPermit(arg.GroupId); err != nil {
		return
	}
	if len(arg.PermitIds) == 0 {
		return
	}
	var t = time.Now()
	var listImport []models.AdminImport
	if listImport, err = dao.GetDefaultOpenImportByMenuIds(arg.PermitIds...); err != nil {
		return
	} else {
		for _, importData := range listImport {
			dt = models.AdminUserGroupPermit{
				GroupId:   arg.GroupId,
				MenuId:    importData.Id,
				PathType:  models.PathTypeApi,
				CreatedAt: t,
				UpdatedAt: t,
			}
			list = append(list, dt)
		}
	}

	for _, pid := range arg.PermitIds {
		dt = models.AdminUserGroupPermit{
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
	res.Result = true
	return
}

func (r *PermitService) AdminGroup(arg *wrappers.ArgAdminGroup) (res *wrappers.ResultAdminGroup, err error) {

	res = &wrappers.ResultAdminGroup{Pager: *response.NewPagerAndDefault(&arg.BaseQuery),}

	var db *gorm.DB
	dao := dao_impl.NewDaoPermit(r.Context)
	// 获取分页数据
	res.Pager.CallGetPagerData(func(pagerObject *response.Pager) (err error) {
		pagerObject.TotalCount, db, err = dao.GetAdminGroupCount(db, arg)
		return
	}, func(pagerObject *response.Pager) (err error) {
		pagerObject.List, err = dao.GetAdminGroupList(db, arg, pagerObject)
		return
	})
	return
}

func (r *PermitService) AdminUser(arg *wrappers.ArgAdminUser) (res *wrappers.ResultAdminUser, err error) {
	res = &wrappers.ResultAdminUser{Pager: *response.NewPagerAndDefault(&arg.BaseQuery),}
	dao := dao_impl.NewDaoPermit(r.Context)

	var db *gorm.DB

	// 分页获取数据
	res.Pager.CallGetPagerData(func(pagerObject *response.Pager) (err error) {
		pagerObject.TotalCount, db, err = dao.GetAdminUserCount(db, arg)
		return
	}, func(pagerObject *response.Pager) (err error) {
		list, err := dao.GetAdminUserList(db, arg, pagerObject)
		if err != nil {
			return
		}
		pagerObject.List, err = r.leftAdminUser(list, dao)
		return
	})
	return
}

func (r *PermitService) getUserGroup(list []models.AdminUser, dao *dao_impl.DaoPermit) (res map[string][]wrappers.AdminUserGroupName, err error) {
	res = map[string][]wrappers.AdminUserGroupName{}
	uIds := make([]string, 0, len(list))
	for _, value := range list {
		uIds = append(uIds, value.UserHid)
	}
	listResult, err := dao.GetUserGroupByUIds(uIds)
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
	groupList, err := dao.GetAdminGroupByIds(gIds)
	if err != nil {
		return
	}

	groupMap := make(map[int]string, len(groupList))
	for _, value := range groupList {
		groupMap[value.Id] = value.Name
	}

	for _, item := range listResult {
		tmp = wrappers.AdminUserGroupName{
			AdminUserGroup: item,
		}
		if _, ok := groupMap[item.GroupId]; ok {
			tmp.GroupName = groupMap[item.GroupId]
		}
		res[item.UserHid] = append(res[item.UserHid], tmp)
	}
	return
}

func (r *PermitService) leftAdminUser(list []models.AdminUser, dao *dao_impl.DaoPermit) (res []wrappers.ResultAdminUserList, err error) {
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
func (r *PermitService) Menu(arg *wrappers.ArgPermitMenu) (res *wrappers.ResultPermitMenuReturn, err error) {
	var groupIds []int
	var isSuperAdmin bool
	res = &wrappers.ResultPermitMenuReturn{
		ResultPermitMenu: wrappers.ResultPermitMenu{
			Children: []wrappers.ResultPermitMenu{},
		},
		RoutParentMap: map[string][]string{},
		Menu:          []wrappers.ResultSystemMenu{},
	}

	dao := dao_impl.NewDaoPermit(r.Context)
	if groupIds, isSuperAdmin, err = r.getUserGroupIds(dao, arg); err != nil {
		return
	}

	if isSuperAdmin { // 如果是超级管理员
		err = r.getGroupMenu(dao, arg, res)
		return
	}
	var menuIds []int
	if menuIds, err = r.getPermitByGroupIds(dao, arg.PathTypes, groupIds...); err != nil { // 普通管理员
		return
	} else if err = r.getGroupMenu(dao, arg, res, menuIds...); err != nil {
		return
	}
	return
}

func (r *PermitService) getGroupMenu(dao *dao_impl.DaoPermit, arg *wrappers.ArgPermitMenu, res *wrappers.ResultPermitMenuReturn, menuIds ...int) (err error) {

	list, err := dao.GetPermitMenuByIds(menuIds...)
	if err != nil {
		return
	}
	var (
		groupPermit map[int][]wrappers.ResultPermitMenu
		permitMap   map[int]wrappers.ResultPermitMenu
	)
	res.Menu, groupPermit, permitMap = r.groupPermit(list)

	if len(res.Menu) <= 0 {
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
func (r *PermitService) orgRoutParentMap(parentId int, list []models.AdminMenu, dataChild *wrappers.ResultPermitMenu) (returnData map[string][]string) {
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
func (r *PermitService) getParentId(mapIdToParent map[int]int, nodeId, parentId int, value *[]int) {
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
func (r *PermitService) groupPermit(list []models.AdminMenu) (systemList []wrappers.ResultSystemMenu, groupPermit map[int][]wrappers.ResultPermitMenu, permitMap map[int]wrappers.ResultPermitMenu) {
	systemList = make([]wrappers.ResultSystemMenu, 0, 30)
	groupPermit = map[int][]wrappers.ResultPermitMenu{}
	l := len(list)
	permitMap = make(map[int]wrappers.ResultPermitMenu, l)
	var data wrappers.ResultPermitMenu
	for _, item := range list {

		if item.ParentId == 0 {
			systemList = append(systemList, wrappers.ResultSystemMenu{
				Id:        item.Id,
				PermitKey: item.PermitKey,
				Label:     item.Label,
				Icon:      item.Icon,
				SortValue: item.SortValue,
				Module:    item.Module,
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

func (r *PermitService) orgPermit(group map[int][]wrappers.ResultPermitMenu, pid int, res *wrappers.ResultPermitMenu, permitMap map[int]wrappers.ResultPermitMenu) {

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
func (r *PermitService) orgResultPermitMenu(item models.AdminMenu) (res wrappers.ResultPermitMenu) {
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

func (r *PermitService) getPermitByGroupIds(dao *dao_impl.DaoPermit, pathType []string, groupIds ...int) (menuIds []int, err error) {
	res, err := dao.GetMenuIdsByPermitByGroupIds(pathType, groupIds...)
	if err != nil {
		return
	}
	menuIds = make([]int, 0, len(res))
	for _, value := range res {
		menuIds = append(menuIds, value.MenuId)
	}
	return
}

func (r *PermitService) getUserGroupIds(dao *dao_impl.DaoPermit, arg *wrappers.ArgPermitMenu) (res []int, superAdmin bool, err error) {
	groups, err := dao.GetGroupByUserId(arg.UserId)
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

func (r *PermitService) Flag(arg *wrappers.ArgFlag) (res *wrappers.ResultFlag, err error) {
	res = &wrappers.ResultFlag{}
	return
}
