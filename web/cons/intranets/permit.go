package intranets

import (
	"github.com/gin-gonic/gin"
)

type ConPermitIntranet interface {
	GetImportPermit(c *gin.Context)

	GetUerImportPermit(c *gin.Context) // 判断用户是否有接口权限

	ValidateUserHavePermit(c *gin.Context) //验证用户是否有客服后台权限
}
