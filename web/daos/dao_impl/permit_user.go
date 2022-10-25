// Package dao_impl /**
package dao_impl

import (
	"context"
	"fmt"
	"github.com/juetun/base-wrapper/lib/base/cache_act"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/pkg/parameters"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl/cache_act_local"
	"time"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"gorm.io/gorm"
)

type DaoPermitUserImpl struct {
	base.ServiceDao
	ctx context.Context
}

func (r *DaoPermitUserImpl) GetAdminUserByIds(arg *base.ArgGetByNumberIds) (res map[int64]*models.AdminUser, err error) {
	res = map[int64]*models.AdminUser{}
	res, err = cache_act_local.NewCacheAdminUserAction(
		cache_act_local.CacheAdminUserActionArg(arg),
		cache_act_local.CacheAdminUserActionGetByIdsFromDb(r.getAdminUserByIdsFromDb),
	).LoadBaseOption(
		cache_act.CacheActionBaseGetCacheKey(r.GetCacheKey),
		cache_act.CacheActionBaseContext(r.Context),
		cache_act.CacheActionBaseCtx(r.ctx), ).
		Action()
	return
}

func (r *DaoPermitUserImpl) GetCacheKey(id interface{}, expireTimeRands ...bool) (res string, timeExpire time.Duration) {

	res = fmt.Sprintf(parameters.CacheKeyAdminUserWithUserHId.Key, id)
	timeExpire = parameters.CacheKeyAdminUserWithUserHId.Expire
	var expireTimeRand bool
	if len(expireTimeRands) > 0 {
		expireTimeRand = expireTimeRands[0]
	}
	if expireTimeRand {
		randNumber, _ := r.RealRandomNumber(60) //打乱缓存时长，防止缓存同一时间过期导致数据库访问压力变大
		timeExpire = timeExpire + time.Duration(randNumber)*time.Second
	}
	return
}
func (r *DaoPermitUserImpl) getAdminUserByIdsFromDb(userHIds ...int64) (resData map[int64]*models.AdminUser, err error) {
	resData = make(map[int64]*models.AdminUser, len(userHIds))

	defer func() {
		if err == nil {
			return
		}
		logContent := map[string]interface{}{
			"userHIds": userHIds,
			"err":      err.Error(),
		}
		r.Context.Error(logContent, "DaoPermitUserImplGetAdminUserByIdsFromDb")
	}()

	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminUser
		var res []*models.AdminUser
		actRes = r.GetDefaultActErrorHandlerResult(m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where("`user_hid` IN (?)", userHIds).
			Find(&res).
			Error
		if actRes.Err != nil {
			return
		}
		for _, item := range res {
			resData[item.UserHid] = item
		}
		return
	})

	return
}
func (r *DaoPermitUserImpl) DeleteAdminUserGroup(userHIds ...int64) (err error) {

	if len(userHIds) == 0 {
		return
	}

	logC := map[string]interface{}{
		"userHIds": userHIds,
	}

	defer func() {
		if err != nil {
			logC["err"] = err.Error()
			r.Context.Error(logC, "DaoPermitUserImplDeleteAdminUserGroup")
			err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		}
	}()

	err = r.ActErrorHandler(func() (actErrorHandlerResult *base.ActErrorHandlerResult) {
		var m models.AdminUserGroup
		actErrorHandlerResult = r.GetDefaultActErrorHandlerResult(&m)
		actErrorHandlerResult.Err = actErrorHandlerResult.Db.Table(actErrorHandlerResult.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where("user_hid IN (?)", userHIds).
			Delete(&m).Error
		return
	})
	return
}

func (r *DaoPermitUserImpl) GetUserHidByGroupIds(groupIds ...int64) (res []int64, err error) {
	res = []int64{}
	if len(groupIds) == 0 {
		return
	}

	var (
		list    []*models.AdminUserGroup
		mapUser map[int64]int64
		l       int
	)

	if list, err = NewDaoPermitGroup(r.Context).
		GetGroupUserByIds(groupIds...); err != nil {
		return
	}

	l = len(list)
	mapUser = make(map[int64]int64, l)
	res = make([]int64, 0, l)

	for _, item := range list {
		if _, ok := mapUser[item.UserHid]; ok {
			continue
		}
		mapUser[item.UserHid] = item.UserHid
		res = append(res, item.UserHid)
	}

	return
}
func (r *DaoPermitUserImpl) GetAdminUserCount(db *gorm.DB, arg *wrappers.ArgAdminUser) (total int64, dba *gorm.DB, err error) {
	var m models.AdminUser
	if db == nil {
		db = r.Context.Db
	}

	if len(arg.GroupIds) > 0 {
		var userHIds []int64
		if userHIds, err = r.GetUserHidByGroupIds(arg.GroupIds...); err != nil {
			return
		}
		if len(userHIds) == 0 {
			return
		}
		arg.UserHIdArray = append(arg.UserHIdArray, userHIds...)
	}

	dba = db.Table(m.TableName()).Scopes(base.ScopesDeletedAt())
	if arg.Name != "" {
		dba = dba.Where("real_name LIKE ?", "%"+arg.Name+"%")
	}
	if len(arg.UserHIdArray) > 0 {
		dba = dba.Where("user_hid IN(?)", arg.UserHIdArray)
	}
	if arg.CannotUse >= 0 {
		dba = dba.Where("can_not_use =?", arg.CannotUse)
	}

	if arg.Mobile != "" {
		dba = dba.Where("mobile =?", arg.Mobile)
	}
	logContent := map[string]interface{}{"arg": arg}

	defer func() {
		if err != nil {
			logContent["err"] = err.Error()
			r.Context.Error(logContent, "DaoPermitUserImplGetAdminUserCount")
			err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		}
	}()

	err = r.ActErrorHandler(func() (actErrorHandlerResult *base.ActErrorHandlerResult) {
		actErrorHandlerResult = r.GetDefaultActErrorHandlerResult(&m)
		actErrorHandlerResult.Err = dba.Count(&total).Error
		return
	})

	return
}

func (r *DaoPermitUserImpl) GetAdminUserList(db *gorm.DB, arg *wrappers.ArgAdminUser, pager *response.Pager) (res []models.AdminUser, err error) {
	res = []models.AdminUser{}
	if err = db.Offset(pager.GetFromAndLimit()).Limit(arg.PageSize).Find(&res).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"arg":   arg,
			"pager": pager,
			"err":   err,
		}, "DaoPermitUserImplGetAdminUserList")
	}
	return
}

