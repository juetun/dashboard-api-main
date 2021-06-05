/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 11:19 下午
 */
package dao_impl

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type DaoPermitImpl struct {
	base.ServiceDao
}

func NewDaoPermit(context ...*base.Context) daos.DaoPermit {
	p := &DaoPermitImpl{}
	p.SetContext(context...)
	p.Context.Db = base.GetDbClient(&base.GetDbClientData{
		Context:     p.Context,
		DbNameSpace: daos.DatabaseAdmin,
	})
	return p
}
func (r *DaoPermitImpl) fetchDb(db *gorm.DB, arg *wrappers.ArgImportList) (dba *gorm.DB) {
	var m models.AdminImport
	dba = db
	if db == nil {
		dba = r.Context.Db.Table(m.TableName())
	}
	return
}
func (r *DaoPermitImpl) GetImportListCount(db *gorm.DB, arg *wrappers.ArgImportList) (totalCount int, dba *gorm.DB, err error) {
	dba = r.fetchDb(db, arg)
	if err = dba.Count(&totalCount).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "daoPermitImplGetImportListCount")
		return
	}
	return
}
func (r *DaoPermitImpl) GetImportListData(db *gorm.DB, arg *wrappers.ArgImportList, pager *response.Pager) (res []models.AdminImport, err error) {

	res = []models.AdminImport{}
	if err = r.fetchDb(db, arg).
		Offset(pager.GetOffset()).
		Order(arg.Order).
		Limit(pager.PageSize).
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "daoPermitImplGetImportListData")
		return
	}

	return
}
func (r *DaoPermitImpl) GetPermitImportByModule(arg *wrappers.ArgPermitMenu) (res []wrappers.Op, err error) {
	var adminUserGroupPermit models.AdminUserGroupPermit
	var adminImport models.AdminImport
	var adminMenu models.AdminMenu
	if arg.IsSuperAdmin { // 如果是超级管理员
		if err = r.Context.Db.Table(adminImport.TableName()).Select("b.permit_key,c.permit_key AS menu_permit_key").
			Joins(fmt.Sprintf("AS b LEFT JOIN  %s as c ON b.menu_id=c.id", adminMenu.TableName())).
			Find(&res).
			Error; err != nil {
			r.Context.Error(map[string]interface{}{
				"arg": arg,
				"err": err.Error(),
			}, "daoPermitImplGetPermitImportByModule0")
			return
		}
		return
	}
	if err = r.Context.Db.Table(adminUserGroupPermit.TableName()).Select("b.permit_key,c.permit_key AS menu_permit_key").
		Joins(fmt.Sprintf("AS a LEFT JOIN  %s as b ON a.menu_id=b.id", adminImport.TableName())).
		Joins(fmt.Sprintf(" LEFT JOIN  %s as c ON b.menu_id=c.id", adminMenu.TableName())).
		Where("a.path_type=? AND a.group_id IN (?)", models.PathTypeApi, arg.GroupId).
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "daoPermitImplGetPermitImportByModule1")
		return
	}

	return
}
func (r *DaoPermitImpl) DeleteImportByIds(id ...int) (err error) {
	if len(id) == 0 {
		return
	}
	var m models.AdminImport
	if err = r.Context.Db.Table(m.TableName()).Unscoped().
		Where("id IN(?)", id).
		Delete(&models.AdminImport{}).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"id":  id,
			"err": err,
		}, "daoPermitDeleteImportByIds0")
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
		return
	}
	return
}
func (r *DaoPermitImpl) GetImportMenuId(menuId ...int) (list []models.AdminImport, err error) {
	list = []models.AdminImport{}
	var m models.AdminImport
	if len(menuId) == 0 {
		return
	}
	if err = r.Context.Db.Table(m.TableName()).Unscoped().
		Where("menu_id IN(?)", menuId).
		Find(&list).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"menuId": menuId,
			"err":    err,
		}, "daoPermitGetImportMenuId")
		return
	}
	return
}

