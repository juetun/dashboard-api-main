// Package dao_impl
/**
* @Author:changjiang
* @Description:
* @File:permit_import
* @Version: 1.0.0
* @Date 2021/6/5 11:35 下午
 */
package dao_impl

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"gorm.io/gorm"
)

type DaoPermitImport struct {
	base.ServiceDao
}

func (r *DaoPermitImport) AddData(adminImport *models.AdminImport) (err error) {
	logContent := map[string]interface{}{
		"adminImport": adminImport,
	}
	defer func() {
		if err != nil {
			logContent["err"] = err.Error()
			r.Context.Error(logContent, "DaoPermitImportAddData")
		}
	}()
	if err = adminImport.InitRegexpString(); err != nil {
		logContent["desc"] = "InitRegexpString"
		return
	}
	if err = r.Context.Db.Create(adminImport).Error; err != nil {
		logContent["desc"] = "Create"
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}

func (r *DaoPermitImport) getColumnName(s string) (res string) {
	li := strings.Split(s, ";")
	res = s
	for _, s2 := range li {
		if s2 == "" {
			return
		}
		li1 := strings.Split(s2, ":")
		if len(li1) > 1 && li1[0] == "column" {
			res = li1[1]
		}
	}
	if res == "primary_key" {
		res = "id"
	}
	return
}
func (r *DaoPermitImport) BatchMenuImport(tableName string, list []models.AdminMenuImport) (err error) {
	if len(list) == 0 {
		return
	}
	var ind = 0
	var keys []string
	var fieldNum int
	var columns []string
	var vals []interface{}
	var vvGroup = make([]string, 0, len(list))
	vals = make([]interface{}, 0, fieldNum)
	tagValue := "gorm"
	var tName string
	var tag string
	for _, data := range list {

		var valueStruct reflect.Value
		if ind == 0 {
			types := reflect.TypeOf(data)
			fieldNum = types.NumField()
			keys = make([]string, 0, fieldNum)
			columns = make([]string, 0, fieldNum)
			for i := 0; i < fieldNum; i++ {
				t := types.Field(i).Tag
				tName = t.Get(tagValue)
				tag = r.getColumnName(tName)
				if tag == "id" {
					continue
				}
				keys = append(keys, tag)
				if tag == "created_at" {
					continue
				}
				columns = append(columns, fmt.Sprintf("`%s`=VALUES(`%s`)", tag, tag))
			}
		}
		types := reflect.TypeOf(data)
		vv := make([]string, 0, fieldNum)
		for i := 0; i < fieldNum; i++ {
			tName = types.Field(i).Tag.Get(tagValue)
			if tag = r.getColumnName(tName); tag == "id" {
				continue
			}
			values := reflect.ValueOf(data)
			valueStruct = values.Field(i)
			switch valueStruct.Kind() {
			case reflect.Interface:
				vals = append(vals, valueStruct.Interface())
			case reflect.Ptr:
				if valueStruct.IsZero() {
					vals = append(vals, nil)
				} else {
					vals = append(vals, valueStruct.Elem().Interface())
				}
			case reflect.Bool:
				vals = append(vals, strconv.FormatBool(valueStruct.Bool()))
			case reflect.String:
				vals = append(vals, valueStruct.String())
			default:
				switch valueStruct.Type().String() {
				case "time.Time":
					vals = append(vals, valueStruct.Interface().(time.Time).Format("2006-01-02 15:04:05"))
				case "time.Duration":
					vals = append(vals, valueStruct.Interface().(time.Duration).String())
				case "int":
					vals = append(vals, fmt.Sprintf("%v", valueStruct.Int()))
				default:
					vals = append(vals, valueStruct.String())
				}
			}
			vv = append(vv, "?")
		}
		vvGroup = append(vvGroup, fmt.Sprintf("(%s)", strings.Join(vv, ",")))
		ind++
	}
	sql := fmt.Sprintf("INSERT INTO `%s`(`"+strings.Join(keys, "`,`")+"`) VALUES "+strings.Join(vvGroup, ",")+
		" ON DUPLICATE KEY UPDATE "+strings.Join(columns, ","), tableName)
	if err = r.Context.Db.Exec(sql, vals...).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"sql":  sql,
			"data": list,
			"vals": vals,
			"err":  err,
		}, "DaoPermitImportBatchMenuImport error")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}
	return
}

func (r *DaoPermitImport) DeleteByCondition(condition interface{}) (res bool, err error) {
	var m models.AdminImport
	if condition == nil {
		err = fmt.Errorf("您没有选择要删除的数据")
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"err":       err.Error(),
		}, "DaoPermitImportDeleteByCondition0")
		return
	}
	if err = r.Context.Db.Table(m.TableName()).Where(condition).
		Delete(&models.AdminImport{}).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"err":       err.Error(),
		}, "DaoPermitImportDeleteByCondition1")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}

