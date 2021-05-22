/**
* @Author:changjiang
* @Description:
* @File:service
* @Version: 1.0.0
* @Date 2021/5/22 5:17 下午
 */
package cons_outernet

import (
	"github.com/gin-gonic/gin"
)

type ConService interface {
	List(c *gin.Context)
	Edit(c *gin.Context)
}
