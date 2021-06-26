/**
* @Author:changjiang
* @Description:
* @File:PermitGroupImpl
* @Version: 1.0.0
* @Date 2021/6/20 10:03 下午
 */
package srv_impl

import (
	"time"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type SrvPermitGroupImpl struct {
	base.ServiceBase
}

func (r *SrvPermitGroupImpl) MenuImportSet(arg *wrappers.ArgMenuImportSet) (res *wrappers.ResultMenuImportSet, err error) {
	res = &wrappers.ResultMenuImportSet{
		Result: false,
	}
	var m models.AdminMenuImport
	t := time.Now()

	var dts = make([]models.AdminMenuImport, 0, len(arg.ImportIds))
	var dt models.AdminMenuImport
	if arg.Type == "delete" {
		dt.DeletedAt = &t
		for _, value := range arg.ImportIds {
			if value == 0 {
				continue
			}
			dt = models.AdminMenuImport{
				MenuId:    arg.MenuId,
				ImportId:  value,
				CreatedAt: t,
				UpdatedAt: t,
				DeletedAt: &t,
			}
			dts = append(dts, dt)
		}
	} else {
		for _, value := range arg.ImportIds {
			if value == 0 {
				continue
			}
			dt = models.AdminMenuImport{
				MenuId:    arg.MenuId,
				ImportId:  value,
				CreatedAt: t,
				UpdatedAt: t,
			}
			dts = append(dts, dt)
		}
	}
	if err = dao_impl.NewPermitImportImpl(r.Context).
		BatchMenuImport(m.TableName(), dts); err != nil {
		return
	}

	res.Result = true
	return
}

func NewSrvPermitGroupImpl(context ...*base.Context) srvs.SrvPermitGroup {
	p := &SrvPermitGroupImpl{}
	p.SetContext(context...)
	return p
}
