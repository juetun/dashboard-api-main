// Package admins /**
package admins

import (
	"github.com/gin-gonic/gin"
)

type Export interface {
	// List 导出列表
	List(c *gin.Context)

	// Init 初始化导出任务
	Init(c *gin.Context)

	// Cancel 取消任务
	Cancel(c *gin.Context)

	// Progress 导出进度
	Progress(c *gin.Context)
}
