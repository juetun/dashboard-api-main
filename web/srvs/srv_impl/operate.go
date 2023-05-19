package srv_impl

import (
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/base-wrapper/lib/plugins/rpc"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
	"github.com/juetun/library/common/app_param"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type SrvOperateImpl struct {
	base.ServiceBase
	dao daos.DaoOperateLog
}

func (r *SrvOperateImpl) OperateLog(arg *wrapper_admin.ArgOperateLog) (res *wrapper_admin.ResultOperateLog, err error) {
	res = &wrapper_admin.ResultOperateLog{Pager: response.NewPager(response.PagerBaseQuery(&arg.PageQuery)), Tabs: base.ModelItemOptions{{Label: "操作日志", Value: "operate_log"}}}
	var (
		actResObj *base.ActErrorHandlerResult
		list      []*models.OperateLog
	)
	err = res.Pager.CallGetPagerData(func(pager *response.Pager) (e error) {
		actResObj, pager.TotalCount, e = r.dao.GetCount(arg)
		return
	}, func(pager *response.Pager) (e error) {
		if list, e = r.dao.GetList(actResObj, arg, pager); e != nil {
			return
		}
		if res.Pager.List, e = r.parseLogList(list); e != nil {
			return
		}
		return
	})
	return
}

func (r *SrvOperateImpl) GetResultUserByUid(userId string, ctx *base.Context) (res *app_param.ResultUser, err error) {
	var value = url.Values{}

	value.Set("user_hid", userId)
	value.Set("data_type", strings.Join([]string{app_param.UserDataTypeMain, app_param.UserDataTypeInfo, app_param.UserDataTypePortrait}, ","))
	ro := rpc.RequestOptions{
		Method:      http.MethodPost,
		AppName:     app_param.AppNameUser,
		URI:         "/user/get_by_uid",
		Header:      http.Header{},
		Value:       value,
		Context:     ctx,
		PathVersion: app_obj.App.AppRouterPrefix.Intranet,
	}
	var data = struct {
		Code int                   `json:"code"`
		Data *app_param.ResultUser `json:"data"`
		Msg  string                `json:"message"`
	}{}
	err = rpc.NewHttpRpc(&ro).
		Send().
		GetBody().
		Bind(&data).Error
	if err != nil {
		return
	}
	res = data.Data
	return
}

func (r *SrvOperateImpl) parseLogList(logs []*models.OperateLog) (res []*wrapper_admin.ResultOperateLogItem, err error) {
	var (
		l          = len(logs)
		dataItem   *wrapper_admin.ResultOperateLogItem
		userMap    = make(map[int64]bool, l)
		userHids   = make([]string, 0, l)
		ResultUser *app_param.ResultUser
	)
	res = make([]*wrapper_admin.ResultOperateLogItem, 0, l)

	for _, item := range logs {
		if _, ok := userMap[item.UserHid]; !ok {
			userHids = append(userHids, strconv.FormatInt(item.UserHid, 10))
			userMap[item.UserHid] = true
		}
	}
	if ResultUser, err = r.GetResultUserByUid(strings.Join(userHids, ","), r.Context); err != nil {
		return
	}
	for _, item := range logs {
		dataItem = &wrapper_admin.ResultOperateLogItem{}
		if err = dataItem.ParseLog(item); err != nil {
			return
		}
		if user, ok := ResultUser.List[item.UserHid]; ok {
			dataItem.Avatar = user.PortraitUrl
			dataItem.NickName = user.NickName
			dataItem.Name = user.RealName
		}
		res = append(res, dataItem)
	}

	return
}

func NewSrvOperate(ctx ...*base.Context) srvs.SrvOperate {
	p := &SrvOperateImpl{}
	p.SetContext(ctx...)
	p.dao = dao_impl.NewDaoOperateLog(p.Context)
	return p
}