func (r *DaoPermitImpl) CreateImport(data *models.AdminImport) (res bool, err error) {
	var m models.AdminImport
	if err = r.Context.Db.Table(m.TableName()).Create(data).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"data": data,
			"err":  err,
		}, "daoPermitCreateImport")
		return
	}

	return
}
func (r *DaoPermitImpl) UpdateAdminImport(condition, data map[string]interface{}) (res bool, err error) {
	var m models.AdminImport
	if len(condition) == 0 {
		return
	}
	if err = r.Context.Db.Table(m.TableName()).
		Where(condition).
		Update(data).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"data":      data,
			"err":       err,
		}, "daoPermitUpdateAdminImport")
		return
	}
	res = true
	return
}
func (r *DaoPermitImpl) GetAdminImportById(id ...int) (res []models.AdminImport, err error) {
	if len(id) == 0 {
		res = []models.AdminImport{}
		return
	}
	var m models.AdminImport
	err = r.Context.Db.Table(m.TableName()).
		Where("id IN(?)", id).
		Find(&res).Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"id":  id,
			"err": err,
		}, "daoPermitGetAdminImportByIdError")
		return
	}
	return
}
func (r *DaoPermitImpl) getColumnName(s string) (res string) {
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
	return
}
func (r *DaoPermitImpl) BatchGroupPermit(tableName string, list []models.AdminUserGroupPermit) (err error) {
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
				tName = types.Field(i).Tag.Get(tagValue)
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
		}, "BatchGroupPermit error")
	}
	return
}
func (r *DaoPermitImpl) DeleteGroupPermit(groupId int, pathType string, menuId ...int) (err error) {
	var m models.AdminUserGroupPermit
	db := r.Context.Db.Table(m.TableName()).
		Where("group_id =?", groupId)
	if pathType != "" {
		db = db.Where("path_type=?", pathType)
	}
	if len(menuId) > 0 {
		db = db.Where("menu_id IN(?)", menuId)
	}
	if err = db.
		Delete(&models.AdminUserGroupPermit{}).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"groupId": groupId,
			"err":     err,
		}, "daoPermitDeleteGroupPermit")
	}
	return
}

// func (r *DaoPermitImpl) DeleteByMenuIds(pageMenuIds ...string) (err error) {
// 	if len(pageMenuIds) == 0 {
// 		return
// 	}
// 	var m models.AdminUserGroupPermit
// 	if err = r.Context.Db.Table(m.TableName()).
// 		Where("path_type =? AND menu_id IN(?)", models.PathTypePage, pageMenuIds).
// 		Delete(&models.AdminUserGroupPermit{}).
// 		Error; err != nil {
// 		r.Context.Error(map[string]interface{}{
// 			"pageMenuIds": pageMenuIds,
// 			"err":         err,
// 		}, "daoPermitDeleteByMenuIds")
// 		return
// 	}
// 	return
// }
func (r *DaoPermitImpl) DeleteGroupPermitByGroupId(groupId int) (err error) {
	if groupId == 0 {
		return
	}
	var m models.AdminUserGroupPermit
	if err = r.Context.Db.Table(m.TableName()).
		Where("group_id=?", groupId).
		Delete(&models.AdminUserGroupPermit{}).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"groupId": groupId,
			"err":     err,
		}, "daoPermitDeleteGroupPermitByGroupId")
	}
	return
}

