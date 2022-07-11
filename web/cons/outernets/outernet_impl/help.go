package outernet_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/cons/outernets"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_outernet"
)

type ConOuterNetsHelpImpl struct {
	base.ControllerBase
}

func (r *ConOuterNetsHelpImpl) Tree(c *gin.Context) {

	var (
		arg wrapper_outernet.ArgTree
		res *wrapper_outernet.ResultTree
		err error
	)

	if err = c.ShouldBind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	if res, err = srv_impl.NewSrvHelpRelate(base.CreateContext(&r.ControllerBase, c)).
		Tree(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	r.Response(c, base.SuccessCode, res)
 	return
}

func NewConOuterNetsHelp() outernets.ConOuterNetsHelp {
	c := &ConOuterNetsHelpImpl{}
	c.ControllerBase.Init()
	return c
}
