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
	"strconv"
	"strings"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/daos"
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

func (r *DaoPermitImpl) MenuImportList(db *gorm.DB, arg *wrapper_admin.ArgMenuImport) (res []wrappers.ResultMenuImportItem, err error) {
	var data []models.AdminImport
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
	iIds := make([]int64, 0, len(data))
	for _, value := range data {
		iIds = append(iIds, value.Id)
	}
	if len(iIds) == 0 {
		return
	}
	var m1 models.AdminMenuImport
	var dt []models.AdminMenuImport
	if err = r.Context.Db.Table(m1.TableName()).
		Where("menu_id=? AND import_id IN (?)", arg.MenuId, iIds).
		Scopes(base.ScopesDeletedAt()).
		Find(&dt).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"arg":  arg,
			"iIds": iIds,
			"err":  err.Error(),
		}, "daoPermitImplMenuImportList1")
		return
	}
	var mv = make(map[int64]models.AdminMenuImport, len(iIds))
	for _, it := range dt {
		mv[it.ImportId] = it
	}

	var dta wrappers.ResultMenuImportItem
	for _, value := range data {
		iIds = append(iIds, value.Id)
		dta = wrappers.ResultMenuImportItem{
			AdminImport: value,
		}
		if dtm, ok := mv[value.Id]; ok {
			dta.Checked = true
			dta.DefaultOpen = dtm.DefaultOpen
		}
		res = append(res, dta)
	}
	return
}

func NewDaoPermit(c ...*base.Context) daos.DaoPermit {
	p := &DaoPermitImpl{}
	p.SetContext(c...)
	return p
}

func (r *DaoPermitImpl) fetchDb(db *gorm.DB, arg *wrappers.ArgImportList) (dba *gorm.DB) {
	var m models.AdminImport
	dba = db
	if db == nil {
		dba = r.Context.Db.Table(m.TableName())
	}
	dba = dba.Scopes(base.ScopesDeletedAt())
	if arg.AppName != "" {
		dba = dba.Where("app_name = ?", arg.AppName)
	}
	if arg.NeedLogin > 0 {
		dba = dba.Where("need_login = ?", arg.NeedLogin)
	}
	if arg.NeedSign > 0 {
		dba = dba.Where("need_sign = ?", arg.NeedSign)
	}
	if arg.DefaultOpen > 0 {
		dba = dba.Where("default_open = ?", arg.DefaultOpen)
	}
	if arg.PermitKey != "" {
		dba = dba.Where("permit_key = ?", arg.PermitKey)
	}
	if arg.UrlPath != "" {
		dba = dba.Where("url_path LIKE  ?", fmt.Sprintf("%%%s%%", arg.UrlPath))
	}
	return
}

func (r *DaoPermitImpl) GetImportListCount(db *gorm.DB, arg *wrappers.ArgImportList) (totalCount int64, dba *gorm.DB, err error) {
	dba = r.fetchDb(db, arg)
	if err = dba.Count(&totalCount).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "daoPermitImplGetImportListCount")
		return
	}
	return
}
func (r *DaoPermitImpl) GetImportListData(db *gorm.DB, arg *wrappers.ArgImportList, pager *response.Pager) (res []models.AdminImport, err error) {

	res = []models.AdminImport{}
	if err = r.fetchDb(db, arg).
		Offset(arg.PageQuery.GetOffset()).
		Order(arg.Order).
		Limit(pager.PageSize).
		Find(&res).
		Error; err == nil {
		return
	}

	r.Context.Error(map[string]interface{}{
		"arg": arg,
		"err": err.Error(),
	}, "daoPermitImplGetImportListData")
	err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	return
}

func (r *DaoPermitImpl) getPermitImportByModuleWithSuperAdmin(arg *wrappers.ArgPermitMenu) (res []wrappers.Op, err error) {
	var (
		adminImport     models.AdminImport
		adminMenu       models.AdminMenu
		adminMenuImport models.AdminMenuImport
	)

	if err = r.Context.Db.Table(adminImport.TableName()).
		Select("b.permit_key,a.permit_key AS menu_permit_key").
		Joins(fmt.Sprintf("AS b LEFT JOIN  %s as c ON c.import_id=b.id LEFT JOIN %s AS a  ON a.id = c.menu_id",
			adminMenuImport.TableName(),
			adminMenu.TableName())).
		Where("b.deleted_at IS NULL AND a.deleted_at IS NULL AND c.deleted_at IS NULL").Unscoped().
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "DaoPermitImplGetPermitImportByModuleWithSuperAdmin")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}
func (r *DaoPermitImpl) GetPermitImportByModule(arg *wrappers.ArgPermitMenu) (res []wrappers.Op, err error) {
	var (
		adminUserGroupPermit models.AdminUserGroupImport
		adminImport          models.AdminImport
		adminMenu            models.AdminMenu
		adminMenuImport      models.AdminMenuImport
	)

	if arg.IsSuperAdmin { // 如果是超级管理员
		res, err = r.getPermitImportByModuleWithSuperAdmin(arg)
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
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "DaoPermitImplGetPermitImportByModule")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}

	return
}
func (r *DaoPermitImpl) GetImportMenuId(menuId ...int64) (list []models.AdminImport, err error) {
	list = []models.AdminImport{}
	var m models.AdminImport
	if len(menuId) == 0 {
		return
	}
	if err = r.Context.Db.Table(m.TableName()).
		Scopes(base.ScopesDeletedAt()).
		Where("menu_id IN(?)", menuId).
		Find(&list).Error; err == nil {
		return
	}

	r.Context.Error(map[string]interface{}{
		"menuId": menuId,
		"err":    err,
	}, "daoPermitGetImportMenuId")
	err = base.NewErrorRuntime(err, base.ErrorSqlCode)

	return
}

