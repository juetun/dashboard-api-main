package daos

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type DaoHelp interface {

	AddHelpDocument(data *models.HelpDocument) (err error)

	GetByPKey(arg *base.ArgGetByStringIds) (res map[string]*models.HelpDocument,err error)

	GetByIds(arg *base.ArgGetByNumberIds) (res map[int64]*models.HelpDocument,err error)


	GetCountByArg(list *wrapper_admin.ArgHelpList, pager *response.Pager) (*base.ActErrorHandlerResult, error)

	GetListByArg(list *wrapper_admin.ArgHelpList, result *base.ActErrorHandlerResult,pagerObject *response.Pager) ([]*models.HelpDocument, error)
}
