package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func getRouteGroup(r *gin.Engine, urlPrefix, routerModule string, otherModule ...string) (g *gin.RouterGroup) {
	if len(otherModule) == 0 {
		g = r.Group(fmt.Sprintf("%s/%s", urlPrefix, routerModule))
	} else {
		g = r.Group(fmt.Sprintf("%s/%s_%s", urlPrefix, routerModule, strings.Join(otherModule, "/")))
	}
	return
}
