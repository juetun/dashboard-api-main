/**
* @Author:changjiang
* @Description:
* @File:export
* @Version: 1.0.0
* @Date 2020/6/10 10:10 上午
 */
package cons_outernet

import (
	"github.com/gin-gonic/gin"
)

type Export interface {
	// 导出列表
	List(c *gin.Context)

	// 初始化导出任务
	Init(c *gin.Context)

	// 取消任务
	Cancel(c *gin.Context)

	// 导出进度
	Progress(c *gin.Context)
}