func (r *DaoPermitImport) UpdateByCondition(condition interface{}, data map[string]interface{}) (res bool, err error) {
	var m models.AdminImport
	if condition == nil || len(data) == 0 {
		err = fmt.Errorf("您没有选择要更新的数据")
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"data":      data,
			"err":       err.Error(),
		}, "DaoPermitImportUpdateByCondition0")
		return
	}
	if err = r.Context.Db.Table(m.TableName()).
		Where(condition).
		Updates(data).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"data":      data,
			"err":       err.Error(),
		}, "DaoPermitImportUpdateByCondition1")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	res = true
	return
}

func (r *DaoPermitImport) GetChildImportByMenuId(menuIds ...int) (list []models.AdminMenuImport, err error) {

	list = []models.AdminMenuImport{}
	if len(menuIds) == 0 {
		return
	}
	var m models.AdminMenuImport
	if err = r.Context.Db.Table(m.TableName()).Scopes(base.ScopesDeletedAt()).
		Where("menu_id IN (?)", menuIds).
		Find(&list).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"iIds": menuIds,
			"err":  err.Error(),
		}, "DaoPermitImportGetChildImportByMenuId")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}
func (r *DaoPermitImport) GetImportMenuByImportIds(iIds ...int) (list []models.AdminMenuImport, err error) {
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
		}, "DaoPermitImportGetImportMenuByImportIds")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}
func (r *DaoPermitImport) UpdateMenuImport(condition string, data map[string]interface{}) (err error) {

	var m models.AdminMenuImport
	var e error
	var fetchData = base.FetchDataParameter{
		SourceDb:  r.Context.Db,
		DbName:    r.Context.DbName,
		TableName: m.TableName(),
	}
	if e = fetchData.SourceDb.Table(fetchData.TableName).
		Where(condition).
		Updates(data).
		Error; e == nil {
		return
	}

	logContent := map[string]interface{}{
		"condition": condition,
		"e":         e.Error(),
	}
	defer func() {
		r.Context.Error(logContent, "DaoPermitImportUpdateMenuImport")
	}()
	if err = r.CreateTableWithError(fetchData.SourceDb, fetchData.TableName, e, &m, base.TableSetOption{"COMMENT": m.GetTableComment()}); err != nil {
		return
	}
	if err = r.UpdateMenuImport(condition, data); err != nil {
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}

	return
}

func (r *DaoPermitImport) BatchAddData(list []models.AdminImport) (err error) {
	var m models.AdminImport
	data := &base.BatchAddDataParameter{
		DbName:    r.Context.DbName,
		Db:        r.Context.Db,
		TableName: m.TableName(),
	}
	for _, adminImport := range list {
		data.Data = append(data.Data, &adminImport)
	}
	var e error
	if e = r.BatchAdd(data); e == nil {
		return
	}
	logContent := map[string]interface{}{
		"list": list,
		"e":    e.Error(),
	}
	defer func() {
		r.Context.Error(logContent, "DaoPermitImportBatchAddData")
	}()
	if err = r.CreateTableWithError(data.Db, data.TableName, e, &m, base.TableSetOption{"COMMENT": m.GetTableComment()}); err != nil {
		logContent["err1"] = err.Error()
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	if err = r.BatchAddData(list); err != nil {
		logContent["err2"] = err.Error()
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}
func (r *DaoPermitImport) DeleteImportByIds(id ...int) (err error) {
	if len(id) == 0 {
		return
	}
	var m models.AdminImport
	if err = r.Context.Db.Table(m.TableName()).Unscoped().
		Where("id IN(?)", id).
		Delete(&models.AdminImport{}).Error; err == nil {
		r.Context.Error(map[string]interface{}{
			"id":  id,
			"err": err,
		}, "daoPermitDeleteImportByIds0")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}

	var m1 models.AdminUserGroupPermit
	if err = r.Context.Db.Table(m1.TableName()).
		Where("menu_id IN(?) AND path_type=?", id, "api").
		Delete(&models.AdminImport{}).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"id":  id,
			"err": err,
		}, "daoPermitDeleteImportByIds1")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}

func (r *DaoPermitImport) GetImportByCondition(condition map[string]interface{}) (list []models.AdminImport, err error) {
	list = []models.AdminImport{}
	if len(condition) == 0 {
		return
	}
	var m models.AdminImport

	if err = r.Context.Db.Table(m.TableName()).Where(condition).Find(&list).Limit(1000).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"err":       err.Error(),
		}, "DaoPermitImportGetImportByCondition")
		return
	}
	return
}
func NewDaoPermitImport(c ...*base.Context) daos.DaoPermitImport {
	p := &DaoPermitImport{}
	p.SetContext(c...)
	s, ctx := p.Context.GetTraceId()
	p.Context.Db, p.Context.DbName, _ = base.GetDbClient(&base.GetDbClientData{
		Context:     p.Context,
		DbNameSpace: daos.DatabaseAdmin,
		CallBack: func(db *gorm.DB, dbName string) (dba *gorm.DB, err error) {
			dba = db.WithContext(context.WithValue(ctx, app_obj.DbContextValueKey, base.DbContextValue{
				TraceId: s,
				DbName:  dbName,
			}))
			return
		},
	})
	return p
}
