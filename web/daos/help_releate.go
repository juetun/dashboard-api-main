package daos

import (
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type DaoHelpRelate interface {
	AddOneHelpRelate(relate *models.HelpDocumentRelate) (err error)
	GetByTopId(bizCode string,topIds ...int64) (res map[int64][]*wrapper_admin.ResultHelpTreeItem, err error)
}
