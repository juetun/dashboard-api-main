// Package dao_impl
/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 11:19 下午
 */
package dao_impl

import (
	"fmt"
	"github.com/juetun/dashboard-api-main/web/daos"
	"strconv"
	"strings"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
	"gorm.io/gorm"
)

type DaoPermitImpl struct {
	base.ServiceDao
}

func (r *DaoPermitImpl) getMenuImportCountDb(db *gorm.DB, arg *wrapper_admin.ArgMenuImport) (dba *gorm.DB) {
	dba = db
	if dba != nil {
		return
	}
	var m models.AdminImport
	dba = r.Context.Db.Table(fmt.Sprintf("%s as a", m.TableName())).
		Scopes(base.ScopesDeletedAt("a"))
	if arg.AppName != "" {
		dba = dba.Where("a.app_name = ?", arg.AppName)
	}
	if arg.MenuId != 0 {
		var b models.AdminMenuImport
		switch arg.HaveSelect {
		case wrapper_admin.ArgMenuImportHaveSelectYes:
			dba = dba.Joins(fmt.Sprintf(" RIGHT JOIN %s as b ON a.id= b.import_id",
				b.TableName())).
				Where("b.menu_id = ?", arg.MenuId).
				Scopes(base.ScopesDeletedAt("b"))
		}
	}
	if arg.UrlPath != "" {
		dba = dba.Where("a.url_path LIKE ?", fmt.Sprintf("%%%s%%", arg.UrlPath))
	}
	return
}

func (r *DaoPermitImpl) MenuImportCount(arg *wrapper_admin.ArgMenuImport, count *int64) (db *gorm.DB, err error) {
	db = r.getMenuImportCountDb(db, arg)
	if err = db.Count(count).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "daoPermitImplMenuImportCount")
		return
	}
	return
}

func (r *DaoPermitImpl) getIIds(data []models.AdminImport) (iIds []int64) {
	iIds = make([]int64, 0, len(data))
	for _, value := range data {
		iIds = append(iIds, value.Id)
	}
	return
}

func (r *DaoPermitImpl) MenuImportList(db *gorm.DB, arg *wrapper_admin.ArgMenuImport) (res []wrappers.ResultMenuImportItem, err error) {
	var (
		data []models.AdminImport
		iIds []int64
	)
	db = r.getMenuImportCountDb(db, arg)
	if err = db.Select("a.*").
		Offset(arg.GetOffset()).
		Limit(arg.PageSize).
		Find(&data).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "daoPermitImplMenuImportList0")
		return
	}
	if arg.MenuId == 0 {
		return
	}

	if iIds = r.getIIds(data); len(iIds) == 0 {
		return
	}

	mv := make(map[int64]models.AdminMenuImport, len(iIds))
	if mv, err = r.getAdminMenuImportByMenuAndImportIds(arg.MenuId, iIds...); err != nil {
		return
	}
	r.orgMenuImportListRes(data, &iIds, mv, &res)
	return
}

func (r *DaoPermitImpl) orgMenuImportListRes(data []models.AdminImport, iIds *[]int64, mv map[int64]models.AdminMenuImport, res *[]wrappers.ResultMenuImportItem) {
	for _, value := range data {
		*iIds = append(*iIds, value.Id)
		dta := wrappers.ResultMenuImportItem{AdminImport: value,}
		if dtm, ok := mv[value.Id]; ok {
			dta.Checked = true
			dta.DefaultOpen = dtm.DefaultOpen
		}
		*res = append(*res, dta)
	}
	return
}

func (r *DaoPermitImpl) getAdminMenuImportByMenuAndImportIds(menuId int, iIds ...int64) (mv map[int64]models.AdminMenuImport, err error) {
	mv = map[int64]models.AdminMenuImport{}
	var dt []models.AdminMenuImport
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m1 *models.AdminMenuImport
		actRes = r.GetDefaultActErrorHandlerResult(m1)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Where("menu_id=? AND import_id IN (?)", menuId, iIds).
			Scopes(base.ScopesDeletedAt()).
			Find(&dt).
			Error
		return
	})
	if err != nil {
		return
	}
	mv = make(map[int64]models.AdminMenuImport, len(dt))
	for _, it := range dt {
		mv[it.ImportId] = it
	}
	return
}

