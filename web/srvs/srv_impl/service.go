/**
* @Author:changjiang
* @Description:
* @File:service
* @Version: 1.0.0
* @Date 2021/5/22 11:30 下午
 */
package srv_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type SrvServiceImpl struct {
	base.ServiceBase
}

func (r *SrvServiceImpl) Edit(arg *wrappers.ArgServiceEdit) (res *wrappers.ResultServiceEdit, err error) {
	res = &wrappers.ResultServiceEdit{}

	res.Result = true
	return
}

func (r *SrvServiceImpl) List(arg *wrappers.ArgServiceList) (res *wrappers.ResultServiceList, err error) {
	res = &wrappers.ResultServiceList{
		Pager: response.NewPager(),
	}

	return
}

func NewSrvServiceImpl(context ...*base.Context) srvs.SrvService {
	p := &SrvServiceImpl{}
	p.SetContext(context...)
	return p
}
