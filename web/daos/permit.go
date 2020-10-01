/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 11:19 下午
 */
package daos

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/juetun/base-wrapper/lib/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
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

func (r *DaoPermit) Save(id int, data *models.AdminMenu) (err error) {
	if id == 0 {
		return
	}
	var m models.AdminMenu
	err = r.Context.Db.Table(m.TableName()).Where("id=?", id).Save(data).Error
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
		Update(map[string]interface{}{"is_del":1,"updated_at":time.Now().Format("2006-01-02 15:04:05")}).
		Error
	return
}
func (r *DaoPermit) Add(data *models.AdminMenu) (err error) {

	err = r.Context.Db.Create(data).Error
	return
}
func (r *DaoPermit) GetAdminUserCount(db *gorm.DB, arg *pojos.ArgAdminUser) (total int, dba *gorm.DB, err error) {
	var m models.AdminUser
	dba = r.Context.Db.Table(m.TableName()).Unscoped().Where("deleted_at = ?", "1970-01-02 08:00:00")
	if arg.Name != "" {
		dba = dba.Where("real_name LIKE ?", "%"+arg.Name+"%")
	}
	if arg.UserHId != "" {
		dba = dba.Where("user_hid=?", arg.UserHId)
	}
	err = dba.Count(&total).Error
	return
}
func (r *DaoPermit) GetAdminUserList(db *gorm.DB, arg *pojos.ArgAdminUser) (res []models.AdminUser, err error) {
	res = []models.AdminUser{}
	err = db.Find(&res).Error
	return
}

func (r *DaoPermit) GetAdminMenuList(arg *pojos.ArgAdminMenu) (res []models.AdminMenu, err error) {
	res = []models.AdminMenu{}
	var m models.AdminMenu
	dba := r.Context.Db.Table(m.TableName()).Where("is_del=?", 0)
	if arg.Label != "" {
		dba = dba.Where("label LIKE ?", "%"+arg.Name+"%")
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
func (r *DaoPermit) GetAdminGroupList(db *gorm.DB, arg *pojos.ArgAdminGroup) (res []models.AdminGroup, err error) {
	res = []models.AdminGroup{}
	err = db.Find(&res).Error
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
