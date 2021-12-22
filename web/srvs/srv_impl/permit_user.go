// Package srv_impl
package srv_impl

import (
	"fmt"
	"strings"
	"time"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/app_param"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
	"gorm.io/gorm"
)

type SrvPermitUserImpl struct {
	base.ServiceDao
}

func (r *SrvPermitUserImpl) AdminUserUpdateWithColumn(arg *wrapper_admin.ArgAdminUserUpdateWithColumn) (res *wrapper_admin.ResultAdminUserUpdateWithColumn, err error) {
	res = &wrapper_admin.ResultAdminUserUpdateWithColumn{}
	dao := dao_impl.NewDaoPermitUser(r.Context)
	switch arg.Column {
	case "can_not_use":
		if err = r.adminUserUpdateWithColumnCanNotUse(dao, arg); err != nil {
			return
		}
	}
	if err = dao.UpdateDataByUserHIds(map[string]interface{}{arg.Column: arg.Value}, arg.UserHidVals...); err != nil {
		return
	}
	res.Result = true
	return
}

func (r *SrvPermitUserImpl) adminUserUpdateWithColumnCanNotUse(dao daos.DaoPermitUser, arg *wrapper_admin.ArgAdminUserUpdateWithColumn) (err error) {
	switch arg.Value {
	case fmt.Sprintf("%d", models.AdminUserCanNotUseYes): //
		if err = dao.DeleteAdminUserGroup(arg.UserHidVals...); err != nil {
			return
		}
	}
	return
}

func (r *SrvPermitUserImpl) AdminUserEdit(arg *wrappers.ArgAdminUserAdd) (res wrappers.ResultAdminUserAdd, err error) {
	res = wrappers.ResultAdminUserAdd{}
	if arg.UserHid == "" {
		err = fmt.Errorf("您没有选择要编辑的用户")
		return
	}

	userHid := strings.TrimSpace(arg.UserHid)

	var user *app_param.ResultUserItem
	if user, err = NewUserService(r.Context).
		GetUserById(userHid); err != nil {
		return
	}
	if user.UserHid == "" {
		err = fmt.Errorf("您要编辑的用户信息不存在")
		return
	}
	t := time.Now()
	if arg.RealName == "" {
		arg.RealName = user.RealName
	}
	if arg.Mobile == "" {
		arg.Mobile = user.Mobile
	}
	daoPermitUser := dao_impl.NewDaoPermitUser(r.Context)
	if arg.Id == 0 { // 如果是添加数据

		adminUser := &models.AdminUser{
			UserHid:   arg.UserHid,
			RealName:  arg.RealName,
			Mobile:    arg.Mobile,
			CreatedAt: t,
			UpdatedAt: t,
		}

		if err = daoPermitUser.AdminUserAdd([]base.ModelBase{adminUser}); err != nil {
			return
		}
		res.Result = true
		return
	}

	data := map[string]interface{}{
		"id":         arg.Id,
		"user_hid":   arg.UserHid,
		"real_name":  arg.RealName,
		"mobile":     arg.Mobile,
		"flag_admin": arg.FlagAdmin,
		"updated_at": t.Format(utils.DateTimeGeneral),
		"deleted_at": nil,
	}
	if err = daoPermitUser.UpdateDataByUserHIds(data, arg.UserHid); err != nil {
		return
	}
	res.Result = true
	return
}

func (r *SrvPermitUserImpl) AdminUserDelete(arg *wrappers.ArgAdminUserDelete) (res wrappers.ResultAdminUserDelete, err error) {
	res = wrappers.ResultAdminUserDelete{}
	dao := dao_impl.NewDaoPermitUser(r.Context)
	if err = dao.DeleteAdminUser(arg.IdString...); err != nil {
		return
	}

	if err = dao_impl.NewDaoPermitGroupImpl(r.Context).
		DeleteUserGroupByUserId(arg.IdString...); err != nil {
		return
	}
	res.Result = true
	return
}

func (r *SrvPermitUserImpl) AdminUser(arg *wrappers.ArgAdminUser) (res *wrappers.ResultAdminUser, err1 error) {

	res = &wrappers.ResultAdminUser{Pager: *response.NewPagerAndDefault(&arg.PageQuery)}

	dao := dao_impl.NewDaoPermitUser(r.Context)

	var db *gorm.DB

	// 分页获取数据
	err1 = res.Pager.CallGetPagerData(func(pagerObject *response.Pager) (err error) {
		pagerObject.TotalCount, db, err = dao.GetAdminUserCount(db, arg)
		return
	}, func(pagerObject *response.Pager) (err error) {
		var list []models.AdminUser
		if list, err = dao.GetAdminUserList(db, arg, pagerObject); err != nil {
			return
		}
		pagerObject.List, err = r.leftAdminUser(list)
		return
	})
	return
}

func (r *SrvPermitUserImpl) getUserGroup(list []models.AdminUser) (res map[string][]wrappers.AdminUserGroupName, err error) {
	res = map[string][]wrappers.AdminUserGroupName{}
	uIds := make([]string, 0, len(list))
	for _, value := range list {
		uIds = append(uIds, value.UserHid)
	}
	dao := dao_impl.NewDaoPermit(r.Context)
	listResult, err := dao.GetUserGroupByUIds(uIds...)
	if err != nil {
		return
	}
	var tmp wrappers.AdminUserGroupName
	gIds := make([]int64, 0, len(listResult))
	gIdsMap := make(map[int64]int64, len(listResult))
	for _, item := range listResult {
		if _, ok := gIdsMap[item.GroupId]; !ok {
			gIdsMap[item.GroupId] = item.GroupId
			gIds = append(gIds, item.GroupId)
		}
		if _, ok := res[item.UserHid]; !ok {
			res[item.UserHid] = []wrappers.AdminUserGroupName{}
		}
	}

	var groupList []models.AdminGroup
	if groupList, err = dao.GetAdminGroupByIds(gIds...); err != nil {
		return
	}

	groupMap := make(map[int64]models.AdminGroup, len(groupList))
	for _, value := range groupList {
		groupMap[value.Id] = value
	}

	for _, item := range listResult {
		tmp = wrappers.AdminUserGroupName{
			AdminUserGroup: item,
		}
		if _, ok := groupMap[item.GroupId]; ok {
			tmp.GroupName = groupMap[item.GroupId].Name
			tmp.IsSuperAdmin = groupMap[item.GroupId].IsSuperAdmin
			tmp.IsAdminGroup = groupMap[item.GroupId].IsAdminGroup
		}
		res[item.UserHid] = append(res[item.UserHid], tmp)
	}
	return
}

func (r *SrvPermitUserImpl) leftAdminUser(list []models.AdminUser) (res []wrappers.ResultAdminUserList, err error) {
	res = []wrappers.ResultAdminUserList{}
	mapGroupPermit, err := r.getUserGroup(list)
	if err != nil {
		return
	}
	var tmp wrappers.ResultAdminUserList
	for _, item := range list {
		tmp = wrappers.ResultAdminUserList{
			AdminUser: item,
		}
		if _, ok := mapGroupPermit[item.UserHid]; ok {
			tmp.Group = mapGroupPermit[item.UserHid]
		}
		res = append(res, tmp)
	}

	return

}

func NewSrvPermitUserImpl(ctx ...*base.Context) (res srvs.SrvPermitUser) {
	p := &SrvPermitUserImpl{}
	p.SetContext(ctx...)
	return p
}
