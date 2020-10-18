/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 12:20 上午
 */
package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
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

func (r *PermitService) AdminUserGroupAdd(arg *pojos.ArgAdminUserGroupAdd) (res pojos.ResultAdminUserGroupAdd, err error) {
	res = pojos.ResultAdminUserGroupAdd{}
	dao := daos.NewDaoPermit(r.Context)

	var args = make([]map[string]interface{}, 0)
	for _, userHId := range arg.UserHIds {
		for _, groupId := range arg.GroupIds {
			args = append(args, map[string]interface{}{
				"group_id": groupId,
				"user_hid": userHId,
				"is_del":   0,
			})
		}
	}
	err = dao.AdminUserGroupAdd(args)
	if err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitService) AdminUserGroupRelease(arg *pojos.ArgAdminUserGroupRelease) (res pojos.ResultAdminUserGroupRelease, err error) {
	res = pojos.ResultAdminUserGroupRelease{}
	dao := daos.NewDaoPermit(r.Context)
	err = dao.AdminUserGroupRelease(arg.IdString)
	if err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitService) AdminUserDelete(arg *pojos.ArgAdminUserDelete) (res pojos.ResultAdminUserDelete, err error) {
	res = pojos.ResultAdminUserDelete{}
	dao := daos.NewDaoPermit(r.Context)
	err = dao.AdminUserDelete(arg.IdString)
	if err != nil {
		return
	}
	res.Result = true
	return
}

func (r *PermitService) AdminUserAdd(arg *pojos.ArgAdminUserAdd) (res pojos.ResultAdminUserAdd, err error) {
	res = pojos.ResultAdminUserAdd{}
	if arg.UserHid == "" {
		err = fmt.Errorf("您没有选择要添加的用户")
		return
	}
	user, err := NewUserService().GetUserById(arg.UserHid)
	if err != nil {
		return
	}
	err = daos.NewDaoPermit(r.Context).
		AdminUserAdd(&models.AdminUser{
			UserHid:  arg.UserHid,
			RealName: user.Name,
			Mobile:   user.Mobile,
		})
	if err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitService) AdminGroupDelete(arg *pojos.ArgAdminGroupDelete) (res pojos.ResultAdminGroupDelete, err error) {
	res = pojos.ResultAdminGroupDelete{}
	dao := daos.NewDaoPermit(r.Context)
	err = dao.DeleteAdminGroupByIds(arg.IdString)
	if err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitService) AdminGroupEdit(arg *pojos.ArgAdminGroupEdit) (res *pojos.ResultAdminGroupEdit, err error) {
	res = &pojos.ResultAdminGroupEdit{}
	dao := daos.NewDaoPermit(r.Context)
	var g models.AdminGroup
	if g, err = dao.FetchByName(arg.Name); err != nil {
		return
	}
	if arg.Id == 0 {
		if g.Name != "" {
			err = fmt.Errorf("您输入的组名已存在")
			return
		}
		if err = dao.InsertAdminGroup(&models.AdminGroup{
			Name:  arg.Name,
			IsDel: 0,
		}); err != nil {
			return
		}
		res.Result = true
		return
	}

	if g.Name != "" && g.Id != arg.Id {
		err = fmt.Errorf("您输入的组名已存在")
		return
	}

	if dao.UpdateAdminGroup(&models.AdminGroup{
		Name:  arg.Name,
		Id:    arg.Id,
		IsDel: 0,
	}); err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitService) MenuAdd(arg *pojos.ArgMenuAdd) (res *pojos.ResultMenuAdd, err error) {
	res = &pojos.ResultMenuAdd{}
	dao := daos.NewDaoPermit(r.Context)
	// timeNow:=time.Now().Format("2006-01-02 15:04:05")
	t := time.Now()
	err = dao.Add(&models.AdminMenu{
		ParentId:   arg.ParentId,
		AppName:    arg.AppName,
		Label:      arg.Label,
		Icon:       arg.Icon,
		IsMenuShow: arg.IsMenuShow,
		AppVersion: arg.AppVersion,
		UrlPath:    arg.UrlPath,
		PathType:   arg.PathType,
		SortValue:  arg.SortValue,
		OtherValue: arg.OtherValue,
		CreatedAt:  t,
		UpdatedAt:  t,
		IsDel:      0,
	})
	return
}
func (r *PermitService) MenuDelete(arg *pojos.ArgMenuDelete) (res *pojos.ResultMenuDelete, err error) {
	res = &pojos.ResultMenuDelete{}
	dao := daos.NewDaoPermit(r.Context)
	err = dao.DeleteByIds(arg.IdValue)
	return
}