func (r *DaoPermitImpl) DeleteGroupPermitByMenuIds(groupId int, module string, pageMenuId, apiMenuId []int) (err error) {
	if len(apiMenuId) == 0 && len(pageMenuId) == 0 {
		return
	}
	var m models.AdminUserGroupPermit
	db := r.Context.Db.Table(m.TableName()).
		Where("group_id =? AND module=?", groupId, module)

	if len(apiMenuId) > 0 {
		parameters := []interface{}{
			apiMenuId,
			models.PathTypeApi,
		}
		sql := "(menu_id IN (?) AND path_type=?)"
		if len(pageMenuId) > 0 {
			parameters = append(parameters, []interface{}{pageMenuId, models.PathTypePage}...)
			sql += " OR (menu_id IN (?) AND path_type=?)"
		}
		db = db.Where(sql, parameters...)
	} else {
		if len(pageMenuId) > 0 {
			db = db.Where("menu_id IN (?) AND path_type=?", pageMenuId, models.PathTypePage)
		}
	}
	db = db.Where(db)
	if err = db.
		Delete(&models.AdminUserGroupPermit{}).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"groupId":    groupId,
			"apiMenuId":  apiMenuId,
			"pageMenuId": pageMenuId,
			"err":        err,
		}, "DeleteGroupPermitByMenuIds")
	}
	return
}
func (r *DaoPermitImpl) GetDefaultOpenImportByMenuIds(menuId ...int) (res []models.AdminImport, err error) {
	res = []models.AdminImport{}
	if len(menuId) == 0 {
		return
	}
	var m models.AdminImport
	if err = r.Context.Db.Table(m.TableName()).
		Where("menu_id IN(?) AND default_open=?", menuId, models.DefaultOpen).
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"menuId": menuId,
			"err":    err,
		}, "GetDefaultOpenImportByMenuIds")
		return
	}
	return
}
func (r *DaoPermitImpl) GetDefaultImportByMenuIds(pageType, module string, menuId ...int) (res []models.AdminImport, err error) {
	res = []models.AdminImport{}
	if len(menuId) == 0 {
		return
	}
	var m models.AdminImport
	db := r.Context.Db.Table(m.TableName())
	switch pageType {
	case models.PathTypePage:
		db = db.Where("menu_id IN(?) ", menuId)
	case models.PathTypeApi:
		db = db.Where("id IN(?) ", menuId)
	default:
		err = fmt.Errorf("pathtype is error")
		return
	}
	if err = db.
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"menuId": menuId,
			"err":    err,
		}, "GetDefaultImportByMenuIds")
		return
	}
	return
}
func (r *DaoPermitImpl) GetImportId(id int) (res models.AdminImport, err error) {
	var m models.AdminImport
	err1 := r.Context.Db.Table(m.TableName()).Where("id=?", id).Find(&res).Error
	if err1 != nil && gorm.IsRecordNotFoundError(err1) {
		err = err1
		r.Context.Error(map[string]interface{}{
			"id":  id,
			"err": err,
		}, "daoPermitGetImportId")
		return
	}
	return
}
func (r *DaoPermitImpl) GetImportCount(arg *wrappers.ArgGetImport, count *int) (db *gorm.DB, err error) {
	if arg.MenuId == 0 {
		return
	}
	var m models.AdminImport
	db = r.Context.Db.Table(m.TableName()).
		Where("`menu_id` = ? AND `deleted_at` IS NULL", arg.MenuId).
		Unscoped()
	err = db.Count(count).Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err,
		}, "daoPermitGetImportCount error")
	}
	return
}
func (r *DaoPermitImpl) GetSelectImportByImportId(groupId int, importId ...int) (res []models.AdminUserGroupPermit, err error) {
	if len(importId) == 0 {
		return
	}
	var m models.AdminUserGroupPermit
	if err = r.Context.Db.Table(m.TableName()).
		Where("menu_id IN(?) AND path_type=? AND group_id=?", importId, models.PathTypeApi, groupId).
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"importId": importId,
			"err":      err,
		}, "daoPermitGetSelectImportByImportId")
		return
	}
	return
}
func (r *DaoPermitImpl) GetImportList(db *gorm.DB, arg *wrappers.ArgGetImport) (res []models.AdminImport, err error) {
	err = db.Offset(arg.BaseQuery.GetOffset()).
		Limit(arg.PageSize).
		Find(&res).Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err,
		}, "daoPermitGetImportList")
		return
	}
	return
}
func (r *DaoPermitImpl) AdminUserGroupAdd(data []map[string]interface{}) (err error) {
	field := make([]string, 0, 10)
	dataMsg := make([]string, 0)
	dataTmp := make([]interface{}, 0)
	duplicate := make([]string, 0)
	dataMsgA := make([]string, 0)
	if len(data) == 0 {
		err = fmt.Errorf("您没有为模板配置参数,请至少配置一个参数")
		return
	}

	for key, value := range data {
		if key == 0 {
			for fieldString, _ := range value {
				field = append(field, fieldString)
			}
		}
		for _, v := range field {
			if key == 0 {
				if v != "created_at" {
					duplicate = append(duplicate, "`"+v+"`=VALUES(`"+v+"`)")
				}
				dataMsgA = append(dataMsgA, "?")
			}
			var v1 interface{}
			if _, ok := value[v]; ok {
				v1 = value[v]
			}
			dataTmp = append(dataTmp, v1)
		}
		dataMsg = append(dataMsg, "("+strings.Join(dataMsgA, ",")+")")
	}
	var m models.AdminUserGroup
	sql := "INSERT INTO `" + m.TableName() + "`(`" + strings.Join(field, "`,`") +
		"`) VALUES" + strings.Join(dataMsg, ",") + " ON DUPLICATE KEY UPDATE " + strings.Join(duplicate, ",")
	err = r.Context.Db.Exec(sql, dataTmp...).Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"data": data,
			"err":  err,
		}, "daoPermitAdminUserGroupAdd")
	}
	return
}
func (r *DaoPermitImpl) AdminUserGroupRelease(ids ...string) (err error) {
	if len(ids) == 0 {
		return
	}
	var m models.AdminUserGroup
	err = r.Context.Db.Table(m.TableName()).
		Where("id IN (?) ", ids).
		Delete(&models.AdminUserGroup{}).
		Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"ids": ids,
			"err": err,
		}, "daoPermitAdminUserGroupRelease")
	}
	return
}
func (r *DaoPermitImpl) AdminUserAdd(arg *models.AdminUser) (err error) {

	fields := []string{"user_hid", "real_name", "updated_at", "deleted_at",}
	var bt bytes.Buffer
	bt.WriteString("ON DUPLICATE KEY UPDATE ")
	for k, value := range fields {
		if k == 0 {
			bt.WriteString(fmt.Sprintf("%s=VALUES(%s)", value, value))
			continue
		}
		bt.WriteString(fmt.Sprintf(",%s=VALUES(%s)", value, value))
	}
	err = r.Context.Db.Set("gorm:insert_option", bt.String()).
		Create(arg).Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err,
		}, "daoPermitAdminUserAdd")
	}
	return
}
func (r *DaoPermitImpl) DeleteAdminUser(ids []string) (err error) {
	if len(ids) == 0 {
		return
	}
	var m models.AdminUser
	err = r.Context.Db.Table(m.TableName()).
		Where("user_hid IN (?) ", ids).
		Update(map[string]interface{}{
			"deleted_at": time.Now().Format("2006-01-02 15:04:05"),
		}).
		Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"ids": ids,
			"err": err,
		}, "daoPermitAdminUserDelete")
	}
	return
}
func (r *DaoPermitImpl) DeleteUserGroupByUserId(userId ...string) (err error) {
	if len(userId) == 0 {
		return
	}
	var m models.AdminUserGroup

	if err = r.Context.Db.Table(m.TableName()).
		Where("user_hid IN (?) ", userId).
		Update(map[string]interface{}{
			"deleted_at": time.Now().Format("2006-01-02 15:04:05"),
		}).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"userId": userId,
			"err":    err,
		}, "daoPermitDeleteUserGroupByUserId")
		return
	}
	return
}
func (r *DaoPermitImpl) DeleteAdminGroupByIds(ids ...string) (err error) {
	if len(ids) == 0 {
		return
	}
	var m models.AdminGroup
	err = r.Context.Db.Table(m.TableName()).
		Where("id IN (?) ", ids).
		Delete(&models.AdminGroup{}).
		Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"ids": ids,
			"err": err,
		}, "daoPermitDeleteAdminGroupByIds")
	}
	return
}
func (r *DaoPermitImpl) DeleteUserGroupPermitByGroupId(ids ...string) (err error) {
	if len(ids) == 0 {
		return
	}
	var m models.AdminUserGroupPermit
	if err = r.Context.Db.Table(m.TableName()).
		Where("group_id IN (?) ", ids).
		Delete(&models.AdminGroup{}).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"ids": ids,
			"err": err,
		}, "daoPermitDeleteUserGroupPermitByGroupId")
		return
	}
	return
}
func (r *DaoPermitImpl) DeleteUserGroupPermit(pathType string, menu_id ...int) (err error) {
	if len(menu_id) == 0 || pathType == "" {
		return
	}
	var m models.AdminUserGroupPermit
	if err = r.Context.Db.Table(m.TableName()).Unscoped().
		Where("menu_id IN (?) AND path_type= ?", menu_id, pathType).
		Delete(&models.AdminGroup{}).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"menu_id":  menu_id,
			"pathType": pathType,
			"err":      err,
		}, "daoPermitDeleteUserGroupPermit")
		return
	}
	return
}
func (r *DaoPermitImpl) FetchByName(name string) (res models.AdminGroup, err error) {
	var m models.AdminGroup
	err1 := r.Context.Db.Table(m.TableName()).
		Where("name = ?", name).
		Find(&res).
		Error
	if gorm.IsRecordNotFoundError(err1) {
		res = models.AdminGroup{}
	} else {
		err = err1
		r.Context.Error(map[string]interface{}{
			"name": name,
			"err":  err,
		}, "daoPermitFetchByName")
	}
	return
}
func (r *DaoPermitImpl) InsertAdminGroup(group *models.AdminGroup) (err error) {
	var m models.AdminGroup
	err = r.Context.Db.Table(m.TableName()).
		Create(group).
		Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"group": group,
			"err":   err,
		}, "daoPermitInsertAdminGroup")
	}
	return
}
func (r *DaoPermitImpl) UpdateAdminGroup(group *models.AdminGroup) (err error) {
	if err = r.Context.Db.Model(group).Where("id=?", group.Id).
		Update(group).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"group": group,
			"err":   err,
		}, "daoPermitUpdateAdminGroup")
	}
	return
}
func (r *DaoPermitImpl) GetAdminGroupByIds(gIds ...int) (res []models.AdminGroup, err error) {
	if len(gIds) == 0 {
		return
	}
	var m models.AdminGroup
	err = r.Context.Db.Table(m.TableName()).
		Where("id IN (?)", gIds).
		Find(&res).
		Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"gIds": gIds,
			"err":  err,
		}, "daoPermitGetAdminGroupByIds")
	}
	return
}
func (r *DaoPermitImpl) GetUserGroupByUIds(uIds ...string) (res []models.AdminUserGroup, err error) {
	if len(uIds) == 0 {
		return
	}
	var m models.AdminUserGroup
	err = r.Context.Db.Table(m.TableName()).
		Where(" user_hid IN (?)", uIds).
		Find(&res).
		Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"uIds": uIds,
			"err":  err,
		}, "daoPermitGetUserGroupByUIds")
	}
	return
}
func (r *DaoPermitImpl) UpdateMenuByCondition(condition interface{}, data map[string]interface{}) (err error) {
	var m models.AdminMenu
	if err = r.Context.Db.
		Table(m.TableName()).
		Where(condition).
		Update(data).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"data":      data,
			"err":       err,
		}, "daoPermitUpdateMenuByCondition")
		return
	}
	return
}
func (r *DaoPermitImpl) Save(id int, data map[string]interface{}) (err error) {
	if id == 0 {
		return
	}
	if err = r.Context.Db.
		Model(&models.AdminMenu{}).
		Where("id=?", id).
		Update(data).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"id":   id,
			"data": data,
			"err":  err,
		}, "daoPermitSave")
	}
	return
}
func (r *DaoPermitImpl) DeleteMenuByIds(ids ...string) (err error) {
	if len(ids) == 0 {
		return
	}
	var m models.AdminMenu
	if err = r.Context.Db.
		Table(m.TableName()).Unscoped().
		Where("id in (?)", ids).
		Delete(&models.AdminMenu{}).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"ids": ids,
			"err": err,
		}, "daoPermitDeleteByIds")
	}
	return
}

