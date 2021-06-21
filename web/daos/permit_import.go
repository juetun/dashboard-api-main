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

type PermitImport interface {

	BatchMenuImport(tableName string, list []models.AdminMenuImport) (err error)

	GetImportMenuByImportIds(iIds ...int) (list []models.AdminMenuImport, err error)
	UpdateByCondition(condition, data map[string]interface{}) (res bool, err error)
	DeleteByCondition(condition map[string]interface{}) (res bool, err error)
}
