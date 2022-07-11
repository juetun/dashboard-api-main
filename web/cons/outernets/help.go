package outernets

import "github.com/gin-gonic/gin"

type (
	ConOuterNetsHelp interface {
		Tree(c *gin.Context)
	}
)
