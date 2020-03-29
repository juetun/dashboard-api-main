/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-04-18
 * Time: 00:05
 */
package validate

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/pojos"
)

type CateStoreV struct {
}

func (cv *CateStoreV) MyValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := common.NewGin(c)
		var json pojos.CateStore
		// 接收各种参数
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusOK, base.Result{
				Code: 400001000,
				Data: nil,
				Msg:  err.Error()})
			c.Abort()
			return
		}
		if b := appG.Validate(&json); !b {
			c.Abort()
			return
		}
		c.Set("json", json)
		c.Next()
	}
}
