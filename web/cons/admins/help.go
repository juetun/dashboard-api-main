package admins

import "github.com/gin-gonic/gin"

//帮助文档接口
type ConsoleHelp interface {


	HelpList(c *gin.Context)

	HelpDetail(c *gin.Context)

	HelpEdit(c *gin.Context)

	HelpTree(c *gin.Context)

	TreeEditNode(c *gin.Context)

	OperateLog(c *gin.Context)
}
