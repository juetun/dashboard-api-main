package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var MiddleWareComponent = make([]gin.HandlerFunc, 0)

func LoadMiddleWare() {
	fmt.Println("Load gin middleWare start")
	MiddleWareComponent = append(MiddleWareComponent, Permission(), )
	fmt.Println("Load gin middleWare over")

}