func (r *PermitService) MenuSave(arg *pojos.ArgMenuSave) (res *pojos.ResultMenuSave, err error) {
	res = &pojos.ResultMenuSave{}
	dao := daos.NewDaoPermit(r.Context)
	t := time.Now()
	err = dao.Save(arg.Id, &models.AdminMenu{
		ParentId:   arg.ParentId,
		AppName:    arg.AppName,
		Label:      arg.Label,
		Icon:       arg.Icon,
		IsMenuShow: arg.IsMenuShow,
		AppVersion: arg.AppVersion,
		UrlPath:    arg.UrlPath,
		PathType:   arg.PathType,
		SortValue:  arg.SortValue,
		OtherValue: arg.OtherValue,
		UpdatedAt:  t,
		IsDel:      0,
	})
	return
}
func (r *PermitService) AdminMenu(arg *pojos.ArgAdminMenu) (res *pojos.ResultAdminMenu, err error) {
	res = &pojos.ResultAdminMenu{List: make([]pojos.AdminMenuObject, 0, 20)}
	dao := daos.NewDaoPermit(r.Context)
	var list []models.AdminMenu
	list, err = dao.GetAdminMenuList(arg)
	r.orgTree(list, 0, &res.List)
	return
}

//
func (r *PermitService) orgTree(list []models.AdminMenu, parentId int, res *[]pojos.AdminMenuObject) () {
	var tmp pojos.AdminMenuObject
	for _, value := range list {
		if value.ParentId != parentId {
			continue
		}
		tmp = pojos.AdminMenuObject{Children: make([]pojos.AdminMenuObject, 0, 20),}

		tmp.ResultAdminMenuSingle = pojos.ResultAdminMenuSingle{
			Id:         value.Id,
			ParentId:   value.ParentId,
			AppName:    value.AppName,
			Title:      value.Label,
			Icon:       value.Icon,
			IsMenuShow: value.IsMenuShow,
			AppVersion: value.AppVersion,
			UrlPath:    value.UrlPath,
			PathType:   value.PathType,
			SortValue:  value.SortValue,
			IsDel:      value.IsDel,
		}
		if value.OtherValue != "" {
			json.Unmarshal([]byte(value.OtherValue), &tmp.ResultAdminMenuSingle.ResultAdminMenuOtherValue)
		} else {
			tmp.ResultAdminMenuSingle.ResultAdminMenuOtherValue = pojos.ResultAdminMenuOtherValue{Expand: true,}
		}
		r.orgTree(list, value.Id, &tmp.Children)
		*res = append(*res, tmp)
	}
}

func (r *PermitService) AdminGroup(arg *pojos.ArgAdminGroup) (res *pojos.ResultAdminGroup, err error) {

	res = &pojos.ResultAdminGroup{Pager: *response.NewPagerAndDefault(&arg.BaseQuery),}
	var db *gorm.DB
	dao := daos.NewDaoPermit(r.Context)
	// 获取分页数据
	res.Pager.CallGetPagerData(func(pagerObject *response.Pager) (err error) {
		pagerObject.TotalCount, db, err = dao.GetAdminGroupCount(db, arg)
		return
	}, func(pagerObject *response.Pager) (err error) {
		pagerObject.List, err = dao.GetAdminGroupList(db, arg)
		return
	})
	return
}

