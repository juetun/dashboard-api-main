package dao_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type DaoOperateLogImpl struct {
	base.ServiceDao
}

func (r *DaoOperateLogImpl) AddLog(list []*models.OperateLog) (err error) {
	var (
		l    = len(list)
		data = make([]base.ModelBase, 0, l)
	)
	if l == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"list": list,
			"err":  err.Error(),
		}, "DaoOperateLogAddLog")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	for _, item := range list {
		data = append(data, item)
	}
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.OperateLog
		actRes = r.GetDefaultActErrorHandlerResult(m)
		if actRes.Err = r.BatchAdd(actRes.ParseBatchAddDataParameter(base.BatchAddDataParameterData(data))); actRes.Err != nil {
			return
		}
		return
	})
	return
}

func (r *DaoOperateLogImpl) GetCount(arg *wrapper_admin.ArgOperateLog) (actResObject *base.ActErrorHandlerResult, total int64, err error) {
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "DaoOperateLogGetCount")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		actRes = r.getActErrorHandlerResult(arg)
		actResObject = actRes
		actResObject.Err = actResObject.Db.Limit(1).
			Count(&total).
			Error
		if actResObject.Err != nil {
			var m *models.OperateLog
			actRes = r.GetDefaultActErrorHandlerResult(m)
			actRes.Err = actResObject.Err
		}

		return
	})
	return
}

func (r *DaoOperateLogImpl) getActErrorHandlerResult(arg *wrapper_admin.ArgOperateLog, actResList ...*base.ActErrorHandlerResult) (actRes *base.ActErrorHandlerResult) {
	var m *models.OperateLog
	if len(actResList) > 0 {
		actRes = actResList[0]
		return
	}
	actRes = r.GetDefaultActErrorHandlerResult(m)
	actRes.Db = actRes.Db.Table(actRes.TableName)
	if !arg.StartTime.IsZero() {
		actRes.Db = actRes.Db.Where("created_at>=?", arg.StartTime.Format(utils.DateTimeGeneral))
	}
	if !arg.OverTime.IsZero() {
		actRes.Db = actRes.Db.Where("created_at<=?", arg.OverTime.Format(utils.DateTimeGeneral))
	}
	return
}

func (r *DaoOperateLogImpl) GetList(actResult *base.ActErrorHandlerResult, arg *wrapper_admin.ArgOperateLog, pager *response.Pager) (list []*models.OperateLog, err error) {
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "DaoOperateLogGetList")
	}()

	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		actRes = r.getActErrorHandlerResult(arg, actResult)
		actRes.Err = actRes.Db.Table(actRes.TableName).Offset(arg.GetOffset()).
			Limit(arg.PageSize).Order("`created_at` DESC").
			Find(&list).
			Error
		return
	})
	return
}

func NewDaoOperateLog(c ...*base.Context) daos.DaoOperateLog {
	p := &DaoOperateLogImpl{}
	p.SetContext(c...)
	return p
}
