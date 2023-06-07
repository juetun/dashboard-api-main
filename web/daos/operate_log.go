package daos

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type DaoOperateLog interface {

	AddLog(list []*models.OperateLog) (err error)

	GetCount(arg *wrapper_admin.ArgOperateLog) (actResObject *base.ActErrorHandlerResult, count int64, err error)

	GetList(actResObject *base.ActErrorHandlerResult, arg *wrapper_admin.ArgOperateLog, pager *response.Pager) (list []*models.OperateLog, err error)
}
