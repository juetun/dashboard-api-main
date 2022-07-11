package dao_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
)

type DaoHelpRelateImpl struct {
	base.ServiceDao
}

func (r *DaoHelpRelateImpl) AddOneHelpRelate(relate *models.HelpDocumentRelate) (err error) {
	if relate == nil {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"relate": relate,
			"err":    err.Error(),
		}, "DaoHelpRelateImplAddOneHelpRelate")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		actRes = r.GetDefaultActErrorHandlerResult(relate)
		actRes.Err = r.AddOneData(actRes.ParseAddOneDataParameter(base.AddOneDataParameterModel(relate)))
		return
	})
	return
}

func NewDaoHelpRelate(c ...*base.Context) daos.DaoHelpRelate {
	p := &DaoHelpRelateImpl{}
	p.SetContext(c...)
	return p
}
