package intranets

import "github.com/gin-gonic/gin"

type ConLogIntranet interface {
	AddLog(c *gin.Context)
}
