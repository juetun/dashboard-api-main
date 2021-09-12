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
	"gorm.io/gorm"
)

type DaoPermit interface {
	GetImportListData(db *gorm.DB, arg *wrappers.ArgImportList, pager *response.Pager) (res []models.AdminImport, err error)

	GetImportListCount(db *gorm.DB, arg *wrappers.ArgImportList) (totalCount int, dba *gorm.DB, err error)

	GetPermitImportByModule(arg *wrappers.ArgPermitMenu) (res []wrappers.Op, err error)
	GetImportMenuId(menuId ...int) (list []models.AdminImport, err error)
	GetImportByCondition(condition map[string]interface{}) (list []models.AdminImport, err error)
	CreateImport(data *models.AdminImport) (res bool, err error)
	UpdateAdminImport(condition, data map[string]interface{}) (res bool, err error)
	GetAdminImportById(id ...int) (res []models.AdminImport, err error)
	BatchGroupPermit(tableName string, list []models.AdminUserGroupPermit) (err error)
	DeleteGroupPermit(groupId int, pathType string, menuId ...int) (err error)
	// DeleteByMenuIds(pageMenuIds ...string) (err error)
	DeleteGroupPermitByGroupId(groupId int) (err error)
	DeleteGroupPermitByMenuIds(groupId int, module string, pageMenuId, apiMenuId []int) (err error)
	GetDefaultOpenImportByMenuIds(menuId ...int) (res []models.AdminImport, err error)
	GetDefaultImportByMenuIds(pageType, module string, menuId ...int) (res []models.AdminImport, err error)
	GetImportId(id int) (res models.AdminImport, err error)
	MenuImportCount(arg *wrappers.ArgMenuImport, count *int) (db *gorm.DB, err error)
	GetImportCount(arg *wrappers.ArgGetImport, count *int) (db *gorm.DB, err error)
	GetSelectImportByImportId(groupId int, importId ...int) (res []models.AdminUserGroupPermit, err error)
	GetImportList(db *gorm.DB, arg *wrappers.ArgGetImport) (res []models.AdminImport, err error)
	MenuImportList(db *gorm.DB, arg *wrappers.ArgMenuImport) (res []wrappers.ResultMenuImportItem, err error)
	AdminUserGroupAdd(data []map[string]interface{}) (err error)
	AdminUserGroupRelease(ids ...string) (err error)
	AdminUserAdd(arg *models.AdminUser) (err error)
	DeleteAdminUser(ids []string) (err error)
	DeleteAdminGroupByIds(ids ...string) (err error)
	FetchByName(name string) (res models.AdminGroup, err error)
	InsertAdminGroup(group *models.AdminGroup) (err error)
	UpdateAdminGroup(group *models.AdminGroup) (err error)
	GetAdminGroupByIds(gIds ...int) (res []models.AdminGroup, err error)
	GetUserGroupByUIds(uIds ...string) (res []models.AdminUserGroup, err error)
	UpdateMenuByCondition(condition interface{}, data map[string]interface{}) (err error)
	Save(id int, data map[string]interface{}) (err error)
	DeleteMenuByIds(ids ...string) (err error)
	GetByCondition(condition map[string]interface{}, orderBy []wrappers.DaoOrderBy, limit int) (res []models.AdminMenu, err error)
	Add(data *models.AdminMenu) (err error)
	GetAdminUserCount(db *gorm.DB, arg *wrappers.ArgAdminUser) (total int, dba *gorm.DB, err error)
	GetAdminUserList(db *gorm.DB, arg *wrappers.ArgAdminUser, pager *response.Pager) (res []models.AdminUser, err error)
	GetMenuByCondition(condition interface{}) (res []models.AdminMenu, err error)
	GetMenu(menuId ...int) (res []models.AdminMenu, err error)

	GetMenuByPermitKey(module string, permitKey ...string) (res []models.AdminMenu, err error)

	GetAdminMenuList(arg *wrappers.ArgAdminMenu) (res []models.AdminMenu, err error)
	GetAdminGroupCount(db *gorm.DB, arg *wrappers.ArgAdminGroup) (total int, dba *gorm.DB, err error)
	GetAdminGroupList(db *gorm.DB, arg *wrappers.ArgAdminGroup, pagerObject *response.Pager) (res []models.AdminGroup, err error)
	GetGroupByUserId(userId string) (res []wrappers.AdminGroupUserStruct, err error)
	GetPermitMenuByIds(module []string, menuIds ...int) (res []models.AdminMenu, err error)
	GetMenuIdsByPermitByGroupIds(module string, pathType []string, groupIds ...int) (res []models.AdminUserGroupPermit, err error)
}
