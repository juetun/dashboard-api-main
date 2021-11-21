// Package dao_impl /**
package dao_impl

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/pkg/parameters"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_intranet"
)

type DaoPermitGroupImpl struct {
	base.ServiceDao
}

func (r *DaoPermitGroupImpl) GetGroupAppPermitImport(groupId int64, appName string, refreshCache ...bool) (res []wrapper_intranet.AdminUserGroupPermit, err error) {
	// 如果不需要强制刷新缓存
	if !r.RefreshCache(refreshCache...) {
		var dataNil bool
		if dataNil, err = r.getGroupAppPermitImportFromCache(groupId, appName, &res); err != nil {
			return
		}
		if !dataNil { // 如果缓存中有数据，则使用缓存中的数据
			return
		}
	}
	if res, err = r.getGroupAppPermitImportFromDb(groupId, appName); err != nil {
		return
	}

	if err = r.setGroupAppPermitImportToCache(groupId, appName, res); err != nil {
		return
	}
	return
}

func (r *DaoPermitGroupImpl) getGroupAppPermitImportFromDb(groupId int64, appName string) (res []wrapper_intranet.AdminUserGroupPermit, err error) {
	var m models.AdminUserGroupImport

	defer func() {
		if err != nil {
			r.Context.Error(map[string]interface{}{
				"groupId": groupId,
				"appName": appName,
			}, "DaoPermitGroupImplGetGroupAppPermitImportFromDb")
			err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		}
	}()

	db := r.Context.Db.Table(m.TableName())
	var mi models.AdminImport
	err = db.Select("a.group_id,a.app_name,a.path_type,b.*").Joins(fmt.Sprintf(" AS a LEFT JOIN  %s AS b ON a.import_id= b.id ",
		mi.TableName())).
		Where("a.group_id = ? AND a.app_name = ?",
			groupId,
			appName).
		Scopes(base.ScopesDeletedAt("a")).
		Scopes(base.ScopesDeletedAt("b")).
		Find(&res).
		Error

	return
}
func (r *DaoPermitGroupImpl) getGroupAppPermitMenuFromDb(groupId int64, appName string) (res []wrapper_intranet.AdminUserGroupPermit, err error) {
	var m models.AdminUserGroupMenu

	defer func() {
		if err != nil {
			r.Context.Error(map[string]interface{}{
				"groupId": groupId,
				"appName": appName,
			}, "DaoPermitGroupImplGetGroupAppPermitMenuFromDb")
			err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		}
	}()

	db := r.Context.Db.Table(m.TableName())

	var mm models.AdminImport
	db = db.Joins(fmt.Sprintf(" AS a LEFT JOIN  %s AS b ON a.menu_id= b.id",
		mm.TableName()))
	if appName != "" {
		db = db.Where("a.module = ?", appName)
	}
	err = db.
		Where("a.group_id = ?", groupId).
		Scopes(base.ScopesDeletedAt("a")).
		Scopes(base.ScopesDeletedAt("b")).
		Find(&res).
		Error

	return
}

func (r *DaoPermitGroupImpl) setGroupAppPermitImportToCache(groupId int64, appName string, res []wrapper_intranet.AdminUserGroupPermit) (err error) {
	key, duration := r.getGroupAppPermitImportCacheKey(groupId, appName)
	defer func() {
		if err != nil {
			r.Context.Error(map[string]interface{}{
				"groupId": groupId,
				"appName": appName,
				"err":     err.Error(),
			}, "DaoPermitGroupImplSetGroupAppPermitImportToCache")
			err = base.NewErrorRuntime(err, base.ErrorRedisCode)
			return
		}
	}()

	op := r.Context.CacheClient.Set(context.TODO(), key, &res, duration)
	if err = op.Err(); err != nil {
		return
	}
	return
}

func (r *DaoPermitGroupImpl) getGroupAppPermitImportFromCache(groupId int64, appName string, data interface{}) (dataNil bool, err error) {
	key, _ := r.getGroupAppPermitImportCacheKey(groupId, appName)

	defer func() {
		if err != nil {
			r.Context.Error(map[string]interface{}{
				"groupId": groupId,
				"appName": appName,
				"err":     err.Error(),
			}, "DaoPermitGroupImplGetGroupAppPermitImport")
			err = base.NewErrorRuntime(err, base.ErrorRedisCode)
			return
		}
	}()

	ctx := context.TODO()
	op := r.Context.CacheClient.Get(ctx, key)

	var e error
	if e = op.Err(); e != nil {
		if e == redis.Nil {
			dataNil = true
			return
		}
		err = e
		return
	}

	if errString := op.Scan(data).Error(); errString != "" {
		err = fmt.Errorf(errString)
		return
	}
	return
}

