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
func (r *UserService) GetUserById(userId int) (*models.ZUsers, error) {
	user := new(models.ZUsers)
	_, err := r.Db.Id(userId).Cols("name", "email").Get(user)
	if err != nil {
		r.Log.Errorln("message", "service.GetUserById", "error", err.Error())
		return user, err
	}
	return user, nil
}

func (r *UserService) UserCnt() (cnt int64, err error) {
	user := new(models.ZUsers)
	cnt, err = r.Db.Count(user)
	return
}
