package outernet

import (
	"html/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
)

func init() {
	app_start.HandleFuncOuterNet = append(app_start.HandleFuncOuterNet, func(r *gin.Engine, urlPrefix string) {
		// web := web_impl.NewControllerWeb()
		r.SetFuncMap(template.FuncMap{
			"rem":    rem,
			"MDate":  mDate,
			"MDate2": mDate2,
		})
		// h := r.Group("web")
		// r.LoadHTMLGlob("template/home/*.tmpl")
		r.Static("/static/home", "./static/home")
		// h.GET("/", web.Index)
		// h.GET("/categories/:name", web.IndexCate)
		// h.GET("/tags/:name", web.IndexTag)
		// h.GET("/detail/:id", web.Detail)
		// h.GET("/archives", web.Archives)
		// h.GET("/404", web.NoFound)

	}, )
}

func rem(divisor int) bool {
	if (divisor+1)%4 == 0 {
		return true
	} else {
		return false
	}
}

func mDate(times time.Time) string {
	return times.Format("2006-01-02 15:04:05")
}

func mDate2(times time.Time) string {
	return times.Format("01-02")
}
