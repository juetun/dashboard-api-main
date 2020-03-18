/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2018-12-16
 * Time: 00:28
 */
package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/app_log"
	"github.com/juetun/app-dashboard/lib/common"
	"net/http"
)

//登录信息验证
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//apiG := api.Gin{C: c}
		token := c.Request.Header.Get("x-auth-token")
		if token == "" {
			msg := "token is null"
			app_log.GetLog().Error(map[string]string{
				"method": "zgh.ginmiddleware.auth",
				"error":  msg,
			})
			c.JSON(http.StatusOK, common.NewHttpResult().SetCode(401).SetMessage(msg))
			c.Abort()
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
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}
