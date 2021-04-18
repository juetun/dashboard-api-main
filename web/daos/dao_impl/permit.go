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
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type DaoPermit struct {
	base.ServiceDao
}

func NewDaoPermit(context ...*base.Context) (p *DaoPermit) {
	p = &DaoPermit{}
	p.SetContext(context...)
	p.Context.Db = base.GetDbClient(&base.GetDbClientData{
		Context:     p.Context,
		DbNameSpace: "admin",
	})
	return
}
func (r *DaoPermit) DeleteImportByIds(id ...int) (err error) {
	if len(id) == 0 {
		return
	}
	var m models.AdminImport
	if err = r.Context.Db.Table(m.TableName()).
		Where("id IN(?)", id).
		Delete(&models.AdminImport{}).Error; err != nil {
		return
	}

	var m1 models.AdminUserGroupPermit
	if err = r.Context.Db.Table(m1.TableName()).
		Where("menu_id IN(?) AND path_type=?", id, "api").
		Delete(&models.AdminImport{}).Error; err != nil {
		return
	}

	return
}
func (r *DaoPermit) GetImportMenuId(menuId int) (list []models.AdminImport, err error) {
	list = []models.AdminImport{}
	var m models.AdminImport
	if menuId == 0 {
		return
	}
	if err = r.Context.Db.Table(m.TableName()).Unscoped().
		Where("menu_id=?", menuId).
		Find(&list).Error; err != nil {
	}
	return
}