// 用户组每个接口权限缓存
func (r *DaoPermitGroupImpl) getGroupAppPermitImportCacheKey(groupId int64, appName string) (res string, duration time.Duration) {
	res = fmt.Sprintf(parameters.CacheKeyUserGroupAppImportWithAppKey, groupId, appName)

	// 生成一个基础时间加上随机时间的时间值避免缓存数据在同一个时间失效
	duration = parameters.CacheKeyUserGroupAppImportTime + time.Duration(rand.Int63n(100))*time.Minute
	return
}

func (r *DaoPermitGroupImpl) UpdateDaoPermitUserGroupByGroupId(groupId int64, data map[string]interface{}) (err error) {
	var m models.AdminUserGroup
	if err = r.Context.Db.Table(m.TableName()).
		Scopes(base.ScopesDeletedAt()).
		Where("group_id = ?", groupId).
		Updates(data).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"groupId": groupId,
			"data":    data,
			"err":     err.Error(),
		}, "DaoPermitGroupImplUpdateDaoPermitUserGroupByGroupId")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	err = r.DeleteUserGroupPermitByGroupId(groupId)
	return
}

func (r *DaoPermitGroupImpl) GetPermitGroupByUid(userHid string, refreshCache ...bool) (res []models.AdminUserGroup, err error) {

	if !r.RefreshCache(refreshCache...) { // 如果不是需要重新刷新缓存
		dataNil, _ := r.getCacheUserPermitGroup(userHid, &res)
		if !dataNil {
			return
		}
	}
	var m models.AdminUserGroup
	if err = r.Context.Db.Table(m.TableName()).
		Where("user_hid = ?", userHid).
		Scopes(base.ScopesDeletedAt()).
		Find(&res).Error; err == nil {
		go func() { // 重新设置缓存
			_ = r.setCacheUserPermitGroupCache(userHid, res)
		}()
		return
	}
	r.Context.Error(map[string]interface{}{
		"userHid": userHid,
		"err":     err.Error(),
	}, "DaoPermitGroupImplGetPermitGroupByUid")
	return
}

func (r *DaoPermitGroupImpl) getCacheUserPermitGroup(userHid string, data interface{}) (dataNil bool, err error) {
	key, _ := r.getUserPermitGroupCacheKey(userHid)
	var e error
	if e = r.Context.CacheClient.Get(context.TODO(), key).Scan(data); e != nil {
		if e == redis.Nil {
			dataNil = true
			return
		}
		r.Context.Error(map[string]interface{}{
			"userHid": userHid,
			"err":     err.Error(),
		}, "DaoPermitGroupImplGetCacheUserPermitGroup")
		err = base.NewErrorRuntime(err, base.ErrorRedisCode)
	}
	return
}

// 用户权限组缓存
func (r *DaoPermitGroupImpl) getUserPermitGroupCacheKey(userHid string) (res string, duration time.Duration) {
	res = fmt.Sprintf(parameters.CacheKeyUserGroupWithAppKey, userHid)
	duration = parameters.CacheKeyUserGroupWithAppKeyTime
	return
}

func (r *DaoPermitGroupImpl) deleteCacheUserPermitGroupCache(userHid string) (err error) {
	key, _ := r.getUserPermitGroupCacheKey(userHid)
	if err = r.Context.CacheClient.Del(context.TODO(), key).Err(); err != nil {
		r.Context.Error(map[string]interface{}{
			"userHid": userHid,
			"err":     err.Error(),
		}, "DaoPermitGroupImplDeleteCacheUserPermitGroupCache")
		err = base.NewErrorRuntime(err, base.ErrorRedisCode)
	}
	return
}

func (r *DaoPermitGroupImpl) setCacheUserPermitGroupCache(userHid string, res []models.AdminUserGroup) (err error) {
	key, duration := r.getUserPermitGroupCacheKey(userHid)
	if err = r.Context.CacheClient.Set(context.TODO(), key, res, duration).Err(); err != nil {
		r.Context.Error(map[string]interface{}{
			"userHid": userHid,
			"res":     res,
			"err":     err.Error(),
		})
		err = base.NewErrorRuntime(err, base.ErrorRedisCode)
		return
	}
	return
}