func (r *DaoPermitImpl) fetchDb(arg *wrappers.ArgImportList) (actRes *base.ActErrorHandlerResult) {
	var m *models.AdminImport
	actRes = r.GetDefaultActErrorHandlerResult(m)
 	db := actRes.Db.Table(actRes.Model.TableName())
	db = db.Scopes(base.ScopesDeletedAt())
	if arg.AppName != "" {
		db = db.Where("app_name = ?", arg.AppName)
	}
	if arg.NeedLogin > 0 {
		db = db.Where("need_login = ?", arg.NeedLogin)
	}
	if arg.NeedSign > 0 {
		db = db.Where("need_sign = ?", arg.NeedSign)
	}
	if arg.DefaultOpen > 0 {
		db = db.Where("default_open = ?", arg.DefaultOpen)
	}
	if arg.PermitKey != "" {
		db = db.Where("permit_key = ?", arg.PermitKey)
	}
	if arg.UrlPath != "" {
		db = db.Where("url_path LIKE  ?", fmt.Sprintf("%%%s%%", arg.UrlPath))
	}
	actRes.Db = db
	return
}

func (r *DaoPermitImpl) GetImportListCount(arg *wrappers.ArgImportList) (totalCount int64, actRes *base.ActErrorHandlerResult, err error) {
	actRes = r.fetchDb(arg)
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "daoPermitImplGetImportListCount")
	}()
	err = r.ActErrorHandler(func() (actResChild *base.ActErrorHandlerResult) {
		actResChild = actRes
		actResChild.Err = actResChild.Db.
			Count(&totalCount).
			Error
		return
	})

	return
}
func (r *DaoPermitImpl) GetImportListData(actRes *base.ActErrorHandlerResult, arg *wrappers.ArgImportList, pager *response.Pager) (res []models.AdminImport, err error) {
	actRes = r.fetchDb(arg)
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"arg":   arg,
			"pager": pager,
			"err":   err.Error(),
		}, "daoPermitImplGetImportListData")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	res = []models.AdminImport{}
	err = r.ActErrorHandler(func() (actResChild *base.ActErrorHandlerResult) {
		actResChild = actRes
		actResChild.Err = actResChild.Db.
			Offset(arg.PageQuery.GetOffset()).
			Order(arg.Order).
			Limit(pager.PageSize).
			Find(&res).
			Error

		return
	})

	return
}

func (r *DaoPermitImpl) getPermitImportByModuleWithSuperAdmin(arg *wrappers.ArgPermitMenu) (res []wrappers.Op, err error) {
	var (
		adminImport     models.AdminImport
		adminMenu       models.AdminMenu
		adminMenuImport models.AdminMenuImport
	)
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "DaoPermitImplGetPermitImportByModuleWithSuperAdmin")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.Context.Db.Table(adminImport.TableName()).
		Select("b.permit_key,a.permit_key AS menu_permit_key").
		Joins(fmt.Sprintf("AS b LEFT JOIN  %s as c ON c.import_id=b.id LEFT JOIN %s AS a  ON a.id = c.menu_id",
			adminMenuImport.TableName(),
			adminMenu.TableName())).
		Where("b.deleted_at IS NULL AND a.deleted_at IS NULL AND c.deleted_at IS NULL").Unscoped().
		Find(&res).
		Error
	return
}