func (r *DaoPermit) CreateImport(data *models.AdminImport) (res bool, err error) {
	var m models.AdminImport
	if err = r.Context.Db.Table(m.TableName()).Create(data).Error; err != nil {
		return
	}

	return
}
func (r *DaoPermit) UpdateAdminImport(condition, data map[string]interface{}) (res bool, err error) {
	var m models.AdminImport
	if len(condition) == 0 {
		return
	}
	if err = r.Context.Db.Table(m.TableName()).Where(condition).Update(data).Error; err != nil {
		return
	}
	res = true
	return
}
func (r *DaoPermit) GetImportId(id int) (res models.AdminImport, err error) {
	var m models.AdminImport
	err1 := r.Context.Db.Table(m.TableName()).Where("id=?", id).Find(&res).Error
	if err1 != nil && gorm.IsRecordNotFoundError(err1) {
		err = err1
		return
	}
	return
}
func (r *DaoPermit) GetImportCount(arg *wrappers.ArgGetImport, count *int) (db *gorm.DB, err error) {
	if arg.MenuId == 0 {
		return
	}
	var m models.AdminImport
	db = r.Context.Db.Table(m.TableName()).Where("menu_id=?", arg.MenuId)
	err = db.Count(count).Error
	return
}
func (r *DaoPermit) GetImportList(db *gorm.DB, arg *wrappers.ArgGetImport) (res []models.AdminImport, err error) {
	err = db.Offset(arg.GetOffset()).
		Limit(arg.PageSize).
		Find(&res).Error
	return
}
func (r *DaoPermit) AdminUserGroupAdd(data []map[string]interface{}) (err error) {
	field := make([]string, 0, 10)
	dataMsg := make([]string, 0)
	dataTmp := make([]interface{}, 0)
	duplicate := make([]string, 0)
	dataMsgA := make([]string, 0)
	if len(data) == 0 {
		return fmt.Errorf("您没有为模板配置参数,请至少配置一个参数")
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
	return r.Context.Db.Exec(sql, dataTmp...).Error
}
func (r *DaoPermit) AdminUserGroupRelease(ids []string) (err error) {
	if len(ids) == 0 {
		return
	}
	var m models.AdminUserGroup
	err = r.Context.Db.Table(m.TableName()).
		Where("id IN (?) ", ids).
		Update(map[string]interface{}{
			"is_del": 1,
		}).
		Error
	return
}
func (r *DaoPermit) AdminUserAdd(arg *models.AdminUser) (err error) {

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
	return
}
func (r *DaoPermit) AdminUserDelete(ids []string) (err error) {
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
	return
}
func (r *DaoPermit) DeleteAdminGroupByIds(ids []string) (err error) {
	if len(ids) == 0 {
		return
	}
	var m models.AdminGroup
	err = r.Context.Db.Table(m.TableName()).
		Where("id IN (?) ", ids).
		Delete(&models.AdminGroup{}).
		Error
	return
}
func (r *DaoPermit) FetchByName(name string) (res models.AdminGroup, err error) {
	var m models.AdminGroup
	err1 := r.Context.Db.Table(m.TableName()).
		Where("name = ?", name).
		Find(&res).
		Error
	if gorm.IsRecordNotFoundError(err1) {
		res = models.AdminGroup{}
	} else {
		err = err1
	}
	return
}
func (r *DaoPermit) InsertAdminGroup(group *models.AdminGroup) (err error) {
	var m models.AdminGroup
	err = r.Context.Db.Table(m.TableName()).
		Create(group).
		Error
	return
}
func (r *DaoPermit) UpdateAdminGroup(group *models.AdminGroup) (err error) {
	var m models.AdminGroup
	err = r.Context.Db.Table(m.TableName()).Where("id=?", group.Id).
		Update(group).
		Error
	return
}
func (r *DaoPermit) GetAdminGroupByIds(gIds []int) (res []models.AdminGroup, err error) {
	if len(gIds) == 0 {
		return
	}
	var m models.AdminGroup
	err = r.Context.Db.Table(m.TableName()).
		Where("id IN (?)", gIds).
		Find(&res).
		Error
	return
}
func (r *DaoPermit) GetUserGroupByUIds(uIds []string) (res []models.AdminUserGroup, err error) {
	if len(uIds) == 0 {
		return
	}
	var m models.AdminUserGroup
	err = r.Context.Db.Table(m.TableName()).
		Where(" user_hid IN (?)", uIds).
		Find(&res).
		Error
	return
}
func (r *DaoPermit) Save(id int, data *models.AdminMenu) (err error) {
	if id == 0 {
		return
	}
	var m models.AdminMenu
	err = r.Context.Db.
		Table(m.TableName()).
		Where("id=?", id).
		Update(data).
		Error
	return
}
func (r *DaoPermit) DeleteByIds(ids []string) (err error) {
	if len(ids) == 0 {
		return
	}
	var m models.AdminMenu
	err = r.Context.Db.
		Table(m.TableName()).
		Where("id in (?)", ids).
		Update(map[string]interface{}{"is_del": 1, "updated_at": time.Now().Format("2006-01-02 15:04:05")}).
		Error
	return
}
func (r *DaoPermit) GetByCondition(condition map[string]interface{}) (res []models.AdminMenu, err error) {
	if len(condition) == 0 {
		return
	}
	var m models.AdminMenu
	err = r.Context.Db.
		Table(m.TableName()).
		Where(condition).
		Find(&res).
		Error
	return
}
func (r *DaoPermit) Add(data *models.AdminMenu) (err error) {
	err = r.Context.Db.Create(data).Error
	return
}
func (r *DaoPermit) GetAdminUserCount(db *gorm.DB, arg *wrappers.ArgAdminUser) (total int, dba *gorm.DB, err error) {
	var m models.AdminUser
	dba = r.Context.Db.Table(m.TableName())
	if arg.Name != "" {
		dba = dba.Where("real_name LIKE ?", "%"+arg.Name+"%")
	}
	if arg.UserHId != "" {
		dba = dba.Where("user_hid=?", arg.UserHId)
	}
	err = dba.Count(&total).Error
	return
}
func (r *DaoPermit) GetAdminUserList(db *gorm.DB, arg *wrappers.ArgAdminUser, pager *response.Pager) (res []models.AdminUser, err error) {
	res = []models.AdminUser{}
	err = db.Offset(pager.GetFromAndLimit()).Limit(arg.PageSize).Find(&res).Error
	return
}

func (r *DaoPermit) GetMenu(menuId int) (res models.AdminMenu, err error) {
	var m models.AdminMenu
	var l []models.AdminMenu
	if err = r.Context.Db.Table(m.TableName()).
		Where("id = ?", menuId).
		Limit(1).
		Find(&l).
		Error; err != nil {
		return
	}
	if len(l) > 0 {
		res = l[0]
	}
	return
}
func (r *DaoPermit) GetAdminMenuList(arg *wrappers.ArgAdminMenu) (res []models.AdminMenu, err error) {
	res = []models.AdminMenu{}
	var m models.AdminMenu
	dba := r.Context.Db.Table(m.TableName()).Where("is_del=?", 0)

	if arg.SystemId > 0 {
		if err = dba.Where("id=?", arg.SystemId).First(&m).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return
			}
			err = fmt.Errorf("你查看权限系统不存在或已删除")
		}
		arg.Module = m.PermitKey
	}
	if arg.Label = strings.TrimSpace(arg.Label); arg.Label != "" {
		dba = dba.Where("label LIKE ?", "%"+arg.Label+"%")
	}
	if arg.ParentId != -1 {
		dba = dba.Where("parent_id = ?", arg.ParentId)
	}
	if arg.Module != "" {
		dba = dba.Where("module = ? OR parent_id=?", arg.Module, 0)
	}
	if arg.AppName != "" {
		dba = dba.Where("app_name = ?", arg.AppName)
	}
	if arg.IsMenuShow != -1 {
		dba = dba.Where("is_menu_show = ?", arg.IsMenuShow)
	}
	if arg.IsDel == -1 || arg.IsDel == 0 {
		dba = dba.Where("is_del = ?", 0)
	} else {
		dba = dba.Where("is_del = ?", arg.IsDel)
	}
	if arg.Id != 0 {
		dba = dba.Where("id = ?", arg.Id)
	}
	err = dba.Order("sort_value desc").Find(&res).Error
	return
}