func (r *DaoPermitImpl) GetByCondition(condition map[string]interface{}, orderBy []wrappers.DaoOrderBy, limit int) (res []models.AdminMenu, err error) {
	if len(condition) == 0 {
		return
	}
	var m models.AdminMenu
	db := r.Context.Db.
		Table(m.TableName()).
		Where(condition)
	var orderString string
	for _, item := range orderBy {
		orderString += fmt.Sprintf("%s %s", item.Column, item.SortFormat)
	}
	if orderString != "" {
		db = db.Order(orderString)
	}
	if limit > 0 {
		db = db.Limit(limit)
	}
	if err = db.
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"err":       err,
		}, "daoPermitGetByCondition")
		return
	}
	return
}
func (r *DaoPermitImpl) Add(data *models.AdminMenu) (err error) {

	if err = r.Context.Db.Create(data).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"data": data,
			"err":  err,
		}, "daoPermitAdd")
	}
	return
}
func (r *DaoPermitImpl) GetAdminUserCount(db *gorm.DB, arg *wrappers.ArgAdminUser) (total int, dba *gorm.DB, err error) {
	var m models.AdminUser
	dba = r.Context.Db.Table(m.TableName()).Unscoped().
		Where("deleted_at IS NULL")
	if arg.Name != "" {
		dba = dba.Where("real_name LIKE ?", "%"+arg.Name+"%")
	}
	if arg.UserHId != "" {
		dba = dba.Where("user_hid=?", arg.UserHId)
	}
	if err = dba.Count(&total).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err,
		}, "daoPermitGetAdminUserCount")
	}

	return
}
func (r *DaoPermitImpl) GetAdminUserList(db *gorm.DB, arg *wrappers.ArgAdminUser, pager *response.Pager) (res []models.AdminUser, err error) {
	res = []models.AdminUser{}
	if err = db.Offset(pager.GetFromAndLimit()).Limit(arg.PageSize).Find(&res).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"arg":   arg,
			"pager": pager,
			"err":   err,
		}, "daoPermitGetAdminUserList")
	}
	return
}

