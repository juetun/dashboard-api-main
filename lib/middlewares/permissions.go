/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-13
 * Time: 22:36
 */
package middlewares

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/app_log"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web"
)

func Permission(routerAsName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiG := common.NewGin(c)
		res := web.CheckPermissions(routerAsName, c.Request.Method)
		var log = app_log.GetLog()
		if !res {
			log.Error(map[string]string{
				"method":      "middleware.Permission",
				"info":        "router permission",
				"router name": routerAsName,
				"httpMethod":  c.Request.Method,
			})
			apiG.Response(400001005, nil)
			return
		}

		token := c.GetHeader("x-auth-token")
		if routerAsName == "console.post.imgUpload" {
			token = c.PostForm("upload-token")
		}

		if token == "" {
			log.Errorln("method", "middleware.Permission", "info", "token null")
			apiG.Response(400001005, nil)
			return
		}

		userId, err := common.ParseToken(token)
		if err != nil {
			log.Errorln("method", "middleware.Permission", "info", "parse token error")
			apiG.Response(400001005, nil)
			return
		}

		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			log.Errorln("method", "middleware.Permission", "info", "strconv token error")
			apiG.Response(400001005, nil)
			return
		}
		c.Set("userId", userIdInt)
		c.Set("token", token)
		// if routerAsName == "" {
		//	apiG.Response(http.StatusOK,0,nil)
		//	return
		// }
		c.Next()
	}
}