func (r *DaoPermit) GetAdminGroupCount(db *gorm.DB, arg *wrappers.ArgAdminGroup) (total int, dba *gorm.DB, err error) {
	var m models.AdminGroup
	dba = r.Context.Db.Table(m.TableName())
	if arg.Name != "" {
		dba = dba.Where("name LIKE ?", "%"+arg.Name+"%")
	}
	if arg.GroupId != 0 {
		dba = dba.Where("id=?", arg.GroupId)
	}
	err = dba.Count(&total).Error
	return
}
func (r *DaoPermit) GetAdminGroupList(db *gorm.DB, arg *wrappers.ArgAdminGroup, pagerObject *response.Pager) (res []models.AdminGroup, err error) {
	res = []models.AdminGroup{}
	err = db.Limit(pagerObject.PageSize).
		Offset(pagerObject.GetFromAndLimit()).
		Order("updated_at desc").
		Find(&res).
		Error
	return
}
func (r *DaoPermit) GetGroupByUserId(userId string) (res []wrappers.AdminGroupUserStruct, err error) {
	if userId == "" {
		res = []wrappers.AdminGroupUserStruct{}
		return
	}
	var m models.AdminUserGroup
	var m1 models.AdminGroup
	err = r.Context.Db.Select("a.*,b.*").Unscoped().
		Table(m.TableName()).
		Joins(fmt.Sprintf("as a left join %s as b  ON  a.group_id=b.id ", m1.TableName())).
		Where(fmt.Sprintf("a.user_hid=? AND a.deleted_at  IS NULL AND  b.deleted_at IS NULL"), userId, ).
		Find(&res).
		Error
	return
}
func (r *DaoPermit) GetPermitMenuByIds(menuIds ...int) (res []models.AdminMenu, err error) {
	var m models.AdminMenu
	db := r.Context.Db.
		Table(m.TableName()).Where("is_del=?", 0)
	// 兼容超级管理员和普通管理员
	if len(menuIds) != 0 {
		db = db.Where("id IN(?)", menuIds)
	}
	err = db.Order("sort_value desc").Find(&res).Error
	return
}
func (r *DaoPermit) GetMenuIdsByPermitByGroupIds(pathType []string, groupIds ...int) (res []models.AdminUserGroupPermit, err error) {
	if len(groupIds) == 0 {
		return
	}
	var m models.AdminUserGroupPermit
	err = r.Context.Db.
		Table(m.TableName()).
		Select("distinct `menu_id`,`group_id`,`id`,`is_del`").
		Where("path_type IN(?)  AND `group_id` in(?) AND `is_del`=?  ", pathType, groupIds, 0).
		Find(&res).
		Error
	return
}