func (r *PermitService) AdminUser(arg *pojos.ArgAdminUser) (res *pojos.ResultAdminUser, err error) {
	res = &pojos.ResultAdminUser{Pager: *response.NewPagerAndDefault(&arg.BaseQuery),}
	dao := daos.NewDaoPermit(r.Context)

	var db *gorm.DB

	// 分页获取数据
	res.Pager.CallGetPagerData(func(pagerObject *response.Pager) (err error) {
		pagerObject.TotalCount, db, err = dao.GetAdminUserCount(db, arg)
		return
	}, func(pagerObject *response.Pager) (err error) {
		list, err := dao.GetAdminUserList(db, arg)
		if err != nil {
			return
		}
		pagerObject.List, err = r.leftAdminUser(list, dao)
		return
	})
	return
}

func (r *PermitService) getUserGroup(list []models.AdminUser, dao *daos.DaoPermit) (res map[string][]pojos.AdminUserGroupName, err error) {
	res = map[string][]pojos.AdminUserGroupName{}
	uIds := make([]string, 0, len(list))
	for _, value := range list {
		uIds = append(uIds, value.UserHid)
	}
	listResult, err := dao.GetUserGroupByUIds(uIds)
	if err != nil {
		return
	}
	var tmp pojos.AdminUserGroupName
	gIds := make([]int, 0, len(listResult))
	gIdsMap := make(map[int]int, len(listResult))
	for _, item := range listResult {
		if _, ok := gIdsMap[item.GroupId]; !ok {
			gIdsMap[item.GroupId] = item.GroupId
			gIds = append(gIds, item.GroupId)
		}
		if _, ok := res[item.UserHid]; !ok {
			res[item.UserHid] = []pojos.AdminUserGroupName{}
		}
	}
	groupList, err := dao.GetAdminGroupByIds(gIds)
	if err != nil {
		return
	}

	groupMap := make(map[int]string, len(groupList))
	for _, value := range groupList {
		groupMap[value.Id] = value.Name
	}

	for _, item := range listResult {
		tmp = pojos.AdminUserGroupName{
			AdminUserGroup: item,
		}
		if _, ok := groupMap[item.GroupId]; ok {
			tmp.GroupName = groupMap[item.GroupId]
		}
		res[item.UserHid] = append(res[item.UserHid], tmp)
	}
	return
}

func (r *PermitService) leftAdminUser(list []models.AdminUser, dao *daos.DaoPermit) (res []pojos.ResultAdminUserList, err error) {
	res = []pojos.ResultAdminUserList{}
	mapGroupPermit, err := r.getUserGroup(list, dao)
	if err != nil {
		return
	}
	var tmp pojos.ResultAdminUserList
	for _, item := range list {
		tmp = pojos.ResultAdminUserList{
			AdminUser: item,
		}
		if _, ok := mapGroupPermit[item.UserHid]; ok {
			tmp.Group = mapGroupPermit[item.UserHid]
		}
		res = append(res, tmp)
	}

	return

}
func (r *PermitService) Menu(arg *pojos.ArgPermitMenu) (res *pojos.ResultPermitMenu, err error) {
	res = &pojos.ResultPermitMenu{}
	dao := daos.NewDaoPermit(r.Context)
	groupIds, err := r.getUserGroupIds(dao, arg)
	if err != nil {
		return
	}
	menuIds, err := r.getPermitByGroupIds(dao, arg.PathType, groupIds...)
	if err != nil {
		return
	}
	res.Menu, err = r.getGroupMenu(dao, menuIds...)
	return
}
func (r *PermitService) getGroupMenu(dao *daos.DaoPermit, menuIds ...int) (res []models.AdminMenu, err error) {
	res = []models.AdminMenu{}
	if len(menuIds) == 0 {
		return
	}
	res, err = dao.GetPermitMenuByIds(menuIds...)
	return
}
func (r *PermitService) getPermitByGroupIds(dao *daos.DaoPermit, pathType string, groupIds ...int) (menuIds []int, err error) {
	res, err := dao.GetMenuIdsByPermitByGroupIds(pathType, groupIds...)
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
