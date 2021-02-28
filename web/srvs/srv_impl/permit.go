/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 12:20 上午
 */
package srv_impl

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type PermitService struct {
	base.ServiceBase
}

func NewPermitService(context ...*base.Context) (p *PermitService) {
	p = &PermitService{}
	p.SetContext(context...)
	return
}

func (r *PermitService) GetMenu(arg *wrappers.ArgGetMenu) (res wrappers.ResultGetMenu, err error) {
	res = wrappers.ResultGetMenu{}
	if arg.MenuId == 0 {
		return
	}
	res.AdminMenu, err = dao_impl.NewDaoPermit(r.Context).GetMenu(arg.MenuId)
	if err != nil {
		return
	}
	return
}
func (r *PermitService) AdminMenuSearch(arg *wrappers.ArgAdminMenu) (res wrappers.ResAdminMenuSearch, err error) {
	res = wrappers.ResAdminMenuSearch{
		List: []models.AdminMenu{},
	}
	dao := dao_impl.NewDaoPermit(r.Context)
	res.List, err = dao.GetAdminMenuList(arg)
	return
}
func (r *PermitService) AdminUserGroupAdd(arg *wrappers.ArgAdminUserGroupAdd) (res wrappers.ResultAdminUserGroupAdd, err error) {
	res = wrappers.ResultAdminUserGroupAdd{}
	dao := dao_impl.NewDaoPermit(r.Context)

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
func (r *PermitService) AdminUserGroupRelease(arg *wrappers.ArgAdminUserGroupRelease) (res wrappers.ResultAdminUserGroupRelease, err error) {
	res = wrappers.ResultAdminUserGroupRelease{}
	dao := dao_impl.NewDaoPermit(r.Context)
	err = dao.AdminUserGroupRelease(arg.IdString)
	if err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitService) AdminUserDelete(arg *wrappers.ArgAdminUserDelete) (res wrappers.ResultAdminUserDelete, err error) {
	res = wrappers.ResultAdminUserDelete{}
	dao := dao_impl.NewDaoPermit(r.Context)
	err = dao.AdminUserDelete(arg.IdString)
	if err != nil {
		return
	}
	res.Result = true
	return
}

func (r *PermitService) AdminUserAdd(arg *wrappers.ArgAdminUserAdd) (res wrappers.ResultAdminUserAdd, err error) {
	res = wrappers.ResultAdminUserAdd{}
	if arg.UserHid == "" {
		err = fmt.Errorf("您没有选择要添加的用户")
		return
	}
	var user *models.UserMain
	if user, err = NewUserService(r.Context).
		GetUserById(strings.TrimSpace(arg.UserHid)); err != nil {
		return
	} else if user.UserHid == "" {
		err = fmt.Errorf("您要添加的用户信息不存在")
		return
	}

	if err = dao_impl.NewDaoPermit(r.Context).
		AdminUserAdd(&models.AdminUser{
			UserHid:  arg.UserHid,
			RealName: user.Name,
			Mobile:   user.Mobile,
			Model: gorm.Model{
				ID:        0,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeletedAt: nil,
			},
		}); err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitService) AdminGroupDelete(arg *wrappers.ArgAdminGroupDelete) (res wrappers.ResultAdminGroupDelete, err error) {
	res = wrappers.ResultAdminGroupDelete{}
	dao := dao_impl.NewDaoPermit(r.Context)
	err = dao.DeleteAdminGroupByIds(arg.IdString)
	if err != nil {
		return
	}
	res.Result = true
	return
}
func (r *PermitService) AdminGroupEdit(arg *wrappers.ArgAdminGroupEdit) (res *wrappers.ResultAdminGroupEdit, err error) {
	res = &wrappers.ResultAdminGroupEdit{}
	dao := dao_impl.NewDaoPermit(r.Context)
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
func (r *PermitService) MenuAdd(arg *wrappers.ArgMenuAdd) (res *wrappers.ResultMenuAdd, err error) {
	res = &wrappers.ResultMenuAdd{}
	dao := dao_impl.NewDaoPermit(r.Context)
	// timeNow:=time.Now().Format("2006-01-02 15:04:05")
	t := time.Now()
	err = dao.Add(&models.AdminMenu{
		ParentId:   arg.ParentId,
		Label:      arg.Label,
		Icon:       arg.Icon,
		IsMenuShow: arg.IsMenuShow,
		UrlPath:    arg.UrlPath,
		SortValue:  arg.SortValue,
		OtherValue: arg.OtherValue,
		CreatedAt:  t,
		UpdatedAt:  t,
		IsDel:      0,
	})
	return
}
func (r *PermitService) MenuDelete(arg *wrappers.ArgMenuDelete) (res *wrappers.ResultMenuDelete, err error) {
	res = &wrappers.ResultMenuDelete{}
	dao := dao_impl.NewDaoPermit(r.Context)
	err = dao.DeleteByIds(arg.IdValue)
	res.Result = true
	return
}

func (r *PermitService) MenuSave(arg *wrappers.ArgMenuSave) (res *wrappers.ResultMenuSave, err error) {
	res = &wrappers.ResultMenuSave{}
	dao := dao_impl.NewDaoPermit(r.Context)
	t := time.Now()
	err = dao.Save(arg.Id, &models.AdminMenu{
		ParentId:   arg.ParentId,
		Label:      arg.Label,
		Icon:       arg.Icon,
		IsMenuShow: arg.IsMenuShow,
		UrlPath:    arg.UrlPath,
		SortValue:  arg.SortValue,
		OtherValue: arg.OtherValue,
		UpdatedAt:  t,
		IsDel:      0,
	})
	return
}
func (r *PermitService) AdminMenu(arg *wrappers.ArgAdminMenu) (res *wrappers.ResultAdminMenu, err error) {
	res = &wrappers.ResultAdminMenu{List: make([]wrappers.AdminMenuObject, 0, 20)}
	dao := dao_impl.NewDaoPermit(r.Context)
	var list []models.AdminMenu
	list, err = dao.GetAdminMenuList(arg)
	r.orgTree(list, wrappers.DefaultPermitParentId, &res.List)
	return
}

//
func (r *PermitService) orgTree(list []models.AdminMenu, parentId int, res *[]wrappers.AdminMenuObject) () {
	var tmp wrappers.AdminMenuObject
	for _, value := range list {
		// 剔除默认数据那条
		if value.Id == parentId || value.ParentId != parentId {
			continue
		}
		tmp = r.orgAdminMenuObject(&value)
		r.orgTree(list, value.Id, &tmp.Children)
		*res = append(*res, tmp)
	}
}

func (r *PermitService) orgAdminMenuObject(value *models.AdminMenu) (res wrappers.AdminMenuObject) {
	res = wrappers.AdminMenuObject{Children: make([]wrappers.AdminMenuObject, 0, 20),}
	res.ResultAdminMenuSingle = wrappers.ResultAdminMenuSingle{
		Id:         value.Id,
		ParentId:   value.ParentId,
		Title:      value.Label,
		Label:      value.Label,
		Icon:       value.Icon,
		IsMenuShow: value.IsMenuShow,
		UrlPath:    value.UrlPath,
		SortValue:  value.SortValue,
		IsDel:      value.IsDel,
	}
	if value.OtherValue != "" {
		json.Unmarshal([]byte(value.OtherValue), &res.ResultAdminMenuSingle.ResultAdminMenuOtherValue)
	} else {
		res.ResultAdminMenuSingle.ResultAdminMenuOtherValue = wrappers.ResultAdminMenuOtherValue{Expand: true,}
	}
	return
}
func (r *PermitService) AdminGroup(arg *wrappers.ArgAdminGroup) (res *wrappers.ResultAdminGroup, err error) {

	res = &wrappers.ResultAdminGroup{Pager: *response.NewPagerAndDefault(&arg.BaseQuery),}

	var db *gorm.DB
	dao := dao_impl.NewDaoPermit(r.Context)
	// 获取分页数据
	res.Pager.CallGetPagerData(func(pagerObject *response.Pager) (err error) {
		pagerObject.TotalCount, db, err = dao.GetAdminGroupCount(db, arg)
		return
	}, func(pagerObject *response.Pager) (err error) {
		pagerObject.List, err = dao.GetAdminGroupList(db, arg, pagerObject)
		return
	})
	return
}

func (r *PermitService) AdminUser(arg *wrappers.ArgAdminUser) (res *wrappers.ResultAdminUser, err error) {
	res = &wrappers.ResultAdminUser{Pager: *response.NewPagerAndDefault(&arg.BaseQuery),}
	dao := dao_impl.NewDaoPermit(r.Context)

	var db *gorm.DB

	// 分页获取数据
	res.Pager.CallGetPagerData(func(pagerObject *response.Pager) (err error) {
		pagerObject.TotalCount, db, err = dao.GetAdminUserCount(db, arg)
		return
	}, func(pagerObject *response.Pager) (err error) {
		list, err := dao.GetAdminUserList(db, arg, pagerObject)
		if err != nil {
			return
		}
		pagerObject.List, err = r.leftAdminUser(list, dao)
		return
	})
	return
}

func (r *PermitService) getUserGroup(list []models.AdminUser, dao *dao_impl.DaoPermit) (res map[string][]wrappers.AdminUserGroupName, err error) {
	res = map[string][]wrappers.AdminUserGroupName{}
	uIds := make([]string, 0, len(list))
	for _, value := range list {
		uIds = append(uIds, value.UserHid)
	}
	listResult, err := dao.GetUserGroupByUIds(uIds)
	if err != nil {
		return
	}
	var tmp wrappers.AdminUserGroupName
	gIds := make([]int, 0, len(listResult))
	gIdsMap := make(map[int]int, len(listResult))
	for _, item := range listResult {
		if _, ok := gIdsMap[item.GroupId]; !ok {
			gIdsMap[item.GroupId] = item.GroupId
			gIds = append(gIds, item.GroupId)
		}
		if _, ok := res[item.UserHid]; !ok {
			res[item.UserHid] = []wrappers.AdminUserGroupName{}
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
		tmp = wrappers.AdminUserGroupName{
			AdminUserGroup: item,
		}
		if _, ok := groupMap[item.GroupId]; ok {
			tmp.GroupName = groupMap[item.GroupId]
		}
		res[item.UserHid] = append(res[item.UserHid], tmp)
	}
	return
}

func (r *PermitService) leftAdminUser(list []models.AdminUser, dao *dao_impl.DaoPermit) (res []wrappers.ResultAdminUserList, err error) {
	res = []wrappers.ResultAdminUserList{}
	mapGroupPermit, err := r.getUserGroup(list, dao)
	if err != nil {
		return
	}
	var tmp wrappers.ResultAdminUserList
	for _, item := range list {
		tmp = wrappers.ResultAdminUserList{
			AdminUser: item,
		}
		if _, ok := mapGroupPermit[item.UserHid]; ok {
			tmp.Group = mapGroupPermit[item.UserHid]
		}
		res = append(res, tmp)
	}

	return

}
func (r *PermitService) Menu(arg *wrappers.ArgPermitMenu) (res *wrappers.ResultPermitMenu, err error) {
	res = &wrappers.ResultPermitMenu{}
	dao := dao_impl.NewDaoPermit(r.Context)
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
func (r *PermitService) getGroupMenu(dao *dao_impl.DaoPermit, menuIds ...int) (res []models.AdminMenu, err error) {
	res = []models.AdminMenu{}
	if len(menuIds) == 0 {
		return
	}
	res, err = dao.GetPermitMenuByIds(menuIds...)
	return
}
func (r *PermitService) getPermitByGroupIds(dao *dao_impl.DaoPermit, pathType string, groupIds ...int) (menuIds []int, err error) {
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

func (r *PermitService) getUserGroupIds(dao *dao_impl.DaoPermit, arg *wrappers.ArgPermitMenu) (res []int, err error) {
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

func (r *PermitService) Flag(arg *wrappers.ArgFlag) (res *wrappers.ResultFlag, err error) {
	res = &wrappers.ResultFlag{}
	return
}
