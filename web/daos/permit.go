/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 11:19 下午
 */
package daos

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/models"
)

type DaoPermit struct {
	base.ServiceDao
}

func NewDaoPermit(context ...*base.Context) (p *DaoPermit) {
	p = &DaoPermit{}
	p.SetContext(context...)
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
func (r *DaoPermit) GetMenuIdsByPermitByGroupIds(groupIds ...int) (res []models.AdminUserGroupPermit, err error) {
	if len(groupIds) == 0 {
		return
	}
	var m models.AdminUserGroupPermit
	err = r.Context.Db.
		Table(m.TableName()).
		Select("distinct `menu_id`,`group_id`,`id`,`is_del`").
		Where("`group_id` in(?) AND `is_del`=?", groupIds, 0).
		Find(&res).
		Error
	return
}
