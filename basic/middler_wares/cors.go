/**
* @Author:changjiang
* @Description:
* @File:cors
* @Version: 1.0.0
* @Date 2020/7/14 4:14 下午
 */
package middler_wares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)
func CorsComponent() gin.HandlerFunc {
	return func(c *gin.Context) {
		status := cors(c)
		if c.Request.Method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
			c.Next()
			return
		}
		// 跨域配置
		if status {
			return
		}
		c.Next()
	}
}

func cors(c *gin.Context) (exitStatus bool) {
 	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE,PATCH")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Auth-Token, X-Auth-UUID, X-Auth-Openid, referrer, Authorization, x-client-id, x-client-version, x-client-type")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	return
}
