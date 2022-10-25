package srv_impl

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
)

type SrvPermitCommon struct {
	base.ServiceBase
}

func (r *SrvPermitCommon) GetAdminUserInfo(getUserArgument *base.ArgGetByNumberIds, userHid int64) (isSuperAdmin bool, err error) {
	var adminUserMap models.AdminUserMap = make(map[int64]*models.AdminUser, len(getUserArgument.Ids))
	if adminUserMap, err = dao_impl.NewDaoPermitUser(r.Context).
		GetAdminUserByIds(getUserArgument); err != nil {
		return
	}
	if dt, ok := adminUserMap[userHid]; !ok {
		err = fmt.Errorf("当前账号没有管理员权限")
		return
	} else if dt.FlagAdmin == models.AdminUserFlagAdminYes {
		isSuperAdmin = true
	}
	return
}
