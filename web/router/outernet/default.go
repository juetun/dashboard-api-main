// Package outernets
/**
* @Author:changjiang
* @Description:
* @File:service
* @Version: 1.0.0
* @Date 2021/5/22 5:10 下午
 */
package outernet

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
)

func init() {
	app_start.HandleFuncOuterNet = append(app_start.HandleFuncOuterNet, func(r *gin.Engine, urlPrefix string) {
	})
}
