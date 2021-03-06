/**
* @Author:changjiang
* @Description:
* @File:user
* @Version: 1.0.0
* @Date 2020/4/19 8:02 下午
 */
package myplugins

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/app/middlewares"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/basic/web_middlerwares/permissions"
)

func PluginUser(arg *app_start.PluginsOperate) (err error) {
	middlewares.MiddleWareComponent = append(middlewares.MiddleWareComponent, func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}
		var res bool
		// apiG := common.NewGin(c)
		s := middlewares.GetRUri(c)
		// 如果是白名单的链接，则直接让过(用户不需要登录就让访问的URL)
		res = permissions.CheckWhite(c, s)
		if res {
			c.Next()
			return
		}
		app_obj.GetLog().Info(c, map[string]interface{}{
			"method":      "PluginUser",
			"info":        "user message",
			"router name": c.Request.RequestURI,
			"httpMethod":  c.Request.Method,
		})

		// 用户登录信息验证
		if exitStatus := middlewares.Auth(c); exitStatus {
			return
		}
		var needValidateUrl = false
		if needValidateUrl {
			res = RequestPathPermit(c, s)
			if !res {
				return
			}
		}

		// 获取当前登录用户信息
		// code, rd := middlewares.UserMessageSet(c, c.Request.RequestURI)
		// if code > 0 {
		// 	apiG.Response(code, rd)
		// 	return
		// }
	})
	return
}

func RequestPathPermit(c *gin.Context, s string) (res bool) {
	res = true
	// 用户登录了的验证权限
	res = permissions.CheckPermissions(c, s)
	// 如果不在白名单范围内，则让过
	if !res {
		app_obj.GetLog().Error(c, map[string]interface{}{
			"method":      "middleware.Permission",
			"info":        "router permission",
			"router name": c.Request.RequestURI,
			"httpMethod":  c.Request.Method,
		})
		obj := base.NewResult()
		obj.Code = http.StatusForbidden
		obj.Data = ""
		obj.Msg = fmt.Sprintf("您没有权限访问本功能(no auth:%s %s)", s, c.Request.Method)
		c.JSON(http.StatusOK, obj)
		c.Abort()
		res = false
		return
	}
	return
}
