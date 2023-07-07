// Package dao_impl
/**
* @Author:ChangJiang
* @Description:
* @File:permit_import
* @Version: 1.0.0
* @Date 2021/6/5 11:35 下午
 */
package dao_impl

import (
	"context"
	"fmt"
	"github.com/juetun/base-wrapper/lib/base/cache_act"
	"github.com/juetun/base-wrapper/lib/common/redis_pkg"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/pkg/parameters"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl/cache_act_local"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"strings"
	"time"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_intranet"
)

type DaoPermitImportImpl struct {
	base.ServiceDao
	ctx context.Context
}

func (r *DaoPermitImportImpl) GetMenuImportByMenuIdAndImportIds(menuId int64, importIds ...int64) (res []*models.AdminMenuImport, err error) {
	if menuId == 0 || len(importIds) == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"menuId":    menuId,
			"importIds": importIds,
			"err":       err.Error(),
		}, "DaoPermitImportImplGetMenuImportByMenuIdAndImportIds")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actErrorHandlerResult *base.ActErrorHandlerResult) {
		var m models.AdminMenuImport
		actErrorHandlerResult = r.GetDefaultActErrorHandlerResult(&m)
		actErrorHandlerResult.Err = actErrorHandlerResult.Db.
			Scopes(base.ScopesDeletedAt()).
			Table(actErrorHandlerResult.TableName).
			Where("menu_id = ? AND import_id IN(?)", menuId, importIds).
			Find(&res).
			Error
		return
	})

	return
}

func (r *DaoPermitImportImpl) GetImportFromDbByIds(ids ...int64) (res map[int64]*models.AdminImport, err error) {

	var (
		l          = len(ids)
		m          models.AdminImport
		data       []*models.AdminImport
		logContent = map[string]interface{}{"ids": ids,}
	)

	res = make(map[int64]*models.AdminImport, l)
	if l == 0 {
		return
	}

	defer func() {
		if err == nil {
			return
		}
		logContent["err"] = err.Error()
		r.Context.Error(logContent, "DaoPermitImportImplGetImportFromDbByIds")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()

	if err = r.ActErrorHandler(func() (actErrorHandlerResult *base.ActErrorHandlerResult) {
		actErrorHandlerResult = r.GetDefaultActErrorHandlerResult(&m)
		actErrorHandlerResult.Err = actErrorHandlerResult.Db.
			Table(actErrorHandlerResult.TableName).
			Where("id IN (?)", ids).
			Scopes(base.ScopesDeletedAt()).
			Find(&data).
			Error
		return
	}); err != nil {
		return
	}
	for _, item := range data {
		res[item.Id] = item
	}
	return
}

func (r *DaoPermitImportImpl) GetImportForGateway(arg *wrapper_intranet.ArgGetImportPermit) (list []models.AdminImport, err error) {
	condition := map[string]interface{}{}
	if arg.AppName != "" {
		condition["app_name"] = arg.AppName
	}
	logContent := map[string]interface{}{"arg": arg}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(logContent, "DaoPermitImportGetImportForGateway")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m models.AdminImport
		actRes = r.GetDefaultActErrorHandlerResult(&m)
		actRes.Err = actRes.Db.Table(actRes.TableName).Where(condition).
			Scopes(base.ScopesDeletedAt()).
			Where("need_login = ? OR need_sign = ?",
				models.NeedLoginNot,
				models.NeedSignNot).
			Find(&list).
			Error
		return
	})
	return
}

func (r *DaoPermitImportImpl) AddData(adminImport *models.AdminImport) (err error) {
	logContent := map[string]interface{}{"adminImport": adminImport,}
	defer func() {
		if err == nil {
			return
		}
		logContent["err"] = err.Error()
		r.Context.Error(logContent, "DaoPermitImportAddData")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	if err = adminImport.InitRegexpString(); err != nil {
		logContent["desc"] = "InitRegexpString"
		return
	}
	if err = r.Context.Db.Create(adminImport).Error; err != nil {
		logContent["desc"] = "Create"
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}

func (r *DaoPermitImportImpl) getColumnName(s string) (res string) {
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
	if res == "primary_key" {
		res = "id"
	}
	return
}

func (r *DaoPermitImportImpl) BatchMenuImport(list []*models.AdminMenuImport) (err error) {
	if len(list) == 0 {
		return
	}
	var dtList []base.ModelBase
	for _, menuImport := range list {
		dtList = append(dtList, menuImport)
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"data": list,
			"err":  err.Error(),
		}, "DaoPermitImportBatchMenuImport")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		actRes = r.GetDefaultActErrorHandlerResult(list[0])
		actRes.Err = r.BatchAdd(actRes.ParseBatchAddDataParameter(base.BatchAddDataParameterData(dtList)))
		return
	})
	return
}

