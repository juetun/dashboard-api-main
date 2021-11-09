package intranet_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/cons/intranets"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_intranet"
)

type ConPermitIntranetImpl struct {
	base.ControllerBase
}

// GetImportPermit 获取接口权限
func (r *ConPermitIntranetImpl) GetImportPermit(c *gin.Context) {
	var args wrapper_intranet.ArgGetImportPermit
	var err error
	var res *wrapper_intranet.ResultGetImportPermit

	if err = c.Bind(&args); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = args.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	srv := srv_impl.NewSrvGatewayImportPermit(base.CreateContext(&r.ControllerBase, c))

	if res, err = srv.GetImportPermit(&args); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}

func NewConPermitIntranet() (res intranets.ConPermitIntranet) {
	p := &ConPermitIntranetImpl{}
	p.ControllerBase.Init()
	return p
}
