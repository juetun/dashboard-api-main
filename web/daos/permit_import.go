// Package daos
/**
* @Author:changjiang
* @Description:
* @File:permit_import
* @Version: 1.0.0
* @Date 2021/6/5 11:35 下午
 */
package daos

import (
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_intranet"
)

type DaoPermitImport interface {

	GetDefaultOpenImportByMenuIds(menuIds ...int64) (res []wrappers.AdminImportWithMenu, err error)

	GetImportByCondition(condition map[string]interface{}) (list []models.AdminImport, err error)

	GetImportForGateway(arg *wrapper_intranet.ArgGetImportPermit) (list []models.AdminImport, err error)

	GetImportFromDbByIds(ids ...int64) (res map[int64]*models.AdminImport, err error)

	AddData(adminImport *models.AdminImport) (err error)

	BatchAddData([]models.AdminImport) (err error)

	DeleteImportByIds(id ...int64) (err error)

	UpdateMenuImport(condition string, data map[string]interface{}) (err error)

	BatchMenuImport(list []*models.AdminMenuImport) (err error)

	GetImportMenuByImportIds(iIds ...int64) (list []models.AdminMenuImport, err error)

	GetChildImportByMenuId(menuIds ...int64) (list []models.AdminMenuImport, err error)

	UpdateByCondition(condition interface{}, data map[string]interface{}) (res bool, err error)

	DeleteByCondition(condition interface{}) (res bool, err error)

	GetMenuImportByMenuIdAndImportIds(menuId int64, importIds ...int64) (res []*models.AdminMenuImport,err error)
}
