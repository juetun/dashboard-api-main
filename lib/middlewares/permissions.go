/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-13
 * Time: 22:36
 */
package middlewares

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/app_log"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/lib/utils"
	"github.com/juetun/app-dashboard/web"
)

func Permission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跨域配置
		if status := cors(c); status {
			return
		}

		var res bool
		apiG := common.NewGin(c)

		// 如果是白名单的链接，则直接让过
		res = web.CheckWhite(c)
		if res {
			c.Next()
			return
		}

		if exitStatus := auth(c); exitStatus {
			return
		}
		// 验证权限
		res = web.CheckPermissions(c)

		// 如果不在白名单范围内，则让过
		if !res {
			app_log.GetLog().Error(map[string]string{
				"method":      "middleware.Permission",
				"info":        "router permission",
				"router name": c.Request.RequestURI,
				"httpMethod":  c.Request.Method,
			})
			obj := utils.NewEmptyObject()
			c.JSON(http.StatusForbidden, obj)
			c.Abort()
			return
		}

		// 获取当前登录用户信息
		code, rd := userMessageSet(c, c.Request.RequestURI)
		if code > 0 {
			apiG.Response(code, rd)
			return
		}

		c.Next()
	}
}

// 用户登录逻辑处理
func auth(c *gin.Context) (exit bool) {
	token := c.Request.Header.Get("x-auth-token")
	if token == "" {
		msg := "token is null"
		app_log.GetLog().Error(map[string]string{
			"method": "zgh.ginmiddleware.auth",
			"error":  msg,
		})
		c.JSON(http.StatusUnauthorized, common.NewHttpResult().SetCode(http.StatusUnauthorized).SetMessage(msg))
		c.Abort()
		exit = true
		return
	}
	userId, err := common.ParseToken(token)
	if err != nil {
		app_log.GetLog().Error(map[string]string{
			"method": "zgh.ginmiddleware.auth",
			"token":  token,
			"error":  err.Error(),
		})
		c.JSON(http.StatusOK, common.NewHttpResult().SetCode(403).SetMessage(err.Error()))
		c.Abort()
		exit = true
		return
	}
	c.Set("userId", userId)
	return
}

func cors(c *gin.Context) (exitStatus bool) {
	fmt.Println("设置跨域")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE,PATCH")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Auth-Token, X-Auth-UUID, X-Auth-Openid, referrer, Authorization, x-client-id, x-client-version, x-client-type")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
		exitStatus = true
	}
	return
}

// 用户信息获取
func userMessageSet(c *gin.Context, routerAsName string) (code int, res interface{}) {
	token := c.GetHeader("x-auth-token")
	if routerAsName == "console.post.imgUpload" { // 如果是上传图片，则用的POST获取用户信息
		token = c.PostForm("upload-token")
	}

	if token == "" {
		app_log.GetLog().Errorln("method", "middleware.Permission", "info", "token null")
		return 400001005, nil
	}

	userId, err := common.ParseToken(token)
	if err != nil {
		app_log.GetLog().Errorln("method", "middleware.Permission", "info", "parse token error")
		return 400001005, nil
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		app_log.GetLog().Errorln("method", "middleware.Permission", "info", "strconv token error")
		return 400001005, nil
	}
	c.Set("userId", userIdInt)
	c.Set("token", token)
	return
}
