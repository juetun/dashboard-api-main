package admins

import "github.com/gin-gonic/gin"

type Statistics interface {
	Index(*gin.Context)
}
