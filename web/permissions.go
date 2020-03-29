/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-13
 * Time: 22:39
 */
package web

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/app_log"
	"github.com/juetun/app-dashboard/lib/common"
)

type HttpPermit struct {
	Method []string `json:"method"`
	Uri    string   `json:"uri"`
}

var PermissionsWhite = []HttpPermit{
	{
		Method: []string{"GET", "POST"},
		Uri:    "console/login",
	},
	{
		Method: []string{"GET", "POST"},
		Uri:    "console/register",
	},
}

// 用户登录后 具备的接口访问权限
var Permissions = []HttpPermit{
	{
		Method: []string{"GET"},
		Uri:    "console/login",
	},
	{
		Method: []string{"GET"},
		Uri:    "console/home",
	},
	{
		Method: []string{"GET"},
		Uri:    "console/post/trash",
	},
}

// 需要验证权限的配置列表
// var Permissions = []string{
// 	"GET/console/login",
// 	"GET/console.post.index",
// 	"GET/console.post.create",
// 	"POST/console.post.store",
// 	"GET/console.post.edit",
// 	"PUT/console.post.update",
// 	"DELETE/console.post.destroy",
// 	"GET/console.post.trash",
// 	"POST/console.post.unTrash",
// 	"POST/console.post.imgUpload",
// 	"GET/console.cate.index",
// 	"GET/console.cate.edit",
// 	"PUT/console.cate.update",
// 	"POST/console.cate.store",
// 	"DELETE/console.cate.destroy",
// 	"GET/console.tag.index",
// 	"POST/console.tag.store",
// 	"GET/console.tag.edit",
// 	"PUT/console.tag.update",
// 	"DELETE/console.tag.destroy",
// 	"GET/console.system.index",
// 	"PUT/console.system.update",
// 	"GET/console.link.index",
// 	"POST/console.link.store",
// 	"GET/console.link.edit",
// 	"PUT/console.link.update",
// 	"DELETE/console.link.destroy",
// 	"DELETE/console.auth.logout",
// 	"GET/console.home.index",
// 	"DELETE/console.auth.cache",
// }

// 不需要验证权限的配置列表
// var PermissionsWhite = []string{
// 	"GET/console/login",
// }

func CheckPermissions(c *gin.Context) (res bool) {
	s := getRUri(c)
	app_log.GetLog().Error(map[string]string{
		"request_Uri": s,
		"info":        "web.permissions.go(CheckPermissions)",
		"router name": c.Request.RequestURI,
		"httpMethod":  c.Request.Method,
	})
	for _, v := range Permissions {
		if res = everyValidateTrueOrFalse(&v.Method, c.Request.Method, v.Uri, s); res {
			return
		}
	}
	return false
}

func everyValidateTrueOrFalse(methodArea *[]string, method, uri, s string) bool {
	var validateMethod bool
	if s == "default" { // 默认 default路径直接让过
		return true
	}
	validateMethod = false
	if len(*methodArea) != 0 {

		// 如果请求方法是返回内的值
		for _, value := range *methodArea {
			if value == method {
				validateMethod = true
			}
		}

		// 如果请求方法是返回内的值 并且请求地址对，则认为对
		if validateMethod == true && uri == s {
			return true
		}
	}

	// 否则 如果权限控制没有设置Method的值 就是表示所有的请求方式都有效，此时只验证URi路径是否正确
	if uri == s {
		return true
	}

	return false
}

// 白名单验证。此部分的接口用户不需要登录即可访问
func CheckWhite(c *gin.Context) (res bool) {
	s := getRUri(c)
	app_log.GetLog().Error(map[string]string{
		"request_Uri": s,
		"info":        "web.permissions.go(CheckWhite)",
		"router name": c.Request.RequestURI,
		"httpMethod":  c.Request.Method,
	})
	for _, v := range PermissionsWhite {
		if res = everyValidateTrueOrFalse(&v.Method, c.Request.Method, v.Uri, s); res {
			return
		}
	}
	return false
}
func getRUri(c *gin.Context) string {
	uri := strings.TrimLeft(c.Request.RequestURI, common.GetAppConfig().AppName+"/"+common.GetAppConfig().AppApiVersion)
	if uri == "" { // 如果是默认页 ，则直接让过
		return "default"
	}
	s1 := strings.Split(uri, "?")
	fmt.Printf("Uri is :'%v'", s1)
	return s1[0]
}
