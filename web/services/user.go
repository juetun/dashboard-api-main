/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-01-14
 * Time: 22:25
 */
package services

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/rpc"
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

// 根据用户HID获取用户信息
func (r *UserService) GetUserById(userId string) (user *models.Users, err error) {
	user = &models.Users{}
	if userId == "" {
		return
	}

	users, err := r.GetUserByIds([]string{userId})
	if err != nil {
		return
	}
	if len(users) > 0 {
		*user = users[0]
	}
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.GetUserById",
			"error":   err.Error(),
		})

	}
	return
}

func (r *UserService) GetUserByIds(userId []string) (users []models.Users, err error) {
	if len(userId) == 0 {
		return
	}
	type ResultUserHttpRpc struct {
		users []models.Users `json:"list"`
	}
	var rpcUser ResultUserHttpRpc
	request := &rpc.RequestOptions{}
	err = rpc.NewHttpRpc(request).
		Send().
		GetBody().
		Bind(&rpcUser).
		Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.GetUserMapByIds",
			"error":   err.Error(),
		})
	}
	users = rpcUser.users
	return
}

func (r *UserService) GetUserMapByIds(userId []string) (user *map[string]models.Users, err error) {

	users, err := r.GetUserByIds(userId)
	for _, value := range users {
		(*user)[value.UserHid] = value
	}
	return
}