func (r *DaoPermitImpl) CreateImport(data *models.AdminImport) (res bool, err error) {
	var m models.AdminImport
	if err = r.Context.Db.Table(m.TableName()).Create(data).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"data": data,
			"err":  err,
		}, "daoPermitCreateImport")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}

	return
}
func (r *DaoPermitImpl) UpdateAdminImport(condition, data map[string]interface{}) (res bool, err error) {
	var m models.AdminImport
	if len(condition) == 0 {
		return
	}
	if err = r.Context.Db.Table(m.TableName()).
		Scopes(base.ScopesDeletedAt()).
		Where(condition).
		Updates(data).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"data":      data,
			"err":       err,
		}, "daoPermitUpdateAdminImport")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
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
	var m models.AdminImport
	err = r.Context.Db.Table(m.TableName()).
		Scopes(base.ScopesDeletedAt()).
		Where("id IN(?)", id).
		Find(&res).Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"id":  id,
			"err": err,
		}, "daoPermitGetAdminImportByIdError")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
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
	var m models.AdminImport
	err1 := r.Context.Db.Table(m.TableName()).Where("id=?", id).Find(&res).Error
	if err1 != nil && err1 != gorm.ErrRecordNotFound {
		err = err1
		r.Context.Error(map[string]interface{}{
			"id":  id,
			"err": err,
		}, "daoPermitGetImportId")
		return
	}
	return
}
func (r *DaoPermitImpl) GetImportCount(arg *wrappers.ArgGetImport, count *int64) (db *gorm.DB, err error) {
	if arg.MenuId == 0 {
		return
	}
	db = r.getImportListDb(nil, arg)
	err = db.Count(count).Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err,
		}, "daoPermitGetImportCount error")
	}
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

// func (r *DaoPermitImpl) AdminUserGroupAdd(data []map[string]interface{}) (err error) {
// 	field := make([]string, 0, 10)
// 	dataMsg := make([]string, 0)
// 	dataTmp := make([]interface{}, 0)
// 	duplicate := make([]string, 0)
// 	dataMsgA := make([]string, 0)
// 	if len(data) == 0 {
// 		err = fmt.Errorf("您没有为模板配置参数,请至少配置一个参数")
// 		return
// 	}
//
// 	for key, value := range data {
// 		if key == 0 {
// 			for fieldString := range value {
// 				field = append(field, fieldString)
// 			}
// 		}
// 		for _, v := range field {
// 			if key == 0 {
// 				if v != "created_at" {
// 					duplicate = append(duplicate, "`"+v+"`=VALUES(`"+v+"`)")
// 				}
// 				dataMsgA = append(dataMsgA, "?")
// 			}
// 			var v1 interface{}
// 			if _, ok := value[v]; ok {
// 				v1 = value[v]
// 			}
// 			dataTmp = append(dataTmp, v1)
// 		}
// 		dataMsg = append(dataMsg, "("+strings.Join(dataMsgA, ",")+")")
// 	}
// 	var m models.AdminUserGroup
// 	sql := "INSERT INTO `" + m.TableName() + "`(`" + strings.Join(field, "`,`") +
// 		"`) VALUES" + strings.Join(dataMsg, ",") + " ON DUPLICATE KEY UPDATE " + strings.Join(duplicate, ",")
// 	err = r.Context.Db.Exec(sql, dataTmp...).Error
// 	if err != nil {
// 		r.Context.Error(map[string]interface{}{
// 			"data": data,
// 			"err":  err,
// 		}, "daoPermitAdminUserGroupAdd")
// 	}
// 	return
// }
func (r *DaoPermitImpl) AdminUserGroupRelease(ids ...string) (err error) {
	if len(ids) == 0 {
		return
	}
	var m models.AdminUserGroup
	err = r.Context.Db.Table(m.TableName()).
		Where("id IN (?) ", ids).
		Delete(&models.AdminUserGroup{}).
		Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"ids": ids,
			"err": err,
		}, "daoPermitAdminUserGroupRelease")
	}
	return
}

