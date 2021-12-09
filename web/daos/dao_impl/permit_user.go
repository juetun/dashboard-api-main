// Package dao_impl /**
package dao_impl

import (
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
}

func (r *DaoPermitUserImpl) GetAdminUserCount(db *gorm.DB, arg *wrappers.ArgAdminUser) (total int64, dba *gorm.DB, err error) {
	var m models.AdminUser
	if db == nil {
		db = r.Context.Db
	}
	dba = db.Table(m.TableName()).Scopes(base.ScopesDeletedAt())
	if arg.Name != "" {
		dba = dba.Where("real_name LIKE ?", "%"+arg.Name+"%")
	}
	if arg.UserHId != "" {
		dba = dba.Where("user_hid=?", arg.UserHId)
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
		actErrorHandlerResult = &base.ActErrorHandlerResult{
			Db:        r.Context.Db,
			DbName:    r.Context.DbName,
			TableName: m.TableName(),
			Model:     &m,
			Err:       dba.Count(&total).Error,
		}
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

func (r *DaoPermitUserImpl) DeleteAdminUser(ids []string) (err error) {
	if len(ids) == 0 {
		return
	}
	var m models.AdminUser
	err = r.Context.Db.Table(m.TableName()).
		Where("user_hid IN (?) ", ids).
		Scopes(base.ScopesDeletedAt()).
		Updates(map[string]interface{}{
			"deleted_at": time.Now().Format("2006-01-02 15:04:05"),
		}).
		Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"ids": ids,
			"err": err,
		}, "DaoPermitUserImplAdminUserDelete")
	}
	return
}
func (r *DaoPermitUserImpl) AdminUserAdd(arg *models.AdminUser) (err error) {
	var data = base.BatchAddDataParameter{
		DbName:    r.Context.DbName,
		TableName: arg.TableName(),
		Data: []base.ModelBase{
			arg,
		},
	}
	if err = r.BatchAdd(&data); err == nil {
		return
	}
	r.Context.Error(map[string]interface{}{
		"arg": arg,
		"err": err,
	}, "DaoPermitUserImplAdminUserAdd")
	return

}

func (r *DaoPermitUserImpl) UpdateDataByUserHIds(data map[string]interface{}, userHIds ...string) (err error) {

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
		var m models.AdminUser
		actErrorHandlerResult = &base.ActErrorHandlerResult{
			Db:        r.Context.Db,
			DbName:    r.Context.DbName,
			TableName: m.TableName(),
			Model:     &m,
		}
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
	return p
}