func (r *DaoPermitGroupImpl) DeleteUserGroupPermitByGroupId(groupIds ...int64) (err error) {
	if len(groupIds) == 0 {
		return
	}
	daoPermitGroupMenu := NewDaoPermitGroupMenu(r.Context)

	if err = daoPermitGroupMenu.DeleteGroupMenuByGroupIds(groupIds...); err != nil {
		r.Context.Error(map[string]interface{}{
			"groupIds": groupIds,
			"err":      err,
		}, "DaoPermitGroupImplDeleteUserGroupPermitByGroupId")
		err = base.NewErrorRuntime(err, base.ErrorRedisCode)
		return
	}
	go func() {
		_ = r.DeleteUserGroupCacheByGroupIds(groupIds...)
	}()

	return
}

// DeleteUserGroupCacheByGroupIds 根据用户组ID删除用户对应的权限组关系操作
func (r *DaoPermitGroupImpl) DeleteUserGroupCacheByGroupIds(groupIds ...int64) (err error) {

	var list []models.AdminUserGroup

	pageQuery := response.PageQuery{}
	pageQuery.PageNo = 1
	pageQuery.PageSize = 1000
	pageQuery.PageType = response.DefaultPageTypeNext
	pageQuery.RequestId = "0"

	pager := response.NewPager(response.PagerBaseQuery(pageQuery))

	for {
		if list, err = r.GetGroupUserByGroupIds(pager, groupIds...); err != nil {
			break
		}

		if len(list) == 0 {
			break
		}

		if len(list) == pager.PageSize {
			pager.IsNext = true
		}
		mReqId := 0
		for k, item := range list {
			if k == 0 || item.Id > mReqId {
				mReqId = item.Id
			}
			_ = r.deleteCacheUserPermitGroupCache(item.UserHid)
		}
		pager.RequestId = fmt.Sprintf("%d", mReqId)

	}

	return
}

func (r *DaoPermitGroupImpl) GetGroupUserByGroupIds(pager *response.Pager, groupIds ...int64) (list []models.AdminUserGroup, err error) {
	var m models.AdminUserGroup
	db := r.Context.Db.Distinct("user_hid").
		Table(m.TableName()).
		Scopes(base.ScopesDeletedAt()).
		Where("group_id IN (?)", groupIds)
	if pager.RequestId != "" {
		db = db.Where("id > ?", pager.RequestId)
	}
	if err = db.Limit(pager.PageSize).
		Order("id ASC").
		Find(list).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"pager":    pager,
			"groupIds": groupIds,
			"err":      err.Error(),
		}, "DaoPermitGroupImplGetGroupUserByGroupIds")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}

func (r *DaoPermitGroupImpl) DeleteUserGroupByUserId(userId ...string) (err error) {
	if len(userId) == 0 {
		return
	}
	var m models.AdminUserGroup
	if err = r.Context.Db.Table(m.TableName()).
		Scopes(base.ScopesDeletedAt()).
		Where("user_hid IN (?) ", userId).
		Updates(map[string]interface{}{
			"deleted_at": time.Now().Format(utils.DateTimeGeneral),
		}).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"userId": userId,
			"err":    err.Error(),
		}, "daoPermitDeleteUserGroupByUserId")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}

	// 清除缓存
	go func() {
		for _, uid := range userId {
			_ = r.deleteCacheUserPermitGroupCache(uid)
		}
	}()
	return
}

func (r *DaoPermitGroupImpl) DeleteUserGroupPermit(menuId ...int64) (err error) {
	if len(menuId) == 0 {
		return
	}
	var m models.AdminUserGroupMenu
	if err = r.Context.Db.Table(m.TableName()).
		Scopes(base.ScopesDeletedAt()).
		Where("menu_id IN (?) ", menuId).
		Delete(&m).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"menu_id": menuId,
			"err":     err,
		}, "DaoPermitGroupImplDeleteUserGroupPermit")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)

		return
	}
	return
}
func (r *DaoPermitGroupImpl) DeleteAdminGroupByIds(ids ...int64) (err error) {
	if len(ids) == 0 {
		return
	}
	var m models.AdminGroup
	err = r.Context.Db.Table(m.TableName()).
		Where("id IN (?) ", ids).
		Delete(&models.AdminGroup{}).
		Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"ids": ids,
			"err": err,
		}, "DaoPermitGroupImplDeleteAdminGroupByIds")
	}
	return
}
func NewDaoPermitGroupImpl(ctx ...*base.Context) (res daos.DaoPermitGroup) {
	p := &DaoPermitGroupImpl{}
	p.SetContext(ctx...)
	return p
}
