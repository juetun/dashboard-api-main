package intranet_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/cons/intranets"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/library/common/app_param/operate_admin"
)

type ConLogIntranetImpl struct {
	base.ControllerBase
}

func (r *ConLogIntranetImpl) AddLog(c *gin.Context) {
	var (
		args operate_admin.ArgAddAdminLog
		err  error
		res  *operate_admin.ResultAddAdminLog
		ctx  = base.CreateContext(&r.ControllerBase, c)
	)

	if r.ParametersAccept(ctx, &args, base.ControllerGetParamTypeBind) {
		return
	}

	if res, err = srv_impl.NewSrvOperate(ctx).
		AddLogMessage(&args); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}

func NewConLogIntranet() (res intranets.ConLogIntranet) {
	p := &ConLogIntranetImpl{}
	p.ControllerBase.Init()
	return p
}