func (r *DaoPermitImpl) GetMenuByCondition(condition interface{}) (res []models.AdminMenu, err error) {
	if condition == nil {
		return
	}
	var m models.AdminMenu
	if err = r.Context.Db.Table(m.TableName()).
		Where(condition).
		Limit(1000).
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"err":       err,
		}, "GetMenuByCondition")
		return
	}
	return
}
func (r *DaoPermitImpl) GetMenu(menuId ...int) (res []models.AdminMenu, err error) {
	var m models.AdminMenu
	l := len(menuId)
	res = make([]models.AdminMenu, 0, l)
	if l == 0 {
		return
	}
	if err = r.Context.Db.Table(m.TableName()).
		Where("id IN (?)", menuId).
		Limit(len(menuId)).
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"menuId": menuId,
			"err":    err,
		}, "daoPermitGetMenu")
		return
	}
	return
}
func (r *DaoPermitImpl) GetMenuByPermitKey(permitKey ...string) (res []models.AdminMenu, err error) {
	if len(permitKey) == 0 {
		return
	}
	var m models.AdminMenu
	if err = r.Context.Db.Table(m.TableName()).
		Where("permit_key IN(?)", permitKey).
		Find(&res).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"permitKey": permitKey,
			"err":       err.Error(),
		}, "daoPermitGetMenuByPermitKey")
		return
	}
	return
}
func (r *DaoPermitImpl) GetAdminMenuList(arg *wrappers.ArgAdminMenu) (res []models.AdminMenu, err error) {
	res = []models.AdminMenu{}
	var m models.AdminMenu
	dba := r.Context.Db.Table(m.TableName()).Unscoped()

	if arg.SystemId > 0 {
		if err = dba.Where("id=?", arg.SystemId).First(&m).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return
			}
			err = fmt.Errorf("你查看权限系统不存在或已删除")
			r.Context.Error(map[string]interface{}{
				"arg": arg,
				"err": err,
			}, "daoPermitGetAdminMenuList0")
			return
		}
		arg.Module = m.PermitKey
	}
	if arg.Label = strings.TrimSpace(arg.Label); arg.Label != "" {
		dba = dba.Where("label LIKE ?", "%"+arg.Label+"%")
	}
	if arg.ParentId > 0 {
		dba = dba.Where("parent_id = ?", arg.ParentId)
	}
	if arg.Module != "" {
		dba = dba.Where("module = ? OR parent_id=?", arg.Module, wrappers.DefaultPermitParentId)
	}
	if arg.AppName != "" {
		dba = dba.Where("app_name = ?", arg.AppName)
	}
	if arg.IsMenuShow != -1 {
		dba = dba.Where("is_menu_show = ?", arg.IsMenuShow)
	}
	if arg.IsDel == -1 || arg.IsDel == 0 {
		dba = dba.Where("deleted_at IS NULL")
	} else {
		dba = dba.Where("deleted_at NOT NULL")
	}
	if arg.Id != 0 {
		dba = dba.Where("id = ?", arg.Id)
	}
	if err = dba.Order("sort_value desc").Find(&res).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "daoPermitGetAdminMenuList2")
		return
	}
	return
}

