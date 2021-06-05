/**
* @Author:changjiang
* @Description:
* @File:permit_import
* @Version: 1.0.0
* @Date 2021/6/5 11:35 下午
 */
package dao_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
)

type PermitImportImpl struct {
	base.ServiceDao
}

func (r *PermitImportImpl) GetImportMenuByImportIds(iIds ...int) (list []models.AdminMenuImport, err error) {
	list = []models.AdminMenuImport{}
	if len(iIds) == 0 {
		return
	}
	var m models.AdminMenuImport
	if err = r.Context.Db.Table(m.TableName()).
		Where("import_id IN (?)", iIds).
		Find(&list).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"iIds": iIds,
			"err":  err.Error(),
		}, "permitImportImplGetImportMenuByImportIds")
		return
	}
	return
}

func NewPermitImportImpl(context ...*base.Context) daos.PermitImport {
	p := &PermitImportImpl{}
	p.SetContext(context...)
	p.Context.Db = base.GetDbClient(&base.GetDbClientData{
		Context:     p.Context,
		DbNameSpace: daos.DatabaseAdmin,
	})
	return p
}
