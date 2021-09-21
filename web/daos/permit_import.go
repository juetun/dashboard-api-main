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
)


type DaoPermitImport interface {

	BatchAddData([]models.AdminImport) (err error)

	DeleteImportByIds(id ...int) (err error)

	UpdateMenuImport(condition string, data map[string]interface{}) (err error)

	BatchMenuImport(tableName string, list []models.AdminMenuImport) (err error)

	GetImportMenuByImportIds(iIds ...int) (list []models.AdminMenuImport, err error)

	GetChildImportByMenuId(menuIds ...int) (list []models.AdminMenuImport, err error)

	UpdateByCondition(condition interface{}, data map[string]interface{}) (res bool, err error)

	DeleteByCondition(condition interface{}) (res bool, err error)
}