func (r *DaoPermitImpl) GetAdminGroupCount(db *gorm.DB, arg *wrappers.ArgAdminGroup) (total int, dba *gorm.DB, err error) {
	var m models.AdminGroup
	dba = r.Context.Db.Table(m.TableName()).Unscoped().Where("deleted_at IS NULL")
	if arg.Name != "" {
		dba = dba.Where("name LIKE ?", "%"+arg.Name+"%")
	}
	var groupId int
	var err1 error
	if arg.GroupId != "" {
		// 如果参数不是整数，则直接返回没有数据
		if groupId, err1 = strconv.Atoi(arg.GroupId); err1 != nil {
			total = 0
			err = fmt.Errorf("您选择的ID格式不正确")
			return
		}
		if groupId != 0 {
			dba = dba.Where("id = ?", groupId)
		}
	}
	if err = dba.Count(&total).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err,
		}, "daoPermitGetAdminGroupCount")
		return
	}
	return
}
func (r *DaoPermitImpl) GetAdminGroupList(db *gorm.DB, arg *wrappers.ArgAdminGroup, pagerObject *response.Pager) (res []models.AdminGroup, err error) {
	res = []models.AdminGroup{}
	if err = db.Limit(pagerObject.PageSize).
		Offset(pagerObject.GetFromAndLimit()).
		Order("updated_at desc").
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"arg":         arg,
			"pagerObject": pagerObject,
			"err":         err,
		}, "daoPermitGetAdminGroupList")
		return
	}
	return
}
func (r *DaoPermitImpl) GetGroupByUserId(userId string) (res []wrappers.AdminGroupUserStruct, err error) {
	if userId == "" {
		res = []wrappers.AdminGroupUserStruct{}
		return
	}
	var m models.AdminUserGroup
	var m1 models.AdminGroup
	if err = r.Context.Db.Select("a.*,b.*").Unscoped().
		Table(m.TableName()).
		Joins(fmt.Sprintf("as a left join %s as b  ON  a.group_id=b.id ", m1.TableName())).
		Where(fmt.Sprintf("a.user_hid=? AND a.deleted_at  IS NULL AND  b.deleted_at IS NULL"), userId, ).
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"userId": userId,
			"err":    err,
		}, "daoPermitGetGroupByUserId")
		return
	}
	return
}
func (r *DaoPermitImpl) GetPermitMenuByIds(module []string, menuIds ...int) (res []models.AdminMenu, err error) {
	var m models.AdminMenu
	db := r.Context.Db.
		Table(m.TableName())
	// 兼容超级管理员和普通管理员
	if len(menuIds) != 0 {
		db = db.Where("id IN(?)", menuIds)
	}
	if len(module) > 0 {
		db = db.Where("module IN(?)", module)
	}
	if err = db.Order("sort_value desc").Limit(2000).Find(&res).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"menuIds": menuIds,
			"err":     err,
		}, "daoPermitGetPermitMenuByIds")
		return
	}
	return
}
func (r *DaoPermitImpl) GetMenuIdsByPermitByGroupIds(module string, pathType []string, groupIds ...int) (res []models.AdminUserGroupPermit, err error) {
	if len(groupIds) == 0 {
		return
	}
	var m models.AdminUserGroupPermit
	if err = r.Context.Db.
		Table(m.TableName()).
		Select("distinct `menu_id`,`group_id`,`id`").
		Where("module = ? AND path_type IN(?)  AND `group_id` in(?) ", module, pathType, groupIds).
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"groupIds": groupIds,
			"pathType": pathType,
			"err":      err,
		}, "daoPermitGetPermitMenuByIds")
		return
	}
	return
}
