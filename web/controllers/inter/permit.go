/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 12:05 上午
 */
package inter

import (
	"github.com/gin-gonic/gin"
)

type Permit interface {
	Menu(c *gin.Context)
	Flag(c *gin.Context)
}
