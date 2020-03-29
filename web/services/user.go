/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-01-14
 * Time: 22:25
 */
package services

import (
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/web/models"
)

type UserService struct {
	base.ServiceBase
}

func NewUserService(context ...*base.Context) (p *UserService) {
	p = &UserService{}
	p.SetContext(context)
	return
}
func (r *UserService) GetUserById(userId int) (user *models.ZUsers, err error) {
	user = new(models.ZUsers)
	err = r.Context.Db.Table((&models.ZUsers{}).TableName()).
		Where("id=?", userId).
		Select("name,email").
		Find(user).
		Error
	if err != nil {
		r.Context.Log.Errorln("message", "service.GetUserById", "error", err.Error())
		return user, err
	}
	return user, nil
}
func (r *UserService) GetUserMapByIds(userId *[]int) (user *map[int]models.ZUsers, err error) {
	user = &map[int]models.ZUsers{}
	if len(*userId) == 0 {
		return
	}
	var users []models.ZUsers
	err = r.Context.Db.Table((&models.ZUsers{}).TableName()).
		Where("id in (?)", *userId).
		Find(&users).
		Error
	if err != nil {
		r.Context.Log.Errorln("message", "service.GetUserMapByIds",
			"error",
			err.Error())
		return
	}
	for _, value := range users {
		(*user)[value.Id] = value
	}
	return
}

func (r *UserService) UserCnt() (cnt int64, err error) {
	err = r.Context.Db.
		Table((&models.ZUsers{}).TableName()).
		Count(&cnt).
		Error
	return
}