func (r *DaoPermitImpl) GetPermitImportByModule(arg *wrappers.ArgPermitMenu) (res []wrappers.Op, err error) {
	var (
		adminUserGroupPermit models.AdminUserGroupImport
		adminImport          models.AdminImport
		adminMenu            models.AdminMenu
		adminMenuImport      models.AdminMenuImport
		desc                 string
	)

	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"arg":  arg,
			"desc": desc,
			"err":  err.Error(),
		}, "DaoPermitImplGetPermitImportByModule")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()

	if arg.IsSuperAdmin { // 如果是超级管理员
		res, err = r.getPermitImportByModuleWithSuperAdmin(arg)
		desc = "getPermitImportByModuleWithSuperAdmin"
		return
	}

	if len(arg.GroupId) == 0 {
		return
	}

	// 普通用户 先通过用户组权限查询权限 再关联 接口 和菜单界面
	if err = r.Context.Db.Table(adminUserGroupPermit.TableName()).
		Select("b.permit_key,d.permit_key AS menu_permit_key").
		Joins(fmt.Sprintf(" AS a LEFT JOIN  %s as b ON a.menu_id=b.id", adminImport.TableName())).
		Joins(fmt.Sprintf(" LEFT JOIN  %s as c ON c.import_id=b.id", adminMenuImport.TableName())).
		Joins(fmt.Sprintf(" LEFT JOIN  %s as d ON c.menu_id=d.id", adminMenu.TableName())).
		Where("  a.group_id IN (?)", arg.GroupId).
		Find(&res).
		Error; err != nil {
		desc = "get_general_group_permit"
		return
	}

	return
}
func (r *DaoPermitImpl) GetImportMenuId(menuId ...int64) (list []models.AdminImport, err error) {
	list = []models.AdminImport{}

	if len(menuId) == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"menuId": menuId,
			"err":    err.Error(),
		}, "daoPermitGetImportMenuId")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminImport
		actRes = r.GetDefaultActErrorHandlerResult(m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where("menu_id IN(?)", menuId).
			Find(&list).
			Error
		return
	})

	return
}

func (r *DaoPermitImpl) CreateImport(data *models.AdminImport) (res bool, err error) {

	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"data": data,
			"err":  err,
		}, "daoPermitCreateImport")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminImport
		actRes = r.GetDefaultActErrorHandlerResult(m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Create(data).
			Error
		return
	})

	return
}
func (r *DaoPermitImpl) UpdateAdminImport(condition, data map[string]interface{}) (res bool, err error) {
	var m *models.AdminImport
	if len(condition) == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"data":      data,
			"err":       err,
		}, "daoPermitUpdateAdminImport")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		actRes = r.GetDefaultActErrorHandlerResult(m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where(condition).
			Updates(data).
			Error
		return
	})
	if err != nil {
		return
	}
	res = true
	return
}
func (r *DaoPermitImpl) GetAdminImportById(id ...int64) (res []models.AdminImport, err error) {
	if len(id) == 0 {
		res = []models.AdminImport{}
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"id":  id,
			"err": err,
		}, "daoPermitGetAdminImportByIdError")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m models.AdminImport
		actRes = r.GetDefaultActErrorHandlerResult(&m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where("id IN(?)", id).
			Find(&res).Error
		return
	})
	return
}
func (r *DaoPermitImpl) getColumnName(s string) (res string) {
	li := strings.Split(s, ";")
	res = s
	for _, s2 := range li {
		if s2 == "" {
			return
		}
		li1 := strings.Split(s2, ":")
		if len(li1) > 1 && li1[0] == "column" {
			res = li1[1]
		}
	}
	return
}

