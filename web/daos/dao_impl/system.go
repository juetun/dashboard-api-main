package dao_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
)

type DaoSystemImpl struct {
	base.ServiceDao
}

func (r *DaoSystemImpl) GetSystemList() (res []*models.ZBaseSys, err error) {
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var system *models.ZBaseSys
		actRes = r.GetDefaultActErrorHandlerResult(system)
		actRes.Err = actRes.Db.Table(actRes.TableName).Find(&res).Error
		return
	})
	return
}

func NewDaoSystem(c ...*base.Context) daos.DaoSystem {
	p := &DaoSystemImpl{}
	p.SetContext(c...)
	return p
}
