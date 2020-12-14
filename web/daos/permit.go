/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 11:19 下午
 */
package daos

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/juetun/base-wrapper/lib/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/pojos"
)

type DaoPermit struct {
	base.ServiceDao
}

func NewDaoPermit(context ...*base.Context) (p *DaoPermit) {
	p = &DaoPermit{}
	p.SetContext(context...)
	p.Context.Db = app_obj.GetDbClient("admin")
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
		Update(map[string]interface{}{
			"is_del": 1,
		}).
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
		Where("is_del=0 AND  id IN (?)", gIds).
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
		Where("is_del=0 AND  user_hid IN (?)", uIds).
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
func (r *DaoPermit) Add(data *models.AdminMenu) (err error) {

	err = r.Context.Db.Create(data).Error
	return
}
func (r *DaoPermit) GetAdminUserCount(db *gorm.DB, arg *pojos.ArgAdminUser) (total int, dba *gorm.DB, err error) {
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
func (r *DaoPermit) GetAdminUserList(db *gorm.DB, arg *pojos.ArgAdminUser, pager *response.Pager) (res []models.AdminUser, err error) {
	res = []models.AdminUser{}
	err = db.Offset(pager.GetFromAndLimit()).Limit(arg.PageSize).Find(&res).Error
	return
}

func (r *DaoPermit) GetAdminMenuList(arg *pojos.ArgAdminMenu) (res []models.AdminMenu, err error) {
	res = []models.AdminMenu{}
	var m models.AdminMenu
	dba := r.Context.Db.Table(m.TableName()).Where("is_del=?", 0)
	if arg.Label != "" {
		dba = dba.Where("label LIKE ?", "%"+arg.Label+"%")
	}
	if arg.ParentId != -1 {
		dba = dba.Where("parent_id = ?", arg.ParentId)
	}
	if arg.AppName != "" {
		dba = dba.Where("app_name = ?", arg.AppName)
	}
	if arg.IsMenuShow != -1 {
		dba = dba.Where("is_menu_show = ?", arg.IsMenuShow)
	}
	if arg.IsDel != -1 {
		dba = dba.Where("is_del = ?", arg.IsDel)
	}
	err = dba.Find(&res).Error
	return
}

func (r *DaoPermit) GetAdminGroupCount(db *gorm.DB, arg *pojos.ArgAdminGroup) (total int, dba *gorm.DB, err error) {
	var m models.AdminGroup
	dba = r.Context.Db.Table(m.TableName()).Where("is_del=?", 0)
	if arg.Name != "" {
		dba = dba.Where("name LIKE ?", "%"+arg.Name+"%")
	}
	if arg.GroupId != 0 {
		dba = dba.Where("id=?", arg.GroupId)
	}
	err = dba.Count(&total).Error
	return
}
func (r *DaoPermit) GetAdminGroupList(db *gorm.DB, arg *pojos.ArgAdminGroup, pagerObject *response.Pager) (res []models.AdminGroup, err error) {
	res = []models.AdminGroup{}
	err = db.Limit(pagerObject.PageSize).
		Offset(pagerObject.GetFromAndLimit()).
		Find(&res).
		Error
	return
}
func (r *DaoPermit) GetGroupByUserId(userId string) (res []models.AdminUserGroup, err error) {
	if userId == "" {
		res = []models.AdminUserGroup{}
		return
	}
	var m models.AdminUserGroup
	err = r.Context.Db.
		Table(m.TableName()).
		Where("user_hid=? AND is_del=?", userId, 0).
		Find(&res).
		Error
	return
}
func (r *DaoPermit) GetPermitMenuByIds(menuIds ...int) (res []models.AdminMenu, err error) {
	if len(menuIds) == 0 {
		return
	}
	var m models.AdminMenu
	err = r.Context.Db.
		Table(m.TableName()).
		Where("id IN(?) AND is_del=?", menuIds, 0).
		Find(&res).
		Error
	return
}
func (r *DaoPermit) GetMenuIdsByPermitByGroupIds(pathType string, groupIds ...int) (res []models.AdminUserGroupPermit, err error) {
	if len(groupIds) == 0 {
		return
	}
	var m models.AdminUserGroupPermit
	err = r.Context.Db.
		Table(m.TableName()).
		Select("distinct `menu_id`,`group_id`,`id`,`is_del`").
		Where("path_type = ? AND `group_id` in(?) AND `is_del`=?  ", pathType, groupIds, 0).
		Find(&res).
		Error
	return
}
