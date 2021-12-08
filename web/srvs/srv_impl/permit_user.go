// Package srv_impl
package srv_impl

import (
	"fmt"
	"strings"
	"time"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/app_param"
	"github.com/juetun/base-wrapper/lib/common/response"
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

	if err = dao_impl.NewDaoPermitUser(r.Context).
		UpdateDataByUserHIds(map[string]interface{}{arg.Column: arg.Value}, arg.UserHidVals...); err != nil {
		return
	}
	res.Result = true
	return
}

func (r *SrvPermitUserImpl) AdminUserAdd(arg *wrappers.ArgAdminUserAdd) (res wrappers.ResultAdminUserAdd, err error) {
	res = wrappers.ResultAdminUserAdd{}
	if arg.UserHid == "" {
		err = fmt.Errorf("您没有选择要添加的用户")
		return
	}

	userHid := strings.TrimSpace(arg.UserHid)

	var user *app_param.ResultUserItem
	if user, err = NewUserService(r.Context).
		GetUserById(userHid); err != nil {
		return
	}
	if user.UserHid == "" {
		err = fmt.Errorf("您要添加的用户信息不存在")
		return
	}
	t := time.Now()
	adminUser := &models.AdminUser{
		UserHid:   arg.UserHid,
		RealName:  user.RealName,
		Mobile:    user.Mobile,
		CreatedAt: base.TimeNormal{Time: t},
		UpdatedAt: base.TimeNormal{Time: t},
	}

	err = dao_impl.NewDaoPermit(r.Context).
		AdminUserAdd(adminUser)
	if err != nil {
		return
	}
	res.Result = true
	return
}

func (r *SrvPermitUserImpl) AdminUserDelete(arg *wrappers.ArgAdminUserDelete) (res wrappers.ResultAdminUserDelete, err error) {
	res = wrappers.ResultAdminUserDelete{}
	dao := dao_impl.NewDaoPermit(r.Context)
	if err = dao.DeleteAdminUser(arg.IdString); err != nil {
		return
	}

	if err = dao_impl.NewDaoPermitGroupImpl(r.Context).DeleteUserGroupByUserId(arg.IdString...); err != nil {
		return
	}
	res.Result = true
	return
}

func (r *SrvPermitUserImpl) AdminUser(arg *wrappers.ArgAdminUser) (res *wrappers.ResultAdminUser, err1 error) {

	res = &wrappers.ResultAdminUser{Pager: *response.NewPagerAndDefault(&arg.PageQuery)}

	dao := dao_impl.NewDaoPermit(r.Context)

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
		pagerObject.List, err = r.leftAdminUser(list, dao)
		return
	})
	return
}

func (r *SrvPermitUserImpl) getUserGroup(list []models.AdminUser, dao daos.DaoPermit) (res map[string][]wrappers.AdminUserGroupName, err error) {
	res = map[string][]wrappers.AdminUserGroupName{}
	uIds := make([]string, 0, len(list))
	for _, value := range list {
		uIds = append(uIds, value.UserHid)
	}
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

func (r *SrvPermitUserImpl) leftAdminUser(list []models.AdminUser, dao daos.DaoPermit) (res []wrappers.ResultAdminUserList, err error) {
	res = []wrappers.ResultAdminUserList{}
	mapGroupPermit, err := r.getUserGroup(list, dao)
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
