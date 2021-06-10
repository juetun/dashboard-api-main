//Package daos
//@Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package daos

import (
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type DaoExport interface {
	Update(model *models.ZExportData) (err error)
	UpdateByHIds(data map[string]interface{}, hIds *[]string) (err error)
	Create(model *models.ZExportData) (res bool, err error)
	Progress(args *wrappers.ArgumentsExportProgress) (res *[]models.ZExportData, err error)
	GetListByUser(userHid string, limit int) (res *[]models.ZExportData, err error)
}
