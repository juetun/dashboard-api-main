// Package daos /**
package daos

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_intranet"
)

type DaoPermitGroup interface {
	base.DaoBatchAdd

	// GetGroupUserCount 获取管理员组的管理员数量
	GetGroupUserCount(groupIds ...int64) (groupIdUserCountMap map[int64]int, err error)

	AdminUserGroupAdd([]base.ModelBase) (err error)

	DeleteUserGroupByUserId(userId ...int64) (err error)

	DeleteUserGroupPermitByGroupId(groupIds ...int64) (err error)

	DeleteUserGroupPermit(menuId ...int64) (err error)

	DeleteAdminGroupByIds(ids ...int64) (err error)

	// UpdateDaoPermitUserGroupByGroupId 更新用户权限组关系
	UpdateDaoPermitUserGroupByGroupId(groupId int64, data map[string]interface{}) (err error)

	// GetPermitGroupByUid 根据用户ID获取用户所属管理员组
	GetPermitGroupByUid(userHid int64, refreshCache ...bool) (res []models.AdminUserGroup, err error)

	// GetGroupAppPermitImport 从缓存中获取指定组的app权限列表
	GetGroupAppPermitImport(groupId int64, appName string, refreshCache ...bool) (res []wrapper_intranet.AdminUserGroupPermit, err error)

	GetGroupByIds(groupIds ...int64) (res []*models.AdminGroup, err error)

	GetGroupUserByIds(groupIds ...int64) (res []*models.AdminUserGroup, err error)

	// GetGroupSystemByGroupId 获取用户组能够使用的系统列表
	GetGroupSystemByGroupId(module string, groupId ...int64) (res []*models.AdminMenu, err error)

	GetAdminMenuByGidAndMenuIds(menuId, groupId []int64) (list []*models.AdminUserGroupMenu, err error)

	GetGroupAdminMenuByGroupIds(module string, groupIds ...int64) ([]*models.AdminMenu, error)
}
