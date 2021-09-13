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
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	cons_admin2 "github.com/juetun/dashboard-api-main/web/cons/admins"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/dashboard-api-main/web/wrappers"
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
	fmt.Println(string(body))
	err = json.Unmarshal(body, data)
	return
}
func (r *ConServiceImpl) Edit(c *gin.Context) {
	var (
		arg wrappers.ArgServiceEdit
		err error
		res *wrappers.ResultServiceEdit
	)
	if err = c.ShouldBind(&arg); err != nil {
		r.Response(c, -1, nil, err.Error())
		return
	}
	srv := srv_impl.NewSrvServiceImpl(base.CreateContext(&r.ControllerBase, c))
	if res, err = srv.Edit(&arg); err != nil {
		r.Response(c, -1, res, err.Error())
		return
	}
	r.Response(c, 0, res, "success")
}

func (r *ConServiceImpl) Detail(c *gin.Context) {
	var (
		arg wrappers.ArgDetail
		err error
		res *wrappers.ResultDetail
	)
	err = c.ShouldBind(&arg)
	if err != nil {
		r.Response(c, -1, nil, err.Error())
		return
	}
	if arg.Id, err = strconv.Atoi(c.Params.ByName("id")); err != nil {
		r.Response(c, -1, nil, "参数格式不正确")
		return
	}
	srv := srv_impl.NewSrvServiceImpl(base.CreateContext(&r.ControllerBase, c))
	if res, err = srv.Detail(&arg); err != nil {
		r.Response(c, -1, res, err.Error())
		return
	}
	r.Response(c, 0, res, "success")
}
func (r *ConServiceImpl) List(c *gin.Context) {
	var (
		arg wrappers.ArgServiceList
		err error
		res *wrappers.ResultServiceList
	)
	err = c.ShouldBind(&arg)
	if err != nil {
		r.Response(c, -1, nil, err.Error())
		return
	}
	srv := srv_impl.NewSrvServiceImpl(base.CreateContext(&r.ControllerBase, c))
	if res, err = srv.List(&arg); err != nil {
		r.Response(c, -1, res, err.Error())
		return
	}
	r.Response(c, 0, res, "success")
}

func NewConServiceImpl() cons_admin2.ConService {
	p := &ConServiceImpl{}
	p.ControllerBase.Init()
	return p
}
