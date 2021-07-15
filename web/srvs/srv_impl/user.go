// Package srv_impl
/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-01-14
 * Time: 22:25
 */
package srv_impl

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/plugins/rpc"
	"github.com/juetun/dashboard-api-main/basic/const_obj"
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

// GetUserById 根据用户HID获取用户信息
func (r *UserService) GetUserById(userId string) (user *models.UserMain, err error) {
	user = &models.UserMain{}
	if userId == "" {
		return
	}

	users, err := r.GetUserByIds([]string{userId})
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"message": "service.GetUserById",
			"userId":  userId,
			"error":   err.Error(),
		})
		return
	}
	if len(users) > 0 {
		*user = users[0]
	}

	return
}

func (r *UserService) GetUserByIds(userId []string) (users []models.UserMain, err error) {
	if len(userId) == 0 {
		return
	}
	type ResultUserHttpRpc struct {
		Code int `json:"code"`
		Data struct {
			List []models.UserMain `json:"list"`
		}
		Message string `json:"message"`
	}
	var rpcUser ResultUserHttpRpc
	var httpHeader = http.Header{}
	httpHeader.Set(app_obj.HttpUserToken, r.Context.GinContext.GetHeader(app_obj.HttpUserToken))

	request := &rpc.RequestOptions{
		Context:     r.Context,
		Method:      "POST",
		AppName:     const_obj.MicroUser,
		PathVersion: "v1",
		Header:      httpHeader,
		URI:         "/user/list",
		Value:       url.Values{},
	}
	request.Value.Set("user_ids", strings.Join(userId, ","))
	var body string

	action := rpc.NewHttpRpc(request).
		Send().
		GetBody()
	body = action.GetBodyAsString()
	r.Context.Info(map[string]interface{}{
		"message": "service.GetUserMapByIds",
		"request": request,
		"body":    body,
	})

	if err = action.Bind(&rpcUser).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"message": "service.GetUserMapByIds",
			"error":   err.Error(),
		})
	}
	users = rpcUser.Data.List
	return
}

func (r *UserService) GetUserMapByIds(userId []string) (user *map[string]models.UserMain, err error) {
	user = &map[string]models.UserMain{}
	users, err := r.GetUserByIds(userId)
	for _, value := range users {
		(*user)[value.UserHid] = value
	}
	return
}
