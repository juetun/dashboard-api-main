// Package dao_impl /**
package dao_impl

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/pkg/parameters"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
)

type DaoPermitGroupImpl struct {
	base.ServiceBase
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
	err = r.DeleteUserGroupPermitByGroupId(fmt.Sprintf("%d", groupId))
	return
}

func (r *DaoPermitGroupImpl) GetPermitGroupByUid(userHid string, refreshCache ...bool) (res []models.AdminUserGroup, err error) {
	var freshCache bool
	if len(refreshCache) > 0 {
		freshCache = refreshCache[0]
	}
	if !freshCache { // 如果不是需要重新刷新缓存
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

func (r *DaoPermitGroupImpl) DeleteUserGroupPermitByGroupId(ids ...string) (err error) {
	if len(ids) == 0 {
		return
	}
	var m models.AdminUserGroupPermit
	if err = r.Context.Db.
		Table(m.TableName()).
		Scopes(base.ScopesDeletedAt()).
		Where("group_id IN (?) ", ids).
		Delete(&models.AdminGroup{}).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"ids": ids,
			"err": err,
		}, "DaoPermitGroupImplDeleteUserGroupPermitByGroupId")
		err = base.NewErrorRuntime(err, base.ErrorRedisCode)
		return
	}
	go func() {
		_ = r.DeleteUserGroupCacheByGroupIds(ids...)
	}()

	return
}

// DeleteUserGroupCacheByGroupIds 根据用户组ID删除用户对应的权限组关系操作
func (r *DaoPermitGroupImpl) DeleteUserGroupCacheByGroupIds(groupIds ...string) (err error) {

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

func (r *DaoPermitGroupImpl) GetGroupUserByGroupIds(pager *response.Pager, groupIds ...string) (list []models.AdminUserGroup, err error) {
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

func (r *DaoPermitGroupImpl) DeleteUserGroupPermit(pathType string, menuId ...int) (err error) {
	if len(menuId) == 0 || pathType == "" {
		return
	}
	var m models.AdminUserGroupPermit
	if err = r.Context.Db.Table(m.TableName()).
		Scopes(base.ScopesDeletedAt()).
		Where("menu_id IN (?) AND path_type = ?", menuId, pathType).
		Delete(&models.AdminUserGroupPermit{}).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"menu_id":  menuId,
			"pathType": pathType,
			"err":      err,
		}, "daoPermitDeleteUserGroupPermit")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)

		return
	}
	return
}

func NewDaoPermitGroupImpl(ctx ...*base.Context) (res daos.DaoPermitGroup) {
	p := &DaoPermitGroupImpl{}
	p.SetContext(ctx...)
	return p
}
