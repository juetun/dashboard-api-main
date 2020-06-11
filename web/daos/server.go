/**
* @Author:changjiang
* @Description:
* @File:server
* @Version: 1.0.0
* @Date 2020/6/11 10:16 上午
 */
package daos

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/models"
)

type DaoAppPath struct {
	base.ServiceDao
}

func NewDaoAppPath(context ...*base.Context) (p *DaoAppPath) {
	p = &DaoAppPath{}
	p.SetContext(context)
	return
}

func (r *DaoAppPath) GetAllApp() (res *[]models.ZAppPath, err error) {
	res= &[]models.ZAppPath{}
	var m models.ZAppPath
	err = r.Context.Db.
		Table(m.TableName()).
		Find(res).
		Error
	return
}
