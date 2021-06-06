/**
* @Author:changjiang
* @Description:
* @File:permit_import
* @Version: 1.0.0
* @Date 2021/6/5 11:35 下午
 */
package dao_impl

import (
	"fmt"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
)

type PermitImportImpl struct {
	base.ServiceDao
}

func (r *PermitImportImpl) DeleteByCondition(condition map[string]interface{}) (res bool, err error) {
	var m models.AdminImport
	if len(condition) == 0 {
		err = fmt.Errorf("您没有选择要删除的数据")
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"err":       err.Error(),
		}, "permitImportImplDeleteByCondition0")
		return
	}
	if err = r.Context.Db.Table(m.TableName()).Where(condition).
		Delete(&models.AdminImport{}).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"err":       err.Error(),
		}, "permitImportImplDeleteByCondition1")
		return
	}
	return
}

func (r *PermitImportImpl) UpdateByCondition(condition, data map[string]interface{}) (res bool, err error) {
	var m models.AdminImport
	if len(condition) == 0 || len(data) == 0 {
		err = fmt.Errorf("您没有选择要更新的数据")
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"data":      data,
			"err":       err.Error(),
		}, "permitImportImplUpdateByCondition0")
		return
	}
	if err = r.Context.Db.Table(m.TableName()).Where(condition).Update(data).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"data":      data,
			"err":       err.Error(),
		}, "permitImportImplUpdateByCondition1")
		return
	}
	res = true
	return
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
