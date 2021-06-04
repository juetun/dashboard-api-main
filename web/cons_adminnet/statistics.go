package cons_adminnet

import "github.com/gin-gonic/gin"

type Statistics interface {
	Index(*gin.Context)
}