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
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
)

type PermitImportImpl struct {
	base.ServiceDao
}

func (r *PermitImportImpl) BatchMenuImport(tableName string, list []models.AdminMenuImport) (err error) {
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
		}, "PermitImportImplBatchMenuImport error")
	}
	return
}
func (r *PermitImportImpl) getColumnName(s string) (res string) {
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

func (r *PermitImportImpl) DeleteByCondition(condition interface{}) (res bool, err error) {
	var m models.AdminImport
	if condition == nil {
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

func (r *PermitImportImpl) UpdateByCondition(condition interface{}, data map[string]interface{}) (res bool, err error) {
	var m models.AdminImport
	if condition == nil || len(data) == 0 {
		err = fmt.Errorf("您没有选择要更新的数据")
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"data":      data,
			"err":       err.Error(),
		}, "permitImportImplUpdateByCondition0")
		return
	}
	if err = r.Context.Db.Table(m.TableName()).
		Where(condition).
		Updates(data).Error; err != nil {
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
func (r *PermitImportImpl) UpdateMenuImport(condition string, data map[string]interface{}) (err error) {

	var m models.AdminMenuImport
	if err = r.Context.Db.Table(m.TableName()).
		Where(condition).
		Updates(data).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"err":       err.Error(),
		}, "permitImportImplUpdateMenuImport")
		return
	}
	return
}

func NewPermitImportImpl(context ...*base.Context) daos.PermitImport {
	p := &PermitImportImpl{}
	p.SetContext(context...)
	p.Context.Db,p.Context.DbName = base.GetDbClient(&base.GetDbClientData{
		Context:     p.Context,
		DbNameSpace: daos.DatabaseAdmin,
	})
	return p
}
