/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-09
 * Time: 21:34
 */
package auth_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/dashboard-api-main/web/controllers/inter"
	"github.com/juetun/dashboard-api-main/web/pojos"
	"github.com/juetun/dashboard-api-main/web/services"
)

type ControllerAuth struct {
	base.ControllerBase
}

func NewControllerAuth() inter.ConsoleAuth {
	controller := &ControllerAuth{}
	controller.ControllerBase.Init()
	return controller
}
func (c *ControllerAuth) Register(ctx *gin.Context) {
	srv := services.NewAuthService(&base.Context{Log: c.Log})
	data, err := srv.GetRegisterCappcha()
	if err != nil {
		c.Response(ctx, -1, nil, err.Error())
		return
	}
	c.Response(ctx, 0, data, "注册")
	return
}
func (c *ControllerAuth) AuthRegister(ctx *gin.Context) {
	srv := services.NewAuthService(&base.Context{Log: c.Log})
	requestJson, exists := ctx.Get("json")
	if !exists {
		c.Log.Error(map[string]string{
			"message": "auth.AuthRegister",
			"error":   "get request_params from context fail",
		})
		c.Response(ctx, 401000004, nil)
		return
	}
	ar, ok := requestJson.(pojos.AuthRegister)
	if !ok {
		c.Log.Error(map[string]string{
			"message": "auth.AuthRegister",
			"error":   "request_params turn to error",
		})
		c.Response(ctx, 400001001, nil)
		return
	}
	res, err := srv.AuthRegister(&ar)
	if err != nil {
		c.Response(ctx, -1, nil, err.Error())
		return
	}
	c.Response(ctx, 0, res, "注册成功")
	return
}
func (c *ControllerAuth) Login(ctx *gin.Context) {
	srv := services.NewAuthService(&base.Context{Log: c.Log})
	data, err := srv.Login()
	if err != nil {
		c.Response(ctx, 407000115, nil)
		return
	}
	c.Response(ctx, 0, data)
	return
}
func (c *ControllerAuth) AuthLogin(ctx *gin.Context) {
	requestJson, exists := ctx.Get("json")
	if !exists {
		c.Log.Error(map[string]string{
			"message": "auth.AuthLogin",
			"error":   "get request_params from context fail",
		})
		c.Response(ctx, 401000004, nil, "get request_params from context fail")
		return
	}
	al, ok := requestJson.(pojos.AuthLogin)
	if !ok {
		c.Log.Error(map[string]string{
			"message": "auth.AuthLogin",
			"error":   "request_params turn to error",
		})
		c.Response(ctx, 400001001, nil, "request_params turn to error", )
		return
	}

	srv := services.NewAuthService(&base.Context{Log: c.Log})
	_, token, err := srv.AuthLogin(&al)
	if err != nil {
		c.Response(ctx, -1, nil, err.Error())
		return
	}
	c.Response(ctx, 0, token, "登录成功")
	return
}
func (c *ControllerAuth) Logout(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist || token == "" {
		c.Log.Error(map[string]string{
			"message": "auth.Logout",
			"error":   "Can not get token",
		})
		c.Response(ctx, 400001004, nil)
		return
	}
	_, err := common.UnsetToken(token.(string))
	if err != nil {
		c.Log.Error(map[string]string{
			"message": "auth.Logout",
			"error":   "Can not get token",
		})
		c.Response(ctx, 407000014, nil)
		return
	}
	c.Response(ctx, 0, token)
	return
}
func (c *ControllerAuth) DelCache(ctx *gin.Context) {
	srv := services.NewAuthService(&base.Context{Log: c.Log})
	srv.DelAllCache()
	c.Response(ctx, 0, nil)
	return
}
