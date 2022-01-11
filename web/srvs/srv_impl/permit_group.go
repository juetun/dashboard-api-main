// Package srv_impl

package srv_impl

import (
	"fmt"
	"strconv"
	"time"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/app_param"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"gorm.io/gorm"
)

type SrvPermitGroupImpl struct {
	base.ServiceBase
}

func (r *SrvPermitGroupImpl) AdminGroupDelete(arg *wrappers.ArgAdminGroupDelete) (res wrappers.ResultAdminGroupDelete, err error) {
	res = wrappers.ResultAdminGroupDelete{}
	daoGroup := dao_impl.NewDaoPermitGroup(r.Context)
	if err = daoGroup.DeleteAdminGroupByIds(arg.IdString...); err != nil {
		return
	}
	// 删除用户组权限
	if err = daoGroup.DeleteUserGroupPermitByGroupId(arg.IdString...); err != nil {
		return
	}

	res.Result = true
	return
}

func (r *SrvPermitGroupImpl) AdminGroup(arg *wrappers.ArgAdminGroup) (res *wrappers.ResultAdminGroup, err error) {

	res = &wrappers.ResultAdminGroup{Pager: *response.NewPager(response.PagerBaseQuery(&arg.PageQuery))}

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

func (r *SrvPermitGroupImpl) getParentGroupMap(dao daos.DaoPermit, list []models.AdminGroup) (parentGroupMap map[int64]models.AdminGroup, ids []int64, err error) {
	l := len(list)
	parentGroupMap = make(map[int64]models.AdminGroup, l)
	pIds := make([]int64, 0, l)
	for _, value := range list {
		pIds = append(pIds, value.ParentId)
		ids = append(ids, value.Id)
	}
	var listG []models.AdminGroup
	if listG, err = dao.GetAdminGroupByIds(pIds...); err != nil {
		return
	}
	for _, value := range listG {
		parentGroupMap[value.Id] = value
	}
	return
}

func (r *SrvPermitGroupImpl) orgGroupList(dao daos.DaoPermit, list []models.AdminGroup) (res []wrappers.AdminGroup, err error) {
	res = make([]wrappers.AdminGroup, 0, len(list))

	var (
		ids               []int64
		getParentGroupMap map[int64]models.AdminGroup
		mapGroupUserCount = map[int64]int{}
	)
	if getParentGroupMap, ids, err = r.getParentGroupMap(dao, list); err != nil {
		return
	}
	daoGroup := dao_impl.NewDaoPermitGroup(r.Context)
	if mapGroupUserCount, err = daoGroup.GetGroupUserCount(ids...); err != nil {
		return
	}

	var dt wrappers.AdminGroup
	for _, it := range list {
		dt = wrappers.AdminGroup{
			AdminGroup: it,
		}
		dt.UpdatedAtString = dt.UpdatedAt.Format("2006.01.02 15:04")
		if tp, ok := mapGroupUserCount[it.Id]; ok {
			dt.UserCount = tp
		}
		if dta, ok := getParentGroupMap[it.ParentId]; ok {
			dt.ParentName = dta.Name
		}
		res = append(res, dt)
	}

	return
}

func (r *SrvPermitGroupImpl) getMapImport(importIds ...int64) (mapImport map[int64]models.AdminImport, err error) {
	mapImport = map[int64]models.AdminImport{}
	dao := dao_impl.NewDaoPermit(r.Context)
	var importList []models.AdminImport
	if importList, err = dao.GetAdminImportById(importIds...); err != nil {
		return
	}
	mapImport = make(map[int64]models.AdminImport, len(importList))
	for _, value := range importList {
		mapImport[value.Id] = value
	}
	return
}

func (r *SrvPermitGroupImpl) getMenuNameWithMenuId(menuId int64) (menuName string, err error) {
	var (
		menus []models.AdminMenu
	)
	if menus, err = dao_impl.NewDaoPermitMenu(r.Context).
		GetMenu(menuId); err != nil {
		return
	} else if len(menus) > 0 {
		menuName = menus[0].Module
	}
	return
}

func (r *SrvPermitGroupImpl) menuImportSetNotDelete(arg *wrappers.ArgMenuImportSet, menuName string, mapImport map[int64]models.AdminImport) (err error) {
	var (
		m   models.AdminMenuImport
		t   = time.Now()
		dtm models.AdminImport
		ok  bool
		dts = make([]*models.AdminMenuImport, 0, len(arg.ImportIds))
		dt  *models.AdminMenuImport
	)

	for _, value := range arg.ImportIds {

		if value == 0 {
			continue
		}

		dt = &models.AdminMenuImport{
			MenuId:        arg.MenuId,
			MenuModule:    menuName,
			ImportId:      value,
			ImportAppName: "",
			CreatedAt:     t,
			UpdatedAt:     t,
		}
		if dtm, ok = mapImport[dt.ImportId]; ok {
			dt.ImportAppName = dtm.AppName
			dt.DefaultOpen = dtm.DefaultOpen
		}
		dts = append(dts, dt)
	}
	if err = dao_impl.NewDaoPermitImport(r.Context).
		BatchMenuImport(m.TableName(), dts); err != nil {
		return
	}

	return
}

func (r *SrvPermitGroupImpl) menuImportSetUpdate(arg *wrappers.ArgMenuImportSet, mapImport map[int64]models.AdminImport) (err error) {
	var (
		m   models.AdminMenuImport
		t   = time.Now()
		dtm models.AdminImport
		ok  bool
		dts = make([]*models.AdminMenuImport, 0, len(arg.ImportIds))
		dt  *models.AdminMenuImport
	)
	var defaultOpen uint8
	switch arg.Column {
	case "default_open":
		var defaultOpen64 uint64
		if defaultOpen64, err = strconv.ParseUint(arg.Value, 10, 64); err != nil {
			err = fmt.Errorf("参数格式不正确(%s:%s)", arg.Column, arg.Value)
			return
		}
		if defaultOpen64 > 127 {
			err = fmt.Errorf("参数格式不正确(%s:%s)", arg.Column, arg.Value)
			return
		}
		defaultOpen = uint8(defaultOpen64)
	default:
		err = fmt.Errorf("当前不支持你选择的字段修改")
		return
	}
	dao := dao_impl.NewDaoPermitImport(r.Context)
	var data []*models.AdminMenuImport
	if data, err = dao.GetMenuImportByMenuIdAndImportIds(arg.MenuId, arg.ImportIds...); err != nil {
		return
	}
	for _, dt = range data {
		dt.UpdatedAt = t
		dt.DeletedAt = nil
		switch arg.Column {
		case "default_open":
			dt.DefaultOpen = defaultOpen
		}
		if dtm, ok = mapImport[dt.ImportId]; ok {
			dt.ImportAppName = dtm.AppName
		}
		dts = append(dts, dt)
	}
	if err = dao.BatchMenuImport(m.TableName(), dts); err != nil {
		return
	}

	return
}

func (r *SrvPermitGroupImpl) menuImportSetDelete(arg *wrappers.ArgMenuImportSet, menuName string, mapImport map[int64]models.AdminImport) (err error) {
	var (
		m   models.AdminMenuImport
		t   = time.Now()
		dtm models.AdminImport
		ok  bool
		dts = make([]*models.AdminMenuImport, 0, len(arg.ImportIds))
		dt  *models.AdminMenuImport
	)
	for _, value := range arg.ImportIds {
		if value == 0 {
			continue
		}
		dt = &models.AdminMenuImport{
			MenuId:        arg.MenuId,
			MenuModule:    menuName,
			ImportId:      value,
			ImportAppName: "",
			CreatedAt:     t,
			UpdatedAt:     t,
			DeletedAt:     &t,
		}
		if dtm, ok = mapImport[dt.ImportId]; ok {
			dt.ImportAppName = dtm.AppName
			dt.DefaultOpen = dtm.DefaultOpen
		}
		dts = append(dts, dt)
	}
	if err = dao_impl.NewDaoPermitImport(r.Context).
		BatchMenuImport(m.TableName(), dts); err != nil {
		return
	}

	return
}

func (r *SrvPermitGroupImpl) MenuImportSet(arg *wrappers.ArgMenuImportSet) (res *wrappers.ResultMenuImportSet, err error) {

	res = &wrappers.ResultMenuImportSet{Result: false}

	var (
		menuName  string
		mapImport map[int64]models.AdminImport
	)
	if menuName, err = r.getMenuNameWithMenuId(arg.MenuId); err != nil {
		return
	}

	if mapImport, err = r.getMapImport(arg.ImportIds...); err != nil {
		return
	}

	switch arg.Type {
	case "delete":
		if err = r.menuImportSetDelete(arg, menuName, mapImport); err != nil {
			return
		}
	case "update":
		if err = r.menuImportSetUpdate(arg, mapImport); err != nil {
			return
		}
	default:
		if err = r.menuImportSetNotDelete(arg, menuName, mapImport); err != nil {
			return
		}
	}
	res.Result = true
	return
}

func (r *SrvPermitGroupImpl) validateUserHid(userHid ...int64) (userMap map[int64]app_param.ResultUserItem, err error) {
	if userMap, err = NewUserService(r.Context).
		GetUserByIds(userHid); err != nil {
		return
	}
	if len(userMap) != len(userHid) {
		err = fmt.Errorf("您选择的管理用户数据异常,请尝试刷新页面重试")
	}
	return
}

func (r *SrvPermitGroupImpl) validateGroupIds(dao daos.DaoPermitGroup, groupIds ...int64) (err error) {
	var groups []*models.AdminGroup
	if groups, err = dao.GetGroupByIds(groupIds...); err != nil {
		return
	}
	if len(groups) != len(groupIds) {
		err = fmt.Errorf("您选择的管理员组数据异常,请尝试刷新页面重试")
	}
	return
}

func (r *SrvPermitGroupImpl) AdminUserGroupAdd(arg *wrappers.ArgAdminUserGroupAdd) (res wrappers.ResultAdminUserGroupAdd, err error) {
	res = wrappers.ResultAdminUserGroupAdd{}

	var (
		data    []base.ModelBase
		dao     = dao_impl.NewDaoPermitGroup(r.Context)
		userMap map[int64]app_param.ResultUserItem
		t       = time.Now()
	)
	if err = r.validateGroupIds(dao, arg.GroupIds...); err != nil {
		return
	}
	if userMap, err = r.validateUserHid(arg.UserHIds...); err != nil {
		return
	}
	var adminUsers = make([]base.ModelBase, 0, len(userMap))
	var dt *models.AdminUser
	for _, item := range userMap {
		dt = &models.AdminUser{
			UserHid:   item.UserHid,
			RealName:  item.RealName,
			Mobile:    item.Mobile,
			CreatedAt: t,
			UpdatedAt: t,
		}
		adminUsers = append(adminUsers, dt)
	}

	if err = dao_impl.NewDaoPermitUser(r.Context).
		AdminUserAdd(adminUsers); err != nil {
		return
	}
	if data, err = r.orgAdminUserGroup(arg); err != nil {
		return
	}

	if err = dao.AdminUserGroupAdd(data); err != nil {
		return
	}
	res.Result = true

	return
}

func (r *SrvPermitGroupImpl) orgAdminUserGroup(arg *wrappers.ArgAdminUserGroupAdd) (data []base.ModelBase, err error) {
	data = make([]base.ModelBase, 0, len(arg.UserHIds)*len(arg.GroupIds))

	var (
		dt *models.AdminUserGroup
		t  = base.GetNowTimeNormal()
	)

	for _, userHId := range arg.UserHIds {
		for _, groupId := range arg.GroupIds {
			dt = &models.AdminUserGroup{
				GroupId:   groupId,
				UserHid:   userHId,
				UpdatedAt: t,
				CreatedAt: t,
			}
			data = append(data, dt)
		}
	}

	return
}

func (r *SrvPermitGroupImpl) AdminUserGroupRelease(arg *wrappers.ArgAdminUserGroupRelease) (res wrappers.ResultAdminUserGroupRelease, err error) {
	res = wrappers.ResultAdminUserGroupRelease{}
	dao := dao_impl.NewDaoPermit(r.Context)
	err = dao.AdminUserGroupRelease(arg.IdString...)
	if err != nil {
		return
	}
	res.Result = true
	return
}

func (r *SrvPermitGroupImpl) AdminGroupEdit(arg *wrappers.ArgAdminGroupEdit) (res *wrappers.ResultAdminGroupEdit, err error) {
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
			Name:         arg.Name,
			IsSuperAdmin: arg.IsSuperAdmin,
			IsAdminGroup: arg.IsAdminGroup,
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
	if dt, err = dao.GetAdminGroupByIds([]int64{arg.Id}...); err != nil {
		return
	}
	if len(dt) == 0 {
		err = fmt.Errorf("您要编辑的组不存在或已删除")
		return
	}
	dta = dt[0]
	dta.Name = arg.Name
	dta.UpdatedAt = time.Now()
	dta.IsAdminGroup = arg.IsAdminGroup
	dta.IsSuperAdmin = arg.IsSuperAdmin
	if err = dao.UpdateAdminGroup(map[string]interface{}{
		"name":           arg.Name,
		"updated_at":     base.TimeNormal{Time: time.Now()}.Format(utils.DateTimeGeneral),
		"is_admin_group": arg.IsAdminGroup,
		"is_super_admin": arg.IsSuperAdmin,
	}, map[string]interface{}{
		"id": arg.Id,
	}); err != nil {
		return
	}
	daoUserGroup := dao_impl.NewDaoPermitGroup(r.Context)
	if err = daoUserGroup.UpdateDaoPermitUserGroupByGroupId(arg.Id, map[string]interface{}{
		"is_super_admin": dta.IsSuperAdmin,
		"is_admin_group": dta.IsAdminGroup,
	}); err != nil {
		return
	}
	res.Result = true
	return
}

func (r *SrvPermitGroupImpl) getGroupIdWithAdminUserGroup(uGroup []models.AdminUserGroup) (isAdmin, superAdmin bool, groupId []int64) {
	groupId = make([]int64, 0, len(uGroup))
	for _, group := range uGroup {
		groupId = append(groupId, group.GroupId)
		if group.IsAdminGroup == models.IsAdminGroupYes { // 判断是否为管理员
			isAdmin = true
		}
		if group.IsSuperAdmin == models.IsSuperAdminYes { // 判断你是否为超级管理员
			superAdmin = true
		}
	}
	return
}

// GetUserGroup 获取管理员所在的用户组ID
func (r *SrvPermitGroupImpl) GetUserGroup(userHid int64) (isAdmin, isSuperAdmin bool, groupIds []int64, err error) {
	groupIds = []int64{}

	if userHid == 0 {
		return
	}

	var uGroup []models.AdminUserGroup

	if uGroup, err = dao_impl.NewDaoPermitGroup(r.Context).
		GetPermitGroupByUid(userHid); err != nil {
		return
	}

	// 如果没有在权限组中
	if len(uGroup) == 0 {
		err = fmt.Errorf(models.GatewayErrMap[models.GatewayNotHavePermit])
		r.Context.Error(map[string]interface{}{"err": err.Error(), "userHid": userHid}, "SrvGatewayImportPermitImplGetUserGroup")
		err = base.NewErrorRuntime(err, models.GatewayNotHavePermit)
		return
	}

	isAdmin, isSuperAdmin, groupIds = r.getGroupIdWithAdminUserGroup(uGroup)

	return
}

func NewSrvPermitGroupImpl(context ...*base.Context) srvs.SrvPermitGroup {
	p := &SrvPermitGroupImpl{}
	p.SetContext(context...)
	return p
}
