/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 12:20 上午
 */
package services

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/pojos"
)

type PermitService struct {
	base.ServiceBase
}

func NewPermitService(context ...*base.Context) (p *PermitService) {
	p = &PermitService{}
	p.SetContext(context...)
	return
}

func (r *PermitService) Menu(arg *pojos.ArgPermitMenu) (res *pojos.ResultPermitMenu, err error) {
	res = &pojos.ResultPermitMenu{}
	dao := daos.NewDaoPermit(r.Context)
	groupIds, err := r.getUserGroupIds(dao, arg)
	if err != nil {
		return
	}
	menuIds, err := r.getPermitByGroupIds(dao, groupIds...)
	if err != nil {
		return
	}
	res.Menu, err = r.getGroupMenu(dao, menuIds...)
	return
}
func (r *PermitService) getGroupMenu(dao *daos.DaoPermit, menuIds ...int) (res []models.AdminMenu, err error) {
	if len(menuIds) == 0 {
		return
	}
	res, err = dao.GetPermitMenuByIds(menuIds...)
	return
}
func (r *PermitService) getPermitByGroupIds(dao *daos.DaoPermit, groupIds ...int) (menuIds []int, err error) {
	res, err := dao.GetMenuIdsByPermitByGroupIds(groupIds...)
	if err != nil {
		return
	}
	menuIds = make([]int, 0, len(res))
	for _, value := range res {
		menuIds = append(menuIds, value.MenuId)
	}
	return
}

func (r *PermitService) getUserGroupIds(dao *daos.DaoPermit, arg *pojos.ArgPermitMenu) (res []int, err error) {
	groups, err := dao.GetGroupByUserId(arg.UserId)
	if err != nil {
		return
	}
	res = make([]int, 0, len(groups))
	for _, group := range groups {
		res = append(res, group.GroupId)
	}
	return
}

func (r *PermitService) Flag(arg *pojos.ArgFlag) (res *pojos.ResultFlag, err error) {
	res = &pojos.ResultFlag{}
	return
}
