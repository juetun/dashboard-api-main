// Package srv_impl
/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-01-14
 * Time: 22:25
 */
package srv_impl

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/app_param"
	"github.com/juetun/base-wrapper/lib/plugins/rpc"
	"github.com/juetun/dashboard-api-main/pkg/parameters"
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
func (r *UserService) GetUserById(userId int64) (user *app_param.ResultUserItem, err error) {
	user = &app_param.ResultUserItem{}
	if userId == 0 {
		return
	}

	users, err := r.GetUserByIds([]int64{userId})
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"message": "service.GetUserById",
			"userId":  userId,
			"error":   err.Error(),
		})
		return
	}
	if dt, ok := users[userId]; ok {
		*user = dt
	}

	return
}

func (r *UserService) GetUserByIds(userId []int64) (users map[int64]app_param.ResultUserItem, err error) {
	if len(userId) == 0 {
		return
	}
	type ResultUserHttpRpc struct {
		Code int `json:"code"`
		Data struct {
			MapList map[int64]app_param.ResultUserItem `json:"list"`
		}
		Message string `json:"message"`
	}

	var httpHeader = http.Header{}
	httpHeader.Set(app_obj.HttpUserToken, r.Context.GinContext.GetHeader(app_obj.HttpUserToken))

	request := &rpc.RequestOptions{
		Context: r.Context,
		Method:  "POST",
		AppName: parameters.MicroUser,
		Header:  httpHeader,
		PathVersion: app_obj.App.AppRouterPrefix.Intranet,
		URI:     "/user/get_by_uid",
		Value:   url.Values{},
	}
	var userIdString = make([]string, 0, len(userId))
	for _, i2 := range userId {
		userIdString = append(userIdString, fmt.Sprintf("%d", i2))
	}
	request.Value.Set("user_hid", strings.Join(userIdString, ","))
	request.Value.Set("data_type", strings.Join([]string{app_param.UserDataTypeMain, app_param.UserDataTypeInfo}, ","))
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
	var rpcUser ResultUserHttpRpc
	if err = action.Bind(&rpcUser).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"message": "service.GetUserMapByIds",
			"error":   err.Error(),
		})
	}
	users = rpcUser.Data.MapList
	return
}

func (r *UserService) GetUserMapByIds(userId []int64) (user map[int64]app_param.ResultUserItem, err error) {
	user = map[int64]app_param.ResultUserItem{}
	if user, err = r.GetUserByIds(userId); err != nil {
		return
	}

	return
}
