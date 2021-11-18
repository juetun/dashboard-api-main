// Package daos /**
package daos

import (
	"github.com/juetun/dashboard-api-main/web/models"
)

type DaoPermitGroup interface {
	DeleteUserGroupByUserId(userId ...string) (err error)

	DeleteUserGroupPermitByGroupId(ids ...string) (err error)

	DeleteUserGroupPermit(pathType string, menuId ...int) (err error)

	// UpdateDaoPermitUserGroupByGroupId 更新用户权限组关系
	UpdateDaoPermitUserGroupByGroupId(groupId int64, data map[string]interface{}) (err error)

	// GetPermitGroupByUid 根据用户ID获取用户所属管理员组
	GetPermitGroupByUid(hid string, refreshCache ...bool) (res []models.AdminUserGroup, err error)
}
