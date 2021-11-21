// Package daos /**
package daos

import (
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_intranet"
)

type DaoPermitGroup interface {
	DeleteUserGroupByUserId(userId ...string) (err error)

	DeleteUserGroupPermitByGroupId(groupIds ...int64) (err error)

	DeleteUserGroupPermit(menuId ...int64) (err error)

	DeleteAdminGroupByIds(ids ...int64) (err error)

	// UpdateDaoPermitUserGroupByGroupId 更新用户权限组关系
	UpdateDaoPermitUserGroupByGroupId(groupId int64, data map[string]interface{}) (err error)

	// GetPermitGroupByUid 根据用户ID获取用户所属管理员组
	GetPermitGroupByUid(hid string, refreshCache ...bool) (res []models.AdminUserGroup, err error)

	// GetGroupAppPermitImport 从缓存中获取指定组的app权限列表
	GetGroupAppPermitImport(groupId int64, appName string, refreshCache ...bool) (res []wrapper_intranet.AdminUserGroupPermit, err error)
}
