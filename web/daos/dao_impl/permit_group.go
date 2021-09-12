/**
* @Author:changjiang
* @Description:
* @File:permit_group
* @Version: 1.0.0
* @Date 2021/9/12 11:40 上午
 */
package dao_impl

import (
	"time"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
)

type DaoPermitGroupImpl struct {
	base.ServiceBase
}

func (r *DaoPermitGroupImpl) DeleteUserGroupPermitByGroupId(ids ...string) (err error) {
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

func (r *DaoPermitGroupImpl) DeleteUserGroupByUserId(userId ...string) (err error) {
	if len(userId) == 0 {
		return
	}
	var m models.AdminUserGroup

	if err = r.Context.Db.Table(m.TableName()).
		Where("user_hid IN (?) ", userId).
		Updates(map[string]interface{}{
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

func (r *DaoPermitGroupImpl) DeleteUserGroupPermit(pathType string, menuId ...int) (err error) {
	if len(menuId) == 0 || pathType == "" {
		return
	}
	var m models.AdminUserGroupPermit
	if err = r.Context.Db.Table(m.TableName()).Unscoped().
		Where("menu_id IN (?) AND path_type= ?", menuId, pathType).
		Delete(&models.AdminGroup{}).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"menu_id":  menuId,
			"pathType": pathType,
			"err":      err,
		}, "daoPermitDeleteUserGroupPermit")
		return
	}
	return
}

func NewDaoPermitGroupImpl(ctx ...*base.Context) (res daos.DaoPermitGroup) {
	p := &DaoPermitGroupImpl{}
	p.SetContext(ctx...)
	return p
}
