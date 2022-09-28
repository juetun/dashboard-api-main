/**
* @Author:changjiang
* @Description:
* @File:service
* @Version: 1.0.0
* @Date 2021/5/22 5:16 下午
 */
package admin_impl

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	cons_admin2 "github.com/juetun/dashboard-api-main/web/cons/admins"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"io/ioutil"
)

type ConServiceImpl struct {
	base.ControllerBase
}

func (r *ConServiceImpl) getJsonArg(c *gin.Context, data interface{}) (err error) {
	var bt bytes.Buffer
	var body []byte
	if body, err = ioutil.ReadAll(c.Request.Body); err != nil {
		return
	}
	bt.Write(body)
	err = json.Unmarshal(body, data)
	return
}

func (r *ConServiceImpl) Edit(c *gin.Context) {
	var (
		arg wrappers.ArgServiceEdit
		res *wrappers.ResultServiceEdit
		err error
		ctx = base.CreateContext(&r.ControllerBase, c)
	)

	if haveErr := r.ParametersAccept(ctx, &arg); haveErr {
		return
	}
	if res, err = srv_impl.NewSrvServiceImpl(ctx).
		Edit(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, 0, res)
}

func (r *ConServiceImpl) Detail(c *gin.Context) {
	var (
		arg wrappers.ArgDetail
		err error
		res *wrappers.ResultDetail
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg); haveErr {
		return
	}
	if res, err = srv_impl.NewSrvServiceImpl(ctx).
		Detail(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res, )
}

func (r *ConServiceImpl) List(c *gin.Context) {
	var (
		arg wrappers.ArgServiceList
		err error
		res *wrappers.ResultServiceList
		ctx = base.CreateContext(&r.ControllerBase, c)
	)

	if haveErr := r.ParametersAccept(ctx, &arg); haveErr {
		return
	}
	srv := srv_impl.NewSrvServiceImpl(ctx)
	if res, err = srv.List(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res, )
}

func NewConServiceImpl() cons_admin2.ConService {
	p := &ConServiceImpl{}
	p.ControllerBase.Init()
	return p
}