func (r *DaoPermitImpl) GetDefaultImportByMenuIds(pageType, module string, menuId ...int64) (res []models.AdminImport, err error) {
	_ = module
	res = []models.AdminImport{}
	if len(menuId) == 0 {
		return
	}
	var m models.AdminImport
	var db *gorm.DB
	switch pageType {
	case models.PathTypePage:
		var ami models.AdminMenuImport
		db = r.Context.Db.Table(ami.TableName()).
			Where("a.menu_id IN(?) AND `a`.`deleted_at` IS NULL", menuId).
			Unscoped().
			Select("b.*").
			Joins(fmt.Sprintf("AS a LEFT JOIN %s AS b ON b.`id` = a.import_id", m.TableName()))
	case models.PathTypeApi:
		db = r.Context.Db.Table(m.TableName()).Where("`id` IN(?) ", menuId)
	default:
		err = fmt.Errorf("pathtype is error")
		return
	}
	if err = db.Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"menuId": menuId,
			"err":    err,
		}, "GetDefaultImportByMenuIds")
		return
	}
	return
}
func (r *DaoPermitImpl) GetImportId(id int) (res models.AdminImport, err error) {

	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"id":  id,
			"err": err,
		}, "daoPermitGetImportId")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminImport
		actRes = r.GetDefaultActErrorHandlerResult(m)
		var list []models.AdminImport
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Where("id=?", id).
			Find(&list).
			Error
		if actRes.Err != nil {
			return
		}
		if len(list) == 0 {
			return
		}
		res = list[0]
		return
	})

	return
}
func (r *DaoPermitImpl) GetImportCount(arg *wrappers.ArgGetImport, count *int64) (db *gorm.DB, err error) {
	if arg.MenuId == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "daoPermitGetImportCount error")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	db = r.getImportListDb(nil, arg)
	err = db.Count(count).Error

	return
}
func (r *DaoPermitImpl) getImportListDb(db *gorm.DB, arg *wrappers.ArgGetImport) (res *gorm.DB) {
	var m models.AdminMenuImport
	var mi models.AdminImport
	res = db
	if db == nil {
		res = r.Context.Db.Table(m.TableName()).
			Joins(fmt.Sprintf("AS a LEFT JOIN %s AS b ON a.import_id=b.id", mi.TableName())).
			Where("a.`menu_id` = ? AND a.`deleted_at` IS NULL", arg.MenuId).
			Unscoped()

		if arg.AppName != "" {
			res = res.Where("b.`app_name` = ?", arg.AppName)
		}
		if arg.NeedLogin > 0 {
			res = res.Where("b.need_login = ?", arg.NeedLogin)
		}
		if arg.NeedSign > 0 {
			res = res.Where("b.need_sign = ?", arg.NeedSign)
		}
		if arg.DefaultOpen > 0 {
			res = res.Where("b.default_open = ?", arg.DefaultOpen)
		}
		if arg.PermitKey != "" {
			res = res.Where("b.permit_key = ?", arg.PermitKey)
		}
		if arg.UrlPath != "" {
			res = res.Where("b.url_path LIKE ?", fmt.Sprintf("%%%s%%", arg.UrlPath))
		}
		// PermitKey   string `json:"permit_key" form:"permit_key"`
		// AppName     string `json:"app_name" form:"app_name"`
		// DefaultOpen uint8  `json:"default_open" form:"default_open"`
		// NeedLogin   uint8  `json:"need_login" form:"need_login"`
		// NeedSign    uint8  `json:"need_sign" form:"need_sign"`
		// UrlPath     string `json:"url_path" form:"url_path"`
	}

	return
}
func (r *DaoPermitImpl) GetImportList(db *gorm.DB, arg *wrappers.ArgGetImport) (res []models.AdminImport, err error) {
	db = r.getImportListDb(db, arg)
	err = db.Select("b.*").Offset(arg.PageQuery.GetOffset()).
		Limit(arg.PageSize).
		Find(&res).Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err,
		}, "daoPermitGetImportList")
		return
	}
	return
}

func (r *DaoPermitImpl) AdminUserGroupRelease(ids ...string) (err error) {
	if len(ids) == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"ids": ids,
			"err": err.Error(),
		}, "daoPermitAdminUserGroupRelease")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminUserGroup
		actRes = r.GetDefaultActErrorHandlerResult(m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Where("id IN (?) ", ids).
			Delete(&models.AdminUserGroup{}).
			Error
		return
	})
	return
}

func (r *DaoPermitImpl) FetchByName(name string) (res models.AdminGroup, err error) {

	res = models.AdminGroup{}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"name": name,
			"err":  err.Error(),
		}, "daoPermitFetchByName")
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminGroup
		actRes = r.GetDefaultActErrorHandlerResult(m)
		var list []models.AdminGroup
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Where("name = ?", name).
			Find(&list).
			Error
		if actRes.Err != nil {
			return
		}
		if len(list) > 0 {
			res = list[0]
		}
		return
	})

	return
}