func (r *DaoPermitImportImpl) DeleteByCondition(condition interface{}) (res bool, err error) {

	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"err":       err.Error(),
		}, "DaoPermitImportDeleteByCondition")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	if condition == nil {
		err = fmt.Errorf("您没有选择要删除的数据")
		return
	}
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m models.AdminImport
		actRes = r.GetDefaultActErrorHandlerResult(&m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where(condition).
			Delete(&models.AdminImport{}).
			Error
		return
	})

	return
}

func (r *DaoPermitImportImpl) UpdateByCondition(condition interface{}, data map[string]interface{}) (res bool, err error) {

	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"data":      data,
			"err":       err.Error(),
		}, "DaoPermitImportImplUpdateByCondition")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	if condition == nil || len(data) == 0 {
		err = fmt.Errorf("您没有选择要更新的数据")
		return
	}

	if err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m models.AdminImport
		actRes = r.GetDefaultActErrorHandlerResult(&m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Where(condition).
			Scopes(base.ScopesDeletedAt()).
			Updates(data).
			Error
		return
	}); err != nil {
		return
	}
	res = true
	return
}

func (r *DaoPermitImportImpl) getChildImportByMenuIdFromDb(menuIds ...int64) (resData map[int64]models.AdminMenuImportCache, err error) {
	var l = len(menuIds)
	resData = make(map[int64]models.AdminMenuImportCache, l)
	if l == 0 {
		return
	}
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminMenuImport
		var list []*models.AdminMenuImport
		actRes = r.GetDefaultActErrorHandlerResult(m)
		if actRes.Err = actRes.Db.Table(actRes.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where("menu_id IN (?)", menuIds).
			Find(&list).
			Error; actRes.Err != nil {
			return
		}

		for _, item := range list {
			if _, ok := resData[item.MenuId]; !ok {
				resData[item.MenuId] = make([]*models.AdminMenuImport, 0, l)
			}
			resData[item.MenuId] = append(resData[item.MenuId], item)
		}
		return
	})

	return
}

func (r *DaoPermitImportImpl) GetChildImportByMenuId(arg *base.ArgGetByNumberIds) (res map[int64]models.AdminMenuImportCache, err error) {
	res = map[int64]models.AdminMenuImportCache{}
	if len(arg.Ids) == 0 {
		return
	}
	res, err = cache_act_local.NewMenuImportAction(
		cache_act_local.CacheMenuImportActionArg(arg),
		cache_act_local.CacheMenuImportActionGetByIdsFromDb(r.getChildImportByMenuIdFromDb),
	).LoadBaseOption(
		cache_act.CacheActionBaseGetCacheKey(r.GetMenuImportCacheKey),
		cache_act.CacheActionBaseContext(r.Context),
		cache_act.CacheActionBaseCtx(r.ctx), ).
		Action()

	return
}

func (r *DaoPermitImportImpl) GetMenuImportCacheKey(productId interface{}, expireTimeRands ...bool) (res string, expire time.Duration, err error) {
	var CacheMenuImportWithAppKey *redis_pkg.CacheProperty
	if CacheMenuImportWithAppKey, err = parameters.GetCacheParamConfig("CacheMenuImportWithAppKey"); err != nil {
		return
	}
	res = fmt.Sprintf(CacheMenuImportWithAppKey.Key, productId)
	expire = CacheMenuImportWithAppKey.Expire
	return
}

func (r *DaoPermitImportImpl) getImportByImportIdFromDb(importIds ...int64) (resData map[int64]*models.AdminMenuImport, err error) {
	var l = len(importIds)
	resData = make(map[int64]*models.AdminMenuImport, l)
	if l == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"importIds": importIds,
			"err":       err.Error(),
		}, "DaoPermitImportGetImportByImportIdFromDb")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m models.AdminMenuImport
		var list []*models.AdminMenuImport
		actRes = r.GetDefaultActErrorHandlerResult(&m)
		if actRes.Err = actRes.Db.Table(actRes.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where("import_id IN (?)", importIds).
			Find(&list).
			Error; actRes.Err != nil {
			return
		}
		for _, item := range list {
			resData[item.ImportId] = item
		}
		return
	})
	return
}

