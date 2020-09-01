/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-01-14
 * Time: 22:25
 */
package services

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/models"
)

type UserService struct {
	base.ServiceBase
}

func NewUserService(context ...*base.Context) (p *UserService) {
	p = &UserService{}
	p.SetContext(context...)
	return
}

func (r *UserService) GetUserById(userId string) (user *models.ZUsers, err error) {
	user = new(models.ZUsers)
	err = r.Context.Db.Table((&models.ZUsers{}).TableName()).
		Where("id=?", userId).
		Select("name,email").
		Find(user).
		Error
	if err != nil {
		r.Context.Log.Logger.Errorln("message", "service.GetUserById", "error", err.Error())
		return user, err
	}
	return user, nil
}

func (r *UserService) GetUserMapByIds(userId *[]string) (user *map[string]models.ZUsers, err error) {
	user = &map[string]models.ZUsers{}
	if len(*userId) == 0 {
		return
	}
	var users []models.ZUsers
	err = r.Context.Db.Table((&models.ZUsers{}).TableName()).
		Where("user_hid in (?)", *userId).
		Find(&users).
		Error
	if err != nil {
		r.Context.Log.Logger.Errorln("message", "service.GetUserMapByIds",
			"error",
			err.Error())
		return
	}
	for _, value := range users {
		(*user)[value.UserHid] = value
	}
	return
}

//
// func (r *UserService) UserCnt() (cnt int64, err error) {
// 	err = r.Context.Db.
// 		Table((&models.ZUsers{}).TableName()).
// 		Count(&cnt).
// 		Error
// 	return
// }