func (r *DaoPermitImpl) InsertAdminGroup(group *models.AdminGroup) (err error) {

	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"group": group,
			"err":   err.Error(),
		}, "daoPermitInsertAdminGroup")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminGroup
		actRes = r.GetDefaultActErrorHandlerResult(m)
		actRes.Err = actRes.Db.Table(actRes.TableName).Create(group).
			Error
		return
	})

	return
}

func (r *DaoPermitImpl) UpdateAdminGroup(data, where map[string]interface{}) (err error) {

	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"data":  data,
			"where": where,
			"err":   err.Error(),
		}, "daoPermitUpdateAdminGroup")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var group *models.AdminGroup
		actRes = r.GetDefaultActErrorHandlerResult(group)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Where(where).
			Scopes(base.ScopesDeletedAt()).
			Updates(data).
			Error
		return
	})

	return
}
func (r *DaoPermitImpl) GetAdminGroupByIds(gIds ...int64) (res []models.AdminGroup, err error) {
	if len(gIds) == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"gIds": gIds,
			"err":  err.Error(),
		}, "daoPermitGetAdminGroupByIds")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminGroup
		actRes = r.GetDefaultActErrorHandlerResult(m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where("id IN (?)", gIds).
			Find(&res).
			Error
		return
	})

	return
}

func (r *DaoPermitImpl) GetUserGroupByUIds(uIds ...int64) (res []models.AdminUserGroup, err error) {
	if len(uIds) == 0 {
		return
	}
	logContent := map[string]interface{}{"uIds": uIds}
	defer func() {
		if err != nil {
			logContent["err"] = err.Error()
			r.Context.Error(logContent, "DaoPermitImplGetUserGroupByUIds")
		}
	}()

	err = r.ActErrorHandler(func() (actErrorHandlerResult *base.ActErrorHandlerResult) {
		var m models.AdminUserGroup
		actErrorHandlerResult = r.GetDefaultActErrorHandlerResult(&m)

		actErrorHandlerResult.Err = actErrorHandlerResult.Db.
			Table(actErrorHandlerResult.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where(" user_hid IN (?)", uIds).
			Find(&res).
			Error
		return
	})
	return
}

func (r *DaoPermitImpl) DeleteMenuByIds(ids ...string) (err error) {
	if len(ids) == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"ids": ids,
			"err": err,
		}, "daoPermitDeleteByIds")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminMenu
		actRes = r.GetDefaultActErrorHandlerResult(m)
		actRes.Err = actRes.Db.
			Table(actRes.TableName).
			Unscoped().
			Where("id in (?)", ids).
			Delete(&models.AdminMenu{}).
			Error
		return
	})
	return
}

func (r *DaoPermitImpl) GetAdminMenuList(arg *wrappers.ArgAdminMenu) (res []models.AdminMenu, err error) {
	res = []models.AdminMenu{}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "daoPermitGetAdminMenuList")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminMenu
		actRes = r.GetDefaultActErrorHandlerResult(m)
		if arg.SystemId > 0 {
			var list []*models.AdminMenu
			if actRes.Err = actRes.Db.Table(actRes.TableName).
				Scopes(base.ScopesDeletedAt()).
				Order("sort_value desc").
				Where("id = ?", arg.SystemId).
				Order("`sort_value` DESC").
				Find(&list).Error; actRes.Err != nil {
				actRes.Err = base.NewErrorRuntime(actRes.Err, base.ErrorSqlCode)
				return
			}
			if len(list) == 0 {
				actRes.Err = fmt.Errorf("你查看权限系统不存在或已删除")
				return
			}
			m = list[0]
			arg.Module = m.PermitKey
		}
		res, actRes.Err = r.getAdminMenuList(arg)
		return
	})
	return
}

