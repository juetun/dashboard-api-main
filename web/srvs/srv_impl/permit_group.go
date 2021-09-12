// Package srv_impl
/**
* @Author:changjiang
* @Description:
* @File:PermitGroupImpl
* @Version: 1.0.0
* @Date 2021/6/20 10:03 下午
 */
package srv_impl

import (
	"fmt"
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

type SrvPermitGroupImpl struct {
	base.ServiceBase
}

func (r *SrvPermitGroupImpl) AdminGroupDelete(arg *wrappers.ArgAdminGroupDelete) (res wrappers.ResultAdminGroupDelete, err error) {
	res = wrappers.ResultAdminGroupDelete{}
	dao := dao_impl.NewDaoPermit(r.Context)

	if err = dao.DeleteAdminGroupByIds(arg.IdString...); err != nil {
		return
	}
	daoGroup := dao_impl.NewDaoPermitGroupImpl(r.Context)
	// 删除用户组权限
	if err = daoGroup.DeleteUserGroupPermitByGroupId(arg.IdString...); err != nil {
		return
	}

	res.Result = true
	return
}

func (r *SrvPermitGroupImpl) AdminGroup(arg *wrappers.ArgAdminGroup) (res *wrappers.ResultAdminGroup, err error) {

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

func (r *SrvPermitGroupImpl) orgGroupList(dao daos.DaoPermit, list []models.AdminGroup) (res []wrappers.AdminGroup, err error) {
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

func (r *SrvPermitGroupImpl) MenuImportSet(arg *wrappers.ArgMenuImportSet) (res *wrappers.ResultMenuImportSet, err error) {
	res = &wrappers.ResultMenuImportSet{
		Result: false,
	}
	var m models.AdminMenuImport
	t := time.Now()
	var menuName string
	var menus []models.AdminMenu
	dao := dao_impl.NewDaoPermit(r.Context)
	if menus, err = dao.GetMenu(arg.MenuId); err != nil {
		return
	} else if len(menus) > 0 {
		menuName = menus[0].Module
	}
	var importList []models.AdminImport
	if importList, err = dao.GetAdminImportById(arg.ImportIds...); err != nil {
		return
	}
	var mapImport = make(map[int]string, len(importList))
	for _, value := range importList {
		mapImport[value.Id] = value.AppName
	}

	var dts = make([]models.AdminMenuImport, 0, len(arg.ImportIds))
	var dt models.AdminMenuImport
	if arg.Type == "delete" {
		dt.DeletedAt = &t
		for _, value := range arg.ImportIds {
			if value == 0 {
				continue
			}
			dt = models.AdminMenuImport{
				MenuId:        arg.MenuId,
				MenuModule:    menuName,
				ImportId:      value,
				ImportAppName: "",
				CreatedAt:     t,
				UpdatedAt:     t,
				DeletedAt:     &t,
			}
			dt.ImportAppName, _ = mapImport[dt.ImportId]
			dts = append(dts, dt)
		}
	} else {
		for _, value := range arg.ImportIds {
			if value == 0 {
				continue
			}
			dt = models.AdminMenuImport{
				MenuId:        arg.MenuId,
				MenuModule:    menuName,
				ImportId:      value,
				ImportAppName: "",
				CreatedAt:     t,
				UpdatedAt:     t,
			}
			dt.ImportAppName, _ = mapImport[dt.ImportId]
			dts = append(dts, dt)
		}
	}
	if err = dao_impl.NewPermitImportImpl(r.Context).
		BatchMenuImport(m.TableName(), dts); err != nil {
		return
	}

	res.Result = true
	return
}
func (r *SrvPermitGroupImpl) AdminUserGroupAdd(arg *wrappers.ArgAdminUserGroupAdd) (res wrappers.ResultAdminUserGroupAdd, err error) {
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

func NewSrvPermitGroupImpl(context ...*base.Context) srvs.SrvPermitGroup {
	p := &SrvPermitGroupImpl{}
	p.SetContext(context...)
	return p
}
