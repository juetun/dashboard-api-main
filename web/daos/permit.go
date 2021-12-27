// Package daos
/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2021/5/31 11:20 下午
 */
package daos

import (
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
	"gorm.io/gorm"
)

type DaoPermit interface {
	GetImportListData(db *gorm.DB, arg *wrappers.ArgImportList, pager *response.Pager) (res []models.AdminImport, err error)

	GetImportListCount(db *gorm.DB, arg *wrappers.ArgImportList) (totalCount int64, dba *gorm.DB, err error)

	GetPermitImportByModule(arg *wrappers.ArgPermitMenu) (res []wrappers.Op, err error)

	GetImportMenuId(menuId ...int64) (list []models.AdminImport, err error)

	CreateImport(data *models.AdminImport) (res bool, err error)

	UpdateAdminImport(condition, data map[string]interface{}) (res bool, err error)

	GetAdminImportById(id ...int64) (res []models.AdminImport, err error)

	GetDefaultImportByMenuIds(pageType, module string, menuId ...int64) (res []models.AdminImport, err error)

	GetImportId(id int) (res models.AdminImport, err error)

	MenuImportCount(arg *wrapper_admin.ArgMenuImport, count *int64) (db *gorm.DB, err error)

	GetImportCount(arg *wrappers.ArgGetImport, count *int64) (db *gorm.DB, err error)

	GetImportList(db *gorm.DB, arg *wrappers.ArgGetImport) (res []models.AdminImport, err error)

	MenuImportList(db *gorm.DB, arg *wrapper_admin.ArgMenuImport) (res []wrappers.ResultMenuImportItem, err error)

	// AdminUserGroupAdd(data []map[string]interface{}) (err error)

	AdminUserGroupRelease(ids ...string) (err error)

	FetchByName(name string) (res models.AdminGroup, err error)

	InsertAdminGroup(group *models.AdminGroup) (err error)

	UpdateAdminGroup(data, where map[string]interface{}) (err error)

	GetAdminGroupByIds(gIds ...int64) (res []models.AdminGroup, err error)

	GetUserGroupByUIds(uIds ...string) (res []models.AdminUserGroup, err error)

	DeleteMenuByIds(ids ...string) (err error)

	GetAdminMenuList(arg *wrappers.ArgAdminMenu) (res []models.AdminMenu, err error)

	GetAdminGroupCount(db *gorm.DB, arg *wrappers.ArgAdminGroup) (total int64, dba *gorm.DB, err error)

	GetAdminGroupList(db *gorm.DB, arg *wrappers.ArgAdminGroup, pagerObject *response.Pager) (res []models.AdminGroup, err error)

	GetPermitMenuByIds(module []string, menuIds ...int64) (res []*models.AdminMenu, err error)
}
