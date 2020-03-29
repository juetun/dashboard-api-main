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

func NewUserService() *UserService {
	return &UserService{}
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

func (r *UserService) UserCnt() (cnt int64, err error) {
	err = r.Context.Db.
		Table((&models.ZUsers{}).TableName()).
		Count(&cnt).
		Error
	return
}
