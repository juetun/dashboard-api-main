// Package daos
/**
* @Author:changjiang
* @Description:
* @File:service
* @Version: 1.0.0
* @Date 2021/5/23 10:35 上午
 */
package daos

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"gorm.io/gorm"
)

type DaoService interface {
	GetByKeys(keys ...string) (res []models.AdminApp, err error)
	Update(condition, data map[string]interface{}) (err error)
	Create(app *models.AdminApp) (err error)
	GetByIds(id ...int) (res []models.AdminApp, err error)
	GetCount(db *gorm.DB, arg *wrappers.ArgServiceList) (total int64, dba *gorm.DB, err error)
	GetList(db *gorm.DB, arg *wrappers.ArgServiceList, page *response.Pager) (list []models.AdminApp, err error)

	GetImportMenuByModule(module string) (res []wrappers.ImportMenu, err error)

	GetByPort(portArgNumber *base.ArgGetByNumberIds) (res map[int64]*models.AdminApp, err error)

	GetByUniqueKey(ids *base.ArgGetByStringIds) (res map[string]*models.AdminApp,err error)
}