func (r *DaoPermitImportImpl) GetImportMenuByImportIds(arg *base.ArgGetByNumberIds) (list map[int64]*models.AdminMenuImport, err error) {
	list = map[int64]*models.AdminMenuImport{}
	if list, err = cache_act_local.NewMenuImportByImportIdAction(
		cache_act_local.CacheMenuImportByImportIdActionArg(arg),
		cache_act_local.CacheMenuImportByImportIdActionGetByIdsFromDb(r.getImportByImportIdFromDb),
	).LoadBaseOption(
		cache_act.CacheActionBaseGetCacheKey(r.GetMenuImportCacheKey),
		cache_act.CacheActionBaseContext(r.Context),
		cache_act.CacheActionBaseCtx(r.ctx), ).
		Action(); err != nil {
		return
	}

	return
}

func (r *DaoPermitImportImpl) GetImportMenuByImportIdsCacheKey(productId interface{}, expireTimeRands ...bool) (res string, expire time.Duration, err error) {
	var CacheMenuImportWithAppKey *redis_pkg.CacheProperty
	if CacheMenuImportWithAppKey, err = parameters.GetCacheParamConfig("CacheMenuImportWithAppKey"); err != nil {
		return
	}
	res = fmt.Sprintf(CacheMenuImportWithAppKey.Key, productId)
	expire = CacheMenuImportWithAppKey.Expire
	return
}

func (r *DaoPermitImportImpl) UpdateMenuImport(condition string, data map[string]interface{}) (err error) {
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"data":      data,
			"err":       err.Error(),
		}, "DaoPermitImportUpdateMenuImport")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m models.AdminMenuImport
		actRes = r.GetDefaultActErrorHandlerResult(&m)
		actRes.Err = actRes.Db.
			Table(actRes.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where(condition).
			Updates(data).
			Error
		return
	})
	return
}

func (r *DaoPermitImportImpl) BatchAddData(list []models.AdminImport) (err error) {

	if len(list) == 0 {
		return
	}
	var (
		dataList []base.ModelBase
	)

	for _, adminImport := range list {
		dataList = append(dataList, &adminImport)
	}
	logContent := map[string]interface{}{"list": list,}
	defer func() {
		if err == nil {
			return
		}
		logContent["err"] = err.Error()
		r.Context.Error(logContent, "DaoPermitImportBatchAddData")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		actRes = r.GetDefaultActErrorHandlerResult(dataList[0])
		actRes.Err = r.BatchAdd(actRes.ParseBatchAddDataParameter(base.BatchAddDataParameterData(dataList)))
		return
	})
	return
}

func (r *DaoPermitImportImpl) DeleteImportByIds(id ...int64) (err error) {
	if len(id) == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"id":  id,
			"err": err.Error(),
		}, "daoPermitDeleteImportByIds")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m models.AdminImport
		actRes = r.GetDefaultActErrorHandlerResult(&m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Scopes(base.ScopesDeletedAt()).
			Unscoped().
			Where("id IN(?)", id).
			Updates(map[string]interface{}{
				"deleted_at": time.Now().Format(utils.DateTimeGeneral),
			}).
			//Delete(&models.AdminImport{}).
			Error
		return
	})
	if err != nil {
		return
	}
	if err = NewDaoPermitGroupImport(r.Context).
		DeleteGroupImportWithMenuId(id...); err != nil {
		return
	}
	return
}

func (r *DaoPermitImportImpl) GetImportByCondition(condition map[string]interface{}) (list []models.AdminImport, err error) {
	list = []models.AdminImport{}
	if len(condition) == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"err":       err.Error(),
		}, "DaoPermitImportImplGetImportByCondition")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m models.AdminImport
		actRes = r.GetDefaultActErrorHandlerResult(&m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Where(condition).
			Scopes(base.ScopesDeletedAt()).
			Find(&list).
			Limit(1000).Error
		return
	})
	return
}

func (r *DaoPermitImportImpl) GetDefaultOpenImportByMenuIds(menuId ...int64) (res []wrappers.AdminImportWithMenu, err error) {
	res = []wrappers.AdminImportWithMenu{}
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
		}, "GetDefaultOpenImportByMenuIds")
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m models.AdminImport
		var ami models.AdminMenuImport
		actRes = r.GetDefaultActErrorHandlerResult(&m)
		err = actRes.Db.Table(ami.TableName()).
			Unscoped().
			Select(`b.*,a.menu_id`).
			Joins(fmt.Sprintf("AS a LEFT JOIN %s AS b  ON  b.`id` = a.import_id", m.TableName())).
			Where("a.menu_id IN(?) AND  b.default_open = ? AND `a`.`deleted_at` IS NULL ",
				menuId,
				models.DefaultOpen).
			Find(&res).
			Error
		return
	})
	return
}
func NewDaoPermitImport(c ...*base.Context) (res daos.DaoPermitImport) {
	p := &DaoPermitImportImpl{}
	p.SetContext(c...)
	p.ctx = context.TODO()
	res = p
	return
}
