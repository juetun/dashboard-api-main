// Package dao_impl /**
package dao_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
)

type DaoPermitUserImpl struct {
	base.ServiceDao
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
