// Package daos /**
package daos

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"gorm.io/gorm"
)

type DaoPermitUser interface {
	UpdateDataByUserHIds(data map[string]interface{}, userHIds ...string) (err error)

	AdminUserAdd(dataUser []base.ModelBase) (err error)

	DeleteAdminUser(userHIds ...string) (err error)

	DeleteAdminUserGroup(userHIds ...string) (err error)

	GetAdminUserCount(db *gorm.DB, arg *wrappers.ArgAdminUser) (total int64, dba *gorm.DB, err error)

	GetAdminUserList(db *gorm.DB, arg *wrappers.ArgAdminUser, pager *response.Pager) (res []models.AdminUser, err error)

	GetGroupByUserId(userId string) (res []wrappers.AdminGroupUserStruct, err error)
}
