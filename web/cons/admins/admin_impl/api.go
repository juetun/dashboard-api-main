// Package admin_impl
/**
 * Created by GoLand.
 * User: zhu
 * Email: ylsc633@gmail.com
 * Date: 2019-05-17
 * Time: 11:18
 */
package admin_impl

import "C"
import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/dashboard-api-main/web"
)

type ApiController struct {
	C *gin.Context
}

func (r *ApiController) Response(httpCode, errCode int, data gin.H) {
	if data == nil {
		panic("常规信息应该设置")
	}
	msg := web.GetMsg(errCode)
	beginTime, _ := strconv.ParseInt(r.C.Writer.Header().Get("X-Begin-Time"), 10, 64)

	duration := time.Now().Sub(time.Unix(0, beginTime))
	milliseconds := float64(duration) / float64(time.Millisecond)
	rounded := float64(int(milliseconds*100+.5)) / 100
	roundedStr := fmt.Sprintf("%.3fms", rounded)
	r.C.Writer.Header().Set("X-Response-time", roundedStr)
	requestUrl := r.C.Request.URL.String()
	requestMethod := r.C.Request.Method
	app_obj.GetLog().Logger.Infoln("message", "Index Response", "Request Url",
		requestUrl, "Request method", requestMethod,
		"code", errCode, "errMsg", msg,
		"took", roundedStr)
	if errCode == 500 {
		r.C.HTML(http.StatusOK, "5xx.tmpl", data)
	} else if errCode == 404 {
		r.C.HTML(http.StatusOK, "4xx.tmpl", data)
	} else if errCode == 0 {
		r.C.HTML(http.StatusOK, "master.tmpl", data)
	} else {
		r.C.HTML(http.StatusOK, "5xx.tmpl", nil)
	}

	r.C.Abort()
	return
}
