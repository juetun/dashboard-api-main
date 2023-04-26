package daos

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type DaoHelpRelate interface {
	AddOneHelpRelate(relate *models.HelpDocumentRelate) (err error)

	GetByTopId(bizCode string, topIds ...int64) (res map[int64][]*wrapper_admin.ResultHelpTreeItem, err error)

	GetByTopHelp() (res []*models.HelpDocumentRelate, err error)

	UpdateById(id int64, data map[string]interface{}) (err error)

	GetByDocKeys(arg *base.ArgGetByStringIds) (res map[string]*models.HelpDocumentRelate, err error)

	GetAllHelpRelate(arg *base.ArgGetByStringIds) (res models.HelpDocumentRelateCaches, err error)
}