func (r *DaoPermitImpl) FetchByName(name string) (res models.AdminGroup, err error) {
	var m models.AdminGroup
	res = models.AdminGroup{}
	err1 := r.Context.Db.Table(m.TableName()).
		Where("name = ?", name).
		Find(&res).
		Error
	if err1 == gorm.ErrRecordNotFound {
		return
	}
	err = err1
	r.Context.Error(map[string]interface{}{
		"name": name,
		"err":  err,
	}, "daoPermitFetchByName")
	return
}
func (r *DaoPermitImpl) InsertAdminGroup(group *models.AdminGroup) (err error) {
	var m models.AdminGroup
	err = r.Context.Db.Table(m.TableName()).
		Create(group).
		Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"group": group,
			"err":   err,
		}, "daoPermitInsertAdminGroup")
	}
	return
}
func (r *DaoPermitImpl) UpdateAdminGroup(data, where map[string]interface{}) (err error) {
	var group *models.AdminGroup

	if err = r.Context.Db.Table(group.TableName()).
		Where(where).
		Scopes(base.ScopesDeletedAt()).
		Updates(data).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"data":  data,
			"where": where,
			"err":   err,
		}, "daoPermitUpdateAdminGroup")
	}
	return
}
func (r *DaoPermitImpl) GetAdminGroupByIds(gIds ...int64) (res []models.AdminGroup, err error) {
	if len(gIds) == 0 {
		return
	}
	var m models.AdminGroup
	err = r.Context.Db.Table(m.TableName()).
		Scopes(base.ScopesDeletedAt()).
		Where("id IN (?)", gIds).
		Find(&res).
		Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"gIds": gIds,
			"err":  err,
		}, "daoPermitGetAdminGroupByIds")
	}
	return
}

func (r *DaoPermitImpl) GetUserGroupByUIds(uIds ...string) (res []models.AdminUserGroup, err error) {
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
		actErrorHandlerResult = &base.ActErrorHandlerResult{
			Db:        r.Context.Db,
			DbName:    r.Context.DbName,
			TableName: m.TableName(),
			Model:     &m,
		}
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
func (r *DaoPermitImpl) UpdateMenuByCondition(condition interface{}, data map[string]interface{}) (err error) {
	var m models.AdminMenu
	if err = r.Context.Db.
		Table(m.TableName()).
		Where(condition).
		Updates(data).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"data":      data,
			"err":       err,
		}, "daoPermitUpdateMenuByCondition")
		return
	}
	return
}
func (r *DaoPermitImpl) Save(id int64, data map[string]interface{}) (err error) {
	if id == 0 {
		return
	}
	if err = r.Context.Db.
		Model(&models.AdminMenu{}).
		Where("id = ?", id).
		Updates(data).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"id":   id,
			"data": data,
			"err":  err,
		}, "daoPermitSave")
	}
	return
}
func (r *DaoPermitImpl) DeleteMenuByIds(ids ...string) (err error) {
	if len(ids) == 0 {
		return
	}
	var m models.AdminMenu
	if err = r.Context.Db.
		Table(m.TableName()).Unscoped().
		Where("id in (?)", ids).
		Delete(&models.AdminMenu{}).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"ids": ids,
			"err": err,
		}, "daoPermitDeleteByIds")
	}
	return
}

func (r *DaoPermitImpl) GetByCondition(condition map[string]interface{}, orderBy []wrappers.DaoOrderBy, limit int) (res []models.AdminMenu, err error) {
	if len(condition) == 0 {
		return
	}
	var m models.AdminMenu
	db := r.Context.Db.
		Table(m.TableName()).
		Where(condition)
	var orderString string
	for _, item := range orderBy {
		orderString += fmt.Sprintf("%s %s", item.Column, item.SortFormat)
	}
	if orderString != "" {
		db = db.Order(orderString)
	}
	if limit > 0 {
		db = db.Limit(limit)
	}
	if err = db.
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"err":       err,
		}, "daoPermitGetByCondition")
		return
	}
	return
}
func (r *DaoPermitImpl) Add(data *models.AdminMenu) (err error) {

	if err = r.Context.Db.Create(data).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"data": data,
			"err":  err,
		}, "daoPermitAdd")
	}
	return
}

