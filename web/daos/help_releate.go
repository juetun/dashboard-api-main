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

	//注:此处的参数ID无须传值，只是传是否刷新缓存使用
	GetAllHelpRelate(arg *base.ArgGetByStringIds) (res models.HelpDocumentRelateCaches, err error)

	GetByIdFromDb(id ...int64) (res []*models.HelpDocumentRelate, err error)
}