func (r *DaoPermitUserImpl) DeleteAdminUser(userHIds ...int64) (err error) {
	if len(userHIds) == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"ids": userHIds,
			"err": err,
		}, "DaoPermitUserImplAdminUserDelete")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminUser
		actRes = r.GetDefaultActErrorHandlerResult(m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Where("user_hid IN (?) ", userHIds).
			Scopes(base.ScopesDeletedAt()).
			Updates(map[string]interface{}{
				"deleted_at": time.Now().Format(utils.DateTimeGeneral),
			}).
			Error
		return
	})

	return
}

func (r *DaoPermitUserImpl) AdminUserAdd(dataUser []base.ModelBase) (err error) {

	logContent := map[string]interface{}{
		"dataUser": dataUser,
	}
	defer func() {
		if err != nil {
			logContent["err"] = err.Error()
			r.Context.Error(logContent, "DaoPermitUserImplAdminUserAdd")
			err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		}
	}()

	err = r.ActErrorHandler(func() (actErrorHandlerResult *base.ActErrorHandlerResult) {
		var dt *models.AdminUser
		actErrorHandlerResult = r.GetDefaultActErrorHandlerResult(dt)
		actErrorHandlerResult.Err = r.BatchAdd(
			actErrorHandlerResult.ParseBatchAddDataParameter(
				base.BatchAddDataParameterData(dataUser),
			),
		)
		return
	})

	return

}

func (r *DaoPermitUserImpl) GetGroupByUserId(userId int64) (res []wrappers.AdminGroupUserStruct, err error) {
	if userId == 0 {
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
	}, "DaoPermitUserImplGetGroupByUserId")
	return
}

func (r *DaoPermitUserImpl) UpdateDataByUserHIds(data map[string]interface{}, userHIds ...int64) (err error) {

	defer func() {
		if err != nil {
			r.Context.Error(map[string]interface{}{
				"data":     data,
				"userHIds": userHIds,
				"err":      err.Error(),
			}, "DaoPermitUserImplUpdateDataByUserHIds")
			err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		}
	}()
	err = r.ActErrorHandler(func() (actErrorHandlerResult *base.ActErrorHandlerResult) {
		var m *models.AdminUser
		actErrorHandlerResult = r.GetDefaultActErrorHandlerResult(m)
		actErrorHandlerResult.Err = actErrorHandlerResult.Db.
			Table(actErrorHandlerResult.TableName).Scopes(base.ScopesDeletedAt()).
			Where("user_hid IN (?)", userHIds).
			Updates(data).Error
		return
	})
	return
}

func NewDaoPermitUser(ctx ...*base.Context) (res daos.DaoPermitUser) {
	p := &DaoPermitUserImpl{}
	p.SetContext(ctx...)
	if p.ctx == nil {
		p.ctx = context.TODO()
	}

	return p
}
