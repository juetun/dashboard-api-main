/**
* @Author:changjiang
* @Description:
* @File:gin
* @Version: 1.0.0
* @Date 2020/3/19 11:19 下午
 */
package common

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/base"
)

type validate interface {
	Message() map[string]int
}
type Gin struct {
	C *gin.Context
}

func NewGin(c *gin.Context) *Gin {
	return &Gin{C: c}
}

func (r *Gin) Response(code int, data interface{}) {
	var o base.ControllerBase
	o.Response(r.C, code, data)
}

func (r *Gin) Validate(obj validate) bool {
	// valid := validation.Validation{}
	// b, err := valid.Valid(obj)
	// var o base.ControllerBase
	// if err != nil {
	//
	// 	app_log.GetLog().Errorln("message", "valid error", "err", err.Error())
	// 	o.Response(r.C, 400000000, nil)
	// 	return false
	// }

	// if !b {
	// 	errorMaps := obj.Message()
	// 	field := valid.Errors[0].Key
	// 	if v, ok := errorMaps[field]; ok {
	// 		o.Response(r.C, v, nil)
	// 		return b
	// 	}
	// 	o.Response(r.C, 100000001, nil)
	// 	return b
	// }
	return true
}