func (r *DaoPermitImpl) getAdminMenuList(arg *wrappers.ArgAdminMenu) (res []models.AdminMenu, err error) {

	var m *models.AdminMenu
	actRes := r.GetDefaultActErrorHandlerResult(m)
	dba := actRes.Db.Table(actRes.TableName)
	if arg.Label = strings.TrimSpace(arg.Label); arg.Label != "" {
		dba = dba.Where("label LIKE ?", "%"+arg.Label+"%")
	}
	if arg.ParentId > 0 {
		dba = dba.Where("parent_id = ?", arg.ParentId)
	}
	if arg.Module != "" {
		dba = dba.Where("`module` = ? OR `parent_id`=?", arg.Module, wrappers.DefaultPermitParentId)
	}
	if arg.AppName != "" {
		dba = dba.Where("`app_name` = ?", arg.AppName)
	}
	if arg.IsMenuShow != -1 {
		dba = dba.Where("`hide_in_menu` = ?", arg.IsMenuShow)
	}
	if arg.IsDel == -1 || arg.IsDel == 0 {
		dba = dba.Where("`deleted_at` IS NULL")
	} else {
		dba = dba.Where("`deleted_at` NOT NULL")
	}
	if arg.Id != 0 {
		dba = dba.Where("`id` = ?", arg.Id)
	}
	err = dba.
		Order("`sort_value` DESC").
		Find(&res).
		Error
	return
}
func (r *DaoPermitImpl) GetAdminGroupCount(db *gorm.DB, arg *wrappers.ArgAdminGroup) (total int64, dba *gorm.DB, err error) {
	var m models.AdminGroup
	if db == nil {
		db = r.Context.Db
	}
	dba = db.Table(m.TableName()).
		Scopes(base.ScopesDeletedAt())
	if arg.Name != "" {
		dba = dba.Where("name LIKE ?", "%"+arg.Name+"%")
	}
	var groupId int
	var err1 error
	if arg.GroupId != "" {
		// 如果参数不是整数，则直接返回没有数据
		if groupId, err1 = strconv.Atoi(arg.GroupId); err1 != nil {
			total = 0
			err = fmt.Errorf("您选择的ID格式不正确")
			return
		}
		if groupId != 0 {
			dba = dba.Where("id = ?", groupId)
		}
	}
	if err = dba.Count(&total).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err,
		}, "daoPermitGetAdminGroupCount")
		return
	}
	return
}
func (r *DaoPermitImpl) GetAdminGroupList(db *gorm.DB, arg *wrappers.ArgAdminGroup, pagerObject *response.Pager) (res []models.AdminGroup, err error) {
	res = []models.AdminGroup{}
	if err = db.Limit(pagerObject.PageSize).
		Offset(pagerObject.GetFromAndLimit()).
		Order("updated_at desc").
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"arg":         arg,
			"pagerObject": pagerObject,
			"err":         err,
		}, "daoPermitGetAdminGroupList")
		return
	}
	return
}

func (r *DaoPermitImpl) GetPermitMenuByIds(module []string, menuIds ...int64) (res []*models.AdminMenu, err error) {
	res = []*models.AdminMenu{}

	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"menuIds": menuIds,
			"err":     err,
		}, "daoPermitGetPermitMenuByIds")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminMenu
		actRes = r.GetDefaultActErrorHandlerResult(m)
		db := actRes.Db.
			Table(actRes.TableName)
		var haveWhere bool
		// 兼容超级管理员和普通管理员
		if len(menuIds) != 0 {
			db = db.Where("id IN(?)", menuIds)
			haveWhere = true
		}
		if len(module) > 0 {
			db = db.Where("module IN(?)", module)
			haveWhere = true
		}
		if !haveWhere {
			return
		}
		actRes.Err = db.Order("sort_value desc").
			Find(&res).Error
		return
	})

	return
}

func NewDaoPermit(c ...*base.Context) daos.DaoPermit {
	p := &DaoPermitImpl{}
	p.SetContext(c...)
	return p
}