func (r *DaoPermitImpl) GetMenuByPermitKey(module string, permitKey ...string) (res []models.AdminMenu, err error) {

	var m models.AdminMenu

	db := r.Context.Db.Table(m.TableName()).Scopes(base.ScopesDeletedAt())
	lPk := len(permitKey)
	var f bool
	if lPk != 0 {
		f = true
		db = db.Where("permit_key IN(?)", permitKey)
	}
	if module != "" {
		f = true
		db = db.Where("module = ?", module)
	}
	if !f {
		return
	}

	if err = db.Limit(len(permitKey)).
		Order("sort_value desc").
		Find(&res).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"permitKey": permitKey,
			"err":       err.Error(),
		}, "daoPermitGetMenuByPermitKey")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}
func (r *DaoPermitImpl) GetAdminMenuList(arg *wrappers.ArgAdminMenu) (res []models.AdminMenu, err error) {
	res = []models.AdminMenu{}
	var m models.AdminMenu
	dba := r.Context.Db.Table(m.TableName())
	if arg.SystemId > 0 {
		if err = r.Context.Db.Table(m.TableName()).
			Scopes(base.ScopesDeletedAt()).
			Order("sort_value desc").
			Where("id = ?", arg.SystemId).
			First(&m).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return
			}
			err = fmt.Errorf("你查看权限系统不存在或已删除")
			r.Context.Error(map[string]interface{}{
				"arg": arg,
				"err": err,
			}, "daoPermitGetAdminMenuList0")
			return
		}
		arg.Module = m.PermitKey
	}
	if arg.Label = strings.TrimSpace(arg.Label); arg.Label != "" {
		dba = dba.Where("label LIKE ?", "%"+arg.Label+"%")
	}
	if arg.ParentId > 0 {
		dba = dba.Where("parent_id = ?", arg.ParentId)
	}
	if arg.Module != "" {
		dba = dba.Where("module = ? OR parent_id=?", arg.Module, wrappers.DefaultPermitParentId)
	}
	if arg.AppName != "" {
		dba = dba.Where("app_name = ?", arg.AppName)
	}
	if arg.IsMenuShow != -1 {
		dba = dba.Where("hide_in_menu = ?", arg.IsMenuShow)
	}
	if arg.IsDel == -1 || arg.IsDel == 0 {
		dba = dba.Where("deleted_at IS NULL")
	} else {
		dba = dba.Where("deleted_at NOT NULL")
	}
	if arg.Id != 0 {
		dba = dba.Where("id = ?", arg.Id)
	}
	if err = dba.Order("sort_value desc").Find(&res).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "daoPermitGetAdminMenuList2")
		return
	}
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
func (r *DaoPermitImpl) GetGroupByUserId(userId string) (res []wrappers.AdminGroupUserStruct, err error) {
	if userId == "" {
		res = []wrappers.AdminGroupUserStruct{}
		return
	}
	var a models.AdminUserGroup
	var b models.AdminGroup
	if err = r.Context.Db.
		Select("a.group_id,a.user_hid,a.is_super_admin as super_admin,a.is_admin_group as admin_group,b.*").
		Unscoped().
		Table(a.TableName()).
		Joins(fmt.Sprintf("as a LEFT JOIN %s as b  ON  a.group_id=b.id ",
			b.TableName())).
		Where(fmt.Sprintf("a.user_hid = ? AND a.deleted_at  IS NULL AND  b.deleted_at IS NULL"),
			userId).
		Find(&res).
		Error; err == nil {
		return
	}
	r.Context.Error(map[string]interface{}{
		"userId": userId,
		"err":    err,
	}, "daoPermitGetGroupByUserId")
	return
}
func (r *DaoPermitImpl) GetPermitMenuByIds(module []string, menuIds ...int64) (res []models.AdminMenu, err error) {
	res = []models.AdminMenu{}
	var m models.AdminMenu
	db := r.Context.Db.
		Table(m.TableName())
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
	if err = db.Order("sort_value desc").
		Find(&res).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"menuIds": menuIds,
			"err":     err,
		}, "daoPermitGetPermitMenuByIds")
		return
	}
	return
}

// func (r *DaoPermitImpl) GetMenuIdsByPermitByGroupIds(module string, pathType []string, groupIds ...int64) (res []models.AdminUserGroupPermit, err error) {
// 	if len(groupIds) == 0 {
// 		return
// 	}
// 	var m models.AdminUserGroupPermit
// 	if err = r.Context.Db.
// 		Table(m.TableName()).
// 		Select("distinct `menu_id`,`group_id`,`id`").
// 		Where("module = ? AND path_type IN(?)  AND `group_id` in(?) ", module, pathType, groupIds).
// 		Find(&res).
// 		Error; err != nil {
// 		r.Context.Error(map[string]interface{}{
// 			"groupIds": groupIds,
// 			"pathType": pathType,
// 			"err":      err,
// 		}, "daoPermitGetPermitMenuByIds")
// 		return
// 	}
// 	return
// }
