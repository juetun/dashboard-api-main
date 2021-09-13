package admins

import (
	"github.com/gin-gonic/gin"
)

type ConApi interface {
	Response(httpCode, errCode int, data gin.H)
}
